package routes

import (
	"context"
	"course-project/app"
	"course-project/entities"
	"course-project/utils"
	"encoding/json"
	"fmt"
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
		var playgroundData []byte
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get playground with id %u", playgroundId), http.StatusBadRequest)
			return
		} else if playgroundData, err = json.Marshal(playground); err != nil {
			http.Error(w, "Error while serializing playground", http.StatusInternalServerError)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "data", string(playgroundData)))
		next.ServeHTTP(w, req)
	})
}

func GetPlayground(a *app.WebApp, w http.ResponseWriter, r *http.Request) {
	playgroundData := r.Context().Value("playground").(string)
	w.WriteHeader(200)
	w.Write([]byte(playgroundData))
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
	playground, _ := utils.ExtractPlayground(r)
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
	playground, _ := utils.ExtractPlayground(r)
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
	user, _ := utils.ExtractUser(r)
	playground, _ := utils.ExtractPlayground(r)
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
