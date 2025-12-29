package dto

// ProductDTO represents the data transfer object for a product
type ProductDTO struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Price      float64 `json:"price"`
	Active     bool    `json:"active"`
	CategoryID *int64  `json:"category,omitempty"`

	Properties map[string]string `json:"properties,omitempty"`
	Images     []string          `json:"images,omitempty"`
	Videos     []string          `json:"videos,omitempty"`
}
