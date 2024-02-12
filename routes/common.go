package routes

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := r.Context().Value("data").(map[string]interface{})
		if !ok {
			data = make(map[string]interface{})
		}
		w.WriteHeader(200)
		err := t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
