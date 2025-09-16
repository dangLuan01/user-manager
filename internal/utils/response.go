package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest 	ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound   	ErrorCode = "NOT_FOUND"
	ErrCodeConflict   	ErrorCode = "CONFLICT"
	ErrCodeInternal   	ErrorCode = "INTERNAL_ERROR_SERVER"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (ae *AppError) Error() string {
	return ""
}

func NewError(code, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func WrapError(code, message string, err error) error {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok :=err.(*AppError);ok {
		status := httpStatusFromCode(ErrorCode(appErr.Code))
		response := gin.H{
			"error":appErr.Message,
			"code":	appErr.Code,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}
		ctx.JSON(status, response)
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code": ErrCodeInternal,
	})
}

func ResponseSuccess(ctx *gin.Context, status int, message string, data ...any) {
	ctx.JSON(status, gin.H{
		"status":"SUCCESS",
		"message": message,
		"data": data,
	})
}

func ResponseSatus(ctx *gin.Context, status int)  {
	ctx.Status(status)
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadGateway, data)
}
func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadGateway
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	default :
		return http.StatusInternalServerError
	}
}
