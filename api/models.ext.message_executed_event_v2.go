package api

// GetEventID returns the MessageExecutedEventV2.EventID of the event
func (e *MessageExecutedEventV2) GetEventID() string {
	return e.EventID
}

// GetCrossChainID returns the MessageExecutedEventV2.CrossChainID of the event
func (e *MessageExecutedEventV2) GetCrossChainID() CrossChainID {
	return e.CrossChainID
}

// GetStatus returns MessageExecutionStatusSuccessful
func (e *MessageExecutedEventV2) GetStatus() MessageExecutionStatus {
	return MessageExecutionStatusSuccessful
}

// GetMeta returns the MessageExecutedEventV2.Meta of the event
func (e *MessageExecutedEventV2) GetMeta() *MessageExecutedEventMetadata {
	return e.Meta
}

// GetFees returns Fees from the event
// If Cost is Fees, then Fees is returned.
// If Cost is Token, then Fees containing a single item is returned where ID equals ID of the event,
// and .Meta.TxID equals .Meta.TxID of the event (if present).
func (e *MessageExecutedEventV2) GetFees() (Fees, error) {
	return createFees(e)
}

func (e *MessageExecutedEventV2) getEventID() string {
	return e.EventID
}

func (e *MessageExecutedEventV2) getTxID() *string {
	if e.Meta != nil && e.Meta.TxID != nil {
		return e.Meta.TxID
	}
	return nil
}

func (e *MessageExecutedEventV2) getCost() Cost {
	return e.Cost
}
