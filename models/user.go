package models

import (
	"encoding/json"
	"time"
)

type User struct {
	Id        int64           `form:"id" primaryKey"yes"`
	Email     string          `form:"email" binding:"required" sql:"varchar(255)"`
	AuthToken string          `form:"auth_token" sql:"varchar(255)"`
	UserInfo  json.RawMessage `form:"user_info"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
}
