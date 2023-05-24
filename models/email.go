package models

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"github.com/aleeXpress/cerca/utils"
)

type Verification struct {
	ID       string
	Username string
	Token    string
}
type Mail struct {
	Host string
	From string
	To   []string
	Body string
}

type AuthEmail struct {
	Identity string
	Username string
	Password string
	Host     string
}
type MailManager struct {
	DB *sql.DB
}

func (mm *MailManager) SendEmailVerification(to []string, data interface{}) error {
	var body bytes.Buffer
	t, err := template.ParseFiles("./template/register.html")
	if err != nil {
		return err
	}
	err = t.Execute(&body, data)
	if err != nil {
		return err
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	authen := smtp.PlainAuth("", "cercanotreply@gmail.com",os.Getenv("smtp"), "smtp.gmail.com")
	mail := Mail{
		Host: "smtp.gmail.com:587",
		From: "cercanotreply@gmail.com",
		To:   to,
		Body: "Subject: Email Verification\r\n" + headers + "\r\n\r\n" + body.String(),
	}
	err = smtp.SendMail(mail.Host, authen, mail.From, mail.To, []byte(mail.Body))
	if err != nil {
		return err
	}

	return nil
}

func (mm *MailManager) Create(UserID string) (string, error) {
	token, err := utils.String(40)
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(time.Hour * 24)
	queryInsert := `insert into email_verification(token, user_id, expires_at) values($1,$2,$3)`
	_, err = mm.DB.Exec(queryInsert, token, UserID, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

// If the token exist 
func (mm *MailManager) VerifyToken(token string) error {
	query := `SELECT users.id, username, email_verification.token
	FROM users
	INNER JOIN email_verification
	ON users.id = email_verification.user_id
	WHERE email_verification.token=$1
	`

	row := mm.DB.QueryRow(query, token)
	var verificationToken Verification
	row.Scan(&verificationToken.ID, &verificationToken.Username, &verificationToken.Token)
	if row.Err() == sql.ErrNoRows {
		return sql.ErrNoRows
	}
	if verificationToken.Username != "" {
		mm.DB.Exec("update users set is_verified =$1, updated_at=$2 where id=$3", true, time.Now(), verificationToken.ID)
	}
	return nil
}
func (mm *MailManager) DeleteToken(token string) error {
	_, err := mm.DB.Exec(`DELETE FROM email_verification WHERE token=$1`,token)
	if err != nil {
			return err
	}
	return nil
}
func (mm *MailManager)GeneriqueEmailSender(path, subject string, To []string, data interface{})(error)  {
	var body bytes.Buffer
	t, err := template.ParseFiles(path)
	if err != nil {
		return err 
	}
	if err := t.Execute(&body, data); err != nil {
		return err
	}
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=UTF-8;"
	authen := smtp.PlainAuth("", "cercanotreply@gmail.com", os.Getenv("smtp"), "smtp.gmail.com")
	mail := Mail{
		Host: "smtp.gmail.com:587",
		From: "cercanotreply@gmail.com",
		To:   To,
		Body: "Subject:"+subject+"\r\n" + headers + "\r\n\r\n" + body.String(),
	}
	err = smtp.SendMail(mail.Host, authen, mail.From, mail.To, []byte(mail.Body))
	if err != nil {
		return err
	}
	return nil
}