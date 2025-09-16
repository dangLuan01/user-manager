package v1service

import (
	"fmt"

	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUser()  ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		
		return nil, utils.WrapError(
			string(utils.ErrCodeInternal), 
			"Faile fetch users.", 
			err,
		)
	}

	return users, nil
}

func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	
	user, found := us.repo.FindBYUUID(uuid);
	if !found {

		return models.User{}, utils.NewError(string(utils.ErrCodeNotFound), "No user")
	}
	
	return user, nil
}

func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormailizeString(user.Email)
	if user, err := us.repo.FindByEmail(user.Email); err != nil {
		
		return models.User{}, utils.NewError(
			string(utils.ErrCodeConflict), 
			fmt.Sprintf("Email: %v already existed.", user.Email),
		)
	}
	user.UUID = uuid.New().String()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		return models.User{}, utils.WrapError(
			string(utils.ErrCodeInternal), 
			"Faile hash password", 
			err,
		)
	}
	user.Password = string(hashPassword)
	if err := us.repo.Create(user); err != nil {

		return models.User{}, utils.WrapError(
			string(utils.ErrCodeInternal), 
			"Faile create user", 
			err,
		)
	}
	
	return user, nil
}
func (us *userService) UpdateUser(uuid string, user models.User) (models.User, error) {
	user.Email = utils.NormailizeString(user.Email)
	if u, err := us.repo.FindByEmail(user.Email); err != nil && u.UUID != uuid{
		
		return models.User{}, utils.NewError(
			string(utils.ErrCodeConflict), 
			fmt.Sprintf("Email: %v already existed.", u.Email),
		)
	}
	currencyUser, found := us.repo.FindBYUUID(uuid)
	if !found {
		return models.User{}, utils.NewError(string(utils.ErrCodeNotFound), "user not found")
	}
	
	currencyUser.Name = user.Name
	currencyUser.Email = user.Email

	if user.Password != "" {
		hashPassword, err :=bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(string(utils.ErrCodeInternal), "Faile hash pass", err)
		}
		currencyUser.Password = string(hashPassword)
		
	}
	if user.Age != 0 {
		currencyUser.Age = user.Age	
	}
	if user.Level != 0 {
		currencyUser.Level = user.Level	
	}
	
	if user.Status != 0 {
		currencyUser.Status = user.Status	
	}
	
	if err := us.repo.Update(uuid, currencyUser); err != nil {
		return models.User{}, utils.WrapError(string(utils.ErrCodeInternal), "Faile update user", err)
	}
	return currencyUser, nil
}

func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.WrapError(string(utils.ErrCodeInternal), "Faile delete user", err)
	}
	
	return nil

}