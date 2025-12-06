package mapper

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"wiki-service/internal/domain/entity"
	"wiki-service/internal/interface/http/dto/response.go"
	"wiki-service/pkg/gateway"
	file_gateway_dto "wiki-service/pkg/gateway/dto/file"
	media_gateway_dto "wiki-service/pkg/gateway/dto/media"
	libs_constant "wiki-service/pkg/libs/constant"
)

func WikiToResponse(
	ctx context.Context,
	wiki *entity.Wiki,
	fileGateway gateway.FileGateway,
	mediaGateway gateway.MediaGateway,
	createdByUser *response.CreatedByUserInfo,
) *response.WikiResponse {
	if wiki == nil {
		return nil
	}

	var imageWiki string
	if wiki.ImageWiki != "" {
		url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
			Key:  wiki.ImageWiki,
			Mode: string(libs_constant.ImageModePublic),
		})
		if err == nil && url != nil {
			imageWiki = *url
		}
	}

	resp := &response.WikiResponse{
		ID:            wiki.ID.Hex(),
		Code:          wiki.Code,
		ImageWiki:     imageWiki,
		Public:        wiki.Public,
		CreatedByUser: createdByUser,
		CreatedAt:     wiki.CreatedAt,
		UpdatedAt:     wiki.UpdatedAt,
	}

	resp.Translation = make([]response.TranslationResponse, 0, len(wiki.Translation))
	for _, tran := range wiki.Translation {
		elements := make([]response.ElementResponse, 0, len(tran.Elements))
		for _, elem := range tran.Elements {
			value := elem.Value // giữ nguyên từ DB
			var imageUrl *string
			var pdfUrl *string

			// Xử lý ảnh / PDF để gán URL hiển thị
			var valueJson *string
			if value != nil && fileGateway != nil {
				switch strings.ToLower(elem.Type) {
				case "large_picture", "banner", "linked_in", "graphic":
					url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
						Key:  *value,
						Mode: string(libs_constant.ImageModePublic),
					})
					if err == nil && url != nil {
						imageUrl = url
						// Tạo JSON string object cho banner và các type image
						jsonObj := map[string]string{
							"key_url":   *value,
							"image_url": *url,
						}
						jsonBytes, _ := json.Marshal(jsonObj)
						jsonStr := string(jsonBytes)
						valueJson = &jsonStr
					} else {
						log.Printf("failed to get image url: %v", err)
					}
				case "document":
					url, err := fileGateway.GetPDFUrl(ctx, file_gateway_dto.GetFileUrlRequest{
						Key:  *value,
						Mode: string(libs_constant.ImageModePublic),
					})
					if err == nil && url != nil {
						pdfUrl = url
					}
				default:
					// Cho các type khác, kiểm tra nếu value là JSON hợp lệ thì lưu vào value_json
					if strings.TrimSpace(*value) != "" && (strings.HasPrefix(*value, "{") || strings.HasPrefix(*value, "[")) {
						var temp interface{}
						if json.Unmarshal([]byte(*value), &temp) == nil {
							valueJson = value // Lưu toàn bộ JSON string
						}
					}
				}
			}

			// Xử lý button / button_url
			var btn *response.ButtonResponse
			var btnUrl *response.ButtonUrlResponse
			if value != nil && *value != "" {
				if strings.EqualFold(elem.Type, "button") {
					_ = json.Unmarshal([]byte(*value), &btn)
					if btn != nil && btn.ButtonIcon != "" {
						btnIconUrl, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
							Key:  btn.ButtonIcon,
							Mode: string(libs_constant.ImageModePublic),
						})
						if err == nil && btnIconUrl != nil {
							btn.ButtonIconUrl = *btnIconUrl
						}
					}
				} else if strings.EqualFold(elem.Type, "button_url") {
					_ = json.Unmarshal([]byte(*value), &btnUrl)
					if btnUrl != nil && btnUrl.ButtonIcon != "" {
						btnIconUrl, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
							Key:  btnUrl.ButtonIcon,
							Mode: string(libs_constant.ImageModePublic),
						})
						if err == nil && btnIconUrl != nil {
							btnUrl.ButtonIconUrl = *btnIconUrl
						}
					}
				}
			}

			var title *response.TitleResponse
			if value != nil && *value != "" {
				if strings.EqualFold(elem.Type, "title") {
					// Thử parse JSON trước
					title = &response.TitleResponse{}
					if err := json.Unmarshal([]byte(*value), title); err != nil {
						// Nếu không phải JSON, coi như plain string và tạo object
						title.Title = *value // Plain string làm title
						title.ImageKey = ""  // Không có image
						title.ImageUrl = ""  // Không có image URL
					} else {
						// Là JSON object, xử lý image như cũ
						if title.ImageKey != "" {
							url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
								Key:  title.ImageKey,
								Mode: string(libs_constant.ImageModePublic),
							})
							if err == nil && url != nil {
								title.ImageUrl = *url
							}
						}
					}
				}
			}
			// Xử lý picture_keys
			var pictureKeysUrl []response.PictureKeyUrl
			var sortedPictureKeys []response.PictureItem
			if strings.EqualFold(elem.Type, "picture") && len(elem.PictureKeys) > 0 && fileGateway != nil {
				// Sort by order trước khi xử lý
				sortedPictureKeys = make([]response.PictureItem, len(elem.PictureKeys))
				for i, item := range elem.PictureKeys {
					sortedPictureKeys[i] = response.PictureItem{
						Key:   item.Key,
						Order: item.Order,
						Title: item.Title,
					}
				}
				// Sort by Order field
				for i := 0; i < len(sortedPictureKeys)-1; i++ {
					for j := i + 1; j < len(sortedPictureKeys); j++ {
						if sortedPictureKeys[i].Order > sortedPictureKeys[j].Order {
							sortedPictureKeys[i], sortedPictureKeys[j] = sortedPictureKeys[j], sortedPictureKeys[i]
						}
					}
				}

				pictureKeysUrl = make([]response.PictureKeyUrl, len(sortedPictureKeys))
				pictureObjects := make([]map[string]interface{}, 0, len(sortedPictureKeys))
				for i, pictureItem := range sortedPictureKeys {
					key := pictureItem.Key
					if key != "" {
						url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
							Key:  key,
							Mode: string(libs_constant.ImageModePublic),
						})
						if err == nil && url != nil {
							pictureKeysUrl[i] = response.PictureKeyUrl{
								Order: pictureItem.Order,
								Url:   *url,
								Title: pictureItem.Title,
							}
							title := ""
							if pictureItem.Title != nil {
								title = *pictureItem.Title
							}
							pictureObjects = append(pictureObjects, map[string]interface{}{
								"key_url":   key,
								"image_url": *url,
								"title":     title,
								"order":     pictureItem.Order,
							})
						} else {
							pictureKeysUrl[i] = response.PictureKeyUrl{
								Order: pictureItem.Order,
								Url:   key, // fallback to key if no URL
								Title: pictureItem.Title,
							}
							title := ""
							if pictureItem.Title != nil {
								title = *pictureItem.Title
							}
							pictureObjects = append(pictureObjects, map[string]interface{}{
								"key_url":   key,
								"image_url": key, // fallback to key if no URL
								"title":     title,
								"order":     pictureItem.Order,
							})
						}
					} else {
						pictureKeysUrl[i] = response.PictureKeyUrl{
							Order: pictureItem.Order,
							Url:   key,
							Title: pictureItem.Title,
						}
						title := ""
						if pictureItem.Title != nil {
							title = *pictureItem.Title
						}
						pictureObjects = append(pictureObjects, map[string]interface{}{
							"key_url":   key,
							"image_url": key,
							"title":     title,
							"order":     pictureItem.Order,
						})
					}
				}
				// Tạo JSON string array của objects cho picture type
				if len(pictureObjects) > 0 {
					jsonBytes, _ := json.Marshal(pictureObjects)
					jsonStr := string(jsonBytes)
					valueJson = &jsonStr
				}
			}

			var videoUrl *string
			if elem.VideoID != nil {
				videoUrl, _ = mediaGateway.GetVideoUrl(ctx, media_gateway_dto.GetVideoUrlRequest{
					VideoID:  *elem.VideoID,
					Language: tran.Language,
				})
			}

			elements = append(elements, response.ElementResponse{
				Number:         elem.Number,
				Type:           elem.Type,
				Status:         elem.Status,
				Value:          value,     // giữ nguyên DB
				ValueJson:      valueJson, // JSON object chứa key và url
				ImageUrl:       imageUrl,  // URL hiển thị
				PdfUrl:         pdfUrl,    // URL PDF nếu có
				PictureKeys:    sortedPictureKeys,
				PictureKeysUrl: pictureKeysUrl,
				Title:          title,
				Button:         btn,
				ButtonUrl:      btnUrl,
				VideoID:        elem.VideoID,
				VideoUrl:       videoUrl,
			})
		}

		resp.Translation = append(resp.Translation, response.TranslationResponse{
			Language: tran.Language,
			Title:    tran.Title,
			Keywords: tran.Keywords,
			Level:    tran.Level,
			Unit:     tran.Unit,
			Elements: elements,
		})
	}

	return resp
}

func ElementsToResponse(elements []entity.Element) []response.ElementResponse {
	if elements == nil {
		return nil
	}

	resp := make([]response.ElementResponse, len(elements))
	for i, elem := range elements {
		// Convert entity.PictureItem to response.PictureItem and sort by Order
		var pictureKeys []response.PictureItem
		if len(elem.PictureKeys) > 0 {
			pictureKeys = make([]response.PictureItem, len(elem.PictureKeys))
			for j, item := range elem.PictureKeys {
				pictureKeys[j] = response.PictureItem{
					Key:   item.Key,
					Order: item.Order,
					Title: item.Title,
				}
			}
			// Sort by Order field
			for j := 0; j < len(pictureKeys)-1; j++ {
				for k := j + 1; k < len(pictureKeys); k++ {
					if pictureKeys[j].Order > pictureKeys[k].Order {
						pictureKeys[j], pictureKeys[k] = pictureKeys[k], pictureKeys[j]
					}
				}
			}
		}

		resp[i] = response.ElementResponse{
			Number:      elem.Number,
			Type:        elem.Type,
			Value:       elem.Value,
			PictureKeys: pictureKeys,
			VideoID:     elem.VideoID,
		}
	}

	return resp
}
