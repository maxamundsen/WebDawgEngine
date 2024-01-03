package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

type PageData struct {
	Person    Person
	Title     string
	IsGreater bool
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	val1, _ := strconv.Atoi(r.FormValue("val1"))
	val2, _ := strconv.Atoi(r.FormValue("val2"))

	sess := globalSession.GetSessionFromCtx(r)
	
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	age := val1 * val2
	isGreater := false

	if sess.Role != "Administrator" {
		isGreater = true
	}

	harambe := Person{
		"Harambe",
		"Monke",
		age,
	}

	t := template.Must(template.ParseFS(viewTemplates, "views/base.html", "views/test.html"))

	pageData := PageData{
		harambe,
		"Title for page",
		isGreater,
	}

	t.Execute(w, pageData)
}
