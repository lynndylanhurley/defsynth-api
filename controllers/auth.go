package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/lynndylanhurley/defsynth-api/models"
	"github.com/martini-contrib/render"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

func AuthLogin(r render.Render, db gorm.DB, params martini.Params, res http.ResponseWriter, req *http.Request) {
	provider, err := gomniauth.Provider(params["provider"])

	if err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	}

	state := gomniauth.NewState("after", "success")

	authUrl, err := provider.GetBeginAuthURL(state, nil)

	if err != nil {
		r.JSON(404, map[string]interface{}{"error": err})
	} else {

		//r.Header().Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X; en-us) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53")

		r.Redirect(authUrl)
	}
}

func AuthCallback(r render.Render, params martini.Params, req *http.Request, db gorm.DB) {
	provider, err := gomniauth.Provider(params["provider"])

	if err != nil {
		panic(err)
	}

	queryParams := objx.Map{}

	for k, v := range req.URL.Query() {
		queryParams[k] = v
	}

	creds, err := provider.CompleteAuth(queryParams)

	if err != nil {
		panic(err)
	}

	user_info, err := provider.GetUser(creds)

	if err != nil {
		fmt.Printf("@-->err %v", err)
		panic(err)
	}

	user := models.User{}

	if db.Where(models.User{Email: user_info.Email()}).FirstOrCreate(&user).Error != nil {
		fmt.Printf("@-->find err %v", err)
		panic(err)
	}

	json_str, err := user_info.Data().JSON()

	if err != nil {
		fmt.Printf("@-->json string err %v", err)
		panic(err)
	}

	fmt.Printf("@-->user %v", user)

	user.Email = user_info.Email()
	user.UserInfo = json_str
	user.AuthToken = signature.RandomKey(64)

	if db.Save(&user).Error != nil {
		fmt.Printf("@-->db save error %v", err)
		panic(err)
	}

	r.HTML(200, "auth_callback", map[string]interface{}{
		"email":      user.Email,
		"auth_token": user.AuthToken,
		"user_info":  user_info.Data(),
	})
}
