package main

import "github.com/SametAvcii/crypto-trade/app/cmd"

// @title xxx API
// @version 1.0
// @description xxx API Documentation

// @host localhost:8000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cmd.StartApp()
}
