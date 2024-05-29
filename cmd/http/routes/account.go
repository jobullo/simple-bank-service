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

type AccountController struct {
	service *service.AccountService
}

func NewAccountController(service *service.AccountService) *AccountController {
	return &AccountController{service: service}
}

// @Summary create an account record
// @api {post} /accounts
// @apiName CreateAccount
// @apiGroup Accounts
// @Param account body database.Account true "create account"
// @Description creates an account record in the DB
// @Tags Accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} database.Account
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /accounts [post]
func (ac *AccountController) Create(ctx *gin.Context) {

	var account database.Account
	if err := ctx.BindJSON(&account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := ac.service.Create(&account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// @Summary delete an account record
// @Description deletes an account record from the DB
// @Tags Accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 204
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /accounts/:id [delete]
func (ac *AccountController) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := ac.service.Delete(uint(id)); err != nil {

		if errors.Is(err, database.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewError(fmt.Sprintf("Account with ID %d not found", id)))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary fetches an account record by id
// @Description fetches an account record by id from the DB
// @Tags Accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} database.Account
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /accounts/:id [get]
func (accountController *AccountController) FetchById(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	account, err := accountController.service.FetchById(uint(id))

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewError(fmt.Sprintf("Account with ID %d not found", id)))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// @Summary list all account records
// @Description lists all account records in the DB
// @Tags Accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {array} database.Account
// @Failure 500 {object} error
// @Router /accounts [get]
func (ac *AccountController) List(ctx *gin.Context) {
	accounts, err := ac.service.List()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// @Summary update an account record
// @Description updates an account record in the DB
// @Tags Accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} database.Account
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /accounts/:id [put]
func (ac *AccountController) Update(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	var account database.Account
	if err := ctx.BindJSON(&account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	account.ID = uint(id)

	if err := ac.service.Update(&account); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewError(fmt.Sprintf("Account with ID %d not found", id)))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
