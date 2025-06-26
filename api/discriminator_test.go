package api_test

import (
	"encoding/json"
	"testing"

	"github.com/axelarnetwork/amplifier-relayer-api/api"
)

func TestTaskDiscriminatorValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupTask     func() api.Task
		asMethod      func(api.Task) (interface{}, error)
		expectedError bool
		expectedType  string
	}{
		{
			name: "GatewayTransactionTask - correct type",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromGatewayTransactionTask(api.GatewayTransactionTask{
					ExecuteData: []byte("test data"),
					Type:        "GATEWAY_TX",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			asMethod: func(task api.Task) (interface{}, error) {
				return task.AsGatewayTransactionTask()
			},
			expectedError: false,
			expectedType:  "GatewayTransactionTask",
		},
		{
			name: "GatewayTransactionTask - wrong type",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromExecuteTask(api.ExecuteTask{
					AvailableGasBalance: api.Token{Amount: "1000", TokenID: nil},
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
					Payload: []byte("test"),
					Type:    "EXECUTE",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			asMethod: func(task api.Task) (interface{}, error) {
				return task.AsGatewayTransactionTask()
			},
			expectedError: true,
			expectedType:  "GatewayTransactionTask",
		},
		{
			name: "ExecuteTask - correct type",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromExecuteTask(api.ExecuteTask{
					AvailableGasBalance: api.Token{Amount: "1000", TokenID: nil},
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
					Payload: []byte("test"),
					Type:    "EXECUTE",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			asMethod: func(task api.Task) (interface{}, error) {
				return task.AsExecuteTask()
			},
			expectedError: false,
			expectedType:  "ExecuteTask",
		},
		{
			name: "ExecuteTask - wrong type",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromGatewayTransactionTask(api.GatewayTransactionTask{
					ExecuteData: []byte("test data"),
					Type:        "GATEWAY_TX",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			asMethod: func(task api.Task) (interface{}, error) {
				return task.AsExecuteTask()
			},
			expectedError: true,
			expectedType:  "ExecuteTask",
		},
		{
			name: "ConstructProofTask - correct type",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromConstructProofTask(api.ConstructProofTask{
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
					Payload: []byte("test"),
					Type:    "CONSTRUCT_PROOF",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			asMethod: func(task api.Task) (interface{}, error) {
				return task.AsConstructProofTask()
			},
			expectedError: false,
			expectedType:  "ConstructProofTask",
		},
		{
			name: "VerifyTask - correct type",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromVerifyTask(api.VerifyTask{
					DestinationChain: "test",
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
					Payload: []byte("test"),
					Type:    "VERIFY",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			asMethod: func(task api.Task) (interface{}, error) {
				return task.AsVerifyTask()
			},
			expectedError: false,
			expectedType:  "VerifyTask",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setupTask()
			_, err := tt.asMethod(task)

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
			discriminator, err := task.Discriminator()
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

func TestTaskValueByDiscriminator(t *testing.T) {
	tests := []struct {
		name          string
		setupTask     func() api.Task
		expectedType  string
		expectedError bool
	}{
		{
			name: "GATEWAY_TX discriminator",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromGatewayTransactionTask(api.GatewayTransactionTask{
					ExecuteData: []byte("test data"),
					Type:        "GATEWAY_TX",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			expectedType:  "GatewayTransactionTask",
			expectedError: false,
		},
		{
			name: "EXECUTE discriminator",
			setupTask: func() api.Task {
				task := api.Task{}
				err := task.FromExecuteTask(api.ExecuteTask{
					AvailableGasBalance: api.Token{Amount: "1000", TokenID: nil},
					Message: api.Message{
						MessageID:          "test",
						SourceChain:        "test",
						SourceAddress:      "test",
						DestinationAddress: "test",
						PayloadHash:        []byte("test"),
					},
					Payload: []byte("test"),
					Type:    "EXECUTE",
				})
				if err != nil {
					t.Fatalf("Failed to setup task: %v", err)
				}
				return task
			},
			expectedType:  "ExecuteTask",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setupTask()
			result, err := task.ValueByDiscriminator()

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

func TestEventValueByDiscriminator(t *testing.T) {
	tests := []struct {
		name          string
		setupEvent    func() api.Event
		expectedType  string
		expectedError bool
	}{
		{
			name: "GAS_CREDIT discriminator",
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
			expectedType:  "GasCreditEvent",
			expectedError: false,
		},
		{
			name: "GAS_REFUNDED discriminator",
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
			expectedType:  "GasRefundedEvent",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.setupEvent()
			result, err := event.ValueByDiscriminator()

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

func TestTaskMergeFunctionality(t *testing.T) {
	t.Run("Merge GatewayTransactionTask", func(t *testing.T) {
		task := api.Task{}

		// First, set up a task
		err := task.FromGatewayTransactionTask(api.GatewayTransactionTask{
			ExecuteData: []byte("initial data"),
			Type:        "GATEWAY_TX",
		})
		if err != nil {
			t.Fatalf("Failed to setup initial task: %v", err)
		}

		// Now merge with new data
		err = task.MergeGatewayTransactionTask(api.GatewayTransactionTask{
			ExecuteData: []byte("merged data"),
			Type:        "GATEWAY_TX",
		})
		if err != nil {
			t.Errorf("Failed to merge task: %v", err)
		}

		// Verify the task can still be cast correctly
		result, err := task.AsGatewayTransactionTask()
		if err != nil {
			t.Errorf("Failed to cast merged task: %v", err)
		} else {
			t.Logf("✅ Successfully cast merged task: %+v", result)
		}
	})
}

func TestEventMergeFunctionality(t *testing.T) {
	t.Run("Merge GasCreditEvent", func(t *testing.T) {
		event := api.Event{}

		// First, set up an event
		err := event.FromGasCreditEvent(api.GasCreditEvent{
			EventID:       "initial-event",
			MessageID:     "initial-message",
			Payment:       api.Token{Amount: "1000", TokenID: nil},
			RefundAddress: "initial-address",
		})
		if err != nil {
			t.Fatalf("Failed to setup initial event: %v", err)
		}

		// Now merge with new data
		err = event.MergeGasCreditEvent(api.GasCreditEvent{
			EventID:       "merged-event",
			MessageID:     "merged-message",
			Payment:       api.Token{Amount: "2000", TokenID: nil},
			RefundAddress: "merged-address",
		})
		if err != nil {
			t.Errorf("Failed to merge event: %v", err)
		}

		// Verify the event can still be cast correctly
		result, err := event.AsGasCreditEvent()
		if err != nil {
			t.Errorf("Failed to cast merged event: %v", err)
		} else {
			t.Logf("✅ Successfully cast merged event: %+v", result)
		}
	})
}

func TestJSONMarshaling(t *testing.T) {
	t.Run("Task JSON marshaling", func(t *testing.T) {
		task := api.Task{}
		err := task.FromGatewayTransactionTask(api.GatewayTransactionTask{
			ExecuteData: []byte("test data"),
			Type:        "GATEWAY_TX",
		})
		if err != nil {
			t.Fatalf("Failed to setup task: %v", err)
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(task)
		if err != nil {
			t.Errorf("Failed to marshal task: %v", err)
		} else {
			t.Logf("✅ Task JSON: %s", string(jsonData))
		}

		// Unmarshal back
		var newTask api.Task
		err = json.Unmarshal(jsonData, &newTask)
		if err != nil {
			t.Errorf("Failed to unmarshal task: %v", err)
		} else {
			t.Logf("✅ Successfully unmarshaled task")
		}

		// Verify discriminator still works
		discriminator, err := newTask.Discriminator()
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
	t.Run("Empty Task discriminator", func(t *testing.T) {
		task := api.Task{}
		discriminator, err := task.Discriminator()
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
