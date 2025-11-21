package repository_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/heaveless/dbz-api/internal/domain/character"
	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
	repo "github.com/heaveless/dbz-api/internal/infrastructure/repositoy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MockDbCollection struct {
	mock.Mock
}

func (m *MockDbCollection) FindOne(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.FindOneOptions],
) (breaker.SingleResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(breaker.SingleResult), args.Error(1)
}

func (m *MockDbCollection) InsertOne(
	ctx context.Context,
	document any,
	opts ...options.Lister[options.InsertOneOptions],
) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v any) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockSingleResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

func TestCharacterRepository_Create_OK(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockDbCollection)

	mockClient.
		On("InsertOne", ctx, mock.Anything).
		Return(&mongo.InsertOneResult{
			InsertedID: int64(1),
		}, nil)

	r := repo.NewCharacterRepository(mockClient)

	err := r.Create(ctx, &domain.CharacterEntity{
		Id:          1,
		Name:        "Goku",
		Ki:          "Ki",
		MaxKi:       "MaxKi",
		Race:        "Race",
		Gender:      "Gender",
		Image:       "Image",
		Affiliation: "Affiliation",
	})

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestCharacterRepository_Create_Error(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockDbCollection)

	mockClient.
		On("InsertOne", ctx, mock.Anything).
		Return((*mongo.InsertOneResult)(nil), errors.New("db error"))

	r := repo.NewCharacterRepository(mockClient)

	err := r.Create(ctx, &domain.CharacterEntity{})

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockClient.AssertExpectations(t)
}

func TestCharacterRepository_Get_OK(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockDbCollection)
	mockResult := new(MockSingleResult)

	entity := domain.CharacterEntity{
		Id:          1,
		Name:        "Goku",
		Ki:          "Ki",
		MaxKi:       "MaxKi",
		Race:        "Race",
		Gender:      "Gender",
		Image:       "Image",
		Affiliation: "Affiliation",
	}

	mockClient.
		On("FindOne", ctx, mock.Anything).
		Return(mockResult, nil)

	mockResult.
		On("Err").
		Return(nil)

	mockResult.
		On("Decode", mock.AnythingOfType("*character.CharacterEntity")).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*domain.CharacterEntity)
			*arg = entity
		}).
		Return(nil)

	r := repo.NewCharacterRepository(mockClient)

	res, err := r.Get(ctx, "Goku")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Id)
	assert.Equal(t, "Goku", res.Name)
	mockClient.AssertExpectations(t)
	mockResult.AssertNotCalled(t, "Err")
}

func TestCharacterRepository_Get_FindError(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockDbCollection)

	mockClient.
		On("FindOne", ctx, mock.Anything).
		Return((*MockSingleResult)(nil), errors.New("find error"))

	r := repo.NewCharacterRepository(mockClient)

	res, err := r.Get(ctx, "Goku")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "find error", err.Error())
	mockClient.AssertExpectations(t)
}

func TestCharacterRepository_Get_DecodeError(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockDbCollection)
	mockResult := new(MockSingleResult)

	mockClient.
		On("FindOne", ctx, mock.Anything).
		Return(mockResult, nil)

	mockResult.
		On("Decode", mock.Anything).
		Return(errors.New("decode error"))

	r := repo.NewCharacterRepository(mockClient)

	res, err := r.Get(ctx, "Goku")

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "decode error", err.Error())
	mockClient.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}
