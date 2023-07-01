package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "user"

func GetToken(w http.ResponseWriter, r *http.Request) *http.Request {
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

	if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
		writeError(w, r, "Missing auth token")
		return r
	}

	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		writeError(w, r, "Invalid/Malformed auth token")
		return r
	}

	tokenPart := splitted[1] //Grab the token part, what we are truly interested in
	tkn := &models.Token{}

	token, err := jwt.ParseWithClaims(tokenPart, tkn, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	if err != nil { //Malformed token, returns with http code 403 as usual
		writeError(w, r, "Malformed authentication token")
		return r
	}

	if !token.Valid { //Token is invalid, maybe not signed on this server
		writeError(w, r, "Invalid token")
		return r
	}

	//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
	ctx := context.WithValue(r.Context(), UserKey, tkn.UserInfo)
	r = r.WithContext(ctx)

	return r
}

func writeError(w http.ResponseWriter, r *http.Request, message string) {
	response := u.Message(false, message)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, response)
}
