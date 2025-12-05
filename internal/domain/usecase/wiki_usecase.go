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

	templateWiki, err := u.wikiRepo.GetTemplates(ctx, "wiki_web")
	if err != nil {
		return nil, err
	}

	if templateWiki == nil {
		return nil, errors.New("template wiki not found")
	}

	if language != nil {
		filterTranslations([]*entity.Wiki{wiki}, language, templateWiki.Elements)
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

	templateWiki, err := u.wikiRepo.GetTemplates(ctx, "wiki_web")
	if err != nil {
		return nil, 0, err
	}

	if templateWiki == nil {
		return nil, 0, errors.New("template wiki not found")
	}

	if language != nil {
		filterTranslations(wikis, language, templateWiki.Elements)
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

	templateWiki, err := u.wikiRepo.GetTemplates(ctx, "wiki_web")
	if err != nil {
		return nil, err
	}

	if templateWiki == nil {
		return nil, errors.New("template wiki not found")
	}

	if language != nil {
		filterTranslations([]*entity.Wiki{wiki}, language, templateWiki.Elements)
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

	wikiRes := mapper.WikiToResponse(ctx, wiki, u.fileGateway, u.mediaGateway, createdByUser)

	return wikiRes, nil
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

func convertElements(reqElements []request.Element, includeValues bool) []entity.Element {
	elements := make([]entity.Element, len(reqElements))
	for i, elem := range reqElements {
		var value *string
		var pictureKeys []entity.PictureItem

		if includeValues {
			if elem.Type == "picture" {
				// For picture type, use the first picture key as value for backward compatibility
				// and store all PictureItems in PictureKeys
				if len(elem.PictureKeys) > 0 {
					value = &elem.PictureKeys[0].Key // First key as main value
					// Convert request.PictureItem to entity.PictureItem
					pictureKeys = make([]entity.PictureItem, len(elem.PictureKeys))
					for j, reqItem := range elem.PictureKeys {
						pictureKeys[j] = entity.PictureItem{
							Key:   reqItem.Key,
							Order: reqItem.Order,
							Title: reqItem.Title,
						}
					}
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

func filterTranslations(wikis []*entity.Wiki, language *int, templateElements []entity.Element) {
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
			filtered = append(filtered, entity.Translation{
				Language: language,
				Title:    nil,
				Keywords: nil,
				Level:    nil,
				Unit:     nil,
				Elements: templateElements,
			})
		}

		wiki.Translation = filtered
	}
}

// pictureKeysEqual compares two slices of picture keys (exact match)
func (u *wikiUseCase) pictureKeysEqual(a, b []entity.PictureItem) bool {
	if len(a) != len(b) {
		return false
	}

	// Sort both slices by Order before comparing
	sortedA := make([]entity.PictureItem, len(a))
	copy(sortedA, a)
	for i := 0; i < len(sortedA)-1; i++ {
		for j := i + 1; j < len(sortedA); j++ {
			if sortedA[i].Order > sortedA[j].Order {
				sortedA[i], sortedA[j] = sortedA[j], sortedA[i]
			}
		}
	}

	sortedB := make([]entity.PictureItem, len(b))
	copy(sortedB, b)
	for i := 0; i < len(sortedB)-1; i++ {
		for j := i + 1; j < len(sortedB); j++ {
			if sortedB[i].Order > sortedB[j].Order {
				sortedB[i], sortedB[j] = sortedB[j], sortedB[i]
			}
		}
	}

	for i, v := range sortedA {
		if i >= len(sortedB) || v.Key != sortedB[i].Key || v.Order != sortedB[i].Order {
			return false
		}
	}
	return true
}

// getKeysToDelete returns keys that exist in current DB but not in new request
func (u *wikiUseCase) getKeysToDelete(currentKeys []entity.PictureItem, newKeys []entity.PictureItem) []string {
	currentKeyMap := make(map[string]bool)
	for _, item := range currentKeys {
		currentKeyMap[item.Key] = true
	}

	var keysToDelete []string
	for _, item := range newKeys {
		if currentKeyMap[item.Key] {
			delete(currentKeyMap, item.Key)
		}
	}

	for key := range currentKeyMap {
		keysToDelete = append(keysToDelete, key)
	}

	return keysToDelete
}

// func validateElements(elements []request.Element) error {
// 	for _, element := range elements {
// 		if element.Number == 0 {
// 			return errors.New("number is required")
// 		}
// 		if element.Type == "" {
// 			return errors.New("type is required")
// 		}
// 	}
// 	return nil
// }

func validateElements(elements []request.Element) error {
	numberSet := make(map[int]bool)

	for _, element := range elements {
		if element.Number <= 0 {
			return fmt.Errorf("element number must be positive, got: %d", element.Number)
		}
		if element.Type == "" {
			return errors.New("type is required")
		}
		if numberSet[element.Number] {
			return fmt.Errorf("duplicate element number: %d", element.Number)
		}
		numberSet[element.Number] = true
	}
	return nil
}

func (u *wikiUseCase) mergeElements(ctx context.Context, translation *entity.Translation, reqElements []request.Element) error {
	// PHASE 1: Collect all file keys being used in the request
	// This ensures we never delete files that are still in use (even if repositioned)
	requestFileKeys := make(map[string]bool)
	for _, reqElem := range reqElements {
		// Collect single file keys (banner, graphic, etc.)
		if reqElem.Value != nil && *reqElem.Value != "" {
			requestFileKeys[*reqElem.Value] = true
		}
		// Collect picture keys
		for _, picItem := range reqElem.PictureKeys {
			if picItem.Key != "" {
				requestFileKeys[picItem.Key] = true
			}
		}
	}

	// PHASE 2: Collect all existing file keys (for cleanup later)
	existingFileKeys := make(map[string]bool)
	for _, elem := range translation.Elements {
		if elem.Value != nil && *elem.Value != "" {
			existingFileKeys[*elem.Value] = true
		}
		for _, picItem := range elem.PictureKeys {
			if picItem.Key != "" {
				existingFileKeys[picItem.Key] = true
			}
		}
	}

	// PHASE 3: Build new elements array from request (completely replace old array)
	// This ensures no elements are lost during position changes
	newElements := make([]entity.Element, len(reqElements))
	for i, reqElem := range reqElements {
		newElem := entity.Element{
			Number: reqElem.Number,
			Type:   reqElem.Type,
		}

		// Handle picture type
		if strings.EqualFold(reqElem.Type, "picture") {
			if len(reqElem.PictureKeys) > 0 {
				// Convert request.PictureItem to entity.PictureItem
				newElem.PictureKeys = make([]entity.PictureItem, len(reqElem.PictureKeys))
				for j, reqItem := range reqElem.PictureKeys {
					newElem.PictureKeys[j] = entity.PictureItem{
						Key:   reqItem.Key,
						Order: reqItem.Order,
						Title: reqItem.Title,
					}
				}
				// Set first key as main value for backward compatibility
				if len(reqElem.PictureKeys) > 0 {
					newElem.Value = &reqElem.PictureKeys[0].Key
				}
			}
		} else {
			// Handle other types (text, banner, graphic, etc.)
			if reqElem.Value != nil {
				newElem.Value = reqElem.Value
			}
		}

		// Set video ID if present
		if reqElem.VideoID != nil {
			newElem.VideoID = reqElem.VideoID
		}

		newElements[i] = newElem
	}

	// PHASE 4: Replace old elements array with new one
	// This is atomic operation - no partial updates
	translation.Elements = newElements

	// PHASE 5: Cleanup unused files
	// Delete files that were in old elements but not in new request
	for fileKey := range existingFileKeys {
		// Skip if file is still being used in request
		if requestFileKeys[fileKey] {
			continue
		}

		// File is no longer used - delete it
		// Determine file type by extension
		if strings.HasSuffix(strings.ToLower(fileKey), ".pdf") {
			if err := u.fileGateway.DeletePDF(ctx, fileKey); err != nil {
				log.Printf("failed to delete unused PDF %s: %v", fileKey, err)
			}
		} else {
			// Assume it's an image
			if err := u.fileGateway.DeleteImage(ctx, fileKey); err != nil {
				log.Printf("failed to delete unused image %s: %v", fileKey, err)
			}
		}
	}

	return nil
}
