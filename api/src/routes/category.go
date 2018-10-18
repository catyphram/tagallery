package routes

import (
	"net/http"
  "github.com/go-chi/render"
	"tagallery.com/api/models"
)

type CategoryPayload struct {
  *models.Category
}

func NewCategoryPayload(category *models.Category) *CategoryPayload {
	resp := &CategoryPayload{Category: category}
	return resp
}

func (category CategoryPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CategoryListResponse(categories []*models.Category) []render.Renderer {
	list := []render.Renderer{}
	for _, category := range categories {
		list = append(list, NewCategoryPayload(category))
	}
	return list
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	render.RenderList(w, r, CategoryListResponse(models.CategoryFixtures))
}
