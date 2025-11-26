package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"wiki-service/internal/domain/entity"
	"wiki-service/internal/domain/repository"
	"wiki-service/internal/interface/http/dto/request"
	"wiki-service/internal/interface/http/dto/response.go"
	"wiki-service/internal/interface/http/mapper"
	"wiki-service/pkg/gateway"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WikiUseCase interface {
	CreateWikiTemplate(ctx context.Context, req request.CreateWikiTemplateRequest, userID string) error
	GetWikis(ctx context.Context, organizationID string, page, limit int, language *int, typeParam, search string) ([]*entity.Wiki, int64, error)
	GetWikiByID(ctx context.Context, id string, language *int) (*response.WikiResponse, error)
	GetWikiByCode(ctx context.Context, code string, language *int, organizationID string, typeParam string) (*response.WikiResponse, error)
	UpdateWiki(ctx context.Context, id string, req request.UpdateWikiRequest) error
	GetStatistics(ctx context.Context, organizationID string, page, limit int, typeParam, search string) ([]*response.WikiStatisticsResponse, error)
	GetTemplate(ctx context.Context, organizationID string, typeParam string) (*entity.WikiTemplate, error)
}

type wikiUseCase struct {
	wikiRepo    repository.WikiRepository
	fileGateway gateway.FileGateway
}

func NewWikiUseCase(
	wikiRepo repository.WikiRepository,
	fileGateway gateway.FileGateway,
) WikiUseCase {
	return &wikiUseCase{
		wikiRepo:    wikiRepo,
		fileGateway: fileGateway,
	}
}

