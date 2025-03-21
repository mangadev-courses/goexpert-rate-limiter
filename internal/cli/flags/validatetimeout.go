package flags

import (
	"fmt"
	"time"
)

var MaxTimeout = 5 * time.Minute

func ValidateTimeout(timeout time.Duration) error {
	if timeout > MaxTimeout {
		return fmt.Errorf("timeout value exceeds the maximum allowed value of %v", MaxTimeout)
	}
	return nil
}
