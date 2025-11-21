package breaker_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"

	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
)

type fakeRoundTripper struct {
	fn    func(req *http.Request) (*http.Response, error)
	calls int
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.calls++
	return f.fn(req)
}

func newTestHttpWithBreaker(rt http.RoundTripper, settings gobreaker.Settings) breaker.ExternalClient {
	client := &http.Client{
		Timeout:   time.Second,
		Transport: rt,
	}
	return breaker.NewHttpWithBreakerRef(client, settings)
}

func TestHttpWithCircuitBreaker_Do_Success(t *testing.T) {
	fakeRT := &fakeRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok")),
				Request:    req,
			}, nil
		},
	}

	settings := gobreaker.Settings{
		Name:        "http-breaker-test-success",
		MaxRequests: 5,
		Timeout:     time.Second,
	}

	cb := newTestHttpWithBreaker(fakeRT, settings)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	res, err := cb.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 1, fakeRT.calls)
}

func TestHttpWithCircuitBreaker_Do_PropagatesError(t *testing.T) {
	fakeRT := &fakeRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	settings := gobreaker.Settings{
		Name:        "http-breaker-test-error",
		MaxRequests: 5,
		Timeout:     time.Second,
	}

	cb := newTestHttpWithBreaker(fakeRT, settings)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	res, err := cb.Do(req)
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "network error")
	assert.Equal(t, 1, fakeRT.calls)
}

func TestHttpWithCircuitBreaker_Do_OpenCircuit(t *testing.T) {
	fakeRT := &fakeRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	settings := gobreaker.Settings{
		Name:        "http-breaker-open-test",
		MaxRequests: 5,
		Timeout:     time.Second,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures >= 1
		},
	}

	cb := newTestHttpWithBreaker(fakeRT, settings)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	res1, err1 := cb.Do(req)
	assert.Nil(t, res1)
	assert.Error(t, err1)
	assert.ErrorContains(t, err1, "network error")
	assert.Equal(t, 1, fakeRT.calls)

	res2, err2 := cb.Do(req)
	assert.Nil(t, res2)
	assert.Error(t, err2)
	assert.ErrorIs(t, err2, gobreaker.ErrOpenState)
	assert.Equal(t, 1, fakeRT.calls)
}
