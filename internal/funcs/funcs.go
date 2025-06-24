package funcs

import (
	"fmt"
)

// Must returns the result if err is nil, panics otherwise
func Must[T any](result T, err error) T {
	if err != nil {
		panic(fmt.Errorf("call should not have failed: %w", err))
	}

	return result
}

// MustNoErr panics if err is not nil
func MustNoErr(err error) {
	if err != nil {
		panic(fmt.Errorf("call should not have failed: %w", err))
	}
}
