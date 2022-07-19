package models

type PositionData struct {
	PositionId int    `json:"PositionId,omitempty"`
	EntryPrice int    `json:"EntryPrice,omitempty"`
	TradeType  bool   `json:"TradeType,omitempty"`
	Quantity   string `json:"Quantity,omitempty"`
}
