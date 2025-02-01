package engine

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/IamNirvan/veritas/pkg/types"
	"github.com/IamNirvan/veritas/pkg/util"
)

type Engine struct {
	rules    []types.Rule
	handlers map[string]*types.ActionHandler
}

func NewEngine() *Engine {
	return &Engine{
		rules:    nil,
		handlers: map[string]*types.ActionHandler{},
	}
}

func (e *Engine) LoadRules(rules json.RawMessage) error {
	return json.Unmarshal(rules, &e.rules)
}

func (e *Engine) RegisterActionHandler(actionType string, handler *types.ActionHandler) {
	e.handlers[actionType] = handler
}

func (e *Engine) EvaluateRules(input *json.RawMessage) error {
	// Ensure the prerequisites are met
	if err := e.validatePrerequisites(); err != nil {
		return err
	}

	// Decode input json intp map
	inputMap, inputMapErr := e.decodeInputToMap(input)
	if inputMapErr != nil {
		return inputMapErr
	}

	// Evaluate each condition in the rule against the input
	var finalResult int
	for _, rule := range e.rules {
		for index, condition := range rule.Conditions {
			conditionResult, conditionErr := e.evaluateCondition(index, inputMap, &condition)
			if conditionErr != nil {
				return conditionErr
			}

			if index == 0 {
				finalResult = conditionResult
			} else {
				if condition.And != nil {
					finalResult &= conditionResult
				} else if condition.Or != nil {
					finalResult |= conditionResult
				}
			}
		}

		// Execute the action if the final result is 1
		if finalResult == 1 {
			if err := e.executeAction(rule.Then); err != nil {
				return fmt.Errorf("failed to execute action: %s", err)
			}
		}
	}
	return nil
}

func (e *Engine) decodeInputToMap(input *json.RawMessage) (*map[string]interface{}, error) {
	var inputMap map[string]interface{}

	if err := json.Unmarshal(*input, &inputMap); err != nil {
		return nil, fmt.Errorf("failed to decode json rules to map: %v", err)
	}

	return &inputMap, nil
}

func (e *Engine) evaluateCondition(conditionCount int, inputMap *map[string]interface{}, condition *types.Condition) (int, error) {
	// The condition in the rule tells me what is expected
	// Compare this with the inputValue in the input
	var inputValue interface{}
	var operationType string
	var conditionValue interface{}

	if conditionCount == 0 {
		inputValue = util.GetValueFromPathV2(inputMap, *(*condition).Path)
		operationType = (*condition).Operator.Type
		conditionValue = (*condition).Operator.Value
	} else {
		if condition.And != nil {
			inputValue = util.GetValueFromPathV2(inputMap, (*condition).And.Path)
			operationType = (*condition).And.Operator.Type
			conditionValue = (*condition).And.Operator.Value
		} else {
			inputValue = util.GetValueFromPathV2(inputMap, (*condition).Or.Path)
			operationType = (*condition).Or.Operator.Type
			conditionValue = (*condition).Or.Operator.Value
		}
	}

	var res bool
	switch operationType {
	case "equals":
		res = reflect.DeepEqual(inputValue, conditionValue)
	case "gt":
		res = util.CompareValues(inputValue, conditionValue) > 0
	case "lt":
		res = util.CompareValues(inputValue, conditionValue) == -1
	default:
		return 0, fmt.Errorf("unsupported operation type: %s", operationType)
	}

	util.LogConditionSummary(inputValue, operationType, conditionValue, res)
	if res {
		return 1, nil
	}
	return 0, nil
}

func (e *Engine) executeAction(action types.Action) error {
	handler, exists := e.handlers[action.Type]
	if !exists {
		return fmt.Errorf("no handler registered for action type: %s", action.Type)
	}
	return (*handler)(action.Params)
}

func (e *Engine) validatePrerequisites() error {
	if e.rules == nil {
		return fmt.Errorf("no rules loaded")
	}

	if len(e.handlers) == 0 {
		return fmt.Errorf("no action handlers registered")
	}

	return nil
}
