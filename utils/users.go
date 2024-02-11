package utils

import (
	"context"
	"course-project/app"
	"course-project/entities"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var Store = sessions.NewCookieStore([]byte("asdkjhasd-asfj-secret-key-hdsf-sdflkshdfl"))

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "sessionID")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userData, ok := session.Values["user"].(string)
		if !ok || userData == "" {
			next.ServeHTTP(w, r)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "user", userData))
		next.ServeHTTP(w, req)
	})
}

func UserAccessRightsMiddleware(next http.Handler) http.Handler {
	return GetUserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, ok := r.Context().Value("user").(string)
		userId, err := GetIdFromRouteParam(w, r, "userId")
		if err != nil {
			return
		}

		forbidden := false
		if !ok {
			forbidden = true
		} else {
			var user entities.User
			err := json.Unmarshal([]byte(userData), &user)
			if err != nil {
				log.Printf("json?\n")
				forbidden = true
			} else if !user.Administrator && user.ID != userId {
				log.Printf("%t %v %v", user.Administrator, user.ID, userId)
				forbidden = true
			}
		}

		if forbidden {
			http.Error(w, "You are not allowed to access information about this user", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func AccessRightsMiddleware(d app.DAO, admin bool, next http.Handler) http.Handler {
	return GetUserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, ok := r.Context().Value("user").(string)
		if !ok {
			//missing value means user is not logged in
			if admin {
				http.Error(w, "You don't have access rights for this page.", http.StatusForbidden)
			} else {
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			}
			return
		}

		var user entities.User
		err := json.Unmarshal([]byte(userData), &user)

		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		if admin && !user.Administrator {
			http.Error(w, "You don't have access rights for this page.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func ExtractUser(r *http.Request) (entities.User, error) {
	var user entities.User
	userData, ok := r.Context().Value("user").(string)
	if !ok {
		return user, fmt.Errorf("No user found attached to context of provided request, here's the context: %v", r.Context())
	}
	err := json.Unmarshal([]byte(userData), &user)
	return user, fmt.Errorf("Could not parse attached user to request context: %v", err)
}
