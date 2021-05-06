package routes

import (
	"github.com/go-chi/chi/v5"

	"go.fromyama/controllers"
	"go.fromyama/middleware"
)

func UserRoutes() *chi.Mux{
	router:= chi.NewRouter()
	router.Post("/signup", controllers.SignUpUser)
	router.Post("/login", controllers.LoginUser)
	router.With(middleware.ProtectedRoute).Get("/details",controllers.GetUserDetails)
	router.With(middleware.ProtectedRoute).Get("/refresh", controllers.RefreshToken)
	return router
}


