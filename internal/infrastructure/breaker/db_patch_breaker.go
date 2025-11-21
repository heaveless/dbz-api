package breaker

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SingleResult interface {
	Decode(v any) error
	Err() error
}

type MongoSingleResultWrapper struct {
	Sr *mongo.SingleResult
}

func (w *MongoSingleResultWrapper) Decode(v any) error {
	return w.Sr.Decode(v)
}

func (w *MongoSingleResultWrapper) Err() error {
	return w.Sr.Err()
}

func WrapMongoSingleResult(sr *mongo.SingleResult) SingleResult {
	return &MongoSingleResultWrapper{Sr: sr}
}

type MongoDbCollection struct {
	col *mongo.Collection
}

func NewMongoDbCollection(col *mongo.Collection) DbCollection {
	return &MongoDbCollection{col: col}
}

func (r *MongoDbCollection) FindOne(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.FindOneOptions],
) (SingleResult, error) {
	sr := r.col.FindOne(ctx, filter, opts...)
	return WrapMongoSingleResult(sr), nil
}

func (r *MongoDbCollection) InsertOne(
	ctx context.Context,
	document any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, error) {
	return r.col.InsertOne(ctx, document, opts...)
}
