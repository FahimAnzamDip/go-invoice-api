package routes

import (
	"net/http"
	"os"

	"github.com/fahimanzamdip/go-invoice-api/controllers"
	"github.com/fahimanzamdip/go-invoice-api/middlewares"
	"github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Configure() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.Respond(w, utils.Message(true, "Welcome to Go-Shop-Api"))
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.Respond(w, utils.Message(false, "404 the requested route is not available"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.Respond(w, utils.Message(false, "405 the provided method is not allowed for the requested route"))
	})

	// Sub router
	ar := chi.NewRouter()

	// Public routes
	ar.Group(func(ar chi.Router) {
		ar.Post("/auth/register", controllers.StoreUserHandler)
		ar.Post("/auth/login", controllers.LoginHandler)
	})

	// Auth routes
	ar.Group(func(ar chi.Router) {
		ar.Use(middlewares.AuthMiddleware)
		// User routes
		ar.Post("/auth/password", controllers.UpdatePasswordHandler)
		ar.Put("/auth/{id}/update", controllers.UpdateUserHandler)
		ar.Get("/auth/validate", controllers.ValidateUserHandler)
	})

	// Admin routes
	ar.Group(func(ar chi.Router) {
		ar.Use(middlewares.AdminMiddleware)
		// User routes
		ar.Get("/users", controllers.IndexUserHandler)
		ar.Get("/users/{id}", controllers.ShowUserHandler)
		ar.Delete("/users/{id}", controllers.DestroyUserHandler)

		// Categories routes
		ar.Get("/categories", controllers.IndexCategoryHandler)
		ar.Post("/categories", controllers.StoreCategoryHandler)
		ar.Get("/categories/{id}", controllers.ShowCategoryHandler)
		ar.Put("/categories/{id}", controllers.UpdateCategoryHandler)
		ar.Delete("/categories/{id}", controllers.DestroyCategoryHandler)

		// Uploads routes
		ar.Post("/uploads", controllers.UploadHandler)

		// Products routes
		ar.Get("/products", controllers.IndexProductHandler)
		ar.Post("/products", controllers.StoreProductHandler)
		ar.Get("/products/{id}", controllers.ShowProductHandler)
		ar.Put("/products/{id}", controllers.UpdateProductHandler)
		ar.Delete("/products/{id}", controllers.DestroyProductHandler)

		// Tags routes
		ar.Get("/tags", controllers.IndexTagHandler)
		ar.Post("/tags", controllers.StoreTagHandler)
		ar.Get("/tags/{id}", controllers.ShowTagHandler)
		ar.Put("/tags/{id}", controllers.UpdateTagHandler)
		ar.Delete("/tags/{id}", controllers.DestroyTagHandler)

		// Clients routes
		ar.Get("/clients", controllers.IndexClientHandler)
		ar.Post("/clients", controllers.StoreClientHandler)
		ar.Get("/clients/{id}", controllers.ShowClientHandler)
		ar.Put("/clients/{id}", controllers.UpdateClientHandler)
		ar.Delete("/clients/{id}", controllers.DestroyClientHandler)

		// Clients routes
		ar.Get("/invoices", controllers.IndexInvoiceHandler)
		ar.Post("/invoices", controllers.StoreInvoiceHandler)
		ar.Get("/invoices/{id}", controllers.ShowInvoiceHandler)
		ar.Delete("/invoices/{id}", controllers.DestroyInvoiceHandler)
	})

	// Mount to main router
	r.Mount(os.Getenv("api_uri"), ar)

	return r
}
