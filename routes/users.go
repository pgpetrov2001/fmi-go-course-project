package routes

import (
	"context"
	"course-project/app"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("asdkjhasd-asfj-secret-key-hdsf-sdflkshdfl"))

func Login(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	ok, err := a.Dao.Authenticate(email, password)
	if err != nil {
		http.Error(w, "Error while trying to extract user data", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Redirect(w, r, fmt.Sprintf("/sign-in?error=%s", "Authentication failed!"), http.StatusSeeOther)
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(_ *app.WebApp, w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func Register(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	err = a.Dao.CreateUser(email, username, password)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/sign-up?error=%s", "Registration failed!"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func SignInDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errVal := r.URL.Query().Get("error")
		data, err := json.Marshal(map[string]interface{}{
			"error": errVal,
		})
		if err != nil {
			log.Fatal(err)
		}
		req := r.WithContext(context.WithValue(r.Context(), "data", string(data)))
		next.ServeHTTP(w, req)
	})
}

func SignUpDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errVal := r.URL.Query().Get("error")
		data, err := json.Marshal(map[string]interface{}{
			"error": errVal,
		})
		if err != nil {
			log.Fatal(err)
		}
		req := r.WithContext(context.WithValue(r.Context(), "data", string(data)))
		next.ServeHTTP(w, req)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "sessionID")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
