package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"tagallery.com/api/models"
)

type ImagePayload struct {
	*models.Image
}

func NewImagePayload(image *models.Image) *ImagePayload {
	resp := &ImagePayload{Image: image}
	return resp
}

func (image ImagePayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ImageListResponse(images []*models.Image) []render.Renderer {
	list := []render.Renderer{}
	for _, image := range images {
		list = append(list, NewImagePayload(image))
	}
	return list
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	var images []*models.Image
	status := r.URL.Query().Get("status")
	categories := r.URL.Query()["categories"]
	lastImage := r.URL.Query().Get("last")

	fmt.Printf("Status: %#v, Categories: %#v, Last Image: %#v\n", status, categories, lastImage)

	switch status {
	case "unprocessed":
		images = getUnprocessedImages(lastImage)
	case "uncategorized":
		images = getUnCategorizedImages(lastImage)
	case "autocategorized":
		images = getAutoCategorizedImages(categories, lastImage)
	default:
		images = getCategorizedImages(categories, lastImage)
	}

	render.RenderList(w, r, ImageListResponse(images))
}

// Get images from a directory that have not been seen/categorized yet
func getUnprocessedImages(lastImage string) []*models.Image {
	return models.ImageFixtures
}

// Check in the DB for images that were unable to be categorized automatically
func getUnCategorizedImages(lastImage string) []*models.Image {
	return models.ImageFixtures
}

// Get images from the DB that have been automatically categorized by the NL
func getAutoCategorizedImages(categories []string, lastImage string) []*models.Image {
	return models.ImageFixtures
}

// Get images from the DB that have been categorized
func getCategorizedImages(categories []string, lastImage string) []*models.Image {
	return models.ImageFixtures
}
