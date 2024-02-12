package utils

import (
	"context"
	"course-project/app"
	"fmt"
	"net/http"
)

func GetPhotoMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		photoId, _ := r.Context().Value("photoId").(uint)
		photo, err := d.GetPhoto(photoId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get photo with id %u", photoId), http.StatusBadRequest)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "photo", photo))
		next.ServeHTTP(w, req)
	})
}
