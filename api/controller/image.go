package controller

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"tagallery.com/api/config"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
)

// GetUnprocessedImages gets {count} images from a directory that have not been seen/categorized yet.
// Subdirectories are ignored and if count <= 0 all images are returned.
func GetUnprocessedImages(count int, lastImage string) ([]model.Image, error) {

	var images []model.Image
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
			}).Warn("A file processed during the dir walk of unprocessed images caused an error.")
			return nil
		}

		// Skip root dir and nested dirs
		if path == imageDir {
			return nil
		} else if info.IsDir() {
			return filepath.SkipDir
		}

		// If lastImage is specified wait till we find it and then get {count} files
		if !selectImages {
			if filepath.Base(path) == lastImage {
				selectImages = true
			}
		} else if count <= 0 || counter < count {
			counter++
			images = append(images, model.Image{File: filepath.Base(path)})
		} else {
			return io.EOF
		}
		return nil
	})

	if err != nil && err != io.EOF {
		log.WithFields(log.Fields{
			"error":  err,
			"images": images,
		}).Warn("An error during the file walk occured when listing unprocessed files.")
		return nil, err
	}

	return images, nil
}

// GetUnCategorizedImages checks in the DB for images that were unable to get categorized automatically.
func GetUnCategorizedImages(count int, lastImage string) ([]model.Image, error) {
	return mongodb.GetImages(count, nil, lastImage)
}

// GetAutoCategorizedImages gets images from the DB that have been automatically categorized.
func GetAutoCategorizedImages(count int, categories []string, lastImage string) ([]model.Image, error) {
	return mongodb.GetImages(count, &model.CategoryMap{
		Proposed: categories,
	}, lastImage)
}

// GetCategorizedImages gets images from the DB that have been categorized.
func GetCategorizedImages(count int, categories []string, lastImage string) ([]model.Image, error) {
	images, err := mongodb.GetImages(count, &model.CategoryMap{
		Assigned: categories,
	}, lastImage)
	return images, err
}

// UpsertImage inserts or updates an image in the db.
func UpsertImage(image model.Image) error {
	return mongodb.UpsertImage(image)
}
