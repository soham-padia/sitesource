package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/solow-crypt/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//returns a user by id
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at 
				from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

//updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update user set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

//authenticates a user
func (m *postgresDBRepo) Authenticate(email, testpassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testpassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
func (m *postgresDBRepo) DoesEmailExist(u models.Registration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int

	checkforEmailQuery := `select id from users where email = $1`

	row := m.DB.QueryRowContext(ctx, checkforEmailQuery, u.Email)

	_ = row.Scan(&id)
	log.Println("ooooooooooooooooooo")
	log.Println(id)
	//log.Println(err)
	log.Println("-------------------")

	if !(id != 0) {
		return false
	} else {
		return true
	}

}

func (m *postgresDBRepo) InsertUser(u models.Registration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into users (first_name,last_name,email,password,access_level, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)

	_, err := m.DB.ExecContext(ctx, stmt,
		u.Firstname,
		u.Lastname,
		u.Email,
		hashedPassword,
		1,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
