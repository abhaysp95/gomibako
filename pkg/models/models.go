package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicateEmail = errors.New("models: provided email already exists")
var ErrInvalidCredentials = errors.New("models: invalid credentials, try again")

type Gomi struct {
	Id int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	Id int
	Name string
	Email string
	Passwd []byte
	Created time.Time
}
