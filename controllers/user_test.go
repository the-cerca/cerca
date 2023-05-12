package controllers

import (
	"net/http"
	"testing"

	m "github.com/aleeXpress/cerca/models"
)

func TestUserC_SignUp(t *testing.T) {
	type fields struct {
		Us  *m.UserManager
		Usm *m.SessionManager
		Ms  *m.MailManager
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usc := &UserC{
				Us:  tt.fields.Us,
				Usm: tt.fields.Usm,
				Ms:  tt.fields.Ms,
			}
			usc.SignUp(tt.args.w, tt.args.r)
		})
	}
}
