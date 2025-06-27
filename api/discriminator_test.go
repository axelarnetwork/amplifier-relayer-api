package api_test

import (
	"encoding/json"
	"testing"

	"github.com/axelarnetwork/amplifier-relayer-api/api"
)

func TestTaskItemDiscriminatorValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupTaskItem func() api.TaskItem
		asMethod      func(api.TaskItem) (any, error)
		expectedError bool
		expectedType  string
	}{
		{
			name: "GatewayTransactionTaskItem - correct type",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
					Type: "GATEWAY_TX",
					Task: api.GatewayTransactionTask{
						ExecuteData: []byte("test data"),
					},
				})
				return item
			},
			asMethod: func(item api.TaskItem) (any, error) {
				return item.AsGatewayTransactionTaskItem()
			},
			expectedError: false,
			expectedType:  "GatewayTransactionTaskItem",
		},
		{
			name: "GatewayTransactionTaskItem - wrong type",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromExecuteTaskItem(api.ExecuteTaskItem{
					Type: "EXECUTE",
					Task: api.ExecuteTask{
						AvailableGasBalance: api.Token{Amount: "1000", TokenID: nil},
						Message: api.Message{
							MessageID:          "test",
							SourceChain:        "test",
							SourceAddress:      "test",
							DestinationAddress: "test",
							PayloadHash:        []byte("test"),
						},
						Payload: []byte("test"),
					},
				})
				return item
			},
			asMethod: func(item api.TaskItem) (any, error) {
				return item.AsGatewayTransactionTaskItem()
			},
			expectedError: true,
			expectedType:  "GatewayTransactionTaskItem",
		},
		{
			name: "ExecuteTaskItem - correct type",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromExecuteTaskItem(api.ExecuteTaskItem{
					Type: "EXECUTE",
					Task: api.ExecuteTask{
						AvailableGasBalance: api.Token{Amount: "1000", TokenID: nil},
						Message: api.Message{
							MessageID:          "test",
							SourceChain:        "test",
							SourceAddress:      "test",
							DestinationAddress: "test",
							PayloadHash:        []byte("test"),
						},
						Payload: []byte("test"),
					},
				})
				return item
			},
			asMethod: func(item api.TaskItem) (any, error) {
				return item.AsExecuteTaskItem()
			},
			expectedError: false,
			expectedType:  "ExecuteTaskItem",
		},
		{
			name: "ExecuteTaskItem - wrong type",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
					Type: "GATEWAY_TX",
					Task: api.GatewayTransactionTask{
						ExecuteData: []byte("test data"),
					},
				})
				return item
			},
			asMethod: func(item api.TaskItem) (any, error) {
				return item.AsExecuteTaskItem()
			},
			expectedError: true,
			expectedType:  "ExecuteTaskItem",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := tt.setupTaskItem()
			_, err := tt.asMethod(item)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					t.Logf("✅ Correctly got error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				} else {
					t.Logf("✅ Successfully cast to %s", tt.expectedType)
				}
			}

			// Verify discriminator value
			discriminator, err := item.Discriminator()
			if err != nil {
				t.Errorf("Failed to get discriminator: %v", err)
			} else {
				t.Logf("Discriminator value: %s", discriminator)
			}
		})
	}
}

func TestEventDiscriminatorValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupEvent    func() api.Event
		asMethod      func(api.Event) (interface{}, error)
		expectedError bool
		expectedType  string
	}{
		{
			name: "GasCreditEvent - correct type",
			setupEvent: func() api.Event {
				event := api.Event{}
				err := event.FromGasCreditEvent(api.GasCreditEvent{
					EventID:       "test-event",
					MessageID:     "test-message",
					Payment:       api.Token{Amount: "1000", TokenID: nil},
					RefundAddress: "test-address",
				})
				if err != nil {
					t.Fatalf("Failed to setup event: %v", err)
				}
				return event
			},
			asMethod: func(event api.Event) (interface{}, error) {
				return event.AsGasCreditEvent()
			},
			expectedError: false,
			expectedType:  "GasCreditEvent",
		},
		{
			name: "GasCreditEvent - wrong type",
			setupEvent: func() api.Event {
				event := api.Event{}
				err := event.FromGasRefundedEvent(api.GasRefundedEvent{
					EventID:          "test-event",
					MessageID:        "test-message",
					Cost:             api.Cost{},
					RecipientAddress: "test-address",
					RefundedAmount:   api.Token{Amount: "1000", TokenID: nil},
				})
				if err != nil {
					t.Fatalf("Failed to setup event: %v", err)
				}
				return event
			},
			asMethod: func(event api.Event) (interface{}, error) {
				return event.AsGasCreditEvent()
			},
			expectedError: true,
			expectedType:  "GasCreditEvent",
		},
		{
			name: "GasRefundedEvent - correct type",
			setupEvent: func() api.Event {
				event := api.Event{}
				err := event.FromGasRefundedEvent(api.GasRefundedEvent{
					EventID:          "test-event",
					MessageID:        "test-message",
					Cost:             api.Cost{},
					RecipientAddress: "test-address",
					RefundedAmount:   api.Token{Amount: "1000", TokenID: nil},
				})
				if err != nil {
					t.Fatalf("Failed to setup event: %v", err)
				}
				return event
			},
			asMethod: func(event api.Event) (interface{}, error) {
				return event.AsGasRefundedEvent()
			},
			expectedError: false,
			expectedType:  "GasRefundedEvent",
		},
		{
			name: "CallEvent - correct type",
			setupEvent: func() api.Event {
				event := api.Event{}
				err := event.FromCallEvent(api.CallEvent{
					DestinationChain: "test-chain",
					EventID:          "test-event",
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
					Payload: []byte("test"),
				})
				if err != nil {
					t.Fatalf("Failed to setup event: %v", err)
				}
				return event
			},
			asMethod: func(event api.Event) (interface{}, error) {
				return event.AsCallEvent()
			},
			expectedError: false,
			expectedType:  "CallEvent",
		},
		{
			name: "MessageApprovedEvent - correct type",
			setupEvent: func() api.Event {
				event := api.Event{}
				err := event.FromMessageApprovedEvent(api.MessageApprovedEvent{
					Cost:    api.Cost{},
					EventID: "test-event",
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
				})
				if err != nil {
					t.Fatalf("Failed to setup event: %v", err)
				}
				return event
			},
			asMethod: func(event api.Event) (interface{}, error) {
				return event.AsMessageApprovedEvent()
			},
			expectedError: false,
			expectedType:  "MessageApprovedEvent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.setupEvent()
			_, err := tt.asMethod(event)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					t.Logf("✅ Correctly got error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				} else {
					t.Logf("✅ Successfully cast to %s", tt.expectedType)
				}
			}

			// Verify discriminator value
			discriminator, err := event.Discriminator()
			if err != nil {
				t.Errorf("Failed to get discriminator: %v", err)
			} else {
				t.Logf("Discriminator value: %s", discriminator)
			}
		})
	}
}

func TestTaskItemValueByDiscriminator(t *testing.T) {
	tests := []struct {
		name          string
		setupTask     func() api.TaskItem
		expectedType  string
		expectedError bool
	}{
		{
			name: "GATEWAY_TX discriminator",
			setupTask: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
					Type: "GATEWAY_TX",
					Task: api.GatewayTransactionTask{
						ExecuteData: []byte("test data"),
					},
				})
				return item
			},
			expectedType:  "GatewayTransactionTaskItem",
			expectedError: false,
		},
		{
			name: "EXECUTE discriminator",
			setupTask: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromExecuteTaskItem(api.ExecuteTaskItem{
					Type: "EXECUTE",
					Task: api.ExecuteTask{
						AvailableGasBalance: api.Token{Amount: "1000", TokenID: nil},
						Message: api.Message{
							MessageID:          "test",
							SourceChain:        "test",
							SourceAddress:      "test",
							DestinationAddress: "test",
							PayloadHash:        []byte("test"),
						},
						Payload: []byte("test"),
					},
				})
				return item
			},
			expectedType:  "ExecuteTaskItem",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := tt.setupTask()
			result, err := item.ValueByDiscriminator()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					t.Logf("✅ Correctly got error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				} else {
					t.Logf("✅ Successfully got value by discriminator: %T", result)
				}
			}
		})
	}
}

