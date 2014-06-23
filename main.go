package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
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

	// init oauth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		github.New(
			"bffafe2fd4db7beb5d42",
			"60d9528dbc1cf06374968a92be749db25c50899b",
			"http://defsynth-api.dev/auth/github/callback",
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
	})

	m.Group("/users", func(r martini.Router) {
		r.Get("", controllers.GetUsers)
		r.Get("/:id", controllers.GetUser)
		r.Delete("/:id", controllers.DeleteUser)
	})

	// start server
	log.Fatal(http.ListenAndServe(":3001", m))
}
