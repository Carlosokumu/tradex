package models

type Mt4Account struct {
	Platform   string  `json:"platform,omitempty"`
	Broker     string  `json:"broker,omitempty"`
	Currency   string  `json:"currency,omitempty"`
	Server     string  `json:"server,omitempty"`
	Balance    float64 `json:"balance"`
	Equity     float64 `json:"equity,omitempty"`
	Margin     float64 `json:"margin,omitempty"`
	FreeMargin float64 `json:"freemargin,omitempty"`
}
