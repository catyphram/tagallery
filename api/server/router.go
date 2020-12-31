package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tagallery.com/api/controller"
	"tagallery.com/api/logger"
	"tagallery.com/api/model"
	"tagallery.com/api/mongodb"
)

// ConfigureRouter creates and sets the routes on the gin router.
func ConfigureRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/category", func(c *gin.Context) {
		if categories, err := mongodb.QueryCategories(); err != nil {
			logger.Logger().Warnw("Unable to query cagegories.", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			logger.Logger().Infow("Categories queried successfully.", "categories", categories)
			c.JSON(http.StatusOK, categories)
		}
	})

	r.POST("/category", func(c *gin.Context) {
		var category model.Category

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			if err := mongodb.UpsertCategory(category); err != nil {
				logger.Logger().Warnw("Unable to upsert cagegory.", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				logger.Logger().Infow("Category upserted successfully.", "category", category)
				c.JSON(http.StatusOK, category)
			}
		}
	})

	r.GET("/image", func(c *gin.Context) {
		var opts model.ImageOptions
		status := c.Query("status")
		count := c.DefaultQuery("count", "15")
		lastImage := c.Query("lastImage")
		categories := c.QueryArray("categories")

		logger.Logger().Infow("Request parameters.",
			"status", status,
			"count", count,
			"lastImage", lastImage,
			"categories", categories,
		)

		if count != "" {
			if value, err := strconv.Atoi(count); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				opts.Count = &value
			}
		}

		if lastImage != "" {
			opts.LastImage = &lastImage
		}

		if images, err := controller.GetImages(status, opts, categories); err != nil {
			logger.Logger().Warnw("Unable to retrieve images.", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			logger.Logger().Infow("Images retrieved succesfully.", "images", images)
			c.JSON(http.StatusOK, images)
		}
	})

	r.POST("/image", func(c *gin.Context) {
		var image model.Image

		if err := c.ShouldBindJSON(&image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			if image.AssignedCategories == nil {
				image.AssignedCategories = []string{}
			}
			if image.ProposedCategories == nil {
				image.ProposedCategories = []string{}
			}

			if updated, err := controller.UpsertImage(image); err != nil {
				logger.Logger().Warnw("Unable to upsert image.", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				logger.Logger().Infow("Image upserted successfully.", "image", updated)
				c.JSON(http.StatusOK, updated)
			}
		}
	})

	return r
}