func (u *wikiUseCase) CreateWikiTemplate(ctx context.Context, req request.CreateWikiTemplateRequest, userID string) error {
	if userID == "" {
		return errors.New("userID is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organizationID is required")
	}

	if len(req.Elements) == 0 {
		return errors.New("elements is required")
	}

	if err := validateElements(req.Elements); err != nil {
		return err
	}

	templateElements := convertElements(req.Elements, false)
	now := time.Now()

	// Save template first
	template := &entity.WikiTemplate{
		OrganizationID: req.OrganizationID,
		Type:           req.Type,
		Elements:       templateElements,
		CreatedBy:      userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.wikiRepo.CreateTemplate(ctx, template); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	// Create 6000 wiki instances
	wikis := make([]entity.Wiki, 6000)

	for i := 0; i < 6000; i++ {
		code := fmt.Sprintf("%04d", i+1)
		wikis[i] = entity.Wiki{
			OrganizationID: req.OrganizationID,
			Type:           req.Type,
			Code:           code,
			Translation: []entity.Translation{
				{
					Language: nil,
					Title:    nil,
					Keywords: nil,
					Level:    nil,
					Unit:     nil,
					Elements: cloneElements(templateElements),
				},
			},
			Icon:      "",
			CreatedBy: userID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	return u.wikiRepo.CreateMany(ctx, wikis, req.Type, req.OrganizationID)
}

func (u *wikiUseCase) GetWikis(ctx context.Context, organizationID string, page, limit int, language *int, typeParam, search string) ([]*entity.Wiki, int64, error) {
	if organizationID == "" {
		return nil, 0, errors.New("organizationID is required")
	}

	if typeParam == "" {
		return nil, 0, errors.New("type is required")
	}

	if page < 1 {
		return nil, 0, errors.New("page must be greater than 0")
	}

	if limit < 1 {
		return nil, 0, errors.New("limit must be greater than 0")
	}

	wikis, total, err := u.wikiRepo.GetWikis(ctx, organizationID, page, limit, typeParam, search)
	if err != nil {
		return nil, 0, err
	}

	if language != nil {
		filterTranslations(wikis, language)
	}

	return wikis, total, nil
}

func (u *wikiUseCase) GetWikiByID(ctx context.Context, id string, language *int) (*response.WikiResponse, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	wiki, err := u.wikiRepo.GetWikiByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if wiki == nil {
		return nil, errors.New("wiki not found")
	}

	if language != nil {
		filterTranslations([]*entity.Wiki{wiki}, language)
	}

	return mapper.WikiToResponse(ctx, wiki, u.fileGateway), nil
}

func (u *wikiUseCase) GetWikiByCode(ctx context.Context, code string, language *int, organizationID string, typeParam string) (*response.WikiResponse, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}

	if organizationID == "" {
		return nil, errors.New("organizationID is required")
	}

	if typeParam == "" {
		return nil, errors.New("type is required")
	}

	wiki, err := u.wikiRepo.GetWikiByCode(ctx, code, organizationID, typeParam)
	if err != nil {
		return nil, err
	}

	if wiki == nil {
		return nil, errors.New("wiki not found")
	}

	if language != nil {
		filterTranslations([]*entity.Wiki{wiki}, language)
	}

	return mapper.WikiToResponse(ctx, wiki, u.fileGateway), nil

}

func (u *wikiUseCase) GetStatistics(ctx context.Context, organizationID string, page, limit int, typeParam, search string) ([]*response.WikiStatisticsResponse, error) {
	if organizationID == "" {
		return nil, errors.New("organizationID is required")
	}

	if page < 1 {
		return nil, errors.New("page must be greater than 0")
	}

	if limit < 1 {
		return nil, errors.New("limit must be greater than 0")
	}

	wikis, total, err := u.wikiRepo.GetWikis(ctx, organizationID, page, limit, typeParam, search)
	if err != nil {
		return nil, err
	}

	if len(wikis) == 0 {
		return []*response.WikiStatisticsResponse{}, nil
	}

	var responses []*response.WikiStatisticsResponse

	for _, wiki := range wikis {
		if wiki == nil {
			continue
		}

		// Maps to track statistics for this wiki
		elementStats := make(map[string]string) // elementKey -> overall status
		languageStats := make(map[string]bool)  // language -> has any value

		// Process translations for this wiki
		for _, translation := range wiki.Translation {
			langStr := "unknown"
			if translation.Language != nil {
				langStr = fmt.Sprintf("%d", *translation.Language)
				languageStats[langStr] = false // Initialize as false, will be set to true if any element has value
			}

			for _, element := range translation.Elements {
				elementKey := fmt.Sprintf("%d_%s", element.Number, element.Type)

				// Check if element has value
				hasValue := element.Value != nil && strings.TrimSpace(*element.Value) != ""

				if hasValue {
					elementStats[elementKey] = "yes"
					if langStr != "unknown" {
						languageStats[langStr] = true
					}
				} else {
					// Only set to "no" if not already set to "yes"
					if elementStats[elementKey] != "yes" {
						elementStats[elementKey] = "no"
					}
				}
			}
		}

		// Convert language stats to string map
		langStatus := make(map[string]string)
		for lang, hasValue := range languageStats {
			if hasValue {
				langStatus[lang] = "yes"
			} else {
				langStatus[lang] = "no"
			}
		}

		// Convert element stats to response structure
		var elements []response.ElementStatistics
		for elementKey, status := range elementStats {
			// Parse element key
			parts := strings.Split(elementKey, "_")
			if len(parts) != 2 {
				continue
			}

			number := 0
			fmt.Sscanf(parts[0], "%d", &number)

			element := response.ElementStatistics{
				Number: number,
				Type:   parts[1],
				Check:  status,
			}
			elements = append(elements, element)
		}

		// Sort elements by number for consistent output
		for i := 0; i < len(elements)-1; i++ {
			for j := i + 1; j < len(elements); j++ {
				if elements[i].Number > elements[j].Number {
					elements[i], elements[j] = elements[j], elements[i]
				}
			}
		}

		// Create response for this wiki
		response := &response.WikiStatisticsResponse{
			Languages: langStatus,
			Code:      wiki.Code,
			Elements:  elements,
		}

		responses = append(responses, response)
	}

	// Add pagination info to the first response (or create a wrapper response)
	if len(responses) > 0 {
		totalPages := int((total + int64(limit) - 1) / int64(limit))
		responses[0].Pagination = &response.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		}
	}

	return responses, nil
}

func (u *wikiUseCase) UpdateWiki(ctx context.Context, id string, req request.UpdateWikiRequest) error {
	if id == "" {
		return errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	wiki, err := u.wikiRepo.GetWikiByID(ctx, objectID)
	if err != nil {
		return errors.New("wiki not found")
	}

	if wiki == nil {
		return errors.New("wiki not found")
	}

	if req.Icon != nil {
		wiki.Icon = *req.Icon
	}

	if req.Language != nil && *req.Language < 0 {
		return errors.New("language must be greater than or equal to 0")
	}

	var translation *entity.Translation
	for i := range wiki.Translation {
		if wiki.Translation[i].Language != nil && req.Language != nil &&
			*wiki.Translation[i].Language == *req.Language {
			translation = &wiki.Translation[i]
			break
		}
		if req.Language != nil && wiki.Translation[i].Language == nil {
			translation = &wiki.Translation[i]
			break
		}
	}

	if translation == nil {
		if len(req.Elements) == 0 {
			return errors.New("elements are required when creating a new translation")
		}
		if err := validateElements(req.Elements); err != nil {
			return err
		}
		wiki.Translation = append(wiki.Translation, entity.Translation{
			Language: req.Language,
		})
		translation = &wiki.Translation[len(wiki.Translation)-1]
	}

	if req.Language != nil {
		translation.Language = req.Language
	}

	if req.Title != nil {
		translation.Title = req.Title
	}

	if req.Keywords != nil {
		translation.Keywords = req.Keywords
	}

	if req.Level != nil {
		translation.Level = req.Level
	}

	if req.Unit != nil {
		translation.Unit = req.Unit
	}

	if len(req.Elements) > 0 {
		if err := validateElements(req.Elements); err != nil {
			return err
		}

		if err := u.mergeElements(ctx, translation, req.Elements); err != nil {
			return err
		}
	}

	wiki.UpdatedAt = time.Now()

	return u.wikiRepo.UpdateWiki(ctx, objectID, wiki)
}

func (u *wikiUseCase) mergeElements(ctx context.Context, translation *entity.Translation, reqElements []request.Element) error {
	// Create a map of existing elements by number for quick lookup
	existingElements := make(map[int]*entity.Element)
	for i := range translation.Elements {
		existingElements[translation.Elements[i].Number] = &translation.Elements[i]
	}

	// Create a map to track which elements are being updated
	updatedNumbers := make(map[int]bool)

	// Process each element from request
	for _, reqElem := range reqElements {
		existingElem, exists := existingElements[reqElem.Number]

		if exists {
			// Element exists, check if value changed
			newValue := reqElem.Value
			if newValue != nil && *newValue != "" {
				// If value changed, delete old image/file
				if existingElem.Value != nil && *existingElem.Value != *newValue {
					if strings.EqualFold(existingElem.Type, "image") {
						if err := u.fileGateway.DeleteImage(ctx, *existingElem.Value); err != nil {
							return fmt.Errorf("failed to delete old image: %w", err)
						}
					} else if strings.EqualFold(existingElem.Type, "file") {
						if err := u.fileGateway.DeletePDF(ctx, *existingElem.Value); err != nil {
							return fmt.Errorf("failed to delete old file: %w", err)
						}
					}
				}
			}

			// Update existing element
			existingElem.Type = reqElem.Type
			if reqElem.Value != nil {
				existingElem.Value = reqElem.Value
			}
			if reqElem.TopicID != nil {
				existingElem.TopicID = reqElem.TopicID
			}
		} else {
			// Element doesn't exist, add new one
			newElem := entity.Element{
				Number: reqElem.Number,
				Type:   reqElem.Type,
			}
			if reqElem.Value != nil {
				newElem.Value = reqElem.Value
			}
			if reqElem.TopicID != nil {
				newElem.TopicID = reqElem.TopicID
			}
			translation.Elements = append(translation.Elements, newElem)
		}

		updatedNumbers[reqElem.Number] = true
	}

	// Remove elements that are not in the request
	// Create new elements slice, only keeping elements that were in the request or updated
	var newElements []entity.Element
	for _, existingElem := range translation.Elements {
		if updatedNumbers[existingElem.Number] {
			// Keep elements that were updated
			newElements = append(newElements, existingElem)
		} else {
			// Delete image/file for elements that are being removed
			if existingElem.Value != nil && *existingElem.Value != "" {
				if strings.EqualFold(existingElem.Type, "image") {
					if err := u.fileGateway.DeleteImage(ctx, *existingElem.Value); err != nil {
						return fmt.Errorf("failed to delete removed image: %w", err)
					}
				} else if strings.EqualFold(existingElem.Type, "file") {
					if err := u.fileGateway.DeletePDF(ctx, *existingElem.Value); err != nil {
						return fmt.Errorf("failed to delete removed file: %w", err)
					}
				}
			}
			// Element is removed, don't add to newElements
		}
	}

	translation.Elements = newElements

	return nil
}

func convertElements(reqElements []request.Element, includeValues bool) []entity.Element {
	elements := make([]entity.Element, len(reqElements))
	for i, elem := range reqElements {
		var value *string
		if includeValues {
			value = elem.Value
		}

		elements[i] = entity.Element{
			Number:  elem.Number,
			Type:    elem.Type,
			Value:   value,
			TopicID: nil,
		}
	}

	return elements
}

func cloneElements(elements []entity.Element) []entity.Element {
	cloned := make([]entity.Element, len(elements))
	copy(cloned, elements)
	return cloned
}

func filterTranslations(wikis []*entity.Wiki, language *int) {
	for _, wiki := range wikis {
		if wiki == nil {
			continue
		}

		filtered := make([]entity.Translation, 0, 1)
		for _, translation := range wiki.Translation {
			if *translation.Language == *language {
				filtered = append(filtered, translation)
				break
			}
		}

		if len(filtered) == 0 {
			for _, translation := range wiki.Translation {
				if translation.Language == nil {
					filtered = append(filtered, translation)
					break
				}
			}
		}

		wiki.Translation = filtered
	}
}

func (u *wikiUseCase) GetTemplate(ctx context.Context, organizationID string, typeParam string) (*entity.WikiTemplate, error) {
	if organizationID == "" {
		return nil, errors.New("organizationID is required")
	}

	if typeParam == "" {
		return nil, errors.New("type is required")
	}

	return u.wikiRepo.GetTemplates(ctx, organizationID, typeParam)
}

func validateElements(elements []request.Element) error {
	for _, element := range elements {
		if element.Number == 0 {
			return errors.New("number is required")
		}
		if element.Type == "" {
			return errors.New("type is required")
		}
	}
	return nil
}
