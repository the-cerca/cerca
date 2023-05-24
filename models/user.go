package models

import (
	"database/sql"
	"errors"


	"net/mail"
	"strings"
	"time"

	"github.com/aleeXpress/cerca/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordDoNotMatch   = errors.New("password do not match")
	ErrEmailDoNotMatch      = errors.New("email not found")
	ErrEmailAlreadyUse      = errors.New("email already used")
	ErrBadAddressFormat     = errors.New("errors on the address mail")
	ErrUsernameNotAvailable = errors.New("username already used")
	ErrUsernameNotFound     = errors.New("username not found")
)

type NewUser struct {
	Firstname string    `json:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email,omitempty"`
	Birthday  time.Time `json:"birthday,omitempty"`
}

type User struct {
	ID             string     `json:"id,omitempty"`
	Firstname      string     `json:"last_name,omitempty"`
	Lastname       string     `json:"username,omitempty"`
	Username       string     `json:"user_name,omitempty"`
	PasswordHashed string     `json:"-"`
	Email          string     `json:"email,omitempty"`
	Birthday       *time.Time `json:"birthday,omitempty"`
	IsVerified     bool       `json:"is_verified,omitempty"`
	IsFreelance    bool       `json:"is_freelance,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type UserManager struct {
	DB *sql.DB
}

// SignUp allows a user to sign up by providing a username, password and email.
// It first checks that the email address provided is in a valid format, then generates a password hash
// and checks whether the email and username already exist in the database.
// If the email or username already exist, the function returns an appropriate error.
// Otherwise, it inserts the new user's information into the database and returns the created user object.
func (um *UserManager) SignUp(firstname, lastname, username, password, email string, birthday time.Time) (*User, error) {
	// Convert email to lowercase and parse email address
	email = strings.ToLower(email)
	lastname = strings.TrimSpace(lastname)
	firstname = strings.TrimSpace(firstname)
	username = strings.TrimSpace(username)
	mail, err := mail.ParseAddress(email)
	if err != nil {
		return nil, ErrBadAddressFormat
	}
	// Generate password hash
	hp, _ := utils.HashPassword(password)
	// Check if email and username already exist in the database
	var emailExists, usernameExists bool
	if err := um.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, mail.Address).Scan(&emailExists); err != nil {
		return nil, err
	}
	if err := um.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username).Scan(&usernameExists); err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailAlreadyUse
	}
	if usernameExists {
		return nil, ErrUsernameNotAvailable
	}
	// Insert new user's information into the database
	var user User
	if err := um.DB.QueryRow(`
		INSERT INTO users (firstname, lastname, username, password_hashed, email,birthday)
		VALUES ($1, $2, $3, $4, $5,$6)
		RETURNING *
	`, firstname, lastname, username, hp, mail.Address, birthday).Scan(
		&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.PasswordHashed,
		&user.Email, &birthday, &user.IsVerified, &user.IsFreelance, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserManager) SignIn(username, password string) (*User, error) {
	var user User
	if err := um.DB.QueryRow(`select * from users where username=$1`, username).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.PasswordHashed, &user.Email, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, ErrUsernameNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHashed), []byte(password)); err != nil {
		return nil, ErrPasswordDoNotMatch
	}
	return &user, nil
}

func (um *UserManager) UpdateUserData(userID string, firstname string, lastname string, username string, email, password string) (*User, error) {
	hp, _ := utils.HashPassword(password)
	query := `
			UPDATE users 
			SET firstname=$1, lastname=$2, username=$3, email=$4, password_hashed=$5, updated_at=NOW() 
			WHERE id=$6 RETURNING *
	`
	var user User
	if err := um.DB.QueryRow(query, firstname, lastname, username, email, hp, userID).Scan(
		&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.PasswordHashed,
		&user.Email, &user.Birthday, &user.IsVerified, &user.IsFreelance, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
func (um *UserManager) UpdatePassword(password, mail string) error {
	hp, _ := utils.HashPassword(password)
	query := `UPDATE users SET password_hashed=$1 where email=$2`
	if _, err := um.DB.Exec(query, hp, mail); err != nil {
		return err
	}
	return nil
}

func (um *UserManager) DeleteAccount(username string) error {
	if _, err := um.DB.Exec(`delete from users where usernname=$1`, username); err != nil {
		return err
	}
	return nil
}
