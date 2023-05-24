package controllers

import (
	"encoding/json"
	"net/http"

	m "github.com/aleeXpress/cerca/models"
	"github.com/aleeXpress/cerca/utils"
	"github.com/go-chi/chi/v5"
)

type UserC struct {
	Us  *m.UserManager
	Usm *m.SessionManager
	Ms  *m.MailManager
}

func (usc *UserC) SignUp(w http.ResponseWriter, r *http.Request) {
	var user m.NewUser
	json.NewDecoder(r.Body).Decode(&user)
	u, err := usc.Us.SignUp(user.Firstname, user.Lastname, user.Username, user.Password, user.Email, user.Birthday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := usc.Usm.Create(u.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	verifyToken, _ := usc.Ms.Create(u.ID)
	send := map[string]string{
		"Username": u.Username,
		"Token":    verifyToken,
	}
	err = usc.Ms.SendEmailVerification([]string{user.Email}, send)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	setCookie(w, CookieSession, session.Token)
	utils.Encode(w, u)
}

func (usc *UserC) SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&credentials)
	user, err := usc.Us.SignIn(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if session, _ := usc.Usm.Create(user.ID); session != nil {
		setCookie(w, CookieSession, session.Token)
	}
	utils.Encode(w, user)
}

// Check if the token exist,
func (usc *UserC) VerifyToken(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	usc.Ms.VerifyToken(params)
	_ = usc.Ms.DeleteToken(params)
}
func (usc *UserC) ForgettenPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	setCookie(w, "email", email)
	if err := usc.Ms.GeneriqueEmailSender("./template/forgetten-password.html", "Reset password ! ", []string{email}, ""); err != nil {
		InternalServerError(w)
	}
}

func (usc *UserC) ResetPassword(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	cookie, err := readCookie(r, "email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := usc.Us.UpdatePassword(password, cookie); err != nil {
		InternalServerError(w)
		return
	}
}

func (usc *UserC) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByContext(r.Context())
	if err != nil {
		InternalServerError(w)
		return
	}
	mail := r.FormValue("email")
	password := r.FormValue("password")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	username := r.FormValue("username")
	u, err := usc.Us.UpdateUserData(user.ID, firstname, lastname, username, mail, password)
	if err != nil {
		InternalServerError(w)
		return
	}
	utils.Encode(w, u)
}

func (usc *UserC) CurrentUser(w http.ResponseWriter, r *http.Request) {
	if u, _ := GetUserByContext(r.Context()); u != nil {
		utils.Encode(w, u)
	}
}
