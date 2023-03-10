package api

import (
	"database/sql"
	"net/http"

	db "github.com/BrunoMoises/go-finance-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	UserID      int32  `json:"user_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreateCategoryParams{
		UserID:      req.UserID,
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
	}

	category, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type getCategoryRequest struct {
	ID int32 `uri:"id" binging:"required"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	category, err := server.store.GetCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type deleteCategoryRequest struct {
	ID int32 `uri:"id" binging:"required"`
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := server.store.DeleteCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type updateCategoryRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var req updateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateCategoryParams{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	}

	category, err := server.store.UpdateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type getCategoriesRequest struct {
	UserID      int32  `json:"user_id" binding:"required"`
	Title       string `json:"title"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
}

func (server *Server) getCategories(ctx *gin.Context) {
	var req getCategoriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var categories []db.Category
	filterUserIdAndType := req.UserID > 0 && len(req.Type) > 0

	filterAsByUserIdAndType := len(req.Description) == 0 && len(req.Title) == 0 && filterUserIdAndType

	if filterAsByUserIdAndType {
		arg := db.GetCategoriesByUserIdAndTypeParams{
			UserID: req.UserID,
			Type:   req.Type,
		}

		categoriesByUserIdAndType, err := server.store.GetCategoriesByUserIdAndType(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByUserIdAndType
	}

	filterAsByUserIdAndTypeAndDescription := len(req.Description) != 0 && len(req.Title) == 0 && filterUserIdAndType
	if filterAsByUserIdAndTypeAndDescription {
		arg := db.GetCategoriesByUserIdAndTypeAndDescriptionParams{
			UserID:      req.UserID,
			Type:        req.Type,
			Description: req.Description,
		}

		categoriesByUserIdAndTypeAndDescription, err := server.store.GetCategoriesByUserIdAndTypeAndDescription(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByUserIdAndTypeAndDescription
	}

	filterAsByUserIdAndTypeAndTitle := len(req.Description) == 0 && len(req.Title) != 0 && filterUserIdAndType
	if filterAsByUserIdAndTypeAndTitle {
		arg := db.GetCategoriesByUserIdAndTypeAndTitleParams{
			UserID: req.UserID,
			Type:   req.Type,
			Title:  req.Title,
		}

		categoriesByUserIdAndTypeAndTitle, err := server.store.GetCategoriesByUserIdAndTypeAndTitle(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByUserIdAndTypeAndTitle
	}

	filterAllParam := len(req.Description) != 0 && len(req.Title) != 0 && filterUserIdAndType
	if filterAllParam {
		arg := db.GetCategoriesParams{
			UserID: req.UserID,
			Type:   req.Type,
			Title:  req.Title,
		}

		categoriesByAllFilters, err := server.store.GetCategories(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByAllFilters
	}

	ctx.JSON(http.StatusOK, categories)
}
