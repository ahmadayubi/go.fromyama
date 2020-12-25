package routes

import (
	"github.com/go-chi/chi"

	"../controllers"
	"../middleware"
)

func CompanyRoutes() *chi.Mux{
	router:= chi.NewRouter()
	router.Post("/register", controllers.RegisterCompany)
	router.With(middleware.ProtectedRoute).Get("/employee/all", controllers.GetAllEmployeeDetails)

	router.With(middleware.ProtectedRoute).Post("/employee/approve/{employeeID}", controllers.LoginUser)
	router.With(middleware.ProtectedRoute).Get("/employee/approved/{employeeID}",controllers.GetUserDetails)
	router.With(middleware.ProtectedRoute).Get("/details", controllers.GetUserDetails)

	router.With(middleware.ProtectedRoute).Post("/payment/add", controllers.GetUserDetails)
	router.With(middleware.ProtectedRoute).Post("/payment/charge", controllers.GetUserDetails)

	router.With(middleware.ProtectedRoute).Post("/parcel/add", controllers.GetUserDetails)
	router.With(middleware.ProtectedRoute).Get("/shipper", controllers.GetUserDetails)

	router.With(middleware.ProtectedRoute).Delete("/unregister",controllers.GetUserDetails)
	return router
}


