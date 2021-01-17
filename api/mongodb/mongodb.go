package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrInvalidObjectID indicates that a provided string is not a valid object id.
var ErrInvalidObjectID = errors.New("the provided string is not a valid ObjectID")

var client *mongo.Client

// Client returns the mongodb client. Make sure to call Connect() beforehand.
func Client() *mongo.Client {
	return client
}

// Connect opens a database connection.
func Connect(ctx context.Context, dbURI string) (*mongo.Client, error) {
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))

	if err != nil {
		return nil, err
	}

	client = dbClient

	return dbClient, nil
}
