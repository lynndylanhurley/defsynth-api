package main

import (
	"fmt"

	"github.com/go-martini/martini"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"

	_ "github.com/lib/pq"
	"github.com/lynndylanhurley/defsynth-api/models"
)

var db gorm.DB

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

	// init server
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token", "X-Requested-With", "Accept", "Accept-Version", "Content-Length", "Content-MD5", "Date", "X-Api-Version", "X-File-Name", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// set up routes
	m.Get("/components", func(r render.Render) {
		components := []models.Component{}
		db.Find(&components)

		r.JSON(200, &components)
	})

	// -- components --
	m.Get("/components/:id", func(params martini.Params, r render.Render) {
		component := models.Component{}
		id := params["id"]

		if err := db.First(&component, id).Error; err != nil {
			panic(fmt.Sprintf("@-->db Find Error: %v", err))
		}

		r.JSON(200, &component)
	})

	m.Post("/components", binding.Bind(models.Component{}), func(r render.Render, component models.Component) {
		if err := db.Save(&component).Error; err != nil {
			panic(fmt.Sprintf("@-->db Save Error: %v", err))
		}
		r.JSON(200, &component)
	})

	m.Put("/components/:id", binding.Bind(models.Component{}), func(r render.Render, params martini.Params, component models.Component) {
		existing := models.Component{}
		id := params["id"]

		if err := db.First(&existing, id).Error; err != nil {
			panic(fmt.Sprintf("@-->db Find Error: %v", err))
		}

		mergo.Merge(&component, existing)

		if err := db.Model(models.Component{}).Updates(&component).Error; err != nil {
			panic(fmt.Sprintf("@-->db Update Error: %v", err))
		}

		r.JSON(200, &component)
	})

	m.Delete("/components/:id", func(r render.Render, params martini.Params) {
		component := models.Component{}
		id := params["id"]

		if err := db.First(&component, id).Error; err != nil {
			panic(fmt.Sprintf("@-->db Find Error: %v", err))
		}

		if err := db.Delete(&component).Error; err != nil {
			panic(fmt.Sprintf("@-->db Delete Error: %v", err))
		}
	})

	// -- synth modules --
	m.Get("/synth_modules", func(r render.Render) {
		records := []models.SynthModule{}
		db.Find(&records)
		r.JSON(200, &records)
	})

	m.Get("/synth_modules/:id", func(params martini.Params, r render.Render) {
		record := models.SynthModule{}
		id := params["id"]

		if err := db.Find(&record, id).Error; err != nil {
			panic(fmt.Sprintf("@-->db Find Error: %+v", err))
		}

		db.Raw(`SELECT
		module_components.id,
		module_components.quantity,
		components.name,
		components.mpn,
		components.newark_s_k_u
		FROM (
			components LEFT JOIN module_components ON module_components.component_id = components.id
		)
		WHERE module_components.synth_module_id = ?`, id).Scan(&record.Components)

		r.JSON(200, &record)
	})

	m.Post("/synth_modules", binding.Bind(models.SynthModule{}), func(r render.Render, record models.SynthModule) {
		if err := db.Save(&record).Error; err != nil {
			panic(fmt.Sprintf("@-->db Save Error: %v", err))
		}

		r.JSON(200, &record)
	})

	m.Put("/synth_modules/:id", binding.Bind(models.SynthModule{}), func(r render.Render, params martini.Params, record models.SynthModule) {
		existing := models.SynthModule{}
		id := params["id"]

		if err := db.First(&existing, id).Error; err != nil {
			panic(fmt.Sprintf("@-->db Find Error: %v", err))
		}

		mergo.Merge(&record, existing)

		if err := db.Save(&record).Error; err != nil {
			panic(fmt.Sprintf("@-->db Update Error: %v", err))
		}

		r.JSON(200, &record)
	})

	m.Delete("/synth_modules/:id", func(r render.Render, params martini.Params) {
		record := models.SynthModule{}
		id := params["id"]

		if err := db.First(&record, id).Error; err != nil {
			panic(fmt.Sprintf("@-->db Find Error: %v", err))
		}

		if err := db.Delete(&record).Error; err != nil {
			panic(fmt.Sprintf("@-->db Delete Error: %v", err))
		}
	})

	m.Post("/module_component", binding.Bind(models.ModuleComponent{}), func(r render.Render, record models.ModuleComponent) {
		if err := db.Save(&record).Error; err != nil {
			panic(fmt.Sprintf("@-->db Save Error: %v", err))
		}
		r.JSON(200, &record)
	})

	// start server
	m.Run()
}
