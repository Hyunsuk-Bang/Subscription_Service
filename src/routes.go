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
	mux.Post("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Get("/activate-account", app.ActivateAccount)

	mux.Post("/login", app.PostLoginPage)
	mux.Post("/logout", app.Logout)
	mux.Post("/register", app.PostRegisterPage)
	return mux
}
