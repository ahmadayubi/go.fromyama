package routes

import (
	"github.com/go-chi/chi"

"../controllers/amazon"
"../middleware"
)


func AmazonRoutes() *chi.Mux{
	router:= chi.NewRouter()

	router.With(middleware.ProtectedApprovedUserRoute).Post("/authorize", amazon.Authorize)
	router.With(middleware.ProtectedApprovedUserRoute).Get("/orders/all", amazon.GetUnfulfilledOrders)

	return router
}
