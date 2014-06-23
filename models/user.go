package models

import "time"

type User struct {
	Id        int64     `form:"id" primaryKey"yes"`
	Email     string    `form:"email" binding:"required" sql:"varchar(255)"`
	AuthToken string    `form:"auth_token" sql:"varchar(255)"`
	UserInfo  string    `form:"user_info" sql:"type:json DEFAULT '{}'"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
