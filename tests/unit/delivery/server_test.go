package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestCharacterHandler_GetOne_Direct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := new(MockCharacterService)
	h := handler.NewCharacterHandler(svc)

	body := bytes.NewBufferString(`{"name":"Goku"}`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/characters", body)
	c.Request.Header.Set("Content-Type", "application/json")

	expected := &domain.CharacterDTO{
		Id:   1,
		Name: "Goku",
	}

	svc.
		On("GetByName", mock.Anything, "Goku").
		Return(expected, nil)

	h.GetOne(c)

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
