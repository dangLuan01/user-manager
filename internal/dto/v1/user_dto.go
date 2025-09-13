package v1dto

import (
	"github.com/dangLuan01/user-manager/internal/models"
)

type UserDTO struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int16 `json:"age"`
	Level  string `json:"level"`
	Status string `json:"status"`
}
type CreateUserInput struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int16  `json:"age" binding:"required,gt=0,lt=127"`
	Status   int8   `json:"status" binding:"required,oneof=1 2"`
	Level    int8   `json:"level" binding:"required,oneof=1 2"`
}

type UpdateUserInput struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Age      int16  `json:"age" binding:"omitempty,gt=0,lt=127"`
	Status   int8   `json:"status" binding:"omitempty,oneof=1 2"`
	Level    int8   `json:"level" binding:"omitempty,oneof=1 2"`
}

func (input * CreateUserInput) MapCreateInputToModel() models.User {
	return models.User{
		Name: input.Name,
		Email: input.Email,
		Password: input.Password,
		Age: input.Age,
		Status: input.Status,
		Level: input.Level,
	}
}

func (input * UpdateUserInput) MapUpdateInputToModel() models.User {
	return models.User{
		Name: input.Name,
		Email: input.Email,
		Password: input.Password,
		Age: input.Age,
		Status: input.Status,
		Level: input.Level,
	}
}

func MapUserDTO(user models.User) *UserDTO {
	return &UserDTO{
		UUID: user.UUID,
		Name: user.Name,
		Email: user.Email,
		Age: user.Age,
		Level: formatLevel(user.Level),
		Status: formatStatus(user.Status),
	}
}

func MapUsersDTO(users []models.User) []UserDTO {
	dtos := make([]UserDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, *MapUserDTO(user))
	}
	return dtos
}

func formatLevel(level int8) string {
	switch level {
	case 1:
		return "Admin"
	default :
		return "Customer"
	}
}

func formatStatus(status int8) string {
	switch status {
	case 1:
		return "Active"
	default :
		return "Hidden"
	}
}