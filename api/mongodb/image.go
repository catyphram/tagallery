package mongodb

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/replaceopt"
	log "github.com/sirupsen/logrus"
	"tagallery.com/api/model"
)

type DBImage struct {
	model.Image
	Id objectid.ObjectID `json:"id" bson:"_id"`
}

// GetImages queries the database for images.
// If count > 0 no more then {count} images will be returned.
// A CategoryMap may be passed to filter only images that are in all of these categories.
// If categories == nil then instead of (auto)-categorized images,
// only images that have no assigned category will be returned.
// With lastImage you get only images after this only. Used for pagination.
func GetImages(count int, categories *model.CategoryMap, lastImage string) ([]model.Image, error) {
	var dbLastImage DBImage
	doc := bson.NewDocument()
	db, ctx, err := GetConnection()

	if err != nil {
		return nil, err
	}

	if lastImage != "" {
		lastImageDoc := bson.NewDocument(bson.EC.String("file", lastImage))
		err = db.Collection("images").FindOne(ctx, lastImageDoc).Decode(&dbLastImage)

		if err != nil {
			log.WithFields(log.Fields{
				"error":             err,
				"lastImageArgument": lastImage,
				"lastImageDatabase": dbLastImage,
			}).Info("Unable to find lastImage in the database.")
		} else {
			doc.Append(bson.EC.SubDocumentFromElements(
				"_id", bson.EC.ObjectID("$gt", dbLastImage.Id),
			))
		}
	}

	// Uncategorized images only
	if categories == nil {

		doc.Append(
			bson.EC.ArrayFromElements(
				"$and", bson.VC.DocumentFromElements(
					bson.EC.ArrayFromElements(
						"$or",
						bson.VC.DocumentFromElements(
							bson.EC.SubDocumentFromElements(
								"assignedCategories", bson.EC.Boolean("$exists", false),
							),
						),
						bson.VC.DocumentFromElements(
							bson.EC.SubDocumentFromElements(
								"assignedCategories", bson.EC.Int32("$size", 0),
							),
						),
						bson.VC.DocumentFromElements(
							bson.EC.Null("assignedCategories"),
						),
					),
				), bson.VC.DocumentFromElements(
					bson.EC.ArrayFromElements(
						"$or",
						bson.VC.DocumentFromElements(
							bson.EC.SubDocumentFromElements(
								"proposedCategories", bson.EC.Boolean("$exists", false),
							),
						),
						bson.VC.DocumentFromElements(
							bson.EC.SubDocumentFromElements(
								"proposedCategories", bson.EC.Int32("$size", 0),
							),
						),
						bson.VC.DocumentFromElements(
							bson.EC.Null("proposedCategories"),
						),
					),
				), bson.VC.DocumentFromElements(
					bson.EC.ArrayFromElements(
						"$or",
						bson.VC.DocumentFromElements(
							bson.EC.SubDocumentFromElements(
								"starredCategory", bson.EC.Boolean("$exists", false),
							),
						),
						bson.VC.DocumentFromElements(
							bson.EC.Null("starredCategory"),
						),
						bson.VC.DocumentFromElements(
							bson.EC.String("starredCategory", ""),
						),
					),
				),
			),
		)
	} else {
		if len(categories.Assigned) > 0 {

			var assignedCategories []*bson.Value
			for _, assignedCategory := range categories.Assigned {
				assignedCategories = append(assignedCategories, bson.VC.String(assignedCategory))
			}

			doc.Append(bson.EC.SubDocumentFromElements(
				"assignedCategories", bson.EC.ArrayFromElements(
					"$all", assignedCategories...,
				),
			))
		}

		if len(categories.Proposed) > 0 {

			var proposedCategories []*bson.Value
			for _, proposedCategory := range categories.Proposed {
				proposedCategories = append(proposedCategories, bson.VC.String(proposedCategory))
			}

			doc.Append(bson.EC.SubDocumentFromElements(
				"proposedCategories", bson.EC.ArrayFromElements(
					"$all", proposedCategories...,
				),
			))
		}

		if categories.Starred != "" {
			doc.Append(bson.EC.String(
				"starredCategory", categories.Starred,
			))
		}
	}

	var opts []findopt.Find

	if count > 0 {
		opts = append(opts, findopt.Limit(int64(count)))
	}

	c, err := db.Collection("images").Find(ctx, doc, opts...)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get cursor for images in database.")
		return nil, err
	}

	defer c.Close(ctx)

	var images []model.Image

	for c.Next(ctx) {

		image := model.Image{}
		err := c.Decode(&image)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to parse image object from bson during select of images.")
			return nil, err
		}
		images = append(images, image)
	}
	if err = c.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Cursor failed during listing of images in database.")
		return nil, err
	}

	return images, nil
}

// UpsertImage inserts or updates an existing image in the db.
func UpsertImage(image model.Image) error {
	var dbImage model.Image
	db, ctx, err := GetConnection()

	if err != nil {
		return err
	}

	doc := bson.NewDocument(bson.EC.String("file", image.File))
	err = db.Collection("images").FindOne(ctx, doc).Decode(&dbImage)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Unable to find image in db. Proceeding to create image.")
	}

	_, err = db.Collection("images").ReplaceOne(
		ctx,
		doc,
		image,
		replaceopt.Upsert(true),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"image": image,
			"error": err,
		}).Error("Failed to update image object in database.")
		return err
	}

	return nil
}
