package character

import "context"

type CharacterApi interface {
	Get(ctx context.Context, name string) (*CharacterEntity, error)
}
