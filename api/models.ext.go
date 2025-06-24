package api

import (
	"errors"
	"fmt"
)

// GeneralizedMessageExecutedEvent is a compatibility interface that generalizes MessageExecutedEvent and MessageExecutedEventV2.
type GeneralizedMessageExecutedEvent interface {
	GetEventID() string
	GetCrossChainID() CrossChainID
	GetStatus() MessageExecutionStatus
	GetFees() (Fees, error)
	GetMeta() *MessageExecutedEventMetadata
}

type eventWithCost interface {
	getTxID() *string
	getEventID() string
	getCost() Cost
}

func createFees[T eventWithCost](event T) (Fees, error) {
	cost := event.getCost()

	token, asTokenErr := cost.AsToken()
	if asTokenErr == nil {
		fee := Fee{
			Token: token,
			ID:    event.getEventID(),
		}

		if txID := event.getTxID(); txID != nil {
			fee.Meta = &FeeMetadata{
				TxID: txID,
			}
		}

		return Fees{fee}, nil
	}

	fees, asFeesErr := cost.AsFees()
	if asFeesErr == nil {
		return fees, nil
	}

	return nil, fmt.Errorf(
		"failed to get fees: %w",
		errors.Join(asTokenErr, asFeesErr),
	)
}
