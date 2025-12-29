// internal/router/router.go
package router

import (
	"fmt"
	"net/http"

	"kenjix.com/persist/cmd/api/handler"
)

type Router struct {
	productHandler *handler.ProductHandler
}

func NewRouter(productHandler *handler.ProductHandler) *Router {
	return &Router{
		productHandler: productHandler,
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
