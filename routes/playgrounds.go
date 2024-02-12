package routes

import (
	"context"
	"course-project/app"
	"course-project/entities"
	"course-project/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetPlaygroundMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playgroundId, err := utils.GetIdFromRouteParam(w, r, "playgroundId")
		if err != nil {
			http.Error(w, "Invalid value for playgroundId", http.StatusBadRequest)
			return
		}

		playground, err := d.GetPlayground(playgroundId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get playground with id %u", playgroundId), http.StatusBadRequest)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "playground", playground))
		next.ServeHTTP(w, req)
	})
}

func GetReviewMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reviewId, err := utils.GetIdFromRouteParam(w, r, "reviewId")
		if err != nil {
			http.Error(w, "Invalid value for reviewId", http.StatusBadRequest)
			return
		}

		review, err := d.GetReview(reviewId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get review with id %u", reviewId), http.StatusBadRequest)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "review", review))
		next.ServeHTTP(w, req)
	})
}

func GetPhotoMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		photoId, err := utils.GetIdFromRouteParam(w, r, "photoId")
		if err != nil {
			http.Error(w, "Invalid value for photoId", http.StatusBadRequest)
			return
		}

		photo, err := d.GetPhoto(photoId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get photo with id %u", photoId), http.StatusBadRequest)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "photo", photo))
		next.ServeHTTP(w, req)
	})
}

func PlaygroundsDataMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playgrounds, err := d.GetPlaygrounds()
		if err != nil {
			http.Error(w, "Could not fetch playgrounds for some reason", http.StatusInternalServerError)
			return
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
		UserId:       user.ID,
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
	playgroundId, err := utils.GetIdFromRouteParam(w, r, "playgroundId")
	if err != nil {
		http.Error(w, "Invalid id specified", http.StatusBadRequest)
		return
	} //actually this parsing of id can be exported to a middleware function, because context supports any type
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

func UploadPlaygroundPhoto(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(entities.User)
	playground, _ := r.Context().Value("playground").(entities.Playground)

	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB max memory
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file into a byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file bytes", http.StatusInternalServerError)
		return
	}
	filename := header.Filename

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
	err = a.Dao.UploadPhoto(&photo, filename, fileBytes)
	if err != nil {
		http.Error(w, "Could not upload photo right now", http.StatusInternalServerError)
		return
	}
	photoData, err := json.Marshal(photo)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(photoData)
}
