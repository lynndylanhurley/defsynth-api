package models

import "time"

type Component struct {
	Id        int64  `form:"id" primaryKey:"yes"`
	Name      string `form:"name" binding:"required" sql:"type:varchar(127)"`
	Mpn       string `form:"mpn" binding:"required" sql:"type:varchar(127)"`
	NewarkSKU string `form:"newarkSKU" binding:"required" sql:"type:varchar(127)"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
