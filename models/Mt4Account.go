package models

type Mt4Account struct {
	Platform   string  `json:"platform,omitempty"`
	Broker     string  `json:"broker,omitempty"`
	Currency   string  `json:"currency,omitempty"`
	Server     string  `json:"server,omitempty"`
	Balance    float32 `json:"balance"`
	Equity     float32 `json:"equity,omitempty"`
	Margin     float32 `json:"margin,omitempty"`
	FreeMargin float32 `json:"freemargin,omitempty"`
}
