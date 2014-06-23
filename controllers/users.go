package controllers

import (
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/lynndylanhurley/defsynth-api/models"
	"github.com/martini-contrib/render"
)

func GetUsers(r render.Render, db gorm.DB) {
	records := []models.User{}
	db.Find(&records)
	r.JSON(200, map[string]interface{}{"data": &records})
}

func GetUser(params martini.Params, r render.Render, db gorm.DB) {
	record := models.User{}
	id := params["id"]

	if err := db.First(&record, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": &record})
	}
}

func DeleteUser(r render.Render, params martini.Params, db gorm.DB) {
	record := models.User{}
	id := params["id"]

	if err := db.First(&record, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	}

	if err := db.Delete(&record).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": "ok"})
	}
}
