package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-redis/redis/v9"
	"github.com/owlint/maestro/internal/errors"
)

type HealthcheckResponse struct {
	OK bool `json:"ok"`
}

// HealthcheckEndpoint creates a endpoint for healthcheck
func HealthcheckEndpoint(client *redis.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, err := client.Ping(ctx).Result()
		if err != nil {
			return HealthcheckResponse{
				OK: false,
			}, errors.RedisError{Origin: err}
		}

		return HealthcheckResponse{
			OK: true,
		}, nil
	}
}
