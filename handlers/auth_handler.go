package handlers

import (
	"gohttp/auth"
	"gohttp/views"
	"net/http"
	"strconv"
)

func LoginHandler(store *auth.MemorySessionStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			views.RenderTemplate(w, "login", nil)
		} else if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			rememberMe, _ := strconv.ParseBool(r.FormValue("rememberMe"))
	
			if username == "admin" && password == "admin" {
				sess := &auth.AuthSession{}
				sess.IsAuthenticated = true
				sess.RememberMe = rememberMe
				sess.Role = "Administrator"
				sess.Username = username
	
				store.PutSession(w, r, sess)
	
				params := r.URL.Query()
				location := params.Get("redirect")
	
				if len(params["redirect"]) > 0 {
					http.Redirect(w, r, location, http.StatusFound)
					return
				}
	
				http.Redirect(w, r, store.Base.LoginPath, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				views.RenderTemplate(w, "error", nil)
			}
		}

	})	
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	MemorySession.DeleteSession(r)
	http.Redirect(w, r, MemorySession.Base.LoginPath, http.StatusFound)
}
