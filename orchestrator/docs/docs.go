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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/calculate": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "calculate"
                ],
                "parameters": [
                    {
                        "description": "Объект, содержащий в себе выражение",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CalculateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CalculateResponse"
                        }
                    },
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.CalculateResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/api/v1/expressions": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "expressions"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ListAllExpressionsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/api/v1/expressions/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "expressions"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID выражения",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GetByIdExpressionResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/internal/task": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "internal"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TaskResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "internal"
                ],
                "parameters": [
                    {
                        "description": "Объект, содержащий в себе результат части выражения",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ApiError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "models.CalculateRequest": {
            "type": "object",
            "required": [
                "expression"
            ],
            "properties": {
                "expression": {
                    "type": "string",
                    "example": "2+2"
                }
            }
        },
        "models.CalculateResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "928b303f-cfcc-46f4-ae24-aabb72bbb7d9"
                }
            }
        },
        "models.Expression": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "928b303f-cfcc-46f4-ae24-aabb72bbb7d9"
                },
                "result": {
                    "type": "number"
                },
                "status": {
                    "type": "string",
                    "example": "DONE"
                }
            }
        },
        "models.GetByIdExpressionResponse": {
            "type": "object",
            "properties": {
                "expression": {
                    "$ref": "#/definitions/models.Expression"
                }
            }
        },
        "models.ListAllExpressionsResponse": {
            "type": "object",
            "properties": {
                "expressions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Expression"
                    }
                }
            }
        },
        "models.TaskRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "928b303f-cfcc-46f4-ae24-aabb72bbb7d9"
                },
                "result": {}
            }
        },
        "models.TaskResponse": {
            "type": "object",
            "properties": {
                "arg1": {
                    "type": "number",
                    "example": 1
                },
                "arg2": {
                    "type": "number",
                    "example": 1
                },
                "id": {
                    "type": "string",
                    "example": "928b303f-cfcc-46f4-ae24-aabb72bbb7d9"
                },
                "operation": {
                    "type": "string",
                    "example": "+"
                },
                "operation_time": {
                    "type": "integer",
                    "example": 1000
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9090",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Orchestrator API",
	Description:      "API documentation for the Calc Orchestrator",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
