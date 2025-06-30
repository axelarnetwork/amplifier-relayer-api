package api_test

import (
	"encoding/json"
	"testing"

	"github.com/axelarnetwork/amplifier-relayer-api/v2/api"
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
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte("initial data"),
			},
		})

		// Now merge with new data
		err := item.MergeGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
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

func TestTaskItemJSONMarshalingCompatibility(t *testing.T) {
	tests := []struct {
		name           string
		setupTaskItem  func() api.TaskItem
		expectedFields []string
		validateJSON   func(t *testing.T, jsonData []byte)
	}{
		{
			name: "GatewayTransactionTaskItem JSON marshaling",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
					Task: api.GatewayTransactionTask{
						ExecuteData: []byte("test execute data"),
					},
					Meta: &api.DestinationChainTaskMetadata{
						ScopedMessages: &[]api.CrossChainID{
							{
								MessageID:   "msg-123",
								SourceChain: "ethereum",
							},
						},
					},
				})
				return item
			},
			expectedFields: []string{"type", "task", "meta"},
			validateJSON: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				err := json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
					return
				}

				// Verify required fields are present
				if result["type"] != "GATEWAY_TX" {
					t.Errorf("Expected type 'GATEWAY_TX', got %v", result["type"])
				}

				// Verify task structure
				task, ok := result["task"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected task to be an object")
				} else {
					if task["executeData"] == nil {
						t.Errorf("Expected executeData field in task")
					}
				}

				// Verify meta structure
				meta, ok := result["meta"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected meta to be an object")
				} else {
					scopedMessages, ok := meta["scopedMessages"].([]interface{})
					if !ok || len(scopedMessages) == 0 {
						t.Errorf("Expected scopedMessages array in meta")
					}
				}
			},
		},
		{
			name: "ExecuteTaskItem JSON marshaling",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromExecuteTaskItem(api.ExecuteTaskItem{
					Task: api.ExecuteTask{
						AvailableGasBalance: api.Token{
							Amount:  "1000000000000000000",
							TokenID: nil,
						},
						Message: api.Message{
							MessageID:          "msg-456",
							SourceChain:        "ethereum",
							SourceAddress:      "0x1234567890123456789012345678901234567890",
							DestinationAddress: "0x0987654321098765432109876543210987654321",
							PayloadHash:        []byte("payload hash"),
						},
						Payload: []byte("execute payload"),
					},
				})
				return item
			},
			expectedFields: []string{"type", "task"},
			validateJSON: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				err := json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
					return
				}

				// Verify required fields are present
				if result["type"] != "EXECUTE" {
					t.Errorf("Expected type 'EXECUTE', got %v", result["type"])
				}

				// Verify task structure
				task, ok := result["task"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected task to be an object")
				} else {
					// Verify availableGasBalance
					gasBalance, ok := task["availableGasBalance"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected availableGasBalance to be an object")
					} else {
						if gasBalance["amount"] != "1000000000000000000" {
							t.Errorf("Expected amount '1000000000000000000', got %v", gasBalance["amount"])
						}
					}

					// Verify message structure
					message, ok := task["message"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected message to be an object")
					} else {
						if message["messageID"] != "msg-456" {
							t.Errorf("Expected messageID 'msg-456', got %v", message["messageID"])
						}
						if message["sourceChain"] != "ethereum" {
							t.Errorf("Expected sourceChain 'ethereum', got %v", message["sourceChain"])
						}
					}

					// Verify payload
					if task["payload"] == nil {
						t.Errorf("Expected payload field in task")
					}
				}
			},
		},
		{
			name: "ConstructProofTaskItem JSON marshaling",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromConstructProofTaskItem(api.ConstructProofTaskItem{
					Task: api.ConstructProofTask{
						Message: api.Message{
							MessageID:          "msg-789",
							SourceChain:        "ethereum",
							SourceAddress:      "0x1234567890123456789012345678901234567890",
							DestinationAddress: "0x0987654321098765432109876543210987654321",
							PayloadHash:        []byte("proof payload hash"),
						},
						Payload: []byte("proof payload"),
					},
				})
				return item
			},
			expectedFields: []string{"type", "task"},
			validateJSON: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				err := json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
					return
				}

				// Verify required fields are present
				if result["type"] != "CONSTRUCT_PROOF" {
					t.Errorf("Expected type 'CONSTRUCT_PROOF', got %v", result["type"])
				}

				// Verify task structure
				task, ok := result["task"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected task to be an object")
				} else {
					// Verify message structure
					message, ok := task["message"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected message to be an object")
					} else {
						if message["messageID"] != "msg-789" {
							t.Errorf("Expected messageID 'msg-789', got %v", message["messageID"])
						}
					}

					// Verify payload
					if task["payload"] == nil {
						t.Errorf("Expected payload field in task")
					}
				}
			},
		},
		{
			name: "RefundTaskItem JSON marshaling",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromRefundTaskItem(api.RefundTaskItem{
					Task: api.RefundTask{
						Message: api.Message{
							MessageID:          "msg-refund",
							SourceChain:        "ethereum",
							SourceAddress:      "0x1234567890123456789012345678901234567890",
							DestinationAddress: "0x0987654321098765432109876543210987654321",
							PayloadHash:        []byte("refund payload hash"),
						},
						RefundRecipientAddress: "0x1111111111111111111111111111111111111111",
						RemainingGasBalance: api.Token{
							Amount:  "500000000000000000",
							TokenID: nil,
						},
					},
					Meta: &api.SourceChainTaskMetadata{
						SourceContext: &api.MessageContext{
							"key1": "value1",
							"key2": "value2",
						},
					},
				})
				return item
			},
			expectedFields: []string{"type", "task", "meta"},
			validateJSON: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				err := json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
					return
				}

				// Verify required fields are present
				if result["type"] != "REFUND" {
					t.Errorf("Expected type 'REFUND', got %v", result["type"])
				}

				// Verify task structure
				task, ok := result["task"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected task to be an object")
				} else {
					// Verify refundRecipientAddress
					if task["refundRecipientAddress"] != "0x1111111111111111111111111111111111111111" {
						t.Errorf("Expected refundRecipientAddress '0x1111111111111111111111111111111111111111', got %v", task["refundRecipientAddress"])
					}

					// Verify remainingGasBalance
					gasBalance, ok := task["remainingGasBalance"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected remainingGasBalance to be an object")
					} else {
						if gasBalance["amount"] != "500000000000000000" {
							t.Errorf("Expected amount '500000000000000000', got %v", gasBalance["amount"])
						}
					}
				}

				// Verify meta structure
				meta, ok := result["meta"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected meta to be an object")
				} else {
					sourceContext, ok := meta["sourceContext"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected sourceContext to be an object")
					} else {
						if sourceContext["key1"] != "value1" {
							t.Errorf("Expected key1 'value1', got %v", sourceContext["key1"])
						}
					}
				}
			},
		},
		{
			name: "VerifyTaskItem JSON marshaling",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromVerifyTaskItem(api.VerifyTaskItem{
					Task: api.VerifyTask{
						DestinationChain: "polygon",
						Message: api.Message{
							MessageID:          "msg-verify",
							SourceChain:        "ethereum",
							SourceAddress:      "0x1234567890123456789012345678901234567890",
							DestinationAddress: "0x0987654321098765432109876543210987654321",
							PayloadHash:        []byte("verify payload hash"),
						},
						Payload: []byte("verify payload"),
					},
					Meta: &api.SourceChainTaskMetadata{
						SourceContext: &api.MessageContext{
							"verifyKey": "verifyValue",
						},
					},
				})
				return item
			},
			expectedFields: []string{"type", "task", "meta"},
			validateJSON: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				err := json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
					return
				}

				// Verify required fields are present
				if result["type"] != "VERIFY" {
					t.Errorf("Expected type 'VERIFY', got %v", result["type"])
				}

				// Verify task structure
				task, ok := result["task"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected task to be an object")
				} else {
					// Verify destinationChain
					if task["destinationChain"] != "polygon" {
						t.Errorf("Expected destinationChain 'polygon', got %v", task["destinationChain"])
					}

					// Verify message structure
					message, ok := task["message"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected message to be an object")
					} else {
						if message["messageID"] != "msg-verify" {
							t.Errorf("Expected messageID 'msg-verify', got %v", message["messageID"])
						}
					}

					// Verify payload
					if task["payload"] == nil {
						t.Errorf("Expected payload field in task")
					}
				}
			},
		},
		{
			name: "ReactToWasmEventTaskItem JSON marshaling",
			setupTaskItem: func() api.TaskItem {
				var item api.TaskItem
				_ = item.FromReactToWasmEventTaskItem(api.ReactToWasmEventTaskItem{
					Task: api.ReactToWasmEventTask{
						Event: api.WasmEvent{
							Type: "wasm.event",
							Attributes: []api.WasmEventAttribute{
								{
									Key:   "key1",
									Value: "value1",
								},
								{
									Key:   "key2",
									Value: "value2",
								},
							},
						},
						Height: 12345,
					},
				})
				return item
			},
			expectedFields: []string{"type", "task"},
			validateJSON: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				err := json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
					return
				}

				// Verify required fields are present
				if result["type"] != "REACT_TO_WASM_EVENT" {
					t.Errorf("Expected type 'REACT_TO_WASM_EVENT', got %v", result["type"])
				}

				// Verify task structure
				task, ok := result["task"].(map[string]interface{})
				if !ok {
					t.Errorf("Expected task to be an object")
				} else {
					// Verify event structure
					event, ok := task["event"].(map[string]interface{})
					if !ok {
						t.Errorf("Expected event to be an object")
					} else {
						if event["type"] != "wasm.event" {
							t.Errorf("Expected event type 'wasm.event', got %v", event["type"])
						}

						attributes, ok := event["attributes"].([]interface{})
						if !ok || len(attributes) != 2 {
							t.Errorf("Expected 2 attributes, got %v", attributes)
						}
					}

					// Verify height
					if task["height"] != float64(12345) {
						t.Errorf("Expected height 12345, got %v", task["height"])
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := tt.setupTaskItem()

			// Marshal to JSON
			jsonData, err := json.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to marshal task item: %v", err)
			}

			t.Logf("✅ Generated JSON: %s", string(jsonData))

			// Validate JSON structure
			tt.validateJSON(t, jsonData)

			// Verify all expected fields are present
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON for field validation: %v", err)
			}

			for _, field := range tt.expectedFields {
				if _, exists := result[field]; !exists {
					t.Errorf("Expected field '%s' to be present in JSON", field)
				}
			}

			// Test round-trip marshaling/unmarshaling
			var newItem api.TaskItem
			err = json.Unmarshal(jsonData, &newItem)
			if err != nil {
				t.Errorf("Failed to unmarshal task item: %v", err)
			}

			// Verify discriminator still works after round-trip
			discriminator, err := newItem.Discriminator()
			if err != nil {
				t.Errorf("Failed to get discriminator after round-trip: %v", err)
			} else {
				t.Logf("✅ Discriminator after round-trip: %s", discriminator)
			}

			// Marshal again and compare
			jsonData2, err := json.Marshal(newItem)
			if err != nil {
				t.Errorf("Failed to marshal task item after round-trip: %v", err)
			}

			if string(jsonData) != string(jsonData2) {
				t.Errorf("JSON changed after round-trip marshaling/unmarshaling")
				t.Logf("Original: %s", string(jsonData))
				t.Logf("After round-trip: %s", string(jsonData2))
			}
		})
	}
}

