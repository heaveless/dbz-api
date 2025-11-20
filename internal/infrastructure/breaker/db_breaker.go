package breaker

import (
	"context"
	"errors"
	"time"

	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DbCollection interface {
	FindOne(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.FindOneOptions],
	) (*mongo.SingleResult, error)

	InsertOne(
		ctx context.Context,
		document any,
		opts ...options.Lister[options.InsertOneOptions],
	) (*mongo.InsertOneResult, error)
}

type DbCollectionWithBreaker struct {
	collection     *mongo.Collection
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewDbCollectionWithBreaker(collection *mongo.Collection, timeout time.Duration) DbCollection {
	settings := gobreaker.Settings{
		Name:        "db-breaker:" + collection.Name(),
		MaxRequests: 5,
		Timeout:     timeout,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			if c.Requests < 10 {
				return false
			}
			errRate := float64(c.TotalFailures) / float64(c.Requests)
			return errRate >= 0.5
		},
	}

	return &DbCollectionWithBreaker{
		collection:     collection,
		circuitBreaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *DbCollectionWithBreaker) FindOne(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.FindOneOptions],
) (*mongo.SingleResult, error) {
	res, err := c.circuitBreaker.Execute(func() (any, error) {
		sr := c.collection.FindOne(ctx, filter, opts...)
		if err := sr.Err(); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return sr, nil
			}
			return nil, err
		}
		return sr, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*mongo.SingleResult), nil
}

func (c *DbCollectionWithBreaker) InsertOne(
	ctx context.Context,
	document any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, error) {
	res, err := c.circuitBreaker.Execute(func() (any, error) {
		return c.collection.InsertOne(ctx, document, opts...)
	})
	if err != nil {
		return nil, err
	}
	return res.(*mongo.InsertOneResult), nil
}
