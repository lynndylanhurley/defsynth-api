package models

import "time"

type ModuleComponent struct {
	Id       int64
	Quantity int64 `form:"quantity" json:"quantity" binding:"required"`

	Component
	SynthModuleId int64 `form:"synth_module_id" json:"-" binding:"required"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