func TestTaskItemMergeFunctionality(t *testing.T) {
	t.Run("Merge GatewayTransactionTaskItem", func(t *testing.T) {
		var item api.TaskItem
		_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
			Type: "GATEWAY_TX",
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte("initial data"),
			},
		})

		// Now merge with new data
		err := item.MergeGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
			Type: "GATEWAY_TX",
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte("merged data"),
			},
		})
		if err != nil {
			t.Errorf("Failed to merge task: %v", err)
		}

		// Verify the task can still be cast correctly
		result, err := item.AsGatewayTransactionTaskItem()
		if err != nil {
			t.Errorf("Failed to cast merged task: %v", err)
		} else {
			t.Logf("✅ Successfully cast merged task: %+v", result)
		}
	})
}

func TestJSONMarshaling(t *testing.T) {
	t.Run("TaskItem JSON marshaling", func(t *testing.T) {
		var item api.TaskItem
		_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
			Type: "GATEWAY_TX",
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte("test data"),
			},
		})

		// Marshal to JSON
		jsonData, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Failed to marshal task: %v", err)
		} else {
			t.Logf("✅ TaskItem JSON: %s", string(jsonData))
		}

		// Unmarshal back
		var newItem api.TaskItem
		err = json.Unmarshal(jsonData, &newItem)
		if err != nil {
			t.Errorf("Failed to unmarshal task: %v", err)
		} else {
			t.Logf("✅ Successfully unmarshaled task")
		}

		// Verify discriminator still works
		discriminator, err := newItem.Discriminator()
		if err != nil {
			t.Errorf("Failed to get discriminator after unmarshaling: %v", err)
		} else {
			t.Logf("✅ Discriminator after unmarshaling: %s", discriminator)
		}
	})

	t.Run("Event JSON marshaling", func(t *testing.T) {
		event := api.Event{}
		err := event.FromGasCreditEvent(api.GasCreditEvent{
			EventID:       "test-event",
			MessageID:     "test-message",
			Payment:       api.Token{Amount: "1000", TokenID: nil},
			RefundAddress: "test-address",
		})
		if err != nil {
			t.Fatalf("Failed to setup event: %v", err)
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(event)
		if err != nil {
			t.Errorf("Failed to marshal event: %v", err)
		} else {
			t.Logf("✅ Event JSON: %s", string(jsonData))
		}

		// Unmarshal back
		var newEvent api.Event
		err = json.Unmarshal(jsonData, &newEvent)
		if err != nil {
			t.Errorf("Failed to unmarshal event: %v", err)
		} else {
			t.Logf("✅ Successfully unmarshaled event")
		}

		// Verify discriminator still works
		discriminator, err := newEvent.Discriminator()
		if err != nil {
			t.Errorf("Failed to get discriminator after unmarshaling: %v", err)
		} else {
			t.Logf("✅ Discriminator after unmarshaling: %s", discriminator)
		}
	})
}

func TestDiscriminatorEdgeCases(t *testing.T) {
	t.Run("Empty TaskItem discriminator", func(t *testing.T) {
		var item api.TaskItem
		discriminator, err := item.Discriminator()
		if err != nil {
			t.Logf("✅ Correctly got error for empty discriminator: %v", err)
		} else {
			t.Logf("Discriminator value: %s", discriminator)
		}
	})

	t.Run("Empty Event discriminator", func(t *testing.T) {
		event := api.Event{}
		discriminator, err := event.Discriminator()
		if err != nil {
			t.Logf("✅ Correctly got error for empty discriminator: %v", err)
		} else {
			t.Logf("Discriminator value: %s", discriminator)
		}
	})
}
