package api

import (
	"errors"
	"fmt"
)

// CostFromToken creates a Cost from a given Token
func CostFromToken(token Token) Cost {
	var cost Cost
	if err := cost.FromToken(token); err != nil {
		panic(fmt.Errorf("failed to create cost from token: %w", err))
	}

	return cost
}

// Validate returns error if Cost isn't valid
//
//goland:noinspection GoMixedReceiverTypes
func (f *Cost) Validate() error {
	_, asTokenErr := f.AsToken()
	if asTokenErr == nil {
		return nil
	}

	fees, asFeesErr := f.AsFees()
	if asFeesErr != nil {
		return fmt.Errorf("cost is neither Fees nor Token: %w", errors.Join(asTokenErr, asFeesErr))
	}

	if len(fees) == 0 {
		return nil
	}

	ids := make(map[string]struct{}, len(fees))
	for _, fee := range fees {
		if _, exists := ids[fee.ID]; exists {
			return fmt.Errorf("duplicate fee ID: %s", fee.ID)
		}
		ids[fee.ID] = struct{}{}
	}

	return nil
}
