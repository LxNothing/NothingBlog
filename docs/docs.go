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
        "/article/all": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "通过该接口可以获得当前的所有文章",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章相关接口"
                ],
                "summary": "获取所有文章的接口",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "security": [
                    {
                        "No": []
                    }
                ],
                "description": "用于用户登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证相关接口"
                ],
                "summary": "用户登录接口",
                "parameters": [
                    {
                        "description": "登录参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/auth/password/modify": {
            "post": {
                "security": [
                    {
                        "No": []
                    }
                ],
                "description": "用于修改用户密码，不需要登录，需要验证旧密码，账户和验证码",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证相关接口"
                ],
                "summary": "修改用户密码",
                "parameters": [
                    {
                        "description": "修改密码的参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ModifyPasswordParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/auth/password/reset": {
            "post": {
                "security": [
                    {
                        "No": []
                    }
                ],
                "description": "用于用户重置密码，需要使用邮箱 - 功能未实现",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证相关接口"
                ],
                "summary": "重置密码接口",
                "parameters": [
                    {
                        "description": "重置密码的参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ResetPasswordParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "security": [
                    {
                        "No": []
                    }
                ],
                "description": "用户注册的接口，需要接收参数",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证相关接口"
                ],
                "summary": "用户注册的接口",
                "parameters": [
                    {
                        "description": "注册参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignUpParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/auth/verifycode": {
            "get": {
                "security": [
                    {
                        "No": []
                    }
                ],
                "description": "通过该接口可以获得基于数字的验证码，目前只支持数字验证码，后续可以更改",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证相关接口"
                ],
                "summary": "获取数字验证码的接口",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.ResponseCodeType": {
            "type": "integer",
            "enum": [
                1000,
                1001,
                1002,
                1003,
                1004,
                1005,
                1006,
                1007,
                1008,
                1009,
                1010,
                1011,
                1012
            ],
            "x-enum-varnames": [
                "CodeSuccess",
                "CodeParameterInvalid",
                "CodeUserNotExist",
                "CodeUserExist",
                "CodePasswordError",
                "CodeServerBusy",
                "CodeTokenInvaild",
                "CodeTokenEmpty",
                "CodeNeedReLogin",
                "CodeCommunityIdInvalid",
                "CodeVerifyCodeInvaild",
                "CodeArticleTitleExisted",
                "CodeArticleNotExisted"
            ]
        },
        "controller.ResponseData": {
            "type": "object",
            "properties": {
                "code": {
                    "$ref": "#/definitions/controller.ResponseCodeType"
                },
                "data": {
                    "description": "omitempty 该字段为空时忽略"
                },
                "msg": {}
            }
        },
        "models.LoginParams": {
            "type": "object",
            "required": [
                "password",
                "username",
                "verify_code"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "verify_code": {
                    "description": "验证码",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.VerifyCodeParams"
                        }
                    ]
                }
            }
        },
        "models.ModifyPasswordParams": {
            "type": "object",
            "required": [
                "new_password",
                "old_password",
                "username",
                "verify_code"
            ],
            "properties": {
                "new_password": {
                    "description": "新密码",
                    "type": "string"
                },
                "old_password": {
                    "description": "旧密码",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                },
                "verify_code": {
                    "description": "验证码",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.VerifyCodeParams"
                        }
                    ]
                }
            }
        },
        "models.ResetPasswordParams": {
            "type": "object",
            "required": [
                "username"
            ],
            "properties": {
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "models.SignUpParams": {
            "type": "object",
            "required": [
                "password",
                "re_password",
                "username",
                "verify_code"
            ],
            "properties": {
                "email": {
                    "description": "邮箱 - 仅在重置密码时，接收验证码使用",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "re_password": {
                    "description": "确认密码",
                    "type": "string"
                },
                "username": {
                    "description": "用户名 - 用户名作为唯一标识，不允许重复",
                    "type": "string"
                },
                "verify_code": {
                    "description": "验证码",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.VerifyCodeParams"
                        }
                    ]
                }
            }
        },
        "models.VerifyCodeParams": {
            "type": "object",
            "required": [
                "code",
                "id"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "id": {
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
	BasePath:         "",
	Schemes:          []string{},
	Title:            "NgBlog",
	Description:      "NgBlog Go博客项目 API 接口文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}