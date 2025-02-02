package docs

import (
	"github.com/swaggo/swag"
)

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
        "/api/events": {
            "get": {
                "description": "Retrieve event logs with optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by status (active/archived)",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by trigger ID",
                        "name": "trigger_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.EventLog"
                            }
                        }
                    }
                }
            }
        },
        // Other paths will be auto-generated here
    },
    "definitions": {
        "models.EventLog": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_test": {
                    "type": "boolean"
                },
                "payload": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string"
                },
                "trigger_id": {
                    "type": "string"
                },
                "triggered_at": {
                    "type": "string"
                }
            }
        },
        // Other definitions will be auto-generated
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:           "Event Trigger API",
	Description:     "API for managing event triggers",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}