package limiter

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	isFrozenFunc  func(ctx context.Context, key string) (bool, float64, error)
	incrementFunc func(ctx context.Context, key string) (int, error)
	freezeFunc    func(ctx context.Context, key string, seconds int) error
}

func (f *fakeRepo) IsFrozen(ctx context.Context, key string) (bool, float64, error) {
	return f.isFrozenFunc(ctx, key)
}

func (f *fakeRepo) IncrementRequestCount(ctx context.Context, key string) (int, error) {
	return f.incrementFunc(ctx, key)
}

func (f *fakeRepo) FreezeRequestCount(ctx context.Context, key string, seconds int) error {
	return f.freezeFunc(ctx, key, seconds)
}

func TestLimiterAllow(t *testing.T) {
	testCases := []struct {
		name                string
		ip                  string
		apiToken            string
		isFrozenFunc        func(ctx context.Context, key string) (bool, float64, error)
		upsertFunc          func(ctx context.Context, key string) (int, error)
		freezeFunc          func(ctx context.Context, key string, seconds int) error
		maxRequestsIP       int
		maxRequestsAPIToken int
		freezeTime          int
		expectedAllowed     bool
		expectedRemaining   int
		expectError         bool
	}{
		{
			name:     "API token frozen",
			ip:       "",
			apiToken: "token1",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return true, 30, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				return 0, nil // Won't be called in this case.
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return nil // Won't be called.
			},
			maxRequestsAPIToken: 10,
			freezeTime:          60,
			expectedAllowed:     false,
			expectedRemaining:   0,
			expectError:         false,
		},
		{
			name:     "IP frozen",
			ip:       "192.168.1.1",
			apiToken: "",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return true, 45, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				return 0, nil
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return nil
			},
			maxRequestsIP:     5,
			freezeTime:        60,
			expectedAllowed:   false,
			expectedRemaining: 0,
			expectError:       false,
		},
		{
			name:     "API token allowed below limit",
			ip:       "",
			apiToken: "token1",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return false, 0, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				// Simulate a current count of 5.
				return 5, nil
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return nil
			},
			maxRequestsAPIToken: 10,
			freezeTime:          60,
			expectedAllowed:     true,
			expectedRemaining:   5,
			expectError:         false,
		},
		{
			name:     "IP allowed below limit",
			ip:       "1.2.3.4",
			apiToken: "",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return false, 0, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				// Simulate a current count of 3.
				return 3, nil
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return nil
			},
			maxRequestsIP:     5,
			freezeTime:        60,
			expectedAllowed:   true,
			expectedRemaining: 2,
			expectError:       false,
		},
		{
			name:     "API token over limit triggers freeze",
			ip:       "",
			apiToken: "token1",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return false, 0, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				// Simulate a current count of 11, which is above the allowed 10.
				return 11, nil
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				// Expect FreezeRequestCount to be called.
				return nil
			},
			maxRequestsAPIToken: 10,
			freezeTime:          60,
			expectedAllowed:     false,
			expectedRemaining:   0,
			expectError:         false,
		},
		{
			name:     "Error in IsFrozen",
			ip:       "1.2.3.4",
			apiToken: "",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return false, 0, fmt.Errorf("isFrozen error")
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				return 0, nil
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return nil
			},
			maxRequestsIP:     5,
			freezeTime:        60,
			expectedAllowed:   false,
			expectedRemaining: 0,
			expectError:       true,
		},
		{
			name:     "Error in UpsertRequestCount",
			ip:       "1.2.3.4",
			apiToken: "",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return false, 0, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				return 0, fmt.Errorf("upsert error")
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return nil
			},
			maxRequestsIP:     5,
			freezeTime:        60,
			expectedAllowed:   false,
			expectedRemaining: 0,
			expectError:       true,
		},
		{
			name:     "Error in FreezeRequestCount",
			ip:       "",
			apiToken: "token1",
			isFrozenFunc: func(ctx context.Context, key string) (bool, float64, error) {
				return false, 0, nil
			},
			upsertFunc: func(ctx context.Context, key string) (int, error) {
				// Simulate a count above the limit.
				return 12, nil
			},
			freezeFunc: func(ctx context.Context, key string, seconds int) error {
				return fmt.Errorf("freeze error")
			},
			maxRequestsAPIToken: 10,
			freezeTime:          60,
			expectedAllowed:     false,
			expectedRemaining:   0,
			expectError:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{
				isFrozenFunc:  tc.isFrozenFunc,
				incrementFunc: tc.upsertFunc,
				freezeFunc:    tc.freezeFunc,
			}

			limiterInstance := &Limiter{
				repository:          repo,
				maxRequestsIP:       tc.maxRequestsIP,
				maxRequestsAPIToken: tc.maxRequestsAPIToken,
				freezeTime:          tc.freezeTime,
			}

			ctx := context.Background()
			var resp AllowResponse
			var err error

			resp, err = limiterInstance.Allow(ctx, tc.ip, tc.apiToken)

			if tc.expectError {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.Equal(t, tc.expectedAllowed, resp.Allowed, "Allowed value mismatch")
				assert.Equal(t, tc.expectedRemaining, resp.Remaining, "Remaining value mismatch")
			}
		})
	}
}
