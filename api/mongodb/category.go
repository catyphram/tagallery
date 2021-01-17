package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tagallery.com/api/config"
	"tagallery.com/api/model"
)

// QueryCategories returns all categories.
func QueryCategories() ([]model.Category, error) {
	client := Client()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(config.Get().Database).Collection("category")

	cur, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	categories := []model.Category{}

	if err := cur.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// UpsertCategory inserts a category or updates it.
// If the provided category has a valid id, then the entire document is updated.
// If the id is missing, then the name is taken as an identifier and everything else is updated.
// You can also provide a valid ObjectId {category.id} for a new category.
// The name is compared case insensitive and must be unique.
func UpsertCategory(category model.Category) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client().Database(config.Get().Database).Collection("category")
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{}

	if category.ID != nil && len(*category.ID) > 0 {
		objectID, err := primitive.ObjectIDFromHex(*category.ID)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		filter["_id"] = objectID
	} else {
		filter["name"] = category.Name
	}
	category.ID = nil

	result, err := collection.ReplaceOne(ctx, filter, category, opts)

	if err != nil {
		return nil, err
	}

	if result.UpsertedID != nil {
		objectID := result.UpsertedID.(primitive.ObjectID).Hex()
		category.ID = &objectID
	}
	return &category, nil
}

// DeleteCategory deletes a category.
func DeleteCategory(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client().Database(config.Get().Database).Collection("category")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidObjectID
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID}, options.Delete())

	return err

}
