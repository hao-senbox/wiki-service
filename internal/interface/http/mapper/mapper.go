package mapper

import (
	"context"
	"encoding/json"
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
			if value != nil && fileGateway != nil {
				switch strings.ToLower(elem.Type) {
				case "large_picture", "banner", "linked_in", "graphic":
					url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
						Key:  *value,
						Mode: string(libs_constant.ImageModePublic),
					})
					if err == nil && url != nil {
						imageUrl = url
					}
				case "document":
					url, err := fileGateway.GetPDFUrl(ctx, file_gateway_dto.GetFileUrlRequest{
						Key:  *value,
						Mode: string(libs_constant.ImageModePublic),
					})
					if err == nil && url != nil {
						pdfUrl = url
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

			// Xử lý picture_keys
			var pictureKeysUrl []string
			if strings.EqualFold(elem.Type, "picture") && len(elem.PictureKeys) > 0 && fileGateway != nil {
				pictureKeysUrl = make([]string, len(elem.PictureKeys))
				for i, key := range elem.PictureKeys {
					if key != "" {
						url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
							Key:  key,
							Mode: string(libs_constant.ImageModePublic),
						})
						if err == nil && url != nil {
							pictureKeysUrl[i] = *url
						} else {
							pictureKeysUrl[i] = key
						}
					} else {
						pictureKeysUrl[i] = key
					}
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
				Value:          value,    // giữ nguyên DB
				ImageUrl:       imageUrl, // URL hiển thị
				PdfUrl:         pdfUrl,   // URL PDF nếu có
				PictureKeys:    elem.PictureKeys,
				PictureKeysUrl: pictureKeysUrl,
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
		resp[i] = response.ElementResponse{
			Number:      elem.Number,
			Type:        elem.Type,
			Value:       elem.Value,
			PictureKeys: elem.PictureKeys,
			VideoID:     elem.VideoID,
		}
	}

	return resp
}
