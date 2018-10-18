package routes

import (
	"net/http"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
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
	var err error
	count := 20
	status := r.URL.Query().Get("status")
	categories := r.URL.Query()["categories"]
	lastImage := r.URL.Query().Get("last")

	switch status {
	case "unprocessed":
		images, err = models.GetUnprocessedImages(count, lastImage)
	case "uncategorized":
		images = models.GetUnCategorizedImages(count, lastImage)
	case "autocategorized":
		images = models.GetAutoCategorizedImages(count, categories, lastImage)
	default:
		images = models.GetCategorizedImages(count, categories, lastImage)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"status":     status,
			"lastImage":  lastImage,
			"categories": categories,
			"error":      err,
		}).Error("An error during the select of images occured")
		render.Render(w, r, ErrInternalServerError(err))
	} else {
		render.RenderList(w, r, ImageListResponse(images))
	}
}
