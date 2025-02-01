package types

import "fmt"

type And struct {
	Path     string   `json:"path"`
	Operator Operator `json:"operator"`
}

func (a And) String() string {
	return fmt.Sprintf("And{Path: %s,  %s}", a.Path, a.Operator)
}

type Or struct {
	Path     string   `json:"path"`
	Operator Operator `json:"operator"`
}

func (o Or) String() string {
	return fmt.Sprintf("Or{Path: %s,  %s}", o.Path, o.Operator)
}

type Condition struct {
	Path     *string   `json:"path"`
	Operator *Operator `json:"operator"`
	And      *And      `json:"and"`
	Or       *Or       `json:"or"`
}

func (c Condition) String() string {
	return fmt.Sprintf("Condition{Path: %s, %s, %v, %v}", getStringValue(c.Path), c.Operator, c.And, c.Or)
}

func getStringValue(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
