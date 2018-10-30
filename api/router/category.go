package router

import (
	"net/http"

	"github.com/go-chi/render"
	"tagallery.com/api/model"
)

type CategoryPayload struct {
	model.Category
}

func NewCategoryPayload(category model.Category) CategoryPayload {
	return CategoryPayload{Category: category}
}

func (category CategoryPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CategoryListResponse(categories []model.Category) []render.Renderer {
	list := []render.Renderer{}
	for _, category := range categories {
		list = append(list, NewCategoryPayload(category))
	}
	return list
}

// GetCategories returns the list of categories images can be assigned to.
func GetCategories(w http.ResponseWriter, r *http.Request) {
	render.RenderList(w, r, CategoryListResponse(model.Categories))
}
