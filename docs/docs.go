// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "lifthus531@gmail.com",
            "email": "lifthus531@gmail.com"
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
        "/session/check/:service/:sid": {
            "post": {
                "description": "checks the hus session in cookie and tells the subservice server if the session is valid with SID.",
                "tags": [
                    "auth"
                ],
                "summary": "accepts sid and service name to check if the session is valid.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "subservice name",
                        "name": "service",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "session id",
                        "name": "sid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok, theclient now should go to subservice's signing endpoint"
                    },
                    "401": {
                        "description": "Unauthorized, the client is not signed in"
                    },
                    "404": {
                        "description": "Not Found, the service is not registered"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/session/revoke": {
            "delete": {
                "description": "can be used to sign out.",
                "tags": [
                    "auth"
                ],
                "summary": "revokes every hus session in cookie from database.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hus session tokens in cookie",
                        "name": "jwt",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok"
                    },
                    "500": {
                        "description": "doesn't have to be handled"
                    }
                }
            }
        },
        "/social/google/{subservice_name}": {
            "post": {
                "description": "validates the google ID token and redirects with hus refresh token to /auth/{token_string}.\nthe refresh token will be expired in 7 days.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "gets google IDtoken and redirect with hus session cookie.",
                "parameters": [
                    {
                        "description": "Google ID token",
                        "name": "jwt",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "301": {
                        "description": "to /error"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.0",
	Host:             "lifthus.com",
	BasePath:         "/auth",
	Schemes:          []string{},
	Title:            "Project-Hus auth server",
	Description:      "This is Project-Hus's root authentication server containing each user's UUID, which is unique for all hus services.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
