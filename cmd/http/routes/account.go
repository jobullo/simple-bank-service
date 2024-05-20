package http

import (
	service "github.com/jobullo/go-api-example/service"
	http "net/http"
	strconv "strconv"

	gin "github.com/gin-gonic/gin"
	database "github.com/jobullo/go-api-example/database"
)

type AccountController struct {
	service *service.AccountService
}

func NewAccountController(service *service.AccountService) *AccountController {
	return &AccountController{service: service}
}

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

func (ac *AccountController) Delete(ctx *gin.Context) {

	if id, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	if err := ac.service.Delete(uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (accountController *AccountController) FetchById(ctx *gin.Context) {

	if id, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	account, err := accountController.service.FetchById(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (ac *AccountController) List(ctx *gin.Context) {
	accounts, err := ac.service.List()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (ac *AccountController) Update(ctx *gin.Context) {
	
	if id, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	
	var account database.Account
	if err := ctx.BindJSON(&account) err != nil {
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
