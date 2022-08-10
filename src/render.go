package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var TEMPLATE_PATH = "templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	DataMap       map[string]interface{}
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	//User *data.User
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", TEMPLATE_PATH),
		fmt.Sprintf("%s/header.partial.gohtml", TEMPLATE_PATH),
		fmt.Sprintf("%s/navbar.partial.gohtml", TEMPLATE_PATH),
		fmt.Sprintf("%s/footer.partial.gohtml", TEMPLATE_PATH),
		fmt.Sprintf("%s/alerts.partial.gohtml", TEMPLATE_PATH),
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", TEMPLATE_PATH, t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	// Used this for data that will be consumer once
	td.Flash = app.Session.PopString(r.Context(), "Flash") // As soon as the data is read, the data is removed
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		td.Authenticated = true
		// TODO - get more user information
	}
	td.Now = time.Now()
	return td
}

func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
