package controllers

import (
	"net/http"

	m "github.com/aleeXpress/cerca/models"
)
type UserMiddleware struct {
	Session *m.SessionManager
}
func (um *UserMiddleware)CurrentUser(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		coo, err := readCookie(r,CookieSession)
		if err != nil {
			next.ServeHTTP(w,r)
			return
		}
		 user, err := um.Session.FindUserByCookie(coo)
		 if err != nil {
			next.ServeHTTP(w,r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		 }
		 ctx = SetContextUser(ctx, user)
		 r = r.WithContext(ctx)

		 next.ServeHTTP(w,r)
		})
}