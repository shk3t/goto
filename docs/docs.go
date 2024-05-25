// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Goto GitHub",
            "url": "http://github.com/shk3t/goto"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/delayed-tasks": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Отложенные задания"
                ],
                "summary": "Список моих отложенных заданий",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Вернуть с",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество возвращаемых элементов",
                        "name": "take",
                        "in": "query"
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
                                    "$ref": "#/definitions/goto_src_model.DelayedTask"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/delayed-tasks/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Отложенные задания"
                ],
                "summary": "Детализация отложенного задания",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор отложенного задания",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.DelayedTask"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Авторизация"
                ],
                "summary": "Логин",
                "parameters": [
                    {
                        "description": "Авторизационные данные",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "login": {
                                    "type": "string"
                                },
                                "password": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.User"
                        }
                    }
                }
            }
        },
        "/projects": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Список моих проектов",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Вернуть с",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество возвращаемых элементов",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Язык",
                        "name": "language",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название модуля",
                        "name": "module",
                        "in": "query"
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
                                    "$ref": "#/definitions/goto_src_model.ProjectMin"
                                }
                            }
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
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Добавление проекта с задачами",
                "parameters": [
                    {
                        "description": "Информация о проекте",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "url": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    {
                        "type": "file",
                        "description": "Архив с проектом",
                        "name": "file",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.DelayedTask"
                        }
                    }
                }
            }
        },
        "/projects/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Детализация моего проекта",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.ProjectPublic"
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
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Обновление проекта с задачами",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Информация о проекте",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "url": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    {
                        "type": "file",
                        "description": "Архив с проектом",
                        "name": "file",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.DelayedTask"
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
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Удаление проекта с задачами",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Авторизация"
                ],
                "summary": "Регистрация",
                "parameters": [
                    {
                        "description": "Авторизационные данные",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "login": {
                                    "type": "string"
                                },
                                "password": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.User"
                        }
                    }
                }
            }
        },
        "/solution/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Решения"
                ],
                "summary": "Детализация моего решения",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор решения",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.Solution"
                        }
                    }
                }
            }
        },
        "/solutions": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Решения"
                ],
                "summary": "Список моих решений",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Вернуть с",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество возвращаемых элементов",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Задача",
                        "name": "taskId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "date-time",
                        "example": "2001-12-31",
                        "description": "Минимальная дата обновления",
                        "name": "dateFrom",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "date-time",
                        "example": "2001-12-31",
                        "description": "Максимальная дата обновления",
                        "name": "dateTo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Статус",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Язык",
                        "name": "language",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название модуля",
                        "name": "module",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Устаревшие",
                        "name": "outdated",
                        "in": "query"
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
                                    "$ref": "#/definitions/goto_src_model.SolutionMin"
                                }
                            }
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
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Решения"
                ],
                "summary": "Отправить решение на проверку",
                "parameters": [
                    {
                        "description": "Решение",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.SolutionInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.Solution"
                        }
                    }
                }
            }
        },
        "/tasks": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Задачи"
                ],
                "summary": "Список задач",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Вернуть с",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество возвращаемых элементов",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Созданные мной",
                        "name": "my",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Язык",
                        "name": "language",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название модуля",
                        "name": "module",
                        "in": "query"
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
                                    "$ref": "#/definitions/goto_src_model.TaskMin"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/tasks/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Задачи"
                ],
                "summary": "Детализация задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор задачи",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/goto_src_model.TaskPrivate"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "goto_src_model.DelayedTask": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "details": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "targetId": {
                    "type": "integer"
                },
                "targetName": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.ProjectMin": {
            "type": "object",
            "properties": {
                "containerization": {
                    "type": "string"
                },
                "failKeywords": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "modules": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "srcDir": {
                    "type": "string"
                },
                "stubDir": {
                    "type": "string"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/goto_src_model.TaskMin"
                    }
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.ProjectPublic": {
            "type": "object",
            "properties": {
                "containerization": {
                    "type": "string"
                },
                "failKeywords": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "modules": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "srcDir": {
                    "type": "string"
                },
                "stubDir": {
                    "type": "string"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/goto_src_model.Task"
                    }
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.Solution": {
            "type": "object",
            "properties": {
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/goto_src_model.SolutionFile"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "result": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "task": {
                    "$ref": "#/definitions/goto_src_model.TaskMin"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.SolutionFile": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "solutionId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.SolutionInput": {
            "type": "object",
            "properties": {
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/goto_src_model.SolutionFile"
                    }
                },
                "taskId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.SolutionMin": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "task": {
                    "$ref": "#/definitions/goto_src_model.TaskMin"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.Task": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/goto_src_model.TaskFile"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "modules": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "projectId": {
                    "type": "integer"
                },
                "runtarget": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "goto_src_model.TaskFile": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "stub": {
                    "type": "string"
                },
                "taskId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.TaskFilePrivate": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "stub": {
                    "type": "string"
                },
                "taskId": {
                    "type": "integer"
                }
            }
        },
        "goto_src_model.TaskMin": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "fileNames": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "modules": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "projectId": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "goto_src_model.TaskPrivate": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/goto_src_model.TaskFilePrivate"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "modules": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "projectId": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "goto_src_model.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Prepend your JWT key with ` + "`" + `Bearer` + "`" + `",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "localhost:3228",
	BasePath:         "/api/",
	Schemes:          []string{"http"},
	Title:            "Goto",
	Description:      "Web app for code challenges with any environments",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
