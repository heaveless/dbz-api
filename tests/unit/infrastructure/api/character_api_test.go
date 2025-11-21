package api_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	api "github.com/heaveless/dbz-api/internal/infrastructure/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCharacterApi_Get_OK(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockExternalClient)

	body := `[{"id":1,"name":"Goku"}]`

	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(resp, nil)

	sut := api.NewCharacterApi("http://test.com", mockClient)

	res, err := sut.Get(ctx, "Goku")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Id)
	assert.Equal(t, "Goku", res.Name)
}

func TestCharacterApi_Get_ClientError(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockExternalClient)

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return((*http.Response)(nil), errors.New("network error"))

	sut := api.NewCharacterApi("http://test.com", mockClient)

	res, err := sut.Get(ctx, "Goku")

	assert.Nil(t, res)
	assert.EqualError(t, err, "service temporarily unavailable, please try again later")
}

func TestCharacterApi_Get_UnexpectedStatus(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockExternalClient)

	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(resp, nil)

	sut := api.NewCharacterApi("http://test.com", mockClient)

	res, err := sut.Get(ctx, "fail")

	assert.Nil(t, res)
	assert.EqualError(t, err, "unexpected status code: 500")
}

func TestCharacterApi_Get_DecodeError(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockExternalClient)

	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("invalid json")),
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(resp, nil)

	sut := api.NewCharacterApi("http://test.com", mockClient)

	res, err := sut.Get(ctx, "Goku")

	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "decoding response")
}

func TestCharacterApi_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockExternalClient)

	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("[]")),
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(resp, nil)

	sut := api.NewCharacterApi("http://test.com", mockClient)

	res, err := sut.Get(ctx, "Unknown")

	assert.Nil(t, res)
	assert.EqualError(t, err, "character not found")
}
