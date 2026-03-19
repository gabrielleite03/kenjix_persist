// internal/router/router.go
package router

import (
	"fmt"
	"net/http"

	"github.com/gabrielleite03/kenjix_persist/cmd/api/handler"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/gabrielleite03/kenjix_persist/internal/repository"
	"github.com/gabrielleite03/kenjix_persist/internal/service"
)

type Router struct {
	productHandler *handler.ProductHandler
}

func NewRouter() *Router {
	dbConection := config.NewDatabaseConfig()
	repo := repository.NewProductRepository(dbConection)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	return &Router{
		productHandler: h,
	}
}

func (r *Router) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/products", r.productHandler.List)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Kenjix Persist API")
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	return mux
}
