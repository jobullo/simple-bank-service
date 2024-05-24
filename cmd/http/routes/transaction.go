package routes

import (
	http "net/http"
	strconv "strconv"

	service "github.com/jobullo/go-api-example/service"

	gin "github.com/gin-gonic/gin"
	database "github.com/jobullo/go-api-example/database"
)

type TransactionController struct {
	service *service.TransactionService
}

func NewTransactionController(service *service.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

// swagger:route POST /transactions transactions createTransaction
// Create a new transaction
// responses:
//
//	200: transactionResponse
//	400: errorResponse
//	500: errorResponse
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

// swagger:route DELETE /transactions/{id} transactions deleteTransaction
// Delete a transaction
// responses:
//
//	200: transactionResponse
//	400: errorResponse
//	500: errorResponse
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

	var transaction database.Transaction
	ctx.JSON(http.StatusOK, transaction)

}

// swagger:route GET /transactions/{id} transactions getTransaction
// Get a transaction by id
// responses:
//
//	200: transactionResponse
//	400: errorResponse
//	500: errorResponse
func (transactionController *TransactionController) FetchById(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	transaction, err := transactionController.service.FetchById(int(id))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)

}

// swagger:route GET /transactions transactions listTransactions
// List all transactions
// responses:
//
//	200: transactionsResponse
//	400: errorResponse
//	500: errorResponse
func (transactionController *TransactionController) List(ctx *gin.Context) {

	transactions, err := transactionController.service.List()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// swagger:route GET /transactions/account/{accountID} transactions listTransactionsByAccount
// List all transactions by account
// responses:
//
//	200: transactionsResponse
//	400: errorResponse
//	500: errorResponse
func (TransactionController *TransactionController) ListByAccount(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("accountID"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	transactions, err := TransactionController.service.ListByAccount(uint(id))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// swagger:route PUT /transactions/{id} transactions updateTransaction
// Update a transaction
// responses:
//
//	200: transactionResponse
//	400: errorResponse
//	500: errorResponse
func (transactionController *TransactionController) Update(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	var transaction database.Transaction

	if err := ctx.BindJSON(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	transaction.ID = uint(id)

	if err := transactionController.service.Update(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
