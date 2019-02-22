package controller

import (
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
)

// GetCategories gets all categories from the DB.
func GetCategories() ([]model.Category, error) {
	categories, err := mongodb.GetCategories()
	return categories, err
}
