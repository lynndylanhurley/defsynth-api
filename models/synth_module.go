package models

import "time"

type SynthModule struct {
	Id          int64  `form:"id" primaryKey:"yes"`
	Name        string `form:"name" binding:"required" sql:"type:varchar(127)"`
	Url         string `form:"url" binding:"required" sql:"type:varchar(127)"`
	Description string `form:"description" sql:"type:text"`

	Components []ModuleComponent `form:"components" json:"components"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
