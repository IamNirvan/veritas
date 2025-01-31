package engine

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/IamNirvan/veritas/pkg/types"
)

type Veritas struct {
	rules    []types.Rule
	handlers map[string]types.ActionHandler
}

func NewVeritas() *Veritas {
	return &Veritas{
		rules:    nil,
		handlers: map[string]types.ActionHandler{},
	}
}

func (v *Veritas) LoadRules(raw []byte) error {
	return json.Unmarshal(raw, &v.rules)
}

func (v *Veritas) RegisterActionHandler(actionType string, handler types.ActionHandler) {
	v.handlers[actionType] = handler
}

func (v *Veritas) EvaluateRules(input interface{}) error {
	// Ensure the prerequisites are met
	if err := v.validatePrerequisites(); err != nil {
		return err
	}

	// Convert input to map if it's JSON
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		// Try to convert from JSON if it's not already a map
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return fmt.Errorf("failed to marshal input: %v", err)
		}

		if err := json.Unmarshal(inputBytes, &inputMap); err != nil {
			return fmt.Errorf("failed to unmarshal input: %v", err)
		}
	}

	// Evaluate each rule
	for _, rule := range v.rules {
		if v.evaluateRule(rule, inputMap) {
			if err := v.executeAction(rule.Then); err != nil {
				return fmt.Errorf("failed to execute action: %s", err)
			}
		}
	}
	return nil
}

// evaluateRule evaluates a single rule against the input
func (v *Veritas) evaluateRule(rule types.Rule, input map[string]interface{}) bool {
	for _, condition := range rule.Conditions {
		if !v.evaluateCondition(condition, input) {
			return false
		}
	}
	return true
}

// evaluateCondition evaluates a single condition
func (v *Veritas) evaluateCondition(condition types.Condition, input map[string]interface{}) bool {
	value := getValueFromPath(input, condition.Path)
	if value == nil {
		return false
	}

	switch condition.Operator {
	case "eq":
		return reflect.DeepEqual(value, condition.Value)
	case "gt":
		return compareValues(value, condition.Value) > 0
	case "lt":
		return compareValues(value, condition.Value) < 0
	case "contains":
		return containsValue(value, condition.Value)
	// Add more operators as needed
	default:
		return false
	}
}

// executeAction executes the specified action
func (v *Veritas) executeAction(action types.Action) error {
	handler, exists := v.handlers[action.Type]
	if !exists {
		return fmt.Errorf("no handler registered for action type: %s", action.Type)
	}
	return handler(action.Params)
}

// Helper functions
func getValueFromPath(input map[string]interface{}, path string) interface{} {
	// Simple path resolution for now - could be enhanced to handle nested paths
	return input[path]
}

func compareValues(a, b interface{}) int {
	// Convert to float64 for numeric comparison
	aFloat, aOk := toFloat64(a)
	bFloat, bOk := toFloat64(b)

	if !aOk || !bOk {
		return 0
	}

	if aFloat > bFloat {
		return 1
	} else if aFloat < bFloat {
		return -1
	}
	return 0
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func containsValue(container, item interface{}) bool {
	switch c := container.(type) {
	case []interface{}:
		for _, v := range c {
			if reflect.DeepEqual(v, item) {
				return true
			}
		}
	case string:
		itemStr, ok := item.(string)
		if !ok {
			return false
		}
		return contains(c, itemStr)
	}
	return false
}

func contains(s string, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}

func (v *Veritas) validatePrerequisites() error {
	if v.rules == nil {
		return fmt.Errorf("no rules loaded")
	}

	if len(v.handlers) == 0 {
		return fmt.Errorf("no action handlers registered")
	}

	return nil
}
