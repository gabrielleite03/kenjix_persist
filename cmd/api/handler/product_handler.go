package handler

import (
	"net/http"

	"kenjix.com/persist/internal/service"
)

type ProductHandler struct {
	ProductService service.ProductService
}

type ProductService interface {
	ListProducts() error
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, "Kenjix Persist API")
	w.Write([]byte("list products"))
}
