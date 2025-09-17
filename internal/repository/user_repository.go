package repository

import (
	"fmt"
	"slices"

	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type SqlUserRepository struct {
	users []models.User
	db *goqu.Database
}

func NewSqlUserRepository(DB *goqu.Database) UserRepository {
	return &SqlUserRepository{
		users : make([]models.User, 0),
		db: DB,
	}
}

func (ur *SqlUserRepository) FindAll() ([]models.User, error){
	
	ds := ur.db.From(goqu.T("users")).
	Select(
		goqu.I("uuid"),
		goqu.I("name"),
		goqu.I("email"),
		goqu.I("age"),
		goqu.I("level"),
		goqu.I("status"),
	)
	var users []models.User
	if err := ds.ScanStructs(&users); err != nil {
		return nil, fmt.Errorf("faile get all user:%v", err)
	}

	return users, nil
}

func (ur *SqlUserRepository) FindBYUUID(uuid uuid.UUID) (models.User, error) {
	ds := ur.db.From(goqu.T("users")).
	Where(
		goqu.C("uuid").Eq(uuid),
	).
	Select(
		goqu.I("uuid"),
		goqu.I("name"),
		goqu.I("email"),
		goqu.I("age"),
		goqu.I("level"),
		goqu.I("status"),
	)
	var user models.User

	found, err := ds.ScanStruct(&user)
	if err != nil || !found {
		return  models.User{}, err
	}

	return user, err
}

func (ur *SqlUserRepository) Create(user models.User) error {
	insertUser := ur.db.Insert("users").Rows(user).Executor()
	if _, err := insertUser.Exec(); err != nil {
       return fmt.Errorf("faile insert rows user")
	}

	return nil
}

func (ur *SqlUserRepository) Update(uuid uuid.UUID, user models.User) error {
	for i, u := range ur.users{
		if u.UUID == uuid {
			ur.users[i] = user
			return nil
		}
	}
	return fmt.Errorf("user not found")
}
func (ur *SqlUserRepository) Delete(uuid uuid.UUID) error {

	for i, u := range ur.users{
		if u.UUID == uuid {
			ur.users = slices.Delete(ur.users, i, i + 1)
			return nil
		}
	}

	return fmt.Errorf("user not found")
}
func (ur *SqlUserRepository) FindByEmail(email string) (models.User, error) {
	
	ds := ur.db.From(goqu.T("users")).Where(
		goqu.C("email").Eq(email),
	).Limit(1)
	
    var user models.User
    found, err := ds.ScanStruct(&user)
	if err != nil {
		return models.User{}, err
	}
	
	if found {
		return user, nil
	}

	return models.User{}, err
}