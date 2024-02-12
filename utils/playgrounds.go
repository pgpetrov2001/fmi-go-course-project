package utils

import (
	"context"
	"course-project/app"
	"course-project/entities"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetPlaygroundMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playgroundId, err := GetIdFromRouteParam(w, r, "playgroundId")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		playground, err := d.GetPlayground(playgroundId)
		var playgroundData []byte
		fail := false
		if err != nil {
			fail = true
		} else if playgroundData, err = json.Marshal(playground); err != nil {
			fail = true
		}
		if fail {
			http.Error(w, "Could not provide data for requested playground", http.StatusInternalServerError)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "playground", playgroundData))
		next.ServeHTTP(w, req)
	})
}

func ExtractPlayground(r *http.Request) (entities.Playground, error) {
	var playground entities.Playground
	playgroundData, ok := r.Context().Value("playground").(string)
	if !ok {
		return playground, fmt.Errorf("No playground found attached to context of provided request, here's the context: %v", r.Context())
	}
	err := json.Unmarshal([]byte(playgroundData), &playground)
	return playground, fmt.Errorf("Could not parse attached playground to request context: %v", err)
}

func ExtractReview(r *http.Request) (entities.PlaygroundReview, error) {
	var review entities.PlaygroundReview
	reviewData, ok := r.Context().Value("review").(string)
	if !ok {
		return review, fmt.Errorf("No review found attached to context of provided request, here's the context: %v", r.Context())
	}
	err := json.Unmarshal([]byte(reviewData), &review)
	return review, fmt.Errorf("Could not parse attached review to request context: %v", err)
}

func ExtractPhoto(r *http.Request) (entities.PlaygroundPhoto, error) {
	var photo entities.PlaygroundPhoto
	photoData, ok := r.Context().Value("photo").(string)
	if !ok {
		return photo, fmt.Errorf("No photo found attached to context of provided request, here's the context: %v", r.Context())
	}
	err := json.Unmarshal([]byte(photoData), &photo)
	return photo, fmt.Errorf("Could not parse attached photo to request context: %v", err)
}
