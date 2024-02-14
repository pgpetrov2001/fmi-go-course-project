package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
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

func ParseJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correctMethod := r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodPut
		if correctMethod && r.Header.Get("Content-Type") == "application/json" {
			var data map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "Invalid JSON in body", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "body", data)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func ComposeMiddlewares(a, b func(next http.Handler) http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a(b(next))
	}
}

var TemplateFuncMap = template.FuncMap{
	"log": func(s interface{}) bool {
		log.Println(s)
		return true
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"add": func(a, b int) int {
		return a + b
	},
	"mod": func(a, b int) int {
		return a % b
	},
	"mul": func(a, b int) int {
		return a * b
	},
	"safeHTMLAttr": func(attr, val string) template.HTMLAttr {
		return template.HTMLAttr(fmt.Sprintf("%s=\"%s\"", attr, val))
	},
	"numRange": func(i, j, k int) []int {
		if k == 0 {
			panic("numRange with step 0 called")
		}
		result := make([]int, 0)
		for l := i; (l < j) != (k < 0); l += k {
			result = append(result, l)
		}
		return result
	},
	"dict": func(values ...interface{}) map[string]interface{} {
		if len(values)%2 != 0 {
			panic("Invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key := values[i].(string)
			dict[key] = values[i+1]
		}
		return dict
	},
	"not": func(val bool) bool {
		return !val
	},
	"asfloat32": func(x int) float32 {
		return float32(x)
	},
	"asfloat64": func(x int) float64 {
		return float64(x)
	},
	"derefBool": func(v *bool) bool {
		return *v
	},
}
