package routes

import (
	"github.com/go-chi/chi"

	"../controllers"
	"../middleware"
)

func CompanyRoutes() *chi.Mux{
	router:= chi.NewRouter()
	router.Post("/register", controllers.RegisterCompany)

	router.With(middleware.ProtectedApprovedUserRoute).Get("/platforms", controllers.GetConnectedPlatforms)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/employee/all", controllers.GetAllEmployeeDetails)

	router.With(middleware.ProtectedApprovedUserRoute).Put("/employee/approve/{employeeID}", controllers.ApproveEmployee)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/employee/approved/{employeeID}",controllers.IsEmployeeApproved)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/details", controllers.GetCompanyDetails)

	router.With(middleware.ProtectedApprovedUserRoute).Post("/add/payment/method", controllers.AddPaymentMethod)
	router.With(middleware.ProtectedApprovedUserRoute).Post("/add/payment/charge", controllers.ChargePaymentAccount)

	router.With(middleware.ProtectedApprovedUserRoute).Post("/add/parcel", controllers.AddParcel)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/shipper", controllers.GetShipper)

	router.With(middleware.ProtectedApprovedUserRoute).Delete("/unregister",controllers.UnregisterCompany)
	return router
}