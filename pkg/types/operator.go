package types

import "fmt"

type Operator struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (o Operator) String() string {
	return fmt.Sprintf("Operator{Type: %s, Value: %v}", o.Type, o.Value)
}
