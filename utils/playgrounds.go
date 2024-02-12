package utils

import (
	"context"
	"course-project/app"
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
		if err != nil {
			http.Error(w, "Could not provide data for requested playground", http.StatusInternalServerError)
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), "playground", playground))
		next.ServeHTTP(w, req)
	})
}
