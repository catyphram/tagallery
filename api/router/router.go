// Package router contains logic for route and error handling
// as well as high level functions for operations.
package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", GetCategories)
	})
	r.Route("/images", func(r chi.Router) {
		r.Get("/", GetImages)
	})
	r.Route("/image", func(r chi.Router) {
		r.Post("/{imageID}", UpsertImage)
	})

	return r
}

// ErrResponse is a general struct to represent an error in the API.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // category-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInternalServerError create a new ErrResponse with HTTP status code 500.
func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error.",
		ErrorText:      err.Error(),
	}
}

// ErrInvalidRequest create a new ErrResponse with HTTP status code 400.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
