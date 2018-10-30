// e. g. curl http://localhost:3333/images?categories=cat2&categories=cat1&status=unprocessed&last=abc123.jpg

package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"tagallery.com/api/config"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/router"
)

func main() {
	err := config.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error while loading the configurations has occured")
	}
	_, _, err = mongodb.Init(mongodb.DatabaseOptions{
		Database: config.GetConfig().Database,
		Host:     config.GetConfig().Database_Host,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error during the startup of the mongodb client occured")
	}
	http.ListenAndServe(fmt.Sprintf(":%v", config.GetConfig().Port), router.CreateRouter())
}
