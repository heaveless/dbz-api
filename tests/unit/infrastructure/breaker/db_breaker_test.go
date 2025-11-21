package breaker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestBreaker_FindOne_OK(t *testing.T) {
	ctx := context.Background()
	mockCol := new(MockMongoCollection)
	mockSR := new(MockSingleResultWrapper)

	mockCol.
		On("FindOne", ctx, mock.Anything).
		Return(mockSR, nil)

	mockSR.
		On("Err").
		Return(nil)

	cb := breaker.NewDbCollectionWithBreaker(mockCol, time.Millisecond*50)

	res, err := cb.FindOne(ctx, map[string]any{"name": "Goku"})
	assert.NoError(t, err)

	mockSR.
		On("Decode", mock.Anything).
		Return(nil)

	err = res.Decode(&struct{}{})
	assert.NoError(t, err)

	mockCol.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestBreaker_FindOne_NoDocuments(t *testing.T) {
	ctx := context.Background()
	mockCol := new(MockMongoCollection)
	mockSR := new(MockSingleResultWrapper)

	mockCol.
		On("FindOne", ctx, mock.Anything).
		Return(mockSR, nil)

	mockSR.
		On("Err").
		Return(mongo.ErrNoDocuments)

	cb := breaker.NewDbCollectionWithBreaker(mockCol, time.Millisecond*50)

	_, err := cb.FindOne(ctx, map[string]any{"name": "Goku"})
	assert.NoError(t, err)

	mockCol.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestBreaker_FindOne_Error(t *testing.T) {
	ctx := context.Background()
	mockCol := new(MockMongoCollection)
	mockSR := new(MockSingleResultWrapper)

	mockCol.
		On("FindOne", ctx, mock.Anything).
		Return(mockSR, nil)

	mockSR.
		On("Err").
		Return(errors.New("db error"))

	cb := breaker.NewDbCollectionWithBreaker(mockCol, time.Millisecond*50)

	res, err := cb.FindOne(ctx, map[string]any{"name": "Goku"})

	assert.Nil(t, res)
	assert.EqualError(t, err, "db error")

	mockCol.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

func TestBreaker_InsertOne_OK(t *testing.T) {
	ctx := context.Background()
	mockCol := new(MockMongoCollection)

	mockCol.
		On("InsertOne", ctx, mock.Anything).
		Return(&mongo.InsertOneResult{InsertedID: "123"}, nil)

	cb := breaker.NewDbCollectionWithBreaker(mockCol, time.Millisecond*50)

	res, err := cb.InsertOne(ctx, map[string]any{})
	assert.NoError(t, err)
	assert.Equal(t, "123", res.InsertedID)

	mockCol.AssertExpectations(t)
}

func TestBreaker_InsertOne_Error(t *testing.T) {
	ctx := context.Background()
	mockCol := new(MockMongoCollection)

	mockCol.
		On("InsertOne", ctx, mock.Anything).
		Return((*mongo.InsertOneResult)(nil), errors.New("insert error"))

	cb := breaker.NewDbCollectionWithBreaker(mockCol, time.Millisecond*50)

	res, err := cb.InsertOne(ctx, map[string]any{})
	assert.Nil(t, res)
	assert.EqualError(t, err, "insert error")

	mockCol.AssertExpectations(t)
}
