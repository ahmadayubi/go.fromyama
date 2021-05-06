package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"go.fromyama/routes"
	"go.fromyama/utils/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)


func CreateRoutes() *chi.Mux{
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.Throttle(20), // limit due to connections the database allows
		)
	router.Route("/", func(r chi.Router){
		r.Mount("/user", routes.UserRoutes())
		r.Mount("/company", routes.CompanyRoutes())
		r.Mount("/postage", routes.PostageRoutes())
		r.Mount("/shopify", routes.ShopifyRoutes())
		r.Mount("/etsy", routes.EtsyRoutes())
		r.Mount("/amazon", routes.AmazonRoutes())
	})
	return router
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error Loading .env File")
		return
	}
	router := CreateRoutes()
	err := database.ConnectToDatabase()
	if err != nil {
		log.Fatal("Error Connecting To Database")
		return
	}
	walkF := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s, %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkF); err != nil {
		log.Fatalf("Logging Error: %s",err.Error())
	}
	log.Fatal(http.ListenAndServe(":8080", router))
}

