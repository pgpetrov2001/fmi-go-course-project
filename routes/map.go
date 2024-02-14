package routes

import (
	"context"
	"course-project/entities"
	"net/http"
)

func MapDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(entities.User)
		data := map[string]interface{}{
			"page": "map",
		}
		if ok {
			data["user"] = user
		}
		req := r.WithContext(context.WithValue(r.Context(), "data", data))
		next.ServeHTTP(w, req)
	})
}
