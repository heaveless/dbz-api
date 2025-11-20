package repositoy

import (
	"context"

	domain "github.com/heaveless/dbz-api/internal/domain/character"
	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type characterRepository struct {
	client breaker.DbCollection
}

func NewCharacterRepository(client breaker.DbCollection) domain.CharacterRepository {
	return &characterRepository{
		client: client,
	}
}

func (repo *characterRepository) Create(ctx context.Context, record *domain.CharacterEntity) error {
	_, err := repo.client.InsertOne(ctx, record)

	return err
}

func (repo *characterRepository) Get(ctx context.Context, name string) (*domain.CharacterEntity, error) {
	res, err := repo.client.FindOne(ctx, bson.M{"name": name})
	if err != nil {
		return nil, err
	}

	var record domain.CharacterEntity
	if err := res.Decode(&record); err != nil {
		return nil, err
	}

	return &record, err
}
