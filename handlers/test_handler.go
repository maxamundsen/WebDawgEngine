package handlers

import (
	"gohttp/auth"
	"gohttp/views"
	"net/http"
)

type testPageModel struct {
	Base     views.ViewBase
	Password string
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	user := sessionStore.GetIdentityFromCtx(r)

	val1 := r.FormValue("val1")

	var password string

	if val1 == "" {
		password = "empty"
	} else {
		password, _ = auth.HashPassword(val1)
	}

	viewData := make(map[string]interface{})

	viewData["Title"] = "Test Page"

	base := views.NewViewBase(user, viewData)

	pageData := testPageModel{
		base,
		password,
	}

	views.RenderTemplate(w, "test", pageData)
}
