package v1handler

import (
	"net/http"

	v1dto "github.com/dangLuan01/user-manager/internal/dto/v1"
	v1service "github.com/dangLuan01/user-manager/internal/service/v1"

	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/internal/validation"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service v1service.UserService
}
type GetUserByUUIDParam struct{
	Uuid string `uri:"uuid" binding:"uuid"`
}
func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}
func (uh *UserHandler) GetAllUser(ctx *gin.Context)  {

	users, err := uh.service.GetAllUser()
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, v1dto.MapUsersDTO(users))
	
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context)  {
	var param GetUserByUUIDParam
	err := ctx.ShouldBindUri(&param)
	if err != nil {

		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return 
	}
	user, err := uh.service.GetUserByUUID(param.Uuid)

	if err != nil {

		utils.ResponseError(ctx, err)
		return
	}
	
	utils.ResponseSuccess(ctx, http.StatusOK, v1dto.MapUserDTO(user))
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {

	var input v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {

		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))

		return
	}
	user := input.MapCreateInputToModel()
	
	createUser, err := uh.service.CreateUser(user)
	if err != nil {

		utils.ResponseError(ctx, err)

		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, v1dto.MapUserDTO(createUser))
}
func (uh *UserHandler) UpdateUser(ctx *gin.Context)  {
	var param GetUserByUUIDParam
	err := ctx.ShouldBindUri(&param)
	if err != nil {

		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return 
	}

	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {

		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))

		return
	}

	user := input.MapUpdateInputToModel()

	updateUser, err := uh.service.UpdateUser(param.Uuid, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK ,updateUser)
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context)  {
	var param GetUserByUUIDParam
	err := ctx.ShouldBindUri(&param)
	if err != nil {

		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return 
	}
	if err := uh.service.DeleteUser(param.Uuid); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	
	utils.ResponseSatus(ctx, http.StatusNoContent)
}
