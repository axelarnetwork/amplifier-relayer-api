package api

// GetEventID returns the MessageExecutedEvent.EventID of the event
func (e *MessageExecutedEvent) GetEventID() string {
	return e.EventID
}

// GetCrossChainID returns CrossChainID constructed with MessageExecutedEvent.SourceChain and MessageExecutedEvent.MessageID of the event
func (e *MessageExecutedEvent) GetCrossChainID() CrossChainID {
	return CrossChainID{
		MessageID:   e.MessageID,
		SourceChain: e.SourceChain,
	}
}

// GetStatus returns the MessageExecutedEvent.Status of the event
func (e *MessageExecutedEvent) GetStatus() MessageExecutionStatus {
	return e.Status
}

// GetMeta returns the MessageExecutedEvent.Meta of the event
func (e *MessageExecutedEvent) GetMeta() *MessageExecutedEventMetadata {
	return e.Meta
}

// GetFees returns Fees from the event
// If Cost is Fees, then Fees is returned.
// If Cost is Token, then Fees containing a single item is returned where ID equals ID of the event,
// and .Meta.TxID equals .Meta.TxID of the event (if present).
func (e *MessageExecutedEvent) GetFees() (Fees, error) {
	return createFees(e)
}

func (e *MessageExecutedEvent) getEventID() string {
	return e.EventID
}

func (e *MessageExecutedEvent) getTxID() *string {
	if e.Meta != nil && e.Meta.TxID != nil {
		return e.Meta.TxID
	}
	return nil
}

func (e *MessageExecutedEvent) getCost() Cost {
	return e.Cost
}
