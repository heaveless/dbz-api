package breaker

import (
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type ExternalClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpWithCircuitBreaker struct {
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewHttpWithBreaker(timeout time.Duration) ExternalClient {
	settings := gobreaker.Settings{
		Name:        "http-breaker",
		MaxRequests: 5,
		Timeout:     5 * time.Second,
	}

	return NewHttpWithBreakerRef(
		&http.Client{Timeout: timeout},
		settings,
	)
}

func NewHttpWithBreakerRef(client *http.Client, settings gobreaker.Settings) ExternalClient {
	return &HttpWithCircuitBreaker{
		client:         client,
		circuitBreaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *HttpWithCircuitBreaker) Do(req *http.Request) (*http.Response, error) {
	res, err := c.circuitBreaker.Execute(func() (any, error) {
		return c.client.Do(req)
	})
	if err != nil {
		return nil, err
	}

	return res.(*http.Response), nil
}
