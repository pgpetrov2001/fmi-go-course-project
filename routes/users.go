package routes

import (
	"context"
	"course-project/app"
	"course-project/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Should be protected with admin acccess rights
func GetUsers(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	users, err := a.Dao.GetUsers()
	if err != nil {
		http.Error(w, "A problem ocurred while fetching users: Database isn't functioning properly, please contact support.", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Could not serialize users, please contact support.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(resp)
}

func GetUser(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userId").(uint)
	user, err := a.Dao.GetUser(userId)
	var userData []byte
	fail := false
	if err != nil {
		fail = true
	} else if userData, err = json.Marshal(user); err != nil {
		fail = true
	}
	if fail {
		http.Error(w, "Could not provide data for requested user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write(userData)
}

func PatchUser(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userId").(uint)
	user, err := a.Dao.GetUser(userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Did not find user with id %d", userId), http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	administrator := r.Form.Get("administrator")
	banned := r.Form.Get("banned")

	if email != "" {
		user.Email = email
	}
	if username != "" {
		user.Username = username
	}
	if password != "" {
		user.PasswordHash, err = utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Something is wrong with the password", http.StatusBadRequest)
			return
		}
	}
	if administrator != "" {
		user.Administrator = administrator == "true"
	}
	if banned != "" {
		user.Banned = banned == "true"
	}

	err = a.Dao.UpdateUser(&user)
	if err != nil {
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	userData, _ := json.Marshal(user)
	w.WriteHeader(200)
	w.Write(userData)
}

func DeleteUser(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userId").(uint)
	user, err := a.Dao.GetUser(userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Did not find user with id %d", userId), http.StatusBadRequest)
		return
	}
	err = a.Dao.DeleteUser(&user)
	if err != nil {
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte{})
}

func PostUser(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	administrator := r.Form.Get("administrator") == "on"
	banned := r.Form.Get("banned") == "on"
	user, err := a.Dao.CreateUser(email, username, password, administrator, banned)

	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	userData, _ := json.Marshal(user)
	w.WriteHeader(200)
	w.Write(userData)
}

func Login(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	session, err := utils.Store.Get(r, "sessionID")
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

	user, err := a.Dao.Authenticate(email, password)
	if err != nil {
		http.Error(w, "Error while trying to extract user data", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Redirect(w, r, fmt.Sprintf("/sign-in?error=%s", "Authentication failed!"), http.StatusSeeOther)
		return
	}

	session.Values["user"] = user
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(_ *app.WebApp, w http.ResponseWriter, r *http.Request) {
	session, err := utils.Store.Get(r, "sessionID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = nil
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
	_, err = a.Dao.CreateUser(email, username, password, false, false)

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
		data := map[string]interface{}{
			"error": errVal,
		}
		req := r.WithContext(context.WithValue(r.Context(), "data", data))
		next.ServeHTTP(w, req)
	})
}

func SignUpDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errVal := r.URL.Query().Get("error")
		data := map[string]interface{}{
			"error": errVal,
		}
		req := r.WithContext(context.WithValue(r.Context(), "data", data))
		next.ServeHTTP(w, req)
	})
}
