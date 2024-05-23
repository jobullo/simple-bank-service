package http

import (
	service "go-api-example/service"
	http "net/http"
	strconv "strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/jobullo/go-api-example/database"
)

type TransactionController struct {
	service *service.TransactionService
}

func NewTransactionController(service *service.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

func (tc *TransactionController) Create(ctx *gin.Context) {
	var transaction database.Transaction

	if err := ctx.BindJSON(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := tc.service.Create(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (tc *TransactionController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := tc.service.Delete(uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)

}

func (transactionController *TransactionController) FetchById(ctx *gin.Context) {

	if id, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	transaction, err := transactionController.service.FetchById(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)

}

func (transactionController *TransactionController) List(ctx *gin.Context) {

	transactions, err := transactionController.service.List()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (transactionController *TransactionController) Update(ctx *gin.Context) {

	if id, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	var transaction database.Transaction

	if err := ctx.BindJSON(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	transaction.ID = uint(id)

	if err := tc.service.Update(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
