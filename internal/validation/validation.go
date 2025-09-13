package validation

import (
	"fmt"
	"strings"

	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}
	RegisterCustomValidation(v)
	return nil
}
func HandlerValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationError {
			root 	:= strings.Split(e.Namespace(), ".")[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root + ".")
			parts 	:= strings.Split(rawPath, ".")
			for i, part := range parts {	
				if strings.Contains("part", "[") {
					idx 	:= strings.Index(part, "[")
					base 	:= utils.CamelToSnakeCase(part[:idx])
					index 	:= part[idx:]
					parts[i] = fmt.Sprintf("%s%s", base, index)
				} else {
					parts[i] = utils.CamelToSnakeCase(part)
				}
			}
			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s phải và uuid %s", fieldPath, e.Param())
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s phải nhỏ hơn %s", fieldPath, e.Param())
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s phải là một slug hợp lệ", fieldPath)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s là trường bắt buộc", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s phải có ít nhất %s ký tự", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s không được vượt quá %s ký tự", fieldPath, e.Param())
			case "url":
				errors[fieldPath] = fmt.Sprintf("%s phải là một URL hợp lệ", fieldPath)
			case "minInt":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn %s", fieldPath, e.Param())
			case "maxInt":
				errors[fieldPath] = fmt.Sprintf("%s không được lớn hơn %s", fieldPath, e.Param())
			case "file_ext":
				exts := strings.Split(e.Param(), " ")
				errors[fieldPath] = fmt.Sprintf("%s phải có phần mở rộng là %s", fieldPath, strings.Join(exts, ", "))
			case "oneof":
				options := strings.Split(e.Param(), " ")
				errors[fieldPath] = fmt.Sprintf("%s phải là một trong các giá trị: %s", fieldPath, strings.Join(options, ", "))
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s phải đúng định dạng %s", fieldPath, fieldPath)
			}
		}
		return gin.H{"errors": errors}
	}
	return gin.H{
		"error": "Validation failed",
		"details": err.Error(),
	}
}