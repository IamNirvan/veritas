package main

import (
	"fmt"

	"github.com/IamNirvan/veritas/pkg/engine"
)

func main() {
	engine := engine.NewVeritas()

	// Define rules in JSON
	rulesJSON := `[
		{
			"conditions": [
				{
					"path": "temperature",
					"operator": "gt",
					"value": 30
				},
				{
					"path": "status",
					"operator": "eq",
					"value": "active"
				}
			],
			"then": {
				"type": "sendEmail",
				"params": {
					"to": "admin@example.com",
					"subject": "High Temperature Alert"
				}
			}
		}
	]`

	// Load rules
	if err := engine.LoadRules([]byte(rulesJSON)); err != nil {
		panic(err)
	}

	// Register action handlers
	engine.RegisterActionHandler("sendEmail", func(params map[string]interface{}) error {
		fmt.Println("sending emails...")
		return nil
	})

	// Evaluate input
	input := map[string]interface{}{
		"temperature": 31,
		"status":      "active",
	}

	if err := engine.EvaluateRules(input); err != nil {
		panic(err)
	}

}
