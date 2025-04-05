package entities

type Log struct {
	Base
	Title   string `json:"title" example:"example title"`
	Message string `json:"message" example:"order created"`
	Entity  string `json:"entity" example:"order"`
	Type    string `json:"type" example:"info"`  // -----> info, error
	Proto   string `json:"proto" example:"http"` // -----> http, grpc
	Ip      string `json:"ip" example:"127.0.0.1"`
}
