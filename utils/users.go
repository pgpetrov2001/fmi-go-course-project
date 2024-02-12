package utils

import (
	"context"
	"course-project/app"
	"course-project/entities"
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

		user, ok := session.Values["user"].(*entities.User)
		if !ok || user == nil {
			next.ServeHTTP(w, r)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "user", *user))
		next.ServeHTTP(w, req)
	})
}

func UserAccessRightsMiddleware(next http.Handler) http.Handler {
	return GetUserMiddleware(GetIdParamMiddleware("userId", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(entities.User) // logged user
		userId, _ := r.Context().Value("userId").(uint)       // requested user

		forbidden := false
		if !ok {
			forbidden = true
		} else if !user.Administrator && user.ID != userId {
			log.Printf("%t %v %v", user.Administrator, user.ID, userId)
			forbidden = true
		}

		if forbidden {
			http.Error(w, "You are not allowed to access information about this user", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})))
}

func AccessRightsMiddleware(d app.DAO, admin bool, next http.Handler) http.Handler {
	return GetUserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(entities.User)
		if !ok {
			//missing value means user is not logged in
			if admin {
				http.Error(w, "You don't have access rights for this page.", http.StatusForbidden)
			} else {
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			}
			return
		}

		if admin && !user.Administrator {
			http.Error(w, "You don't have access rights for this page.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}
