package models

type PositionData struct {
	AccountID                   int    `json:"accountId,omitempty"`
	AccountNumber               int    `json:"accountNumber,omitempty"`
	Live                        bool   `json:"live,omitempty"`
	BrokerName                  string `json:"brokerName,omitempty"`
	BrokerTitle                 string `json:"brokerTitle,omitempty"`
	DepositCurrency             string `json:"depositCurrency,omitempty"`
	TraderRegistrationTimestamp int64  `json:"traderRegistrationTimestamp,omitempty"`
	TraderAccountType           string `json:"traderAccountType,omitempty"`
}
