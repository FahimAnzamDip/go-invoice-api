package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fahimanzamdip/go-invoice-api/app"
	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := &models.AuthUser{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := models.Login(user.EmailOrMobile, user.Password)
	u.Respond(w, res)
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := user.UpdatePassword()
	u.Respond(w, res)
}

func ValidateUserHandler(w http.ResponseWriter, r *http.Request) {
	userInfoValue := r.Context().Value(app.UserKey)
	var userInfo models.UserInfo
	if userInfoValue == nil {
		return
	} else {
		userInfo = userInfoValue.(models.UserInfo)
	}
	user := &models.User{}

	tokenHeader := r.Header.Get("Authorization") // Grab the token from the header
	splitted := strings.Split(tokenHeader, " ")  // The token normally comes in format `Bearer {token-body}`,
	// we check if the retrieved token matched this requirement
	user.Token = splitted[1]

	res := user.ValidateUser(userInfo.UserId)
	u.Respond(w, res)
}
