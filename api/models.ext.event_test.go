//go:build unit_test

package api_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axelarnetwork/amplifier-relayer-api/api"
	"github.com/axelarnetwork/amplifier-relayer-api/internal/funcs"
)

func TestEvent_EventID(t *testing.T) {
	eventID := uuid.NewString()
	eventBase := api.EventBase{
		EventID: eventID,
	}

	eventBaseJSON := funcs.Must(json.Marshal(eventBase))

	var event api.Event
	funcs.MustNoErr(
		json.Unmarshal(eventBaseJSON, &event),
	)

	result := event.EventID()
	assert.Equal(t, eventID, result)
}

func TestEvent_Validate_WhenValid(t *testing.T) {
	var validTokenCost, validFeesCost api.Cost

	funcs.MustNoErr(
		validTokenCost.FromToken(api.Token{
			Amount: "123",
		}),
	)

	funcs.MustNoErr(
		validFeesCost.FromFees([]api.Fee{
			{
				ID: "fee1",
				Token: api.Token{
					Amount: "123",
				},
			},
			{
				ID: "fee2",
				Token: api.Token{
					Amount: "456",
				},
			},
		}),
	)

	type testCase struct {
		Description string        `json:"-"`
		Type        api.EventType `json:"type"`
		Cost        *api.Cost     `json:"cost"`
	}

	eventWithoutCost := func(t api.EventType) testCase {
		return testCase{
			Description: "no cost",
			Type:        t,
		}
	}

	eventWithTokenCost := func(t api.EventType) testCase {
		return testCase{
			Description: "cost is token",
			Type:        t,
			Cost:        &validTokenCost,
		}
	}

	eventWithFeesCost := func(t api.EventType) testCase {
		return testCase{
			Description: "cost is fees",
			Type:        t,
			Cost:        &validFeesCost,
		}
	}

	validEventsWithCost := []testCase{
		eventWithoutCost(api.EventTypeCannotExecuteTask),
		eventWithTokenCost(api.EventTypeCannotExecuteTask),
		eventWithFeesCost(api.EventTypeCannotExecuteTask),
		eventWithTokenCost(api.EventTypeGasRefunded),
		eventWithFeesCost(api.EventTypeGasRefunded),
		eventWithTokenCost(api.EventTypeMessageApproved),
		eventWithFeesCost(api.EventTypeMessageApproved),
		eventWithTokenCost(api.EventTypeMessageExecuted),
		eventWithFeesCost(api.EventTypeMessageExecuted),
		eventWithTokenCost(api.EventTypeMessageExecutedV2),
		eventWithFeesCost(api.EventTypeMessageExecutedV2),
	}

	validEventsWithoutCost := []testCase{
		eventWithoutCost(api.EventTypeCall),
		eventWithoutCost(api.EventTypeCannotExecuteMessage),
		eventWithoutCost(api.EventTypeCannotExecuteMessageV2),
		eventWithoutCost(api.EventTypeSignersRotated),
		eventWithoutCost(api.EventTypeAppInterchainTransferReceived),
		eventWithoutCost(api.EventTypeAppInterchainTransferSent),
		eventWithoutCost(api.EventTypeCannotRouteMessage),
		eventWithoutCost(api.EventTypeGasCredit),
		eventWithoutCost(api.EventTypeITSInterchainTransfer),
		eventWithoutCost(api.EventTypeITSLinkTokenStarted),
		eventWithoutCost(api.EventTypeITSInterchainTokenDeploymentStarted),
		eventWithoutCost(api.EventTypeITSInterchainTokenDeploymentStarted),
	}
	validEvents := slices.Concat(validEventsWithCost, validEventsWithoutCost)

	for _, tc := range validEvents {
		t.Run(fmt.Sprintf("when %s and %s", tc.Type, tc.Description), func(t *testing.T) {
			eventJSON := funcs.Must(json.Marshal(tc))

			var event api.Event
			funcs.MustNoErr(
				json.Unmarshal(eventJSON, &event),
			)

			result := event.Validate()

			require.NoError(t, result)
		})
	}
}

func TestEvent_Validate_WhenInvalid(t *testing.T) {
	var invalidFeesCost api.Cost

	funcs.MustNoErr(
		invalidFeesCost.FromFees([]api.Fee{
			{
				ID: "same_fee",
				Token: api.Token{
					Amount: "123",
				},
			},
			{
				ID: "same_fee",
				Token: api.Token{
					Amount: "456",
				},
			},
		}),
	)

	type testCase struct {
		Type api.EventType `json:"type"`
		Cost *api.Cost     `json:"cost"`
	}

	eventWithInvalidFeesCost := func(t api.EventType) testCase {
		return testCase{
			Type: t,
			Cost: &invalidFeesCost,
		}
	}

	invalidEvents := []testCase{
		eventWithInvalidFeesCost(api.EventTypeCannotExecuteTask),
		eventWithInvalidFeesCost(api.EventTypeGasRefunded),
		eventWithInvalidFeesCost(api.EventTypeMessageApproved),
		eventWithInvalidFeesCost(api.EventTypeMessageExecuted),
		eventWithInvalidFeesCost(api.EventTypeMessageExecutedV2),
	}

	for _, e := range invalidEvents {
		t.Run(fmt.Sprintf("when %s", e.Type), func(t *testing.T) {

			eventJSON := funcs.Must(json.Marshal(e))

			var event api.Event
			funcs.MustNoErr(
				json.Unmarshal(eventJSON, &event),
			)

			result := event.Validate()

			assert.Error(t, result)
		})
	}
}