func TestTaskItemJSONMarshalingEdgeCases(t *testing.T) {
	t.Run("Empty TaskItem JSON marshaling", func(t *testing.T) {
		var item api.TaskItem
		jsonData, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Failed to marshal empty task item: %v", err)
		} else {
			t.Logf("✅ Empty TaskItem JSON: %s", string(jsonData))
		}

		// Should only contain the type field
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Errorf("Failed to unmarshal empty task item JSON: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("Expected empty task item to have only 1 field, got %d", len(result))
		}

		if result["type"] != "" {
			t.Errorf("Expected empty type field, got %v", result["type"])
		}
	})

	t.Run("TaskItem with minimal data JSON marshaling", func(t *testing.T) {
		var item api.TaskItem
		_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte{},
			},
		})

		jsonData, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Failed to marshal minimal task item: %v", err)
		} else {
			t.Logf("✅ Minimal TaskItem JSON: %s", string(jsonData))
		}

		// Verify it can be unmarshaled back
		var newItem api.TaskItem
		err = json.Unmarshal(jsonData, &newItem)
		if err != nil {
			t.Errorf("Failed to unmarshal minimal task item: %v", err)
		}

		// Verify discriminator works
		discriminator, err := newItem.Discriminator()
		if err != nil {
			t.Errorf("Failed to get discriminator for minimal task item: %v", err)
		} else {
			t.Logf("✅ Minimal TaskItem discriminator: %s", discriminator)
		}
	})

	t.Run("TaskItem with nil optional fields JSON marshaling", func(t *testing.T) {
		var item api.TaskItem
		_ = item.FromGatewayTransactionTaskItem(api.GatewayTransactionTaskItem{
			Task: api.GatewayTransactionTask{
				ExecuteData: []byte("test"),
			},
			Meta: nil, // Explicitly nil
		})

		jsonData, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Failed to marshal task item with nil meta: %v", err)
		} else {
			t.Logf("✅ TaskItem with nil meta JSON: %s", string(jsonData))
		}

		// Verify meta field is not present in JSON
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Errorf("Failed to unmarshal task item with nil meta: %v", err)
		}

		if _, exists := result["meta"]; exists {
			t.Errorf("Expected meta field to be omitted when nil")
		}
	})
}
