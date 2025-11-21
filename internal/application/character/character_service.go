package character

import (
	"context"
	"errors"
	"log"
	"time"

	domain "github.com/heaveless/dbz-api/internal/domain/character"
	utils "github.com/heaveless/dbz-api/internal/utils"
)

type CharacterService struct {
	repo domain.CharacterRepository
	api  domain.CharacterApi
}

func NewCharacterService(dr domain.CharacterRepository, hr domain.CharacterApi) *CharacterService {
	return &CharacterService{
		repo: dr,
		api:  hr,
	}
}

func (s *CharacterService) GetByName(ctx context.Context, name string) (*domain.CharacterDTO, error) {
	chr, err := utils.WithFallback(ctx,
		func(ctx context.Context) (*domain.CharacterEntity, error) {
			return s.repo.Get(ctx, name)
		},
		func(ctx context.Context) (*domain.CharacterEntity, error) {
			return s.api.Get(ctx, name)
		},
		func(err error) bool {
			return errors.Is(err, context.DeadlineExceeded) || true
		},
	)

	if err != nil {
		return nil, err
	}

	go func(c *domain.CharacterEntity) {
		saveCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		saveErr := s.repo.Create(saveCtx, c)
		if saveErr != nil {
			log.Printf("[DB] failed to save user %d from api: %v", c.Id, saveErr)
		}
	}(chr)

	return &domain.CharacterDTO{
		Id:          chr.Id,
		Name:        chr.Name,
		Ki:          chr.Ki,
		MaxKi:       chr.MaxKi,
		Race:        chr.Race,
		Gender:      chr.Gender,
		Image:       chr.Image,
		Affiliation: chr.Affiliation,
	}, nil
}
