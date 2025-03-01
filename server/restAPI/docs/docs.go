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
        "/job/result/{id}": {
            "get": {
                "description": "get job results by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "get job results",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.JobResultData"
                        }
                    },
                    "400": {
                        "description": "Invalid ID parameter",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "No results for this ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed during DB connection",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/job/{id}": {
            "get": {
                "description": "get all data from job by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "get job data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "job ID",
                        "name": "id",
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
                                "$ref": "#/definitions/types.JobData"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid ID parameter",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "It wasn't possible to find the job",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed during DB connection",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "delete all data related to this job id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "delete job data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Job ID",
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
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid ID parameter",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "No results for this ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed during DB connection",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/jobs": {
            "get": {
                "description": "get all data from jobs",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "get jobs data",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Last id(order) gotten from db",
                        "name": "cursor",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.JobData"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed during DB connection",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/plugin/{name}": {
            "post": {
                "description": "add plugin by name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "plugins"
                ],
                "summary": "add plugin",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Plugin name as shown in the github org",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.JobData": {
            "type": "object",
            "properties": {
                "finish_time": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "order": {
                    "type": "integer"
                },
                "qasm": {
                    "type": "string"
                },
                "result_types": {
                    "$ref": "#/definitions/types.JobResultTypes"
                },
                "results": {
                    "$ref": "#/definitions/types.JobResultData"
                },
                "start_time": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "submission_date": {
                    "type": "string"
                },
                "target_simulator": {
                    "type": "string"
                }
            }
        },
        "types.JobResultData": {
            "type": "object",
            "properties": {
                "counts": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                },
                "expval": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "id": {
                    "type": "string"
                },
                "job_id": {
                    "type": "string"
                },
                "quasi_dist": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                }
            }
        },
        "types.JobResultTypes": {
            "type": "object",
            "properties": {
                "counts": {
                    "type": "boolean"
                },
                "expval": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "job_id": {
                    "type": "string"
                },
                "quasi_dist": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
