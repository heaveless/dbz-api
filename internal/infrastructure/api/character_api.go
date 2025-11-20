package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	domain "github.com/heaveless/dbz-api/internal/domain/character"
	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
)

type characterApi struct {
	baseURL string
	client  breaker.ExternalClient
}

func NewCharacterApi(baseURL string, client breaker.ExternalClient) domain.CharacterApi {
	return &characterApi{
		baseURL: baseURL,
		client:  client,
	}
}

func (api *characterApi) Get(ctx context.Context, name string) (*domain.CharacterEntity, error) {
	endpoint := fmt.Sprintf("%s/api/characters?name=%s", api.baseURL, url.QueryEscape(name))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	res, err := api.client.Do(req)
	if err != nil {
		return nil, errors.New("service temporarily unavailable, please try again later")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var characters []domain.CharacterEntity
	if err := json.NewDecoder(res.Body).Decode(&characters); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(characters) == 0 {
		return nil, fmt.Errorf("character not found")
	}

	return &characters[0], nil
}
