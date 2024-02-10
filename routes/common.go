package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	data string
}

func EmptyDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type DataMiddleware func(http.ResponseWriter, *http.Request) (interface{}, error)

func RenderTemplate(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawData, ok := r.Context().Value("data").(string)
		var data map[string]interface{}
		if ok {
			err := json.Unmarshal([]byte(rawData), &data)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			data = make(map[string]interface{})
		}
		w.WriteHeader(200)
		err := t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
