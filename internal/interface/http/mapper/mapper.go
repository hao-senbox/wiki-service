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
) *response.WikiResponse {
	if wiki == nil {
		return nil
	}

	resp := &response.WikiResponse{
		ID:             wiki.ID.Hex(),
		Code:           wiki.Code,
		Icon:           wiki.Icon,
		CreatedBy:      wiki.CreatedBy,
		CreatedAt:      wiki.CreatedAt,
		UpdatedAt:      wiki.UpdatedAt,
	}

	resp.Translation = make([]response.TranslationResponse, 0, len(wiki.Translation))
	for _, tran := range wiki.Translation {
		elements := make([]response.ElementResponse, 0, len(tran.Elements))
		for _, elem := range tran.Elements {
			value := elem.Value
			if value != nil &&
				strings.EqualFold(elem.Type, "image") &&
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
				strings.EqualFold(elem.Type, "file") &&
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

			elements = append(elements, response.ElementResponse{
				Number:  elem.Number,
				Type:    elem.Type,
				Value:   value,
				TopicID: elem.TopicID,
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
		ID:             template.ID.Hex(),
		Type:           template.Type,
		Elements:       elements,
		CreatedBy:      template.CreatedBy,
		CreatedAt:      template.CreatedAt,
		UpdatedAt:      template.UpdatedAt,
	}
}
