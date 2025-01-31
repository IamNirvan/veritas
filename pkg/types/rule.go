package types

type Rule struct {
	Conditions []Condition `json:"conditions"`
	Then       Action      `json:"action"`
}
