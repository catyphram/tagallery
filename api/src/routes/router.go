package routes

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

  return r;
}

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

// func ErrInvalidRequest(err error) render.Renderer {
// 	return &ErrResponse{
// 		Err:            err,
// 		HTTPStatusCode: 400,
// 		StatusText:     "Invalid request.",
// 		ErrorText:      err.Error(),
// 	}
// }

// func ErrRender(err error) render.Renderer {
// 	return &ErrResponse{
// 		Err:            err,
// 		HTTPStatusCode: 422,
// 		StatusText:     "Error rendering response.",
// 		ErrorText:      err.Error(),
// 	}
// }

// var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
