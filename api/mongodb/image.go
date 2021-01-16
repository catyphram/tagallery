package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tagallery.com/api/config"
	"tagallery.com/api/logger"
	"tagallery.com/api/model"
)

// DBImage extends a model.Image by an Id.
type DBImage struct {
	model.Image
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

// GetImages queries the database for images.
// If count is set then no more then {ops.count} images will be returned.
// A *CategoryMap may be passed to filter only images that are in all of these categories.
// If categories == nil then instead of (auto)-categorized images,
// only images that have no assigned category will be returned.
// With lastImage you get only images after this one. Used for pagination.
func GetImages(opts model.ImageOptions, categories *model.CategoryMap) ([]model.Image, error) {
	var dbLastImage DBImage
	doc := bson.D{}

	collection := Client().Database(config.Get().Database).Collection("image")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if opts.LastImage != nil {
		err := collection.FindOne(ctx, bson.M{"file": opts.LastImage}).Decode(&dbLastImage)

		if err != nil {
			logger.Logger().Warnw("Unable to find lastImage in the database.",
				"lastImage", *opts.LastImage,
				"error", err,
			)
		} else {
			doc = append(doc, bson.E{
				Key: "_id", Value: bson.M{
					"$gt": dbLastImage.ID,
				},
			})
		}
	}

	// Uncategorized images only
	if categories == nil {
		doc = append(doc, bson.E{Key: "$or", Value: bson.A{
			bson.M{"assignedCategories": bson.M{"$size": 0}},
			bson.M{"assignedCategories": bson.M{"$eq": nil}},
		}}, bson.E{Key: "$or", Value: bson.A{
			bson.M{"proposedCategories": bson.M{"$size": 0}},
			bson.M{"proposedCategories": bson.M{"$eq": nil}},
		}}, bson.E{Key: "$or", Value: bson.A{
			bson.M{"starredCategory": bson.M{"$in": bson.A{nil, ""}}},
		}})
	} else {
		if categories.Assigned != nil {
			if len(categories.Assigned) > 0 {
				doc = append(doc, bson.E{Key: "assignedCategories", Value: bson.M{
					"$all": categories.Assigned},
				})
			} else {
				doc = append(doc, bson.E{Key: "$and", Value: bson.A{
					bson.M{"assignedCategories": bson.M{"$ne": nil}},
					bson.M{"assignedCategories": bson.M{"$ne": []string{}}},
				}})
			}
		}

		if categories.Proposed != nil {
			if len(categories.Proposed) > 0 {
				doc = append(doc, bson.E{Key: "proposedCategories", Value: bson.M{
					"$all": categories.Proposed},
				})
			} else {
				doc = append(doc, bson.E{Key: "$and", Value: bson.A{
					bson.M{"proposedCategories": bson.M{"$ne": nil}},
					bson.M{"proposedCategories": bson.M{"$ne": []string{}}},
				}})
			}
		}

		if categories.Starred != nil {
			doc = append(doc, bson.E{Key: "starredCategory", Value: categories.Starred})
		}
	}

	cur, err := collection.Find(ctx, doc, options.Find().SetLimit(int64(*opts.Count)))

	if err != nil {
		return nil, err
	}

	images := []model.Image{}

	if err := cur.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// UpsertImage inserts or updates an existing image in the db.
func UpsertImage(image model.Image) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client().Database(config.Get().Database).Collection("image")

	opts := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(ctx, bson.M{"file": image.File}, image, opts)

	return err
}
