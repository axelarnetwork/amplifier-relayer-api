package api

// GetFees returns Fees from the event
// If Cost is Fees, then Fees is returned.
// If Cost is Token, then Fees containing a single item is returned where ID equals ID of the event,
// and .Meta.TxID equals .Meta.TxID of the event (if present).
func (e *MessageApprovedEvent) GetFees() (Fees, error) {
	return createFees(e)
}

func (e *MessageApprovedEvent) getEventID() string {
	return e.EventID
}

func (e *MessageApprovedEvent) getTxID() *string {
	if e.Meta != nil && e.Meta.TxID != nil {
		return e.Meta.TxID
	}
	return nil
}

func (e *MessageApprovedEvent) getCost() Cost {
	return e.Cost
}
