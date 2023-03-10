package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/BrunoMoises/go-finance-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	UserID      int32     `json:"user_id" binding:"require"`
	CategoryID  int32     `json:"category_id" binding:"require"`
	Title       string    `json:"title" binding:"require"`
	Type        string    `json:"type" binding:"require"`
	Description string    `json:"description" binding:"require"`
	Value       int32     `json:"value" binding:"require"`
	Date        time.Time `json:"date" binding:"require"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var CategoryId = req.CategoryID
	var accountType = req.Type

	category, err := server.store.GetCategory(ctx, CategoryId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	}
	compareType := category.Type != accountType
	if compareType {
		ctx.JSON(http.StatusBadRequest, "Account type is different of Category type")
		return
	}

	arg := db.CreateAccountParams{
		UserID:      req.UserID,
		CategoryID:  CategoryId,
		Title:       req.Title,
		Type:        accountType,
		Description: req.Description,
		Value:       req.Value,
		Date:        req.Date,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int32 `uri:"id" binging:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int32 `uri:"id" binging:"required"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type updateAccountRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       int32  `json:"value"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateAccountParams{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Value:       req.Value,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CategoryID  int32     `json:"category_id"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
		CategoryID:  req.CategoryID,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type getAccountsGraphRequest struct {
	UserID int32  `json:"user_id" binding:"required"`
	Type   string `json:"type" binding:"required"`
}

func (server *Server) getAccountGraph(ctx *gin.Context) {
	var req getAccountsGraphRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountGraphParams{
		UserID: req.UserID,
		Type:   req.Type,
	}

	countGraph, err := server.store.GetAccountGraph(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, countGraph)
}

type getAccountsReportsRequest struct {
	UserID int32  `json:"user_id" binding:"required"`
	Type   string `json:"type" binding:"required"`
}

func (server *Server) getAccountReports(ctx *gin.Context) {
	var req getAccountsReportsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountReportsParams{
		UserID: req.UserID,
		Type:   req.Type,
	}

	sumReport, err := server.store.GetAccountReports(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sumReport)
}
