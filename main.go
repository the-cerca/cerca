package main

import (
	"database/sql"
	"log"
	"net/http"


	c "github.com/aleeXpress/cerca/controllers"
	m "github.com/aleeXpress/cerca/models"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	cfg := m.PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "example",
		Database: "postgres",
		SSLMode:  "disable",
	}

	db, err := sql.Open("postgres", cfg.String())
	checkError(err)
	err = db.Ping()
	checkError(err)
	log.Println("Dababase connected 🚀")
	userC := c.UserC{
		Us:  &m.UserManager{DB: db},
		Usm: &m.SessionManager{DB: db},
		Ms:  &m.MailManager{DB: db},
	}
	serviceC := c.ServiveC{
		Im:      &m.ImageManager{DB: db},
		Cs:      &m.CategoryManager{DB: db},
		Sm:      &m.ServiceManager{DB: db},
		Session: &m.SessionManager{DB: db},
	}
	AuthenticationMiddleware := c.UserMiddleware{
		Session: &m.SessionManager{DB: db},
	}
	r := chi.NewRouter()
	r.Post("/sign-up", userC.SignUp)
	r.Post("/sign-in", userC.SignIn)
	r.Get("/verify/{id}", userC.VerifyToken)
	r.Post("/forgetten-password", userC.ForgettenPassword)
	r.Post("/reset-password", userC.ResetPassword)
	r.Route("/", func(r chi.Router) {
		r.Use(AuthenticationMiddleware.CurrentUser)
		r.Post("/update-user", userC.UpdateUserData)
		r.Post("/create-service", serviceC.CreateService)
		r.Post("/update-service", serviceC.UpdateService)
	})
	log.Fatal(http.ListenAndServe(":8000", r))
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}
