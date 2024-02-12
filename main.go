package main

import (
	"course-project/app"
	"course-project/dao"
	"course-project/entities"
	"course-project/routes"
	"course-project/utils"
	"fmt"
	"github.com/gorilla/mux"
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
	playgroundsTemplate := template.Must(template.ParseFiles("./templates/playgrounds.tmpl.html"))

	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("/map", http.StatusPermanentRedirect))
	r.HandleFunc("/map", routes.RenderTemplate(mapTemplate)).Methods("GET")
	r.Handle("/playgrounds", routes.PlaygroundsDataMiddleware(w.Dao, routes.RenderTemplate(playgroundsTemplate))).Methods("GET")
	r.Handle("/sign-in", routes.SignInDataMiddleware(routes.RenderTemplate(loginTemplate))).Methods("GET")
	r.Handle("/sign-up", routes.SignUpDataMiddleware(routes.RenderTemplate(registerTemplate))).Methods("GET")
	r.Handle("/api/login", w.WebappWrapper(routes.Login)).Methods("POST")
	r.Handle("/api/register", w.WebappWrapper(routes.Register)).Methods("POST")
	r.Handle("/api/logout", w.WebappWrapper(routes.Logout)).Methods("POST")
	r.Handle("/api/users", utils.AccessRightsMiddleware(w.Dao, true, w.WebappWrapper(routes.GetUsers))).Methods("GET")
	r.Handle("/api/users/{userId}", utils.UserAccessRightsMiddleware(w.WebappWrapper(routes.GetUser))).Methods("GET")
	r.Handle("/api/users/{userId}", utils.UserAccessRightsMiddleware(w.WebappWrapper(routes.PatchUser))).Methods("PATCH")
	r.Handle("/api/users/{userId}", utils.UserAccessRightsMiddleware(w.WebappWrapper(routes.DeleteUser))).Methods("DELETE")
	r.Handle("/api/users", utils.AccessRightsMiddleware(w.Dao, true, w.WebappWrapper(routes.PostUser))).Methods("POST")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.GetPlayground)))).Methods("GET")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.PatchPlayground)))).Methods("PATCH")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.DeletePlayground)))).Methods("DELETE")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.PostPlayground)))).Methods("POST")
	r.Handle("/api/playground/{playgroundId}/review", utils.AccessRightsMiddleware(w.Dao, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.ReviewPlayground)))).Methods("POST")
	r.Handle("/api/playground/{playgroundId}/gallery", utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.PlaygroundGallery))).Methods("GET")
	r.Handle("/api/playground/{playgroundId}/upload", utils.AccessRightsMiddleware(w.Dao, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.UploadPlaygroundPhoto)))).Methods("POST")
	r.Handle("/api/pending_photos", utils.AccessRightsMiddleware(w.Dao, true, w.WebappWrapper(routes.PendingPhotos))).Methods("GET")
	r.Handle("/api/approve/{photoId}", utils.AccessRightsMiddleware(w.Dao, true, w.WebappWrapper(routes.ApprovePhoto))).Methods("POST")
	r.Handle("/api/review/{reviewId}/vote", utils.AccessRightsMiddleware(w.Dao, false, utils.GetReviewMiddleware(w.Dao, w.WebappWrapper(routes.VoteReview)))).Methods("POST")
	r.Handle("/api/photo/{photoId}/vote", utils.AccessRightsMiddleware(w.Dao, false, utils.GetPhotoMiddleware(w.Dao, w.WebappWrapper(routes.VotePhoto)))).Methods("POST")
	r.Handle("/img/{photoId}", utils.GetPhotoMiddleware(w.Dao, w.WebappWrapper(routes.GetPhoto))).Methods("GET")

	w.Mux.Handle("/", r)

	fileServer := http.FileServer(http.Dir(path.Join(workDir, "static")))
	w.Mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Printf("Listening on port %s\n", w.Server.Addr)
	if err := w.Server.ListenAndServe(); err != nil {
		return fmt.Errorf("Failed to start service on port %s:%w", w.Server.Addr, err)
	}
	return nil
}

func main() {
	multiplexer := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: multiplexer,
	}

	dsn := "root:root@tcp(127.0.0.1:3306)/CPNS?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	err = db.AutoMigrate(&entities.User{}, &entities.Playground{}, &entities.PlaygroundPhoto{}, &entities.PlaygroundReview{})
	if err != nil {
		log.Fatalf("Auto-Migration error: %v", err)
	}

	playgroundsApp := app.WebApp{
		Server: server,
		Mux:    multiplexer,
		Dao:    &dao.CPNS{Db: db, FSStoragePath: path.Join(workDir, "data")},
	}
	err = playgroundsApp.Dao.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = CPNSRun(&playgroundsApp)
	if err != nil {
		log.Fatalf("Could intialize CPNS: %v", err)
	}
}
