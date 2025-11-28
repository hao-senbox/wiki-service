package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	GetTemplate(ctx context.Context, typeParam string) (*entity.WikiTemplate, error)
	GetStatistics(ctx context.Context, page, limit int, typeParam, search string) ([]*response.WikiStatisticsResponse, error)
	GetWikiByCode(ctx context.Context, code string, language *int, typeParam string) (*response.WikiResponse, error)
	GetWikis(ctx context.Context, page, limit int, language *int, typeParam, search string) ([]*response.WikiResponse, int64, error)
	GetWikiByID(ctx context.Context, id string, language *int) (*response.WikiResponse, error)
	UpdateWiki(ctx context.Context, id string, req request.UpdateWikiRequest) error
}

type wikiUseCase struct {
	wikiRepo     repository.WikiRepository
	fileGateway  gateway.FileGateway
	userGateway  gateway.UserGateway
	mediaGateway gateway.MediaGateway
}

func NewWikiUseCase(
	wikiRepo repository.WikiRepository,
	fileGateway gateway.FileGateway,
	userGateway gateway.UserGateway,
	mediaGateway gateway.MediaGateway,
) WikiUseCase {
	return &wikiUseCase{
		wikiRepo:     wikiRepo,
		fileGateway:  fileGateway,
		userGateway:  userGateway,
		mediaGateway: mediaGateway,
	}
}

func (u *wikiUseCase) CreateWikiTemplate(ctx context.Context, req request.CreateWikiTemplateRequest, userID string) error {
	if userID == "" {
		return errors.New("userID is required")
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
		Type:      req.Type,
		Elements:  templateElements,
		CreatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.wikiRepo.CreateTemplate(ctx, template); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	// Create 6000 wiki instances
	wikis := make([]entity.Wiki, 6000)

	for i := 0; i < 6000; i++ {
		code := fmt.Sprintf("%04d", i+1)
		wikis[i] = entity.Wiki{
			Type:   req.Type,
			Code:   code,
			Public: 1,
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
			ImageWiki: "",
			CreatedBy: userID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	err := u.wikiRepo.CreateMany(ctx, wikis, req.Type)
	if err != nil {
		fmt.Println("InsertMany error:", err)
		return err
	}

	return nil
}

func (u *wikiUseCase) GetTemplate(ctx context.Context, typeParam string) (*entity.WikiTemplate, error) {
	if typeParam == "" {
		return nil, errors.New("type is required")
	}

	return u.wikiRepo.GetTemplates(ctx, typeParam)
}

func (u *wikiUseCase) GetStatistics(ctx context.Context, page, limit int, typeParam, search string) ([]*response.WikiStatisticsResponse, error) {
	if page < 1 {
		return nil, errors.New("page must be greater than 0")
	}

	if limit < 1 {
		return nil, errors.New("limit must be greater than 0")
	}

	wikis, total, err := u.wikiRepo.GetWikis(ctx, page, limit, typeParam, search)
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
			if _, err := fmt.Sscanf(parts[0], "%d", &number); err != nil {
				continue // Skip invalid element
			}

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

func (u *wikiUseCase) GetWikiByCode(ctx context.Context, code string, language *int, typeParam string) (*response.WikiResponse, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}

	if typeParam == "" {
		return nil, errors.New("type is required")
	}

	wiki, err := u.wikiRepo.GetWikiByCode(ctx, code, typeParam)
	if err != nil {
		return nil, err
	}

	if wiki == nil {
		return nil, errors.New("wiki not found")
	}

	if language != nil {
		filterTranslations([]*entity.Wiki{wiki}, language)
	}

	// Get user info for created_by
	var createdByUser *response.CreatedByUserInfo
	if user, err := u.userGateway.GetCurrentUser(ctx); err == nil && user != nil {
		createdByUser = &response.CreatedByUserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Fullname: user.Fullname,
			Email:    user.Email,
			Avatar:   user.AvatarURL,
		}
	}

	return mapper.WikiToResponse(ctx, wiki, u.fileGateway, u.mediaGateway, createdByUser), nil

}

func (u *wikiUseCase) GetWikis(ctx context.Context, page, limit int, language *int, typeParam, search string) ([]*response.WikiResponse, int64, error) {
	if typeParam == "" {
		return nil, 0, errors.New("type is required")
	}

	if page < 1 {
		return nil, 0, errors.New("page must be greater than 0")
	}

	if limit < 1 {
		return nil, 0, errors.New("limit must be greater than 0")
	}

	wikis, total, err := u.wikiRepo.GetWikis(ctx, page, limit, typeParam, search)
	if err != nil {
		return nil, 0, err
	}

	if language != nil {
		filterTranslations(wikis, language)
	}

	// Get current user info once for all wikis
	var currentUser *response.CreatedByUserInfo
	if user, err := u.userGateway.GetCurrentUser(ctx); err == nil && user != nil {
		currentUser = &response.CreatedByUserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Fullname: user.Fullname,
			Email:    user.Email,
			Avatar:   user.AvatarURL,
		}
	}

	// Convert wikis to responses
	responses := make([]*response.WikiResponse, len(wikis))
	for i, wiki := range wikis {
		responses[i] = mapper.WikiToResponse(ctx, wiki, u.fileGateway, u.mediaGateway, currentUser)
	}

	return responses, total, nil
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

	// Get user info for created_by
	var createdByUser *response.CreatedByUserInfo
	if user, err := u.userGateway.GetCurrentUser(ctx); err == nil && user != nil {
		createdByUser = &response.CreatedByUserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Fullname: user.Fullname,
			Email:    user.Email,
			Avatar:   user.AvatarURL,
		}
	}

	return mapper.WikiToResponse(ctx, wiki, u.fileGateway, u.mediaGateway, createdByUser), nil
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

	if req.ImageWiki != nil {
		wiki.ImageWiki = *req.ImageWiki
	}

	if req.Public != nil {
		wiki.Public = *req.Public
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
        elem := &translation.Elements[i]
        existingElements[elem.Number] = elem
    }

    // Create a map to track which elements are being updated
    updatedNumbers := make(map[int]bool)

    // Process each element from request
    for _, reqElem := range reqElements {
        existingElem, exists := existingElements[reqElem.Number]
        if exists {
            // Element exists, update it
            if strings.EqualFold(reqElem.Type, "picture") {
                if len(reqElem.PictureKeys) > 0 {
                    keysChanged := !u.pictureKeysEqual(existingElem.PictureKeys, reqElem.PictureKeys)
                    if keysChanged {
                        for _, oldKey := range existingElem.PictureKeys {
                            if err := u.fileGateway.DeleteImage(ctx, oldKey); err != nil {
                                log.Printf("failed to delete old picture: %v", err)
                            }
                        }
                    }
                    existingElem.PictureKeys = reqElem.PictureKeys
                }
            } else {
                newValue := reqElem.Value
                if newValue != nil && *newValue != "" {
                    if existingElem.Value != nil && *existingElem.Value != *newValue {
                        // Delete old file based on type
                        if err := u.deleteOldFile(ctx, existingElem); err != nil {
                            log.Printf("failed to delete old file: %v", err)
                        }
                    }
                }
                existingElem.Value = reqElem.Value
            }

            existingElem.Type = reqElem.Type
            if reqElem.VideoID != nil {
                existingElem.VideoID = reqElem.VideoID
            }
        } else {
            // Element doesn't exist, add new one
            newElem := entity.Element{
                Number: reqElem.Number,
                Type:   reqElem.Type,
            }
            if strings.EqualFold(reqElem.Type, "picture") {
                if len(reqElem.PictureKeys) > 0 {
                    newElem.PictureKeys = reqElem.PictureKeys
                    newElem.Value = &reqElem.PictureKeys[0]
                }
            } else {
                if reqElem.Value != nil {
                    newElem.Value = reqElem.Value
                }
            }
            if reqElem.VideoID != nil {
                newElem.VideoID = reqElem.VideoID
            }
            translation.Elements = append(translation.Elements, newElem)
        }
        updatedNumbers[reqElem.Number] = true
    }

    // Remove elements that are not in the request
    var newElements []entity.Element
    for i := range translation.Elements { // SỬ DỤNG INDEX
        existingElem := &translation.Elements[i] // LẤY POINTER
        if updatedNumbers[existingElem.Number] {
            newElements = append(newElements, *existingElem) // DEREFERENCE để append giá trị đã update
        } else {
            // Delete files for removed elements
            if err := u.deleteElementFiles(ctx, existingElem); err != nil {
                log.Printf("failed to delete removed element files: %v", err)
            }
        }
    }
    translation.Elements = newElements

    return nil
}

