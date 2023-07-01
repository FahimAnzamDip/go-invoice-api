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

		// Check if the user is authenticated and is_admin field is true
		if user, ok := r.Context().Value(app.UserKey).(models.UserInfo); ok && user.IsAdmin {
			// User is an admin, proceed to the next handler
			next.ServeHTTP(w, r)
		} else {
			// User is not an admin, return an error response or redirect as needed
			writeError(w, r, "Access denied. You must be admin to access.")
			// Alternatively, you can redirect to a specific page:
			// http.Redirect(w, r, "/unauthorized", http.StatusFound)
		}
	})
}

func writeError(w http.ResponseWriter, r *http.Request, message string) {
	response := u.Message(false, message)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, response)
}
