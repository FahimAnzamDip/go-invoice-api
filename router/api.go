package router

import (
	"net/http"
	"os"

	"github.com/fahimanzamdip/go-invoice-api/handlers"
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
		ar.Post("/auth/register", handlers.StoreUserHandler)
		ar.Post("/auth/login", handlers.LoginHandler)
	})

	// Auth routes
	ar.Group(func(ar chi.Router) {
		ar.Use(middlewares.AuthMiddleware)
		// User routes
		ar.Post("/auth/password", handlers.UpdatePasswordHandler)
		ar.Put("/auth/{id}/update", handlers.UpdateUserHandler)
		ar.Get("/auth/validate", handlers.ValidateUserHandler)
		ar.Get("/auth/user", handlers.GetUserByIdHandler)
	})

	// Admin routes
	ar.Group(func(ar chi.Router) {
		ar.Use(middlewares.AdminMiddleware)
		// User routes
		ar.Get("/users", handlers.IndexUserHandler)
		ar.Get("/users/{id}", handlers.ShowUserHandler)
		ar.Put("/users/{id}", handlers.UpdateUserHandler)
		ar.Delete("/users/{id}", handlers.DestroyUserHandler)

		// Categories routes
		ar.Get("/categories", handlers.IndexCategoryHandler)
		ar.Post("/categories", handlers.StoreCategoryHandler)
		ar.Get("/categories/{id}", handlers.ShowCategoryHandler)
		ar.Put("/categories/{id}", handlers.UpdateCategoryHandler)
		ar.Delete("/categories/{id}", handlers.DestroyCategoryHandler)

		// Uploads routes
		ar.Post("/uploads", handlers.UploadHandler)

		// Products routes
		ar.Get("/products", handlers.IndexProductHandler)
		ar.Post("/products", handlers.StoreProductHandler)
		ar.Get("/products/{id}", handlers.ShowProductHandler)
		ar.Put("/products/{id}", handlers.UpdateProductHandler)
		ar.Delete("/products/{id}", handlers.DestroyProductHandler)

		// Tags routes
		ar.Get("/tags", handlers.IndexTagHandler)
		ar.Post("/tags", handlers.StoreTagHandler)
		ar.Get("/tags/{id}", handlers.ShowTagHandler)
		ar.Put("/tags/{id}", handlers.UpdateTagHandler)
		ar.Delete("/tags/{id}", handlers.DestroyTagHandler)

		// Clients routes
		ar.Get("/clients", handlers.IndexClientHandler)
		ar.Post("/clients", handlers.StoreClientHandler)
		ar.Get("/clients/{id}", handlers.ShowClientHandler)
		ar.Put("/clients/{id}", handlers.UpdateClientHandler)
		ar.Delete("/clients/{id}", handlers.DestroyClientHandler)

		// Taxes routes
		ar.Get("/taxes", handlers.IndexTaxHandler)
		ar.Post("/taxes", handlers.StoreTaxHandler)
		ar.Get("/taxes/{id}", handlers.ShowTaxHandler)
		ar.Put("/taxes/{id}", handlers.UpdateTaxHandler)
		ar.Delete("/taxes/{id}", handlers.DestroyTaxHandler)

		// Invoices routes
		ar.Get("/invoices", handlers.IndexInvoiceHandler)
		ar.Post("/invoices", handlers.StoreInvoiceHandler)
		ar.Get("/invoices/{id}", handlers.ShowInvoiceHandler)
		ar.Put("/invoices/{id}", handlers.UpdateInvoiceHandler)
		ar.Delete("/invoices/{id}", handlers.DestroyInvoiceHandler)

		// Payments routes
		ar.Get("/payments", handlers.IndexPaymentHandler)
		ar.Post("/payments", handlers.StorePaymentHandler)
		ar.Get("/payments/{id}", handlers.ShowPaymentHandler)
		ar.Put("/payments/{id}", handlers.UpdatePaymentHandler)
		ar.Delete("/payments/{id}", handlers.DestroyPaymentHandler)
	})

	// Mount to main router
	r.Mount(os.Getenv("api_uri"), ar)

	return r
}