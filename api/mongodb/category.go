package mongodb

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	log "github.com/sirupsen/logrus"
	"tagallery.com/api/model"
)

// DBCategory encapsulates a model.Category and it's id in the db
type DBCategory struct {
	model.Category
	Id objectid.ObjectID `json:"id" bson:"_id"`
}

// GetCategories queries the database for categories.
func GetCategories() ([]model.Category, error) {
	db, ctx, err := GetConnection()

	if err != nil {
		return nil, err
	}

	c, err := db.Collection("categories").Find(ctx, nil)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get cursor for categories in database.")
		return nil, err
	}

	defer c.Close(ctx)

	var categories []model.Category

	for c.Next(ctx) {
		category := model.Category{}
		err := c.Decode(&category)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to parse category object from bson during select of categories.")
			return nil, err
		}
		categories = append(categories, category)
	}
	if err = c.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Cursor failed during listing of categories in database.")
		return nil, err
	}

	return categories, nil
}
