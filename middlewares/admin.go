package middlewares

import (
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/app"
	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = app.GetToken(w, r)

		if user, ok := r.Context().Value(app.UserKey).(models.UserInfo); ok && user.IsAdmin {
			next.ServeHTTP(w, r)
		} else {
			writeError(w, r, "Access denied, You must be an admin to access!")
		}
	})
}

func writeError(w http.ResponseWriter, r *http.Request, message string) {
	response := u.Message(false, message)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, response)
}
