package models

import "time"

type User struct {
	Id        int64  `form:"id" json:"id" primaryKey:"yes"`
	Email     string `form:"email" json:"email" binding:"required" sql:"varchar(255)"`
	Name      string `form:"name" json:"name" sql:"varchar(255)"`
	AuthToken string `form:"auth_token" json:"auth_token" sql:"varchar(255)"`
	Nickname  string `form:"nickname" json:"nickname" sql:"varchar(127)"`
	AvatarURL string `form:"avatar_url" json:"avatar_url" sql:"varchar(255)"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
