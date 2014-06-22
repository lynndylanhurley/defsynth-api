package controllers

import (
	"github.com/go-martini/martini"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"github.com/lynndylanhurley/defsynth-api/models"
	"github.com/martini-contrib/render"
)

func GetComponents(r render.Render, db gorm.DB) {
	components := []models.Component{}
	db.Find(&components)

	r.JSON(200, map[string]interface{}{"data": &components})
}

func GetComponent(params martini.Params, r render.Render, db gorm.DB) {
	component := models.Component{}
	id := params["id"]

	if err := db.First(&component, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": &component})
	}
}

func NewComponent(r render.Render) {
	component := models.Component{}
	r.JSON(200, map[string]interface{}{"data": &component})
}

func CreateComponent(r render.Render, component models.Component, db gorm.DB) {
	if err := db.Save(&component).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": &component})
	}
}

func UpdateComponent(r render.Render, params martini.Params, component models.Component, db gorm.DB) {
	existing := models.Component{}
	id := params["id"]

	if err := db.First(&existing, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	}

	mergo.Merge(&component, existing)

	if err := db.Save(&component).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": &component})
	}
}

func DeleteComponent(r render.Render, params martini.Params, db gorm.DB) {
	component := models.Component{}
	id := params["id"]

	if err := db.First(&component, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	}

	if err := db.Delete(&component).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": "ok"})
	}
}
