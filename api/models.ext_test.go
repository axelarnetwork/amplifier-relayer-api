package api_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axelarnetwork/amplifier-relayer-api/v2/api"
	"github.com/axelarnetwork/amplifier-relayer-api/v2/internal/funcs"
)

func TestTokenManagerType(t *testing.T) {
	type testCase struct {
		chainTokenManagerValue uint8
		apiTokenManagerValue   api.TokenManagerType
		err                    error
	}

	testCases := []testCase{
		{0, api.TokenManagerNativeInterchainToken, nil},
		{1, api.TokenManagerMintBurnFrom, nil},
		{2, api.TokenManagerLockUnlock, nil},
		{3, api.TokenManagerLockUnlockFee, nil},
		{4, api.TokenManagerMintBurn, nil},
		{5, "", errors.New("invalid TokenManagerType: 5")},
	}

	for _, tc := range testCases {
		tmt, err := api.TokenManagerTypeFromSolidityEnum(tc.chainTokenManagerValue)
		if tc.err != nil {
			assert.Error(t, err)
			assert.Equal(t, tc.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.apiTokenManagerValue, tmt)
		}
	}
}

func TestGetFees(t *testing.T) {
	t.Run("when Cost is nil", func(t *testing.T) {
		events := []interface {
			GetFees() (api.Fees, error)
		}{
			&api.CannotExecuteTaskEvent{
				Cost: nil,
			},
		}

		for _, event := range events {
			t.Run(fmt.Sprintf("when %T", event), func(t *testing.T) {
				fees, err := event.GetFees()
				assert.NoError(t, err)
				assert.Empty(t, fees)
			})
		}
	})

	t.Run("when Cost is Token", func(t *testing.T) {
		var (
			cost    api.Cost
			eventID = uuid.NewString()
		)
		funcs.MustNoErr(
			cost.FromToken(api.Token{
				Amount: "123",
			}),
		)

		events := []interface {
			GetFees() (api.Fees, error)
		}{
			&api.CannotExecuteTaskEvent{
				EventID: eventID,
				Cost:    &cost,
			},
			&api.MessageApprovedEvent{
				EventID: eventID,
				Cost:    cost,
			},
			&api.MessageExecutedEvent{
				EventID: eventID,
				Cost:    cost,
			},
			&api.MessageExecutedEventV2{
				EventID: eventID,
				Cost:    cost,
			},
			&api.GasRefundedEvent{
				EventID: eventID,
				Cost:    cost,
			},
		}

		for _, event := range events {
			t.Run(fmt.Sprintf("when %T", event), func(t *testing.T) {
				fees, err := event.GetFees()
				require.NoError(t, err)
				require.Len(t, fees, 1)
				assert.EqualValues(t, api.Fee{
					Token: api.Token{
						Amount: "123",
					},
					ID: eventID,
				}, fees[0])
			})
		}
	})

	t.Run("when Cost is Fees", func(t *testing.T) {
		description1 := uuid.NewString()
		fee1 := api.Fee{
			Token: api.Token{
				Amount: "123",
			},
			ID:          uuid.NewString(),
			Description: &description1,
		}

		txID2 := uuid.NewString()
		fee2 := api.Fee{
			Token: api.Token{
				Amount: "456",
			},
			ID: uuid.NewString(),
			Meta: &api.FeeMetadata{
				TxID: &txID2,
			},
		}
		fee3 := api.Fee{
			Token: api.Token{
				Amount: "789",
			},
			ID: uuid.NewString(),
		}

		var cost api.Cost
		funcs.MustNoErr(
			cost.FromFees(api.Fees{fee1, fee2, fee3}),
		)

		events := []interface {
			GetFees() (api.Fees, error)
		}{
			&api.CannotExecuteTaskEvent{
				Cost: &cost,
			},
			&api.MessageApprovedEvent{
				Cost: cost,
			},
			&api.MessageExecutedEvent{
				Cost: cost,
			},
			&api.MessageExecutedEventV2{
				Cost: cost,
			},
			&api.GasRefundedEvent{
				Cost: cost,
			},
		}

		for _, event := range events {
			t.Run(fmt.Sprintf("when %T", event), func(t *testing.T) {
				fees, err := event.GetFees()
				require.NoError(t, err)
				require.Len(t, fees, 3)
				assert.EqualValues(t, api.Fees{fee1, fee2, fee3}, fees)
			})
		}
	})
}

func TestTaskEnvelopeHelpers(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	taskID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("TaskEnvelope helper methods", func(t *testing.T) {
		// Create a TaskItem with an ExecuteTaskItem
		var taskItem api.TaskItem
		executeTaskItem := api.ExecuteTaskItem{
			Task: api.ExecuteTask{
				AvailableGasBalance: api.Token{Amount: "1000"},
				Message: api.Message{
					MessageID:          "msg-123",
					SourceChain:        "ethereum",
					SourceAddress:      "0x123",
					DestinationAddress: "0x456",
					PayloadHash:        []byte("hash"),
				},
				Payload: []byte("payload"),
			},
		}

		err := taskItem.FromExecuteTaskItem(executeTaskItem)
		require.NoError(t, err)

		taskEnvelope := api.TaskEnvelope{
			Chain:     "ethereum",
			ID:        taskID,
			Task:      taskItem,
			Timestamp: fixedTime,
		}

		// Test helper methods
		require.Equal(t, taskID, taskEnvelope.GetTaskItemID())
		require.Equal(t, "ethereum", taskEnvelope.GetChain())
		require.Equal(t, fixedTime, taskEnvelope.GetTimestamp())
		require.Equal(t, api.TaskTypeExecute, taskEnvelope.GetTaskType())
	})

	t.Run("TaskEnvelope with GatewayTransactionTaskItem", func(t *testing.T) {
		// Create a TaskItem with a GatewayTransactionTaskItem
		var taskItem api.TaskItem
		gatewayTaskItem := api.GatewayTransactionTaskItem{
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte("execute data"),
			},
			Meta: &api.DestinationChainTaskMetadata{
				ScopedMessages: &[]api.CrossChainID{
					{
						MessageID:   "msg-123",
						SourceChain: "ethereum",
					},
				},
			},
		}

		err := taskItem.FromGatewayTransactionTaskItem(gatewayTaskItem)
		require.NoError(t, err)

		taskEnvelope := api.TaskEnvelope{
			Chain:     "polygon",
			ID:        taskID,
			Task:      taskItem,
			Timestamp: fixedTime,
		}

		// Test helper methods
		require.Equal(t, taskID, taskEnvelope.GetTaskItemID())
		require.Equal(t, "polygon", taskEnvelope.GetChain())
		require.Equal(t, fixedTime, taskEnvelope.GetTimestamp())
		require.Equal(t, api.TaskTypeGatewayTransaction, taskEnvelope.GetTaskType())
	})
}
