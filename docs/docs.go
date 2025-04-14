// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Add Symbol Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbol Endpoints"
                ],
                "summary": "Add Symbol",
                "parameters": [
                    {
                        "description": "Add Symbol Request",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.AddSymbolReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/exchanges": {
            "get": {
                "description": "Retrieves all exchanges",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchanges"
                ],
                "summary": "Get all exchanges",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new exchange",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchanges"
                ],
                "summary": "Add exchange",
                "parameters": [
                    {
                        "description": "Exchange information",
                        "name": "exchange",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.AddExchangeReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully created exchange",
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.AddExchangeRes"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/exchanges/{id}": {
            "get": {
                "description": "Retrieves an exchange by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchanges"
                ],
                "summary": "Get exchange by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Exchange ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an existing exchange with the provided information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchanges"
                ],
                "summary": "Update an exchange",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Exchange ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Exchange update information",
                        "name": "exchange",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateExchangeReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully updated exchange",
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateExchangeRes"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes an exchange by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchanges"
                ],
                "summary": "Delete an exchange",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Exchange ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/signal/interval": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get All Signal Intervals Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Signal Endpoints"
                ],
                "summary": "Get All Signal Intervals",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Add Signal Interval Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Signal Endpoints"
                ],
                "summary": "Add Signal Interval",
                "parameters": [
                    {
                        "description": "Add Symbol Request",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.AddSignalIntervalReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/signal/interval/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get Signal Interval by ID Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Signal Endpoints"
                ],
                "summary": "Get Signal Interval by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Signal ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update Signal Interval Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Signal Endpoints"
                ],
                "summary": "Update Signal Interval",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Signal ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Signal Request",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateSignalIntervalReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete Signal Interval Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Signal Endpoints"
                ],
                "summary": "Delete Signal Interval",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Signal ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/symbol": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get All Symbols Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbol Endpoints"
                ],
                "summary": "Get All Symbols",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.GetSymbolRes"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/symbol/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get Symbol by ID Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbol Endpoints"
                ],
                "summary": "Get Symbol by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Symbol ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.GetSymbolRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update Symbol Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbol Endpoints"
                ],
                "summary": "Update Symbol",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Symbol ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Symbol Request",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateSymbolReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete Symbol Route",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Symbol Endpoints"
                ],
                "summary": "Delete Symbol",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Symbol ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "gin.H": {
            "type": "object",
            "additionalProperties": {}
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.AddExchangeReq": {
            "type": "object",
            "properties": {
                "name": {
                    "description": "Binance",
                    "type": "string"
                },
                "ws_url": {
                    "description": "wss://ws-api.binance.com:443/ws-api/v3",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.AddExchangeRes": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "1",
                    "type": "string"
                },
                "name": {
                    "description": "Binance",
                    "type": "string"
                },
                "ws_url": {
                    "description": "wss://ws-api.binance.com:443/ws-api/v3",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.AddSignalIntervalReq": {
            "type": "object",
            "properties": {
                "exchange_id": {
                    "description": "1",
                    "type": "string"
                },
                "interval": {
                    "description": "1m, 5m, 15m, 1h, 4h, 1d",
                    "type": "string"
                },
                "symbol": {
                    "description": "BTCUSDT",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.AddSymbolReq": {
            "type": "object",
            "properties": {
                "exchange_id": {
                    "type": "string"
                },
                "symbol": {
                    "description": "BTCUSDT",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.GetSymbolRes": {
            "type": "object",
            "properties": {
                "exchange_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_active": {
                    "description": "1 active, 2 passive",
                    "type": "integer"
                },
                "symbol": {
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateExchangeReq": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "description": "Binance",
                    "type": "string"
                },
                "ws_url": {
                    "description": "wss://ws-api.binance.com:443/ws-api/v3",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateExchangeRes": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "description": "Binance",
                    "type": "string"
                },
                "ws_url": {
                    "description": "wss://ws-api.binance.com:443/ws-api/v3",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateSignalIntervalReq": {
            "type": "object",
            "properties": {
                "exchange_id": {
                    "description": "1",
                    "type": "string"
                },
                "id": {
                    "description": "1",
                    "type": "string"
                },
                "interval": {
                    "description": "1m, 5m, 15m, 1h, 4h, 1d",
                    "type": "string"
                },
                "symbol": {
                    "description": "BTCUSDT",
                    "type": "string"
                }
            }
        },
        "github_com_SametAvcii_crypto-trade_pkg_dtos.UpdateSymbolReq": {
            "type": "object",
            "properties": {
                "exchange_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8001",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "Crypto Trade API",
	Description:      "Crypto Trade API Documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
