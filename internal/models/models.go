package models

import (
	"time"
)

type Registration struct {
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Password2 string
}

//Users is a users model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
