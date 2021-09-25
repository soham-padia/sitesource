package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(Nosurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/phone", handlers.Repo.Phone)
	mux.Get("/pc", handlers.Repo.Pc)
	mux.Get("/laptop", handlers.Repo.Laptop)
	mux.Get("/downloads", handlers.Repo.Download)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/docs", handlers.Repo.Docs)
	mux.Get("/donate", handlers.Repo.Donate)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Get("/user/register", handlers.Repo.Registration)
	mux.Post("/user/register", handlers.Repo.PostRegistration)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/users-new", handlers.Repo.AdminNewUsers)
		mux.Get("/users-all", handlers.Repo.AdminAllUsers)
		mux.Get("/donation", handlers.Repo.AdminDonationInfo)
	})

	return mux
}
