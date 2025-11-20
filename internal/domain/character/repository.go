package character

import "context"

type CharacterRepository interface {
	Get(ctx context.Context, name string) (*CharacterEntity, error)
	Create(ctx context.Context, c *CharacterEntity) error
}
