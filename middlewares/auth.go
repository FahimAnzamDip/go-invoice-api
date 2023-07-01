package middlewares

import (
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/app"
	"github.com/fahimanzamdip/go-invoice-api/models"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = app.GetToken(w, r)

		// Check if the user is authenticated and is_admin field is true
		if user, ok := r.Context().Value(app.UserKey).(models.UserInfo); ok && user.UserId > 0 {
			// User is an admin, proceed to the next handler
			next.ServeHTTP(w, r)
		} else {
			// User is not logged in return error
			writeError(w, r, "Access denied. You must be authenticated to access.")
			// Alternatively, you can redirect to a specific page:
			// http.Redirect(w, r, "/unauthorized", http.StatusFound)
		}
	})
}
