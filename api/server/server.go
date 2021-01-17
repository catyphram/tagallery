package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"tagallery.com/api/config"
	"tagallery.com/api/logger"
	"tagallery.com/api/mongodb"
)

// StartServer loads the config, sets up the logger, establishes a connection to the db and starts the router.
func StartServer() {
	config := config.Load()

	log := logger.Setup(config.Debug)

	if client, err := mongodb.Connect(context.Background(), fmt.Sprintf(`mongodb://%s`, config.DatabaseHost)); err != nil {
		log.Fatalw("Unable to connect to database.", "error", err)
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err = client.Disconnect(ctx); err != nil {
				log.Fatalw("Unable to disconnect from database.", "error", err)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatalw("Ping to database not successful.", "error", err)
		}
		if err := setupDatabase(ctx, client); err != nil {
			log.Fatalw("Unable to setup the database.", "error", err)
		}

	}

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	defer log.Sync()

	r := ConfigureRouter()
	r.Run(fmt.Sprintf(":%d", config.Port))
}

// setupDatabase creates a unique index on the category name.
func setupDatabase(ctx context.Context, client *mongo.Client) error {
	categoryCollection := client.Database(config.Get().Database).Collection("category")

	// Creating MongoDb indexes is an idempotent operation.
	// Therefore we don't have to check if the index already exists.
	_, err := categoryCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName("name_unique").SetUnique(true).SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2,
		}),
	})

	return err
}
