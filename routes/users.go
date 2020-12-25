package routes

import (
	"github.com/go-chi/chi"

	"../controllers"
	"../middleware"
)

func Routes() *chi.Mux{
	router:= chi.NewRouter()
	router.Post("/signup", controllers.SignUpUser)
	router.Post("/login", controllers.LoginUser)
	router.With(middleware.ProtectedRoute).Post("/details",controllers.GetUserDetails)
//	router.Get("/details")
	return router
}


