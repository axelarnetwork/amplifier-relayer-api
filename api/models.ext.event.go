package api

import (
	"encoding/json"
	"errors"

	"github.com/axelarnetwork/amplifier-relayer-api/internal/funcs"
)

// EventID returns id of the underlying event.
// The codegen library doesn't provide a way to access this field out of the box.
//
//goland:noinspection GoMixedReceiverTypes
func (e *Event) EventID() string {
	var obj struct {
		EventID string `json:"eventID"`
	}
	funcs.MustNoErr(
		json.Unmarshal(e.union, &obj),
	)
	return obj.EventID
}

// Validate returns error if Event isn't valid
//
//goland:noinspection GoMixedReceiverTypes
func (e *Event) Validate() error {
	switch e.Type {
	case EventTypeCannotExecuteTask:
		return e.validateFees(false)
	case
		EventTypeGasRefunded,
		EventTypeMessageApproved,
		EventTypeMessageExecuted,
		EventTypeMessageExecutedV2:
		return e.validateFees(true)
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (e *Event) validateFees(mandatoryCost bool) error {
	var obj struct {
		Cost *Cost `json:"cost"`
	}

	if err := json.Unmarshal(e.union, &obj); err != nil {
		return err
	}

	if obj.Cost == nil {
		if mandatoryCost {
			return errors.New("cost is required")
		}

		return nil
	}

	return obj.Cost.Validate()
}
