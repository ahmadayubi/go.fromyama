package routes

import (
	"github.com/go-chi/chi"

	"../controllers"
	"../middleware"
)

func UserRoutes() *chi.Mux{
	router:= chi.NewRouter()
	router.Post("/signup", controllers.SignUpUser)
	router.Post("/login", controllers.LoginUser)
	router.With(middleware.ProtectedRoute).Post("/details",controllers.GetUserDetails)
	return router
}


