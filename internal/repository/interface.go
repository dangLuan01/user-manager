package repository

import (
	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindBYUUID(uuid uuid.UUID) (models.User, error)
	Create(user models.User) error
	Update(uuid uuid.UUID, user models.User) error
	Delete(uuid uuid.UUID) error
	FindByEmail(email string) (models.User, error)
}