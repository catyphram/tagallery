package router

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"tagallery.com/api/controller"
	"tagallery.com/api/model"
)

type ImagePayload struct {
	model.Image
}

func NewImagePayload(image model.Image) ImagePayload {
	return ImagePayload{Image: image}
}

func (image ImagePayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (i *ImagePayload) Bind(r *http.Request) error {
	return nil
}

func ImageListResponse(images []model.Image) []render.Renderer {
	list := []render.Renderer{}
	for _, image := range images {
		list = append(list, NewImagePayload(image))
	}
	return list
}

// GetImages returns a list of images filtered by
// count, categories, status and lastImage for pagination.
func GetImages(w http.ResponseWriter, r *http.Request) {
	var images []model.Image
	var err error
	count := 15
	status := r.URL.Query().Get("status")
	categories := r.URL.Query()["categories"]
	lastImage := r.URL.Query().Get("last")

	if i, err := strconv.Atoi(r.URL.Query().Get("count")); err == nil && i > 0 {
		count = i
	}

	switch status {
	case "unprocessed":
		images, err = controller.GetUnprocessedImages(count, lastImage)
	case "uncategorized":
		images, err = controller.GetUnCategorizedImages(count, lastImage)
	case "autocategorized":
		images, err = controller.GetAutoCategorizedImages(count, categories, lastImage)
	default:
		images, err = controller.GetCategorizedImages(count, categories, lastImage)
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

// UpsertImage creates or updates an image in the database.
func UpsertImage(w http.ResponseWriter, r *http.Request) {
	data := &ImagePayload{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	image := data.Image
	err := controller.UpsertImage(image)

	if err != nil {
		log.WithFields(log.Fields{
			"data":  data,
			"error": err,
		}).Error("An error during the upsert of an image occured")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, 200)
}
