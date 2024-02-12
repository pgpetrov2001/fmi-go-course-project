package utils

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func GetIdParamMiddleware(param string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idVal, err := strconv.ParseUint(vars[param], 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Specified value for %s is not a valid id", param), http.StatusBadRequest)
			return
		}
		id := uint(idVal)
		req := r.WithContext(context.WithValue(r.Context(), param, id))
		next.ServeHTTP(w, req)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
