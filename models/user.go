package models

import "time"

type User struct {
	Id        int64  `form:"id" primaryKey"yes"`
	Email     string `form:"email" binding:"required" sql:"varchar(255)"`
	Name      string `form:"name" sql:"varchar(255)"`
	AuthToken string `form:"auth_token" sql:"varchar(255)"`
	Nickname  string `form:"nickname" sql:"varchar(127)"`
	AvatarURL string `form:"avatar_url" sql:"varchar(255)"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
