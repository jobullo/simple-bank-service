package routes

import (
	"errors"
	"fmt"
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

// @Summary create a transaction record
// @Description allows a transaction to be created in the database if the account exists
// @Tags Tranasctions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} database.Transaction
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /transactions [post]

func (tc *TransactionController) Create(ctx *gin.Context) {
	var transaction database.Transaction

	if err := ctx.BindJSON(&transaction); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := tc.service.Create(&transaction); err != nil {

		switch {
		case errors.Is(err, database.ErrParentNotFound):
			ctx.AbortWithStatusJSON(http.StatusConflict, NewError(fmt.Sprintf("Account with ID %d not found", transaction.AccountID)))
		case errors.Is(err, database.ErrInvalidType):
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(fmt.Sprintf("Invalid transaction type %s", transaction.Type)))
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

// @Summary delete a transaction record
// @Description allows a transaction to be deleted from the database
// @Tags Tranasctions
// @Security ApiKeyAuth
// @Accept  json
// @Produce no content
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /transactions/:id [delete]

func (tc *TransactionController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := tc.service.Delete(uint(id)); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewError(fmt.Sprintf("Transaction with ID %d not found", id)))
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)

}

// @Summary delete a transaction record
// @Description allows a transaction to be deleted from the database
// @Tags Tranasctions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 204 {object} database.Transaction
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /transactions/:id [get]
func (transactionController *TransactionController) FetchById(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	transaction, err := transactionController.service.FetchById(int(id))

	if err != nil {

		if errors.Is(err, database.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewError(fmt.Sprintf("Transaction with ID %d not found", id)))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)

}

// @Summary list all transaction records
// @Description lists all transaction records in the DB
// @Tags Tranasctions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {array} database.Transaction
// @Failure 500
// @Router /transactions [get]
func (transactionController *TransactionController) List(ctx *gin.Context) {

	transactions, err := transactionController.service.List()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// @Summary update the amount of a transaction record
// @Description update the amount of a transaction record
// @Tags Tranasctions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} database.Transaction
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /transactions/:id [put]
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

		if errors.Is(err, database.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewError(fmt.Sprintf("Transaction with ID %d not found", id)))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
