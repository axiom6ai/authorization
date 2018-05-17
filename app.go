package main

import (
	"log"
	"net/http"

	"authorization/controller"
	"authorization/model"

	_ "github.com/lib/pq"

	"github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Controller interface {
	InitializeRoutes()
}

type App struct {
	Router        *mux.Router
	DB            *gorm.DB
	err           error
	JwtMiddleware *jwtmiddleware.JWTMiddleware
	controllers   []Controller
}

func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.initializeDB()
	app.JwtMiddleware = InitJwtMiddleware()
	app.controllers = append(app.controllers, controller.UserCtrl{app.Router, app.DB, app.JwtMiddleware})
	app.controllers = append(app.controllers, controller.CampCtrl{app.Router, app.DB, app.JwtMiddleware})
	app.controllers = append(app.controllers, controller.StudentCtrl{app.Router, app.DB, app.JwtMiddleware})
	for _, controller := range app.controllers {
		controller.InitializeRoutes()
	}
}

func (app *App) Run(addr string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Content-Type", "X-Auth-Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(originsOk, headersOk, methodsOK)(app.Router)))
}

func (app *App) initializeDB() {
	app.DB, app.err = gorm.Open("postgres", "user=postgres password=postgres port=5432 dbname=postgres sslmode=disable")
	if app.err != nil {
		panic(app.err.Error())
	}

	app.DB.AutoMigrate(
		&model.Campus{},
		&model.Student{},
		&model.User{},
	)
}
