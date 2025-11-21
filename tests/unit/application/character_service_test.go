package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	app "github.com/heaveless/dbz-api/internal/application/character"
	domain "github.com/heaveless/dbz-api/internal/domain/character"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCharacterRepository struct {
	mock.Mock
}

func (m *MockCharacterRepository) Get(ctx context.Context, name string) (*domain.CharacterEntity, error) {
	args := m.Called(ctx, name)

	var chr *domain.CharacterEntity
	if v := args.Get(0); v != nil {
		chr = v.(*domain.CharacterEntity)
	}

	return chr, args.Error(1)
}

func (m *MockCharacterRepository) Create(ctx context.Context, c *domain.CharacterEntity) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

type MockCharacterApi struct {
	mock.Mock
}

func (m *MockCharacterApi) Get(ctx context.Context, name string) (*domain.CharacterEntity, error) {
	args := m.Called(ctx, name)

	var chr *domain.CharacterEntity
	if v := args.Get(0); v != nil {
		chr = v.(*domain.CharacterEntity)
	}

	return chr, args.Error(1)
}

func TestCharacterService_GetByName_FromRepo(t *testing.T) {
	ctx := context.Background()
	repo := new(MockCharacterRepository)
	api := new(MockCharacterApi)

	svc := app.NewCharacterService(repo, api)

	entity := &domain.CharacterEntity{
		Id:          1,
		Name:        "Goku",
		Ki:          "9000",
		MaxKi:       "10000",
		Race:        "Saiyan",
		Gender:      "Male",
		Image:       "img",
		Affiliation: "Z",
	}

	repo.
		On("Get", mock.Anything, "Goku").
		Return(entity, nil)

	repo.
		On("Create", mock.Anything, entity).
		Return(nil)

	api.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)

	dto, err := svc.GetByName(ctx, "Goku")
	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, entity.Id, dto.Id)
	assert.Equal(t, entity.Name, dto.Name)
	assert.Equal(t, entity.Ki, dto.Ki)
	assert.Equal(t, entity.MaxKi, dto.MaxKi)
	assert.Equal(t, entity.Race, dto.Race)
	assert.Equal(t, entity.Gender, dto.Gender)
	assert.Equal(t, entity.Image, dto.Image)
	assert.Equal(t, entity.Affiliation, dto.Affiliation)

	time.Sleep(10 * time.Millisecond)

	repo.AssertCalled(t, "Create", mock.Anything, entity)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetByName_FallbackToApi(t *testing.T) {
	ctx := context.Background()
	repo := new(MockCharacterRepository)
	api := new(MockCharacterApi)

	svc := app.NewCharacterService(repo, api)

	entityFromApi := &domain.CharacterEntity{
		Id:          2,
		Name:        "Vegeta",
		Ki:          "8500",
		MaxKi:       "9500",
		Race:        "Saiyan",
		Gender:      "Male",
		Image:       "img2",
		Affiliation: "Z",
	}

	repo.
		On("Get", mock.Anything, "Vegeta").
		Return((*domain.CharacterEntity)(nil), errors.New("db error"))

	api.
		On("Get", mock.Anything, "Vegeta").
		Return(entityFromApi, nil)

	repo.
		On("Create", mock.Anything, entityFromApi).
		Return(nil)

	dto, err := svc.GetByName(ctx, "Vegeta")
	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, entityFromApi.Id, dto.Id)
	assert.Equal(t, entityFromApi.Name, dto.Name)

	time.Sleep(10 * time.Millisecond)

	repo.AssertCalled(t, "Create", mock.Anything, entityFromApi)
	repo.AssertExpectations(t)
	api.AssertExpectations(t)
}

func TestCharacterService_GetByName_ErrorWhenRepoAndApiFail(t *testing.T) {
	ctx := context.Background()
	repo := new(MockCharacterRepository)
	api := new(MockCharacterApi)

	svc := app.NewCharacterService(repo, api)

	repo.
		On("Get", mock.Anything, "Piccolo").
		Return((*domain.CharacterEntity)(nil), errors.New("db error"))

	api.
		On("Get", mock.Anything, "Piccolo").
		Return((*domain.CharacterEntity)(nil), errors.New("api error"))

	dto, err := svc.GetByName(ctx, "Piccolo")
	assert.Nil(t, dto)
	assert.Error(t, err)
	assert.Equal(t, "api error", err.Error())

	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	api.AssertExpectations(t)
}
