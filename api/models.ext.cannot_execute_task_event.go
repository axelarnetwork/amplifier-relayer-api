package api

import "fmt"

// GetFees returns Fees from the event
// If Cost is Fees, then Fees is returned.
// If Cost is Token, then Fees containing a single item is returned where ID equals ID of the event,
// and .Meta.TxID equals .Meta.TxID of the event (if present).
func (e *CannotExecuteTaskEvent) GetFees() (Fees, error) {
	return createFees(e)
}

func (e *CannotExecuteTaskEvent) getEventID() string {
	return e.EventID
}

func (e *CannotExecuteTaskEvent) getTxID() *string {
	if e.Meta != nil && e.Meta.TxID != nil {
		return e.Meta.TxID
	}
	return nil
}

func (e *CannotExecuteTaskEvent) getCost() Cost {
	if e.Cost != nil {
		return *e.Cost
	}

	var cost Cost
	if err := cost.FromFees(Fees{}); err != nil {
		panic(fmt.Errorf("failed to create Cost: %w", err))
	}

	return cost
}