// Helper function to delete old file
func (u *wikiUseCase) deleteOldFile(ctx context.Context, elem *entity.Element) error {
    if elem.Value == nil || *elem.Value == "" {
        return nil
    }
    
    elemType := strings.ToLower(elem.Type)
    switch elemType {
    case "banner", "large_picture", "graphic", "linked_in":
        return u.fileGateway.DeleteImage(ctx, *elem.Value)
    default:
        return u.fileGateway.DeletePDF(ctx, *elem.Value)
    }
}

// Helper function to delete element files
func (u *wikiUseCase) deleteElementFiles(ctx context.Context, elem *entity.Element) error {
    if strings.EqualFold(elem.Type, "picture") {
        for _, key := range elem.PictureKeys {
            if key != "" {
                if err := u.fileGateway.DeleteImage(ctx, key); err != nil {
                    return err
                }
            }
        }
    } else if elem.Value != nil && *elem.Value != "" {
        return u.deleteOldFile(ctx, elem)
    }
    return nil
}

func convertElements(reqElements []request.Element, includeValues bool) []entity.Element {
	elements := make([]entity.Element, len(reqElements))
	for i, elem := range reqElements {
		var value *string
		var pictureKeys []string

		if includeValues {
			if elem.Type == "picture" {
				// For picture type, use the first picture key as value for backward compatibility
				// and store all keys in PictureKeys
				if len(elem.PictureKeys) > 0 {
					value = &elem.PictureKeys[0] // First key as main value
					pictureKeys = elem.PictureKeys
				}
			} else {
				value = elem.Value
			}
		}

		elements[i] = entity.Element{
			Number:      elem.Number,
			Type:        elem.Type,
			Value:       value,
			PictureKeys: pictureKeys,
			VideoID:     elem.VideoID,
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
			if translation.Language == nil {
				continue
			}
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

// pictureKeysEqual compares two slices of picture keys
func (u *wikiUseCase) pictureKeysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if i >= len(b) || v != b[i] {
			return false
		}
	}
	return true
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
