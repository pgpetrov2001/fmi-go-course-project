package utils

import (
	"context"
	"course-project/app"
	"fmt"
	"net/http"
)

func GetPlaygroundMiddleware(d app.DAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playgroundId, _ := r.Context().Value("playgroundId").(uint)
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
		reviewId, _ := r.Context().Value("reviewId").(uint)
		review, err := d.GetReview(reviewId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get review with id %u", reviewId), http.StatusBadRequest)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "review", review))
		next.ServeHTTP(w, req)
	})
}
