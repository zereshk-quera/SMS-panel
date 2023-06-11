// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/account/{accountID}/phone-books/": {
            "get": {
                "description": "Get all phone books for a given account ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "phonebook"
                ],
                "summary": "Get all phone books",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "accountID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.PhoneBookResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/account/{accountID}/phone-books/{phoneBookID}": {
            "get": {
                "description": "Get a phone book by ID for a given account ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "phonebook"
                ],
                "summary": "Get a phone book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "accountID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Phone Book ID",
                        "name": "phoneBookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.PhoneBookResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a phone book by ID for a given account ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "phonebook"
                ],
                "summary": "Update a phone book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "accountID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Phone Book ID",
                        "name": "phoneBookID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Phone Book object",
                        "name": "phoneBook",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.PhoneBookRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.PhoneBookResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a phone book by ID for a given account ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "phonebook"
                ],
                "summary": "Delete a phone book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "accountID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Phone Book ID",
                        "name": "phoneBookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/phonebook": {
            "post": {
                "description": "Create a new phone book entry",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "phonebook"
                ],
                "summary": "Create a phone book entry",
                "parameters": [
                    {
                        "description": "Phone book entry data",
                        "name": "phoneBook",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.PhoneBookRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.PhoneBookResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "handlers.PhoneBookRequest": {
            "type": "object",
            "required": [
                "accountID",
                "name"
            ],
            "properties": {
                "accountID": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.PhoneBookResponse": {
            "type": "object",
            "properties": {
                "accountID": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "SMS-PANEL",
	Description:      "Simple SMS-PANEL server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
