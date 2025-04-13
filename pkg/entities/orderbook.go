package entities

type OrderBook struct {
	Base
	Symbol     string `json:"symbol"`
	ExchangeId string `json:"exchange_id"`
	Price      string `json:"price"`
	Amount     string `json:"amount"`
	Side       string `json:"side"`   // bid or ask
	Status     string `json:"status"` // open, closed
}

func (o *OrderBook) FromDto() {

}
