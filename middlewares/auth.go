package middlewares

import (
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/app"
	"github.com/fahimanzamdip/go-invoice-api/models"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = app.GetToken(w, r)

		if user, ok := r.Context().Value(app.UserKey).(models.UserInfo); ok && user.UserId > 0 {
			next.ServeHTTP(w, r)
		} else {
			writeError(w, r, "Please login to access!")
		}
	})
}
