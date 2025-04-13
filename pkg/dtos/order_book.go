package dtos

type OrderBook struct {
	EventType     string     `json:"e"` // Event type
	EventTime     int64      `json:"E"` // Event time
	Symbol        string     `json:"s"` // Symbol
	FirstUpdateID int64      `json:"U"` // First update ID
	LastUpdateID  int64      `json:"u"` // Last update ID
	Bids          [][]string `json:"b"` // Bids
	Asks          [][]string `json:"a"` // Asks
}
