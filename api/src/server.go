// e. g. curl http://localhost:3333/images?categories=cat2&categories=cat1&status=unprocessed&last=abc123.jpg

package main

import (
	"net/http"

	"tagallery.com/api/routes"
)

func main() {
	http.ListenAndServe(":3333", routes.CreateRouter())
}
