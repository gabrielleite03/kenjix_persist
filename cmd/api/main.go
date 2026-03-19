package main

import (
	"net/http"

	"github.com/gabrielleite03/kenjix_persist/cmd/api/router"
)

func main() {

	r := router.NewRouter()

	http.ListenAndServe(":8080", r.RegisterRoutes())
}
