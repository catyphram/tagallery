// Package mongodb contains init and access functions for mongodb
// as well as some CRUD methods to work with our collections.
package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// DatabaseOptions structs a database name and host string
type DatabaseOptions struct {
	Database string
	Host     string
}

var (
	_db  *mongo.Database
	_ctx context.Context
)

type ctxKey string

// Init establishes a new connection to a database.
func Init(options DatabaseOptions) (*mongo.Database, context.Context, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey("database"), options.Database)
	ctx = context.WithValue(ctx, ctxKey("host"), options.Host)

	db, err := configDB(ctx)

	if err != nil {
		Close()
		return nil, nil, err
	}

	_ctx = ctx
	_db = db

	return _db, _ctx, nil
}

// Close closes the current connection to the db, if there is one open.
func Close() error {
	if _db != nil {
		err := _db.Client().Disconnect(_ctx)
		_db = nil
		_ctx = nil
		return err
	}
	return nil
}

// GetConnection returns the connection to the db if there is one open.
// If no connection is open, an error will be returned.
func GetConnection() (*mongo.Database, context.Context, error) {
	if _db == nil {
		return nil, nil, errors.New("The database has to be loaded first. Call Init() to do so.")
	}

	return _db, _ctx, nil
}

// configDB opens the db connection and selects the database.
func configDB(ctx context.Context) (*mongo.Database, error) {
	uri := fmt.Sprintf(`mongodb://%s`,
		ctx.Value(ctxKey("host")),
	)
	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database(ctx.Value(ctxKey("database")).(string))
	return db, nil
}
