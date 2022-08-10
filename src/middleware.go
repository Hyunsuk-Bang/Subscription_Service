package main

import "net/http"

func (app *Config) sessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}
