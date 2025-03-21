package limiter

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

type Repository interface {
	IncrementRequestCount(ctx context.Context, key string) (int, error)
	FreezeRequestCount(ctx context.Context, key string, seconds int) error
	IsFrozen(ctx context.Context, key string) (bool, float64, error)
}

type Limiter struct {
	repository          Repository
	maxRequestsIP       int
	maxRequestsAPIToken int
	freezeTime          int
}

type Request struct {
	IPAddress string
	APIKey    string
}

type AllowResponse struct {
	Allowed   bool
	Remaining int
}

const (
	defaultMaxRequestsPerSecondIP       = 10
	defaultMaxRequestsPerSecondAPIToken = 20
	defaultFreezeTimeInSeconds          = 60
)

func New(repository Repository) *Limiter {
	maxRequestsIPStr := os.Getenv("MAX_REQUESTS_PER_SECOND_IP")
	maxRequestsIP, err := strconv.Atoi(maxRequestsIPStr)
	if err != nil {
		maxRequestsIP = defaultMaxRequestsPerSecondIP
	}

	maxRequestsAPITokenStr := os.Getenv("MAX_REQUESTS_PER_SECOND_API_TOKEN")
	maxRequestsAPIToken, err := strconv.Atoi(maxRequestsAPITokenStr)
	if err != nil {
		maxRequestsAPIToken = defaultMaxRequestsPerSecondAPIToken
	}

	freezeTimeStr := os.Getenv("FREZEE_TIME_IN_SECONDS")
	freezeTime, err := strconv.Atoi(freezeTimeStr)
	if err != nil {
		freezeTime = defaultFreezeTimeInSeconds
	}

	return &Limiter{
		repository:          repository,
		maxRequestsIP:       maxRequestsIP,
		maxRequestsAPIToken: maxRequestsAPIToken,
		freezeTime:          freezeTime,
	}
}

func (l *Limiter) Allow(ctx context.Context, IP, APIToken string) (AllowResponse, error) {
	switch {
	case APIToken != "":
		isAllowed, remaining, err := l.validate(ctx, APIToken, l.maxRequestsAPIToken)
		if err != nil {
			return AllowResponse{}, fmt.Errorf("Error validating API token: %w", err)
		}
		return AllowResponse{
			Allowed:   isAllowed,
			Remaining: remaining,
		}, nil

	case IP != "":
		isAllowed, remaining, err := l.validate(ctx, IP, l.maxRequestsIP)
		if err != nil {
			return AllowResponse{}, fmt.Errorf("Error validating IP: %w", err)
		}
		return AllowResponse{
			Allowed:   isAllowed,
			Remaining: remaining,
		}, nil

	default:
		return AllowResponse{}, fmt.Errorf("IP or API token must be provided")
	}
}

func (l *Limiter) validate(ctx context.Context, key string, maxRequests int) (bool, int, error) {
	frozen, _, err := l.repository.IsFrozen(ctx, key)
	if err != nil {
		return false, 0, fmt.Errorf("error checking if key is frozen: %w", err)
	}
	if frozen {
		return false, 0, nil
	}

	var requests int

	requests, err = l.repository.IncrementRequestCount(ctx, key)
	if err != nil {
		return false, 0, fmt.Errorf("Error incrementing request count by key: %w", err)
	}

	if requests > maxRequests {
		err := l.repository.FreezeRequestCount(ctx, key, l.freezeTime)
		if err != nil {
			return false, 0, fmt.Errorf("Error freezing requests by key: %w", err)
		}
		return false, 0, nil
	}

	return true, maxRequests - requests, nil
}
