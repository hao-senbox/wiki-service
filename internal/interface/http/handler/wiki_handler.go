package handler

import (
	"context"
	"strconv"
	"wiki-service/internal/domain/usecase"
	"wiki-service/internal/interface/http/dto/request"
	libs_constant "wiki-service/pkg/libs/constant"
	libs_helper "wiki-service/pkg/libs/helper"

	"github.com/gofiber/fiber/v2"
)

type WikiHandler struct {
	wikiUseCase usecase.WikiUseCase
}

func NewWikiHandler(wikiUseCase usecase.WikiUseCase) *WikiHandler {
	return &WikiHandler{
		wikiUseCase: wikiUseCase,
	}
}

func (h *WikiHandler) CreateWikiTemplate(c *fiber.Ctx) error {
	var req request.CreateWikiTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		libs_helper.SendError(c, fiber.StatusBadRequest, err, libs_helper.ErrInvalidRequest)
		return nil
	}

	userID, exists := c.Locals("user_id").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing userID")
		return nil
	}

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	err := h.wikiUseCase.CreateWikiTemplate(ctx, req, userID)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}

	return libs_helper.SendSuccess(c, fiber.StatusCreated, "Wiki template created successfully", nil)
}

func (h *WikiHandler) GetTemplate(c *fiber.Ctx) error {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing organizationID")
		return nil
	}

	typeParam := c.Query("type")
	if typeParam == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing type parameter")
		return nil
	}

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	templates, err := h.wikiUseCase.GetTemplate(ctx, organizationID, typeParam)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}

	return libs_helper.SendSuccess(c, fiber.StatusOK, "Templates fetched successfully", templates)
}

func (h *WikiHandler) GetStatistics(c *fiber.Ctx) error {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing organizationID")
		return nil
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid page parameter")
		return nil
	}

	limitParam := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid limit parameter")
		return nil
	}

	typeParam := c.Query("type", "")
	searchParam := c.Query("search")

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	statistics, err := h.wikiUseCase.GetStatistics(ctx, organizationID, page, limit, typeParam, searchParam)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}

	return libs_helper.SendSuccess(c, fiber.StatusOK, "Statistics fetched successfully", statistics)
}

func (h *WikiHandler) GetWikiByCode(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing code")
		return nil
	}

	var language *int
	if langParam := c.Query("language"); langParam != "" {
		lang, err := strconv.Atoi(langParam)
		if err != nil || lang < 0 {
			libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid language parameter")
			return nil
		}
		language = &lang
	}

	organizationID := c.Query("organization_id")
	if organizationID == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing organizationID")
		return nil
	}

	typeParam := c.Query("type")
	if typeParam == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing type parameter")
		return nil
	}

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	wiki, err := h.wikiUseCase.GetWikiByCode(ctx, code, language, organizationID, typeParam)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}

	return libs_helper.SendSuccess(c, fiber.StatusOK, "Wiki fetched successfully", wiki)
}

func (h *WikiHandler) GetWikis(c *fiber.Ctx) error {
	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid page parameter")
		return nil
	}

	limitParam := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid limit parameter")
		return nil
	}

	organizationID := c.Query("organization_id")
	if organizationID == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing organizationID")
		return nil
	}

	var language *int
	if langParam := c.Query("language"); langParam != "" {
		lang, err := strconv.Atoi(langParam)
		if err != nil || lang < 0 {
			libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid language parameter")
			return nil
		}
		language = &lang
	}

	typeParam := c.Query("type")
	if typeParam == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing type parameter")
		return nil
	}

	searchParam := c.Query("search")

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	wikis, total, err := h.wikiUseCase.GetWikis(ctx, organizationID, page, limit, language, typeParam, searchParam)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := fiber.Map{
		"items":       wikis,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	}

	return libs_helper.SendSuccess(c, fiber.StatusOK, "Wikis fetched successfully", response)
}

func (h *WikiHandler) GetWikiByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing id")
		return nil
	}

	var language *int
	if langParam := c.Query("language"); langParam != "" {
		lang, err := strconv.Atoi(langParam)
		if err != nil || lang < 0 {
			libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Invalid language parameter")
			return nil
		}
		language = &lang
	}

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	wiki, err := h.wikiUseCase.GetWikiByID(ctx, id, language)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}

	return libs_helper.SendSuccess(c, fiber.StatusOK, "Wiki fetched successfully", wiki)
}

func (h *WikiHandler) UpdateWiki(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		libs_helper.SendError(c, fiber.StatusBadRequest, nil, "Missing id")
		return nil
	}

	token, exists := c.Locals("token").(string)
	if !exists {
		libs_helper.SendError(c, fiber.StatusUnauthorized, nil, "Missing token")
		return nil
	}

	ctx := context.WithValue(c.Context(), libs_constant.Token, token)

	var req request.UpdateWikiRequest
	if err := c.BodyParser(&req); err != nil {
		libs_helper.SendError(c, fiber.StatusBadRequest, err, libs_helper.ErrInvalidRequest)
		return nil
	}

	err := h.wikiUseCase.UpdateWiki(ctx, id, req)
	if err != nil {
		libs_helper.SendError(c, fiber.StatusInternalServerError, err, libs_helper.ErrInternal)
		return nil
	}

	return libs_helper.SendSuccess(c, fiber.StatusOK, "Wiki updated successfully", nil)
}
