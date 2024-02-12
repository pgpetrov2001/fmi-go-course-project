package routes

import (
	"course-project/app"
	"course-project/entities"
	"encoding/json"
	"fmt"
	"net/http"
)

func PendingPhotos(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	photos, err := a.Dao.PendingPhotos()
	if err != nil {
		http.Error(w, "Could not fetch pending photos", http.StatusInternalServerError)
		return
	}
	photosData, err := json.Marshal(photos)
	if err != nil {
		http.Error(w, "Could not serialize photos", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(photosData)
}

func ApprovePhoto(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	photo, _ := r.Context().Value("photo").(entities.PlaygroundPhoto)
	*photo.Approved = true // TODO: beware of nil pointer dereference
	err := a.Dao.UpdatePhoto(&photo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not update photo: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte{})
}

func GetPhoto(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	photo, _ := r.Context().Value("photo").(entities.PlaygroundPhoto)
	rawUser := r.Context().Value("user")
	var user *entities.User
	if rawUser != nil {
		tmp := rawUser.(entities.User)
		user = &tmp
	}
	allow := false
	if user != nil && user.Administrator {
		allow = true
	}
	if user != nil && photo.UserId == user.ID {
		allow = true
	}
	if photo.Approved != nil && *photo.Approved {
		allow = true
	}
	if !allow {
		http.Error(w, "Not allowed to view photo", http.StatusForbidden)
		return
	}
	data, err := a.Dao.GetPhotoContents(&photo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not retrieve photo with id %u"), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
