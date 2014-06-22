package controllers

import (
	"github.com/go-martini/martini"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"github.com/lynndylanhurley/defsynth-api/models"
	"github.com/martini-contrib/render"
)

func GetSynthModules(r render.Render, db gorm.DB) {
	records := []models.SynthModule{}
	db.Find(&records)
	r.JSON(200, map[string]interface{}{"data": &records})
}

func NewSynthModule(r render.Render) {
	r.JSON(200, map[string]interface{}{"data": models.SynthModule{}})
}

func GetSynthModule(params martini.Params, r render.Render, db gorm.DB) {
	record := models.SynthModule{}
	id := params["id"]

	if err := db.Find(&record, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	} else {
		db.Raw(`SELECT
			module_components.id,
			module_components.quantity,
			components.name,
			components.mpn,
			components.newark_s_k_u
			FROM (
				components LEFT JOIN module_components ON module_components.component_id = components.id
			)
			WHERE module_components.synth_module_id = ?
		`, id).Scan(&record.Components)

		r.JSON(200, map[string]interface{}{"data": &record})
	}
}

func PostSynthModule(r render.Render, record models.SynthModule, db gorm.DB) {
	if err := db.Save(&record).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": &record})
	}
}

func DeleteSynthModule(r render.Render, params martini.Params, db gorm.DB) {
	synth_module := models.SynthModule{}

	id := params["id"]

	if err := db.Find(&synth_module, id).Error; err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	}

	if err := db.Where("synth_module_id = ?", id).Delete(models.ModuleComponent{}).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	}

	if err := db.Delete(&synth_module).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": "ok"})
	}
}

func UpdateSynthModule(r render.Render, params martini.Params, db gorm.DB, record models.SynthModule) {
	existing := models.SynthModule{}
	id := params["id"]

	if err := db.First(&existing, id).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	}

	mergo.Merge(&record, existing)

	if err := db.Save(&record).Error; err != nil {
		r.JSON(500, map[string]interface{}{"error": err})
	} else {
		r.JSON(200, map[string]interface{}{"data": &record})
	}
}
