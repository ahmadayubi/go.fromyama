package routes

import (
	"github.com/go-chi/chi"

	"../controllers"
	"../middleware"
)

func PostageRoutes() *chi.Mux{
	router:= chi.NewRouter()

	router.With(middleware.ProtectedApprovedUserRoute).Post("/buy/canadapost", controllers.BuyCanadaPostPostageLabel)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/employee/approved/{employeeID}",controllers.IsEmployeeApproved)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/details", controllers.GetCompanyDetails)

	return router
}