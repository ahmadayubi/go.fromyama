package routes

import (
	"github.com/go-chi/chi"

	"../controllers"
	"../middleware"
)

func PostageRoutes() *chi.Mux{
	router:= chi.NewRouter()

	router.With(middleware.ProtectedApprovedUserRoute).Post("/buy/canadapost", controllers.BuyCanadaPostPostageLabel)

	return router
}