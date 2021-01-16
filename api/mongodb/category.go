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

// UpsertCategory inserts a category or updates it if a category with the same name already exists.
// The name is compared case insensitive.
func UpsertCategory(category model.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client().Database(config.Get().Database).Collection("category")

	opts := options.Replace().SetCollation(&options.Collation{Strength: 2, Locale: "en"}).SetUpsert(true)
	_, err := collection.ReplaceOne(ctx, bson.M{"name": category.Name}, category, opts)

	return err
}

// DeleteCategory deletes a category.
func DeleteCategory(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client().Database(config.Get().Database).Collection("category")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID}, options.Delete())

	return err

}
