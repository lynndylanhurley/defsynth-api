package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/signature"

	_ "github.com/lib/pq"
	"github.com/lynndylanhurley/defsynth-api/controllers"
	"github.com/lynndylanhurley/defsynth-api/models"
)

var db gorm.DB
var m *martini.ClassicMartini

func main() {
	// init db
	db, err := gorm.Open("postgres", "user=root dbname=defsynth-api sslmode=disable")

	if err != nil {
		panic(fmt.Sprintf("db Error: %v", err))
	}

	db.LogMode(true)

	// migrate db
	db.AutoMigrate(models.Component{})
	db.AutoMigrate(models.ModuleComponent{})
	db.AutoMigrate(models.SynthModule{})
	db.AutoMigrate(models.User{})

	// read keys from env
	envName := martini.Env
	envFileName := strings.Join([]string{".env.", envName}, "")
	envMap, err := godotenv.Read(envFileName)

	if err != nil {
		panic(fmt.Sprintf("@-->env vars load error: %v", err))
	}

	// init oauth providers
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		github.New(
			envMap["AUTH_GITHUB_KEY"],
			envMap["AUTH_GITHUB_SECRET"],
			envMap["AUTH_GITHUB_CALLBACK"],
		),
		facebook.New(
			envMap["AUTH_FACEBOOK_KEY"],
			envMap["AUTH_FACEBOOK_SECRET"],
			envMap["AUTH_FACEBOOK_CALLBACK"],
		),
		google.New(
			envMap["AUTH_GOOGLE_KEY"],
			envMap["AUTH_GOOGLE_SECRET"],
			envMap["AUTH_GOOGLE_CALLBACK"],
		),
	)

	// init server
	m = martini.Classic()
	m.Use(render.Renderer())
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token", "X-Requested-With", "Accept", "Accept-Version", "Content-Length", "Content-MD5", "Date", "X-Api-Version", "X-File-Name", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	m.Map(db)

	m.Group("/components", func(r martini.Router) {
		r.Get("", controllers.GetComponents)
		r.Get("/new", controllers.NewComponent)
		r.Get("/:id", controllers.GetComponent)
		r.Put("/:id", binding.Bind(models.Component{}), controllers.UpdateComponent)
		r.Delete("/:id", binding.Bind(models.Component{}), controllers.DeleteComponent)
		r.Post("", binding.Bind(models.Component{}), controllers.CreateComponent)
	})

	m.Group("/synth_modules", func(r martini.Router) {
		r.Get("", controllers.GetSynthModules)
		r.Get("/new", controllers.NewSynthModule)
		r.Get("/:id", controllers.GetSynthModule)
		r.Post("", binding.Bind(models.SynthModule{}), controllers.PostSynthModule)
		r.Delete("/:id", binding.Bind(models.SynthModule{}), controllers.DeleteSynthModule)
		r.Put("/:id", binding.Bind(models.SynthModule{}), controllers.UpdateSynthModule)
	})

	m.Group("/auth", func(r martini.Router) {
		r.Get("/:provider/login", controllers.AuthLogin)
		r.Get("/:provider/callback", controllers.AuthCallback)
		r.Post("/validate_token", binding.Bind(models.User{}), controllers.ValidateToken)
	})

	m.Group("/users", func(r martini.Router) {
		r.Get("", controllers.GetUsers)
		r.Get("/:id", controllers.GetUser)
		r.Delete("/:id", controllers.DeleteUser)
	})

	// start server
	log.Fatal(http.ListenAndServe(":3001", m))
}
