package main

import (
	"encoding/json"
	"fmt"

	"github.com/IamNirvan/veritas/pkg/engine"
	"github.com/IamNirvan/veritas/pkg/types"
)

func main() {
	engine := engine.NewEngine()

	rulesJSON := json.RawMessage(`
	[
		{
			"conditions": [
				{
					"path": "temperature",
					"operator": {
						"type": "lt",
						"value": 100
					}
				},
				{
					"and": {
						"path": "size",
						"operator": {
							"type": "lt",
							"value": 100
						}
					}
				}
			],
			"then": {
				"type": "sendSMS",
				"params": {
					"contact": "0724454572"
				}
			}
    	}
	]
	`)

	// Load rules
	if err := engine.LoadRules(rulesJSON); err != nil {
		panic(err)
	}

	// Register action handlers
	sendSMSHandler := func(params map[string]interface{}) error {
		fmt.Printf("\nsending SMS to contact %s \n", params["contact"])
		return nil
	}

	engine.RegisterActionHandler("sendSMS", (*types.ActionHandler)(&sendSMSHandler))

	// Evaluate input
	input := json.RawMessage(`
		{
			"temperature": 90,
			"size": 20
		}
	`)

	if err := engine.EvaluateRules(&input); err != nil {
		panic(err)
	}
}
