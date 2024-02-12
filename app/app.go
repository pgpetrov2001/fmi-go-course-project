package app

import (
	"course-project/entities"
	"net/http"
)

type DAO interface {
	Init() error
	CreateUser(email string, username string, pass string, administrator bool, banned bool) (*entities.User, error)
	Authenticate(email string, pass string) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	GetUser(userId uint) (entities.User, error)
	UpdateUser(user *entities.User) error
	DeleteUser(user *entities.User) error
	CreatePlayground(playground *entities.Playground) error
	GetPlayground(playgroundId uint) (entities.Playground, error)
	GetPlaygrounds() ([]entities.Playground, error)
	UpdatePlayground(playground *entities.Playground) error
	DeletePlayground(playground *entities.Playground) error
	ReviewPlayground(review *entities.PlaygroundReview) error
	PlaygroundGallery(playgroundId uint) ([]entities.PlaygroundPhoto, error)
	PendingPhotos() ([]entities.PlaygroundPhoto, error)
	UploadPhoto(photo *entities.PlaygroundPhoto, filename string, data []byte) error
	GetPhoto(photoId uint) (entities.PlaygroundPhoto, error)
	UpdatePhoto(photo *entities.PlaygroundPhoto) error
	GetReview(reviewId uint) (entities.PlaygroundReview, error)
	UpdateReview(review *entities.PlaygroundReview) error
	VoteReview(review *entities.ReviewVote) error
	VotePhoto(photo *entities.PhotoVote) error
}

type WebApp struct {
	Server *http.Server
	Mux    *http.ServeMux
	Dao    DAO
}

type WebappHandler func(*WebApp, http.ResponseWriter, *http.Request)

func (a *WebApp) WebappWrapper(f WebappHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(a, w, r)
	})
}
