package app

import (
	"course-project/entities"
	"net/http"
)

type DAO interface {
	CreateUser(email string, username string, pass string) error
	Authenticate(email string, pass string) (bool, error)
	GetUsers() ([]entities.User, error)
	GetUser(id int) (entities.User, error)
	GetPlayground(id int) (entities.Playground, error)
	GetPlaygrounds() ([]entities.PlaygroundReview, error)
}

type WebApp struct {
	Server *http.Server
	Mux    *http.ServeMux
	Dao    DAO
}

type WebappHandler func(*WebApp, http.ResponseWriter, *http.Request)
type Handler func(http.ResponseWriter, *http.Request)

type WebappMiddleware func(*WebApp, http.ResponseWriter, *http.Request)

func (a *WebApp) WebappWrapper(f WebappHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		f(a, w, r)
	}
}
