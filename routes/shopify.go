package routes

import (
	"github.com/go-chi/chi"

"../controllers/shopify"
"../middleware"
)

func ShopifyRoutes() *chi.Mux{
	router:= chi.NewRouter()

	router.Get("/callback", shopify.Callback)
	router.With(middleware.ProtectedApprovedUserRoute).Post("/generate-link",shopify.GenerateAuthURL)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/locations", shopify.GetLocations)

	router.With(middleware.ProtectedApprovedUserRoute).Get("/orders/all", shopify.GetUnfulfilledOrders)

	return router
}