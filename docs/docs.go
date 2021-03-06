// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Hemakshi Sachdev",
            "email": "sachdev.hemakshi@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/login/github/oauth2": {
            "get": {
                "description": "Performs a Github OAuth by using the ` + "`" + `code` + "`" + ` provided in query param and then uses the received accessToken to fetch user information from Github. If it is a new user then user is also saved in the database. It also creates a user session for future authenticated calls.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Login via Github OAuth",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.UserInfo"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/repositories": {
            "get": {
                "description": "Get complete repository information for all repositories from the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "repository"
                ],
                "summary": "Get all repositories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Repository"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/repository/{repoID}/update/lastEventAt": {
            "put": {
                "description": "Updates ` + "`" + `lastEventAt` + "`" + ` time for the given ` + "`" + `repoID` + "`" + ` in the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "repository"
                ],
                "summary": "Updates ` + "`" + `lastEventAt` + "`" + ` time for the given ` + "`" + `repoID` + "`" + `",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Repository ID for which ` + "`" + `lastEventAt` + "`" + ` needs to be updated",
                        "name": "repoID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "` + "`" + `lastEventAt` + "`" + ` time (with timezone) at which last event for the repository occurred. Format ` + "`" + `lastEventAt` + "`" + `: ` + "`" + `2006-01-02 15:04:05-07:00` + "`" + `",
                        "name": "lastEventAt",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.lastEventAt"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/authenticated": {
            "get": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Get user information from Github for authenticated user via session token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user information from Github for authenticated user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.UserInfo"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/logout": {
            "get": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Performs the logout function by deleting user session",
                "tags": [
                    "user"
                ],
                "summary": "Logout from user session",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/subscription/add": {
            "post": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Create new subscription for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Create new subscription for the authenticated user",
                "parameters": [
                    {
                        "description": "Repository Name and list of Labels to create a new subscription",
                        "name": "subscriptions",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Subscription"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/subscription/remove": {
            "delete": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Deletes existing subscription for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Deletes existing subscription for the authenticated user",
                "parameters": [
                    {
                        "description": "Repository Name and list of Labels which needs to be deleted from the existing subscription",
                        "name": "subscriptions",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Subscription"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/subscription/update": {
            "put": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Update existing subscription for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Update existing subscription for the authenticated user",
                "parameters": [
                    {
                        "description": "Repository Name and list of Labels which needs to be added to the existing subscription",
                        "name": "subscriptions",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Subscription"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/subscription/view": {
            "get": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Get all subscriptions for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Get all subscriptions for the authenticated user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Subscription"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/{repoName}/subscription/labels": {
            "get": {
                "security": [
                    {
                        "Github OAuth": []
                    }
                ],
                "description": "Get all subscriptions for the given ` + "`" + `userID` + "`" + ` and ` + "`" + `repoName` + "`" + ` of the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Get all subscriptions for the given ` + "`" + `userID` + "`" + ` and ` + "`" + `repoName` + "`" + ` of the authenticated user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Repository Name for which subscription data needs to be fetched. Format ` + "`" + `facebook/react` + "`" + `",
                        "name": "repoName",
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
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/models.Label"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/{repoID}/subscription/view": {
            "get": {
                "description": "Get all subscriptions for the given ` + "`" + `repoID` + "`" + `",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Get all subscriptions for the given ` + "`" + `repoID` + "`" + `",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Repository ID for which subscription data needs to be fetched",
                        "name": "repoID",
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
                                "$ref": "#/definitions/models.UserIDLabel"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Label": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.Repository": {
            "type": "object",
            "properties": {
                "lastEventAt": {
                    "type": "string"
                },
                "repoID": {
                    "type": "string"
                },
                "repoName": {
                    "type": "string"
                }
            }
        },
        "models.Subscription": {
            "type": "object",
            "properties": {
                "labels": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Label"
                    }
                },
                "repoName": {
                    "type": "string"
                }
            }
        },
        "models.UserIDLabel": {
            "type": "object",
            "properties": {
                "label": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "routes.UserInfo": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "avatarImg": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "routes.lastEventAt": {
            "type": "object",
            "properties": {
                "lastEventAt": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8001",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Github Issue-Notifier API",
	Description: "APIs for the Github Issue Notifier Project. https://github.com/issue-notifier",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
