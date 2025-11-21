package breaker_test

import (
	"context"

	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) Name() string {
	return "characters"
}

func (m *MockMongoCollection) FindOne(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.FindOneOptions],
) (breaker.SingleResult, error) {

	args := m.Called(ctx, filter)

	var sr breaker.SingleResult
	if v := args.Get(0); v != nil {
		sr = v.(breaker.SingleResult)
	}

	return sr, args.Error(1)
}

func (m *MockMongoCollection) InsertOne(
	ctx context.Context,
	document any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, error) {

	args := m.Called(ctx, document)

	var res *mongo.InsertOneResult
	if v := args.Get(0); v != nil {
		res = v.(*mongo.InsertOneResult)
	}

	return res, args.Error(1)
}
