package mapper

import (
	"context"
	"strings"
	"wiki-service/internal/domain/entity"
	"wiki-service/internal/interface/http/dto/response.go"
	"wiki-service/pkg/gateway"
	file_gateway_dto "wiki-service/pkg/gateway/dto/file"
	libs_constant "wiki-service/pkg/libs/constant"
)

func WikiToResponse(
	ctx context.Context,
	wiki *entity.Wiki,
	fileGateway gateway.FileGateway,
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
			value := elem.Value
			// topicID := elem.TopicID
			if value != nil &&
				strings.EqualFold(elem.Type, "picture") &&
				fileGateway != nil &&
				*value != "" {
				url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
					Key:  *value,
					Mode: string(libs_constant.ImageModePublic),
				})
				if err == nil && url != nil {
					value = url
				}
			}

			if value != nil &&
				strings.EqualFold(elem.Type, "large_picture") &&
				fileGateway != nil &&
				*value != "" {
				url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
					Key:  *value,
					Mode: string(libs_constant.ImageModePublic),
				})
				if err == nil && url != nil {
					value = url
				}
			}

			if value != nil &&
				strings.EqualFold(elem.Type, "banner") &&
				fileGateway != nil &&
				*value != "" {
				url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
					Key:  *value,
					Mode: string(libs_constant.ImageModePublic),
				})
				if err == nil && url != nil {
					value = url
				}
			}

			if value != nil &&
				strings.EqualFold(elem.Type, "linked_in") &&
				fileGateway != nil &&
				*value != "" {
				url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
					Key:  *value,
					Mode: string(libs_constant.ImageModePublic),
				})
				if err == nil && url != nil {
					value = url
				}
			}

			if value != nil &&
				strings.EqualFold(elem.Type, "graphic") &&
				fileGateway != nil &&
				*value != "" {
				url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
					Key:  *value,
					Mode: string(libs_constant.ImageModePublic),
				})
				if err == nil && url != nil {
					value = url
				}
			}

			if value != nil &&
				strings.EqualFold(elem.Type, "document") &&
				fileGateway != nil &&
				*value != "" {
				url, err := fileGateway.GetPDFUrl(ctx, file_gateway_dto.GetFileUrlRequest{
					Key:  *value,
					Mode: string(libs_constant.ImageModePublic),
				})
				if err == nil && url != nil {
					value = url
				}
			}

			// if value != nil &&
			// 	strings.EqualFold(elem.Type, "button_url") &&
			// 	fileGateway != nil &&
			// 	*value != "" {
			// 	url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
			// 		Key:  *value,
			// 		Mode: string(libs_constant.ImageModePublic),
			// 	})
			// }

			// if topicID != nil &&
			// 	strings.EqualFold(elem.Type, "video") &&
			// 	fileGateway != nil &&
			// 	*topicID != "" {
			// 	url, err := fileGateway.GetVideoUrl(ctx, file_gateway_dto.GetFileUrlRequest{
			// 		Key:  *topicID,
			// 		Mode: string(libs_constant.ImageModePublic),
			// 	})
			// }

			// Handle picture keys for picture type
			var pictureKeys []string
			if strings.EqualFold(elem.Type, "picture") && len(elem.PictureKeys) > 0 {
				pictureKeys = make([]string, len(elem.PictureKeys))
				for i, key := range elem.PictureKeys {
					if key != "" && fileGateway != nil {
						url, err := fileGateway.GetImageUrl(ctx, file_gateway_dto.GetFileUrlRequest{
							Key:  key,
							Mode: string(libs_constant.ImageModePublic),
						})
						if err == nil && url != nil {
							pictureKeys[i] = *url
						} else {
							pictureKeys[i] = key // fallback to original key if URL generation fails
						}
					} else {
						pictureKeys[i] = key
					}
				}
			}

			elements = append(elements, response.ElementResponse{
				Number:      elem.Number,
				Type:        elem.Type,
				Value:       value,
				PictureKeys: pictureKeys,
				TopicID:     elem.TopicID,
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

func WikiTemplateToResponse(template *entity.WikiTemplate) *response.WikiTemplateResponse {
	if template == nil {
		return nil
	}

	elements := make([]response.ElementResponse, 0, len(template.Elements))
	for _, elem := range template.Elements {
		elements = append(elements, response.ElementResponse{
			Number:  elem.Number,
			Type:    elem.Type,
			Value:   elem.Value,
			TopicID: elem.TopicID,
		})
	}

	return &response.WikiTemplateResponse{
		ID:        template.ID.Hex(),
		Type:      template.Type,
		Elements:  elements,
		CreatedBy: template.CreatedBy,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}
}
