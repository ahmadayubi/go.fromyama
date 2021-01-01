package routes

import (
"github.com/go-chi/chi"

"../controllers/etsy"
"../middleware"
)

func EtsyRoutes() *chi.Mux{
	router:= chi.NewRouter()

	router.Get("/callback", etsy.Callback)
	router.With(middleware.ProtectedApprovedUserRoute).Post("/generate-link", etsy.GenerateAuthURL)

	router.With(middleware.ProtectedApprovedUserRoute).Post("/fulfill/{orderID}",etsy.FulfillOrder)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/orders/all", etsy.GetUnfulfilledOrders)

	return router
}
