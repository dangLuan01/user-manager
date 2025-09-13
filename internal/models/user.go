package models

type User struct {
	UUID     string `db:"uuid"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Age      int16  `db:"age"`
	Level    int8   `db:"level"`
	Status   int8   `db:"status"`
}