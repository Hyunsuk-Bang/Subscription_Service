package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(app.sessionLoad)

	mux.Get("/", app.MainPage)
	mux.Get("/login", app.LoginPage)
	mux.Get("/register", app.RegisterPage)
	mux.Get("/activate", app.ActivateAccount)
	mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Encryption:  "none",
			FromAddress: "info@mycompany.com",
			FromName:    "info",
			ErrorChan:   make(chan error),
		}

		msg := Message{
			To:      "me@here.com",
			Subject: "Test Email",
			Data:    "Hello World.",
		}

		m.sendMail(msg, make(chan error))
	})
	mux.Get("/plans", app.ChooseSubscription)

	mux.Get("/logout", app.Logout)
	mux.Post("/login", app.PostLoginPage)
	mux.Post("/register", app.PostRegisterPage)
	return mux
}
