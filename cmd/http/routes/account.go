package routes

import (
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

// Swagger:route POST /accounts accounts createAccount
// Create a new account
// responses:
//
//	200: accountResponse
//	400: errorResponse
//	500: errorResponse
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

// swagger:route DELETE /accounts/{id} accounts deleteAccount
// Delete an account
// responses:
//
//	200: accountResponse
//	400: errorResponse
//	500: errorResponse
func (ac *AccountController) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := ac.service.Delete(uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	var account database.Account
	ctx.JSON(http.StatusOK, account)
}

// swagger:route GET /accounts/{id} accounts fetchAccount
// Fetch an account by id
// responses:
//
//	200: accountResponse
//	400: errorResponse
//	500: errorResponse
func (accountController *AccountController) FetchById(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	account, err := accountController.service.FetchById(uint(id))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// swagger:route GET /accounts accounts listAccounts
// List all accounts
// responses:
//
//	200: accountResponse
//	400: errorResponse
//	500: errorResponse
func (ac *AccountController) List(ctx *gin.Context) {
	accounts, err := ac.service.List()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// swagger:route PUT /accounts/{id} accounts updateAccount
// Update an account
// responses:
//
//	200: accountResponse
//	400: errorResponse
//	500: errorResponse
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
