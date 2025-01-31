package types

type Condition struct {
	Path     string       `json:"path"`
	Operator string       `json:"operator"`
	Value    *interface{} `json:"value"`
}
