package types

type Action struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}
