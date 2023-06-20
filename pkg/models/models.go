package models

import "time"

type Gomi struct {
	Id int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}
