package main

import (
	"net/http"

	"kenjix.com/persist/cmd/api/handler"
	"kenjix.com/persist/cmd/api/router"
	"kenjix.com/persist/internal/config"
	"kenjix.com/persist/internal/repository"
	"kenjix.com/persist/internal/service"
)

func main() {

	dbConection := config.NewDatabaseConfig()
	repo := repository.NewProductRepository(dbConection)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	r := router.NewRouter(h)

	http.ListenAndServe(":8080", r.RegisterRoutes())
}
