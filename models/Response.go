package models

type Response struct {
	Data []struct {
		AccountID                   int    `json:"accountId,omitempty"`
		AccountNumber               int    `json:"accountNumber,omitempty"`
		Live                        bool   `json:"live,omitempty"`
		BrokerName                  string `json:"brokerName,omitempty"`
		BrokerTitle                 string `json:"brokerTitle,omitempty"`
		DepositCurrency             string `json:"depositCurrency,omitempty"`
		TraderRegistrationTimestamp int64  `json:"traderRegistrationTimestamp,omitempty"`
		TraderAccountType           string `json:"traderAccountType,omitempty"`
		Leverage                    int    `json:"leverage,omitempty"`
		LeverageInCents             int    `json:"leverageInCents,omitempty"`
		Balance                     int    `json:"balance,omitempty"`
		Deleted                     bool   `json:"deleted,omitempty"`
		AccountStatus               string `json:"accountStatus,omitempty"`
		SwapFree                    bool   `json:"swapFree,omitempty"`
		MoneyDigits                 int    `json:"moneyDigits,omitempty"`
	} `json:"data,omitempty"`
}
