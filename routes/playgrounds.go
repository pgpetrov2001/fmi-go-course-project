package routes

import (
	"context"
	"course-project/app"
	"course-project/entities"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func PlaygroundsDataMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUser := r.Context().Value("user")
		var user *entities.User
		if rawUser != nil {
			tmp := rawUser.(entities.User)
			user = &tmp
		}
		playgrounds, err := d.GetPlaygrounds()
		if err != nil {
			http.Error(w, "Could not fetch playgrounds for some reason", http.StatusInternalServerError)
			return
		}
		for i, playground := range playgrounds {
			playgrounds[i].SelectedPhotos = make([]entities.PlaygroundPhoto, 0)
			for _, photo := range playground.Photos {
				if (user != nil && user.Administrator) || (photo.Approved != nil && *photo.Approved && photo.Selected) {
					playgrounds[i].SelectedPhotos = append(playgrounds[i].SelectedPhotos, photo)
				}
			}
		}
		templateData := map[string]interface{}{
			"playgrounds": playgrounds,
		}
		req := r.WithContext(context.WithValue(r.Context(), "data", templateData))
		next.ServeHTTP(w, req)
	})
}

func GetPlayground(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	playground := r.Context().Value("playground").(entities.Playground)
	playgroundData, err := json.Marshal(playground)
	if err != nil {
		http.Error(w, "Could not serialize playground", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(playgroundData)
}

func parsePlaygroundForm(r *http.Request, playground *entities.Playground) error {
	siteNumber := r.Form.Get("sitenumber")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	area := r.Form.Get("area")
	location := r.Form.Get("location")
	ownership := r.Form.Get("ownership")

	var errors map[string]error

	if siteNumber != "" {
		playground.SiteNumber = siteNumber
	}
	if latitude != "" {
		var val float64
		val, errors["Latitude"] = strconv.ParseFloat(latitude, 64)
		playground.Latitude = val
	}
	if longitude != "" {
		var val float64
		val, errors["Longitude"] = strconv.ParseFloat(longitude, 64)
		playground.Longitude = val
	}
	if area != "" {
		var val int
		val, errors["Area"] = strconv.Atoi(area)
		playground.Area = val
	}
	if location != "" {
		playground.Location = location
	}
	if ownership != "" {
		playground.Ownership = ownership
	}

	if len(errors) > 0 {
		errorStr := "Error parsing fields for playground.\nErrors:\n"
		for key, val := range errors {
			errorStr += fmt.Sprintf("- Error parsing %v: %v", key, val)
		}
		return fmt.Errorf(errorStr)
	}
	return nil
}

func PatchPlayground(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	playground, _ := r.Context().Value("playground").(entities.Playground)
	err := parsePlaygroundForm(r, &playground)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = a.Dao.UpdatePlayground(&playground)
	if err != nil {
		http.Error(w, "Could not update playground", http.StatusInternalServerError)
		return
	}

	playgroundData, err := json.Marshal(playground)
	w.WriteHeader(200)
	w.Write(playgroundData)
}

func DeletePlayground(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	playground, _ := r.Context().Value("playground").(entities.Playground)
	err := a.Dao.DeletePlayground(&playground)
	if err != nil {
		http.Error(w, "Could not delete playground", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Successfully deleted playground."))
}

func PostPlayground(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	var playground entities.Playground
	err := parsePlaygroundForm(r, &playground)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.Dao.CreatePlayground(&playground)

	if err != nil {
		http.Error(w, "Could not create playground", http.StatusInternalServerError)
		return
	}

	playgroundData, err := json.Marshal(playground)
	w.WriteHeader(200)
	w.Write(playgroundData)
}

func ReviewPlayground(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(entities.User)
	playground, _ := r.Context().Value("playground").(entities.Playground)
	starsVal := r.Form.Get("stars")
	stars, err := strconv.Atoi(starsVal)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid value for stars %d", stars), http.StatusBadRequest)
		return
	}
	review := entities.PlaygroundReview{
		PlaygroundID: playground.ID,
		Playground:   playground,
		UserID:       user.ID,
		User:         user,
		Stars:        stars,
		Content:      r.Form.Get("content"),
	}
	err = a.Dao.ReviewPlayground(&review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func PlaygroundGallery(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	playgroundId, _ := r.Context().Value("playgroundId").(uint)
	photos, err := a.Dao.PlaygroundGallery(playgroundId)
	if err != nil {
		http.Error(w, "Could not fetch gallery for playground", http.StatusInternalServerError)
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

func VoteReview(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(entities.User)
	review, _ := r.Context().Value("review").(entities.PlaygroundReview)
	up := r.Form.Get("up") == "true"
	vote := entities.ReviewVote{
		Up:                 up,
		Review:             review,
		User:               user,
		PlaygroundReviewID: review.ID,
		UserID:             user.ID,
	}
	err := a.Dao.VoteReview(&vote)
	if err != nil {
		http.Error(w, "Could not place vote on review", http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte{})
}

func VotePhoto(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(entities.User)
	photo, _ := r.Context().Value("photo").(entities.PlaygroundPhoto)
	up := r.Form.Get("up") == "true"
	vote := entities.PhotoVote{
		Up:                up,
		Photo:             photo,
		User:              user,
		PlaygroundPhotoID: photo.ID,
		UserID:            user.ID,
	}
	err := a.Dao.VotePhoto(&vote)
	if err != nil {
		http.Error(w, "Could not place vote on review", http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte{})
}

func UploadPlaygroundPhotos(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(entities.User)
	playground, _ := r.Context().Value("playground").(entities.Playground)

	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB max memory
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]

	var photos []entities.PlaygroundPhoto

	for _, file := range files {
		uploadedFile, err := file.Open()
		if err != nil {
			http.Error(w, "Unable to open uploaded file", http.StatusBadRequest)
			return
		}
		defer uploadedFile.Close()

		fileBytes, err := io.ReadAll(uploadedFile)
		if err != nil {
			http.Error(w, "Error reading file bytes", http.StatusInternalServerError)
			return
		}

		approved := new(bool)
		*approved = false
		photo := entities.PlaygroundPhoto{
			Playground:   playground,
			PlaygroundID: playground.ID,
			User:         user,
			UserId:       user.ID,
			Approved:     approved,
			Selected:     false,
		}
		log.Printf("Uploading photo %s with %d bytes\n", file.Filename, len(fileBytes))
		err = a.Dao.UploadPhoto(&photo, file.Filename, fileBytes)
		if err != nil {
			http.Error(w, "Could not upload one of the photos right now", http.StatusInternalServerError)
			return
		}
		photos = append(photos, photo)
	}

	photosData, err := json.Marshal(photos)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(photosData)
}
