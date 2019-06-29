package controller

import (
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
)

// GetCategories returns all categories from the db.
func GetCategories() ([]model.Category, error) {
	return mongodb.GetCategories()
}
