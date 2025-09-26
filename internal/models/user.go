package models

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID `db:"uuid"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Age      int16  `db:"age" goqu:"omitempty"`
	Level    int8   `db:"level"`
	Status   int8   `db:"status"`
}