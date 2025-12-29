package mapper

import (
	"kenjix.com/persist/internal/dto"
	"kenjix.com/persist/internal/model"
)

func ProductToDTO(p *model.Product) *dto.ProductDTO {
	if p == nil {
		return nil
	}

	out := &dto.ProductDTO{
		ID:         p.ID,
		Name:       p.Name,
		SKU:        p.SKU,
		Price:      p.Price,
		Active:     p.Active,
		CategoryID: p.CategoryID,
	}

	/*
		// Properties → map
		if len(p.Properties) > 0 {
			out.Properties = make(map[string]string)
			for _, prop := range p.Properties {
				out.Properties[prop.Name] = prop.Value
			}
		}

		// Images → []string (urls)
		for _, img := range p.Images {
			out.Images = append(out.Images, img.URL)
		}

		// Videos → []string (urls)
		for _, vid := range p.Videos {
			out.Videos = append(out.Videos, vid.URL)
		}
	*/
	return out
}
