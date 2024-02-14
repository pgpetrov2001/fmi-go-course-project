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
	mapTemplate := template.Must(template.New("index.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/index.tmpl.html", "./templates/common.tmpl.html"))
	loginTemplate := template.Must(template.New("login.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/login.tmpl.html", "./templates/common.tmpl.html"))
	registerTemplate := template.Must(template.New("register.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/register.tmpl.html", "./templates/common.tmpl.html"))
	playgroundsTemplate := template.Must(template.New("playgrounds.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/playgrounds.tmpl.html", "./templates/playground.tmpl.html", "./templates/reviews.tmpl.html", "./templates/common.tmpl.html"))
	usersTemplate := template.Must(template.New("users.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/users.tmpl.html", "./templates/common.tmpl.html"))
	playgroundGalleryTemplate := template.Must(template.New("playground_gallery.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/playground_gallery.tmpl.html", "./templates/common.tmpl.html"))
	profileTemplate := template.Must(template.New("profile.tmpl.html").Funcs(utils.TemplateFuncMap).ParseFiles("./templates/profile.tmpl.html", "./templates/common.tmpl.html"))

	r := mux.NewRouter()
	r.Use(utils.ComposeMiddlewares(utils.LoggingMiddleware, utils.ParseJSONMiddleware))

	r.Handle("/", http.RedirectHandler("/map", http.StatusPermanentRedirect))
	r.Handle("/map", utils.GetUserMiddleware(w.Dao, routes.MapDataMiddleware(routes.RenderTemplate(mapTemplate)))).Methods("GET")
	r.Handle("/playgrounds", utils.GetUserMiddleware(w.Dao, routes.PlaygroundsDataMiddleware(w.Dao, routes.RenderTemplate(playgroundsTemplate)))).Methods("GET")
	r.Handle("/playground/{playgroundId}/gallery", utils.GetUserMiddleware(w.Dao, utils.GetPlaygroundMiddleware(w.Dao, routes.PlaygroundGalleryDataMiddleware(w.Dao, routes.RenderTemplate(playgroundGalleryTemplate))))).Methods("GET")
	r.Handle("/users", utils.AccessRightsMiddleware(w.Dao, true, false, routes.UsersDataMiddleware(w.Dao, routes.RenderTemplate(usersTemplate)))).Methods("GET")
	r.Handle("/profile", utils.AccessRightsMiddleware(w.Dao, false, true, routes.ProfileDataMiddleware(w.Dao, routes.RenderTemplate(profileTemplate)))).Methods("GET")
	r.Handle("/sign-in", utils.GetUserMiddleware(w.Dao, routes.SignInDataMiddleware(routes.RenderTemplate(loginTemplate)))).Methods("GET")
	r.Handle("/sign-up", utils.GetUserMiddleware(w.Dao, routes.SignUpDataMiddleware(routes.RenderTemplate(registerTemplate)))).Methods("GET")
	r.Handle("/api/login", w.WebappWrapper(routes.Login)).Methods("POST")
	r.Handle("/api/register", w.WebappWrapper(routes.Register)).Methods("POST")
	r.Handle("/api/logout", w.WebappWrapper(routes.Logout)).Methods("POST")
	r.Handle("/api/users", utils.AccessRightsMiddleware(w.Dao, true, false, w.WebappWrapper(routes.GetUsers))).Methods("GET")
	r.Handle("/api/users/{userId}", utils.UserAccessRightsMiddleware(w.Dao, w.WebappWrapper(routes.GetUser))).Methods("GET")
	r.Handle("/api/users/{userId}", utils.UserAccessRightsMiddleware(w.Dao, w.WebappWrapper(routes.PatchUser))).Methods("PATCH")
	r.Handle("/api/users/{userId}", utils.UserAccessRightsMiddleware(w.Dao, w.WebappWrapper(routes.DeleteUser))).Methods("DELETE")
	r.Handle("/api/users", utils.AccessRightsMiddleware(w.Dao, true, false, w.WebappWrapper(routes.PostUser))).Methods("POST")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.GetPlayground)))).Methods("GET")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.PatchPlayground)))).Methods("PATCH")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.DeletePlayground)))).Methods("DELETE")
	r.Handle("/api/playground/{playgroundId}", utils.AccessRightsMiddleware(w.Dao, true, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.PostPlayground)))).Methods("POST")
	r.Handle("/api/playground/{playgroundId}/review", utils.AccessRightsMiddleware(w.Dao, false, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.ReviewPlayground)))).Methods("POST")
	r.Handle("/api/playground/{playgroundId}/upload", utils.AccessRightsMiddleware(w.Dao, false, false, utils.GetPlaygroundMiddleware(w.Dao, w.WebappWrapper(routes.UploadPlaygroundPhotos)))).Methods("POST")
	r.Handle("/api/pending_photos", utils.AccessRightsMiddleware(w.Dao, true, false, w.WebappWrapper(routes.PendingPhotos))).Methods("GET")
	r.Handle("/api/approve/{photoId}", utils.AccessRightsMiddleware(w.Dao, true, false, utils.GetPhotoMiddleware(w.Dao, w.WebappWrapper(routes.ApprovePhoto)))).Methods("POST")
	r.Handle("/api/review/{reviewId}/vote", utils.AccessRightsMiddleware(w.Dao, false, false, utils.GetReviewMiddleware(w.Dao, w.WebappWrapper(routes.VoteReview)))).Methods("POST")
	r.Handle("/api/photo/{photoId}/vote", utils.AccessRightsMiddleware(w.Dao, false, false, utils.GetPhotoMiddleware(w.Dao, w.WebappWrapper(routes.VotePhoto)))).Methods("POST")
	r.Handle("/api/photo/{photoId}", utils.AccessRightsMiddleware(w.Dao, true, false, utils.GetPhotoMiddleware(w.Dao, w.WebappWrapper(routes.PatchPhoto)))).Methods("PATCH")
	r.Handle("/img/{photoId}", utils.GetUserMiddleware(w.Dao, utils.GetPhotoMiddleware(w.Dao, w.WebappWrapper(routes.GetPhoto)))).Methods("GET")

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

	err = db.AutoMigrate(&entities.User{}, &entities.Playground{}, &entities.PlaygroundPhoto{}, &entities.PlaygroundReview{}, &entities.PhotoVote{}, &entities.ReviewVote{})
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

	//the following lines ingest the playground data from the geojson into the database:
	//err = IngestPlaygrounds(playgroundsApp.Dao, filepath.Join(workDir, "static", "detski_ploshtadki_old_new_26_sofpr_20190418.geojson"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = CPNSRun(&playgroundsApp)
	if err != nil {
		log.Fatalf("Could intialize CPNS: %v", err)
	}
}
