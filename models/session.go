package models

import (
	"database/sql"
	"time"

	"github.com/aleeXpress/cerca/utils"
)

type SessionManager struct {
	DB *sql.DB
}

type Session struct {
	Id        string
	UserId    string
	Token     string
	create_at *time.Time
}
type NewSession struct {
	UserId string
	Token  string
}

func (sm *SessionManager) Create(userID string) (*Session, error) {
	token, err := utils.String(32)
	if err != nil {
		return nil, err
	}

	NewSession := Session{
		UserId: userID,
		Token:  token,
	}
	queryInsert := `insert into sessions (user_id, token) values ($1, $2 )`
	_, err = sm.DB.Exec(queryInsert, NewSession.UserId, NewSession.Token)
	if err != nil {
		return nil, err
	}
	querySelect := `select * from sessions where token=$1`
	row := sm.DB.QueryRow(querySelect, NewSession.Token)
	var session Session
	err = row.Scan(&session.Id, &session.UserId, &session.Token, &session.create_at)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (sm *SessionManager) FindUserByCookie(token string) (*User, error) {
	var user User
	if err := sm.DB.QueryRow(`select users.id, users.firstname,users.lastname,users.username,users.password_hashed,
	 users.email, users.birthday, users.is_verified,users.created_at, 
	 users.updated_at from users inner join sessions
	on users.id = sessions.user_id where 
	sessions.token=$1`, token).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.PasswordHashed, &user.Email, &user.Birthday, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
