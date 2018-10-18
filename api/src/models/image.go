package models

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"tagallery.com/api/config"
)

type Image struct {
	Name               string      `json:"name"`
	AssignedCategories []*Category `json:"assignedCategories"`
	ProposedCategories []*Category `json:"proposedCategories"`
	StarredCategory    *Category   `json:"starredCategory"`
}

var ImageFixtures = []*Image{
	{Name: "Image 1", AssignedCategories: []*Category{CategoryFixtures[0]}, ProposedCategories: []*Category{}, StarredCategory: nil},
	{Name: "Image 2", AssignedCategories: []*Category{}, ProposedCategories: []*Category{CategoryFixtures[0]}, StarredCategory: nil},
	{Name: "Image 3", AssignedCategories: []*Category{}, ProposedCategories: []*Category{CategoryFixtures[1]}, StarredCategory: nil},
	{Name: "Image 4", AssignedCategories: []*Category{CategoryFixtures[0], CategoryFixtures[1]}, ProposedCategories: []*Category{}, StarredCategory: nil},
	{Name: "Image 5", AssignedCategories: []*Category{}, ProposedCategories: []*Category{}, StarredCategory: CategoryFixtures[1]},
}

// Get images from a directory that have not been seen/categorized yet
func GetUnprocessedImages(count int, lastImage string) ([]*Image, error) {

	var images []*Image
	var selectImages bool
	counter := 0
	imageDir := config.GetConfig().Unprocessed_Images

	if lastImage != "" {
		selectImages = false
	} else {
		selectImages = true
	}

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"path":  path,
			}).Warn("A file processed during the dir walk of unprocessed images caused an error")
		} else if !selectImages && path == lastImage {
			selectImages = true
		} else if selectImages && path != imageDir && counter < count {
			counter++
			// @TODO: The path has to be switched to a URL if needed to be accessed from a different PC
			// To do so add a static web server that serves from the imageDir and
			// replace the path's directory with the server URL
			images = append(images, &Image{Name: path})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}

// Check in the DB for images that were unable to be categorized automatically
func GetUnCategorizedImages(count int, lastImage string) []*Image {
	return ImageFixtures
}

// Get images from the DB that have been automatically categorized by the NL
func GetAutoCategorizedImages(count int, categories []string, lastImage string) []*Image {
	return ImageFixtures
}

// Get images from the DB that have been categorized
func GetCategorizedImages(count int, categories []string, lastImage string) []*Image {
	return ImageFixtures
}
