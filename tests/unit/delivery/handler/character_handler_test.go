package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/heaveless/dbz-api/internal/delivery/http/handler"
	domain "github.com/heaveless/dbz-api/internal/domain/character"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCharacterService struct {
	mock.Mock
}

func (m *MockCharacterService) GetByName(ctx context.Context, name string) (*domain.CharacterDTO, error) {
	args := m.Called(ctx, name)

	var chr *domain.CharacterDTO
	if v := args.Get(0); v != nil {
		chr = v.(*domain.CharacterDTO)
	}

	return chr, args.Error(1)
}

func setupRouter(h *handler.CharacterHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/characters", h.GetOne)
	return r
}

func TestCharacterHandler_GetOne_InvalidBody(t *testing.T) {
	svc := new(MockCharacterService)
	h := handler.NewCharacterHandler(svc)
	router := setupRouter(h)

	body := bytes.NewBufferString(`{}`)

	req, _ := http.NewRequest(http.MethodPost, "/characters", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "The data submitted is invalid.", resp["message"])

	svc.AssertNotCalled(t, "GetByName", mock.Anything, mock.Anything)
}

func TestCharacterHandler_GetOne_ServiceError(t *testing.T) {
	svc := new(MockCharacterService)
	h := handler.NewCharacterHandler(svc)
	router := setupRouter(h)

	svc.
		On("GetByName", mock.Anything, "Goku").
		Return((*domain.CharacterDTO)(nil), errors.New("character not found"))

	body := bytes.NewBufferString(`{"name":"Goku"}`)

	req, _ := http.NewRequest(http.MethodPost, "/characters", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "character not found", resp["message"])

	svc.AssertExpectations(t)
}

func TestCharacterHandler_GetOne_OK(t *testing.T) {
	svc := new(MockCharacterService)
	h := handler.NewCharacterHandler(svc)
	router := setupRouter(h)

	expected := &domain.CharacterDTO{
		Id:   1,
		Name: "Goku",
	}

	svc.
		On("GetByName", mock.Anything, "Goku").
		Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Goku"}`)

	req, _ := http.NewRequest(http.MethodPost, "/characters", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)

	var resp struct {
		Data *domain.CharacterDTO `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, expected.Id, resp.Data.Id)
	assert.Equal(t, expected.Name, resp.Data.Name)

	svc.AssertExpectations(t)
}
