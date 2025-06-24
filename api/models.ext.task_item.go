package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrUnknownTaskType is an error when TaskType has an unrecognised value
var ErrUnknownTaskType = errors.New("unknown task type")

// SetTaskFromJSON set TaskItem.Task from a given JSON string as specified TaskType
func (t *TaskItem) SetTaskFromJSON(taskType TaskType, taskJSON string) error {
	switch taskType {
	case TaskTypeConstructProof:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromConstructProofTask); err != nil {
			return err
		}
	case TaskTypeExecute:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromExecuteTask); err != nil {
			return err
		}
	case TaskTypeGatewayTransaction:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromGatewayTransactionTask); err != nil {
			return err
		}
	case TaskTypeReactToExpiredSigningSession:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromReactToExpiredSigningSessionTask); err != nil {
			return err
		}
	case TaskTypeReactToWasmEvent:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromReactToWasmEventTask); err != nil {
			return err
		}
	case TaskTypeRefund:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromRefundTask); err != nil {
			return err
		}
	case TaskTypeVerify:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromVerifyTask); err != nil {
			return err
		}
	case TaskTypeReactToRetriablePoll:
		if err := setTaskFromJSONWithSetter(taskJSON, t.Task.FromReactToRetriablePollTask); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownTaskType, taskType)
	}

	return nil
}

func setTaskFromJSONWithSetter[T any](taskJSON string, taskSetter func(T) error) error {
	var task T
	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		return err
	}
	if err := taskSetter(task); err != nil {
		return err
	}
	return nil
}
