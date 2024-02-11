package routes

import (
	"course-project/app"
	"course-project/utils"
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
	id, err := utils.GetIdFromRouteParam(w, r, "photoId")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo, err := a.Dao.GetPhoto(id)
	if err != nil {
		http.Error(w, "Could not fetch photo info", http.StatusInternalServerError)
		return
	}
	*photo.Approved = true // TODO: beware of nil pointer dereference
	err = a.Dao.UpdatePhoto(&photo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not update photo: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte{})
}
