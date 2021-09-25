package repository

import "github.com/solow-crypt/bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testpassword string) (int, string, error)
	InsertUser(u models.Registration) error
	DoesEmailExist(u models.Registration) bool
}
