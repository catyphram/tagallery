package controller

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"tagallery.com/api/config"
	"tagallery.com/api/logger"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
	"tagallery.com/api/util"
)

// GetUnprocessedImages returns unprocessed images from a file directory.
// Subdirectories are ignored.
// Options can be passed to limit the number of images returned and
// to start from a specific image, useful for pagination.
func GetUnprocessedImages(opts model.ImageOptions) ([]model.Image, error) {
	imageDir := filepath.Join(config.Get().Images, config.Get().UnprocessedImagesFolder)
	images := []model.Image{}
	selectImages := true
	counter := 0

	if opts.LastImage != nil {
		selectImages = false
	}

	// Create the unprocessed image directory if it does not exist
	_ = os.MkdirAll(imageDir, 0755)

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Logger().Warnw("A file processed during the dir walk of unprocessed images caused an error.", "error", err)
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
			if filepath.Base(path) == *opts.LastImage {
				selectImages = true
			}
		} else if opts.Count == nil || counter < *opts.Count {
			counter++
			images = append(images, model.Image{
				File:               filepath.Join(config.Get().UnprocessedImagesFolder, filepath.Base(path)),
				AssignedCategories: []string{},
				ProposedCategories: []string{},
			})
		} else {
			return io.EOF
		}
		return nil
	})

	if err != nil && err != io.EOF {
		return nil, err
	}

	return images, nil
}

// GetImages returns a list of images filtered by
// count, categories, status and lastImage for pagination.
func GetImages(
	status string, opts model.ImageOptions, categories []string,
) ([]model.Image, error) {

	switch status {
	case "unprocessed":
		return GetUnprocessedImages(opts)
	case "uncategorized":
		return mongodb.GetImages(opts, nil)
	case "autocategorized":
		return mongodb.GetImages(opts, &model.CategoryMap{
			Proposed: categories,
		})
	default:
		return mongodb.GetImages(opts, &model.CategoryMap{
			Assigned: categories,
		})
	}
}

// UpsertImage inserts or updates an existing image.
func UpsertImage(image model.Image) (*model.Image, error) {
	if strings.HasSuffix(filepath.Dir(image.File), config.Get().UnprocessedImagesFolder) {
		imageDir := filepath.Join(config.Get().Images, config.Get().ProcessedImagesFolder)

		// Create the processed image directory if it does not exist
		_ = os.MkdirAll(imageDir, 0755)

		fileName := filepath.Base(image.File)
		newPath := filepath.Join(imageDir, fileName)

		if _, err := os.Stat(newPath); err == nil {
			return nil, errors.New("file already exists")
		} else if os.IsNotExist(err) {
			if err := os.Rename(filepath.Join(config.Get().Images, image.File), newPath); err != nil {
				return nil, err
			}
			image.File = filepath.Join(config.Get().ProcessedImagesFolder, fileName)
		}
	}

	// Automatically add the starred category to the assigned categories.
	if exists := util.ContainsString(image.AssignedCategories, *image.StarredCategory, false); !exists {
		image.AssignedCategories = append(image.AssignedCategories, *image.StarredCategory)
	}

	return &image, mongodb.UpsertImage(image)
}
