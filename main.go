package main

import (
	"course-project/app"
	"course-project/dao"
	"course-project/entities"
	"course-project/routes"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

var workDir, _ = os.Getwd()

func CPNSRun(w *app.WebApp) error {
	mapTemplate := template.Must(template.ParseFiles("./templates/index.tmpl.html"))
	loginTemplate := template.Must(template.ParseFiles("./templates/login.tmpl.html"))
	registerTemplate := template.Must(template.ParseFiles("./templates/register.tmpl.html"))

	w.Mux.Handle("/", http.RedirectHandler("/map", http.StatusPermanentRedirect))
	w.Mux.HandleFunc("/map", routes.RenderTemplate(mapTemplate))
	w.Mux.HandleFunc("/api/login", w.WebappWrapper(routes.Login))
	w.Mux.HandleFunc("/api/register", w.WebappWrapper(routes.Register))
	w.Mux.HandleFunc("/api/logout", w.WebappWrapper(routes.Logout))
	w.Mux.Handle("/sign-in", routes.SignInDataMiddleware(routes.RenderTemplate(loginTemplate)))
	w.Mux.Handle("/sign-up", routes.SignUpDataMiddleware(routes.RenderTemplate(registerTemplate)))

	fileServer := http.FileServer(http.Dir(path.Join(workDir, "static")))
	w.Mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Printf("Listening on port %s\n", w.Server.Addr)
	if err := w.Server.ListenAndServe(); err != nil {
		return fmt.Errorf("Failed to start service on port %s:%w", w.Server.Addr, err)
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	dsn := "root:root@tcp(127.0.0.1:3306)/CPNS?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Playground{})
	db.AutoMigrate(&entities.PlaygroundPhoto{})
	db.AutoMigrate(&entities.PlaygroundReview{})

	playgroundsApp := app.WebApp{
		Server: server,
		Mux:    mux,
		Dao:    &dao.CPNS{Db: db},
	}

	err = CPNSRun(&playgroundsApp)
	if err != nil {
		log.Fatalf("Could intialize CPNS: %v", err)
	}
}
