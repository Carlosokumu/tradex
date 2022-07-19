package models

type PositionData struct {
	PositionId int     `json:"PositionId,omitempty"`
	EntryPrice float32 `json:"EntryPrice,omitempty"`
	TradeType  string  `json:"TradeType,omitempty"`
	Quantity   float32 `json:"Quantity,omitempty"`
}
