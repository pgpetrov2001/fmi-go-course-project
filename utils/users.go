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

func GetUserMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "sessionID")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userIdVal := session.Values["userId"]

		if userIdVal == nil {
			next.ServeHTTP(w, r)
			return
		}

		userId := userIdVal.(uint)
		user, err := d.GetUser(userId)
		if err != nil {
			log.Printf("Dangling session for user id %d\n", userId)
			next.ServeHTTP(w, r)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "user", user))
		next.ServeHTTP(w, req)
	})
}

func UserAccessRightsMiddleware(d app.DAO, next http.Handler) http.Handler {
	return GetUserMiddleware(d, GetIdParamMiddleware("userId", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func AccessRightsMiddleware(d app.DAO, admin bool, redirectOnFail bool, next http.Handler) http.Handler {
	return GetUserMiddleware(d, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(entities.User)
		if !ok {
			log.Printf("not ok\n")
			//missing value means user is not logged in
			if admin {
				http.Error(w, "You don't have access rights for this page.", http.StatusForbidden)
			} else if redirectOnFail {
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			} else {
				http.Error(w, "You need to be logged in to perform this action", http.StatusForbidden)
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
