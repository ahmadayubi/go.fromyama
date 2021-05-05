package routes

import (
	"github.com/go-chi/chi"

	"go.fromyama/controllers"
	"go.fromyama/middleware"
)

func PostageRoutes() *chi.Mux{
	router:= chi.NewRouter()

	router.With(middleware.ProtectedApprovedUserRoute).Post("/buy/canadapost", controllers.BuyCanadaPostPostageLabel)
	router.With(middleware.ProtectedApprovedUserRoute).Post("/rates/canadapost", controllers.GetCanadaPostRate)

	return router
}