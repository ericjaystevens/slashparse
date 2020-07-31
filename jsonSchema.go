// THIS  FILE IS GENERATED, DO NOT EDIT, INSTEAD UPDATE schema.json and run generate/generateFromStatic.go
	
package slashparse

const jsonSchemaContent ="{\r\n  \"$id\": \"https://example.com/person.schema.json\",\r\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\r\n  \"title\": \"SlashCommand\",\r\n  \"type\": \"object\",\r\n  \"properties\": {\r\n    \"name\": {\r\n      \"type\": \"string\",\r\n      \"description\": \"The Name of the Slash Command.\"\r\n    },\r\n    \"description\": {\r\n      \"type\": \"string\",\r\n      \"description\": \"A description of what the slash command does\"\r\n    },\r\n    \"arguments\": {\r\n      \"type\": \"array\",\r\n      \"description\": \"Pass these to your slash command\",\r\n      \"properties\": {\r\n        \"name\": {\r\n          \"type\": \"string\",\r\n          \"description\": \"Name of argument of Slash command\"\r\n        },\r\n        \"argtype\": {\r\n          \"type\": \"string\",\r\n          \"description\": \"SlashParse built-in argument types\",\r\n          \"enum\": [\"word\", \"number\", \"quoted text\", \"date\", \"time\", \"remaining text\"]\r\n        },\r\n        \"description\": {\r\n          \"type\": \"string\",\r\n          \"description\": \"Description of the argument being passed\"\r\n        },\r\n        \"errorMsg\": {\r\n          \"type\": \"string\",\r\n          \"description\": \"custom error message if argument does not meet requirements\"\r\n        },\r\n        \"position\": {\r\n          \"type\": \"number\",\r\n          \"description\": \"poition of the argument relative to the slash command\"\r\n        },\r\n        \"required\": {\r\n         \"type\": \"boolean\",\r\n         \"description\": \"If the arguemnt is required\" \r\n        }\r\n      },\r\n      \"required\": [\"name\", \"argtype\", \"description\"]\r\n    },\r\n    \"subcommands\": {\r\n      \"type\": \"array\",\r\n      \"description\": \"A Sub command of the slash command, often a noun\",\r\n      \"properties\": {\r\n        \"name\": {\r\n          \"type\": \"string\",\r\n          \"description\": \"Name of sub command\"\r\n        },\r\n        \"description\": {\r\n          \"type\": \"string\",\r\n          \"description\": \"description of sub command\"\r\n        },\r\n        \"arguments\": {\r\n          \"description\": \"Pass these to your slash sub command\",\r\n          \"properties\": {\r\n            \"name\": {\r\n              \"type\": \"string\",\r\n              \"description\": \"Name of argument of Slash sub command\"\r\n            },\r\n            \"argtype\": {\r\n              \"type\": \"string\",\r\n              \"description\": \"SlashParse built-in argument types\",\r\n              \"enum\": [\"word\", \"number\", \"quoted text\", \"date\", \"time\", \"remaining text\"]\r\n            },\r\n            \"description\": {\r\n              \"type\": \"string\",\r\n              \"description\": \"Description of the argument being passed to the sub command\"\r\n            },\r\n            \"errorMsg\": {\r\n              \"type\": \"string\",\r\n              \"description\": \"custom error message if argument does not meet requirements\"\r\n            },\r\n            \"position\": {\r\n              \"type\": \"number\",\r\n              \"description\": \"poition of the argument relative to the slash sub command\"\r\n            },\r\n            \"required\": {\r\n            \"type\": \"boolean\",\r\n            \"description\": \"Is the arguemnt required?\" \r\n            }\r\n          },\r\n          \"required\": [\"name\", \"argtype\", \"description\"]\r\n        },\r\n        \"subcommands\": {\r\n          \"type\": \"array\",\r\n          \"description\": \"a sub sub command\",\r\n          \"properties\": {\r\n            \"name\": {\r\n              \"type\": \"string\",\r\n              \"description\": \"Name of sub sub command, often an action word\"\r\n            },\r\n            \"description\": {\r\n              \"type\": \"string\",\r\n              \"description\": \"description of a sub sub command\"\r\n            },\r\n            \"arguments\": {\r\n              \"description\": \"Pass these to your slash sub-sub command\",\r\n              \"properties\": {\r\n                \"name\": {\r\n                  \"type\": \"string\",\r\n                  \"description\": \"Name of argument of Slash sub-sub command\"\r\n                },\r\n                \"argtype\": {\r\n                  \"type\": \"string\",\r\n                  \"description\": \"SlashParse built-in argument types\",\r\n                  \"enum\": [\"word\", \"number\", \"quoted text\", \"date\", \"time\", \"remaining text\"]\r\n                },\r\n                \"description\": {\r\n                  \"type\": \"string\",\r\n                  \"description\": \"Description of the argument being passed to the sub-sub command\"\r\n                },\r\n                \"errorMsg\": {\r\n                  \"type\": \"string\",\r\n                  \"description\": \"custom error message if argument does not meet requirements\"\r\n                },\r\n                \"position\": {\r\n                  \"type\": \"number\",\r\n                  \"description\": \"poition of the argument relative to the slash sub-sub command\"\r\n                },\r\n                \"required\": {\r\n                \"type\": \"boolean\",\r\n                \"description\": \"If the arguemnt is required\" \r\n                }\r\n              },\r\n              \"required\": [\"name\", \"argtype\", \"description\"]\r\n            }\r\n          },\r\n          \"required\": [\"name\", \"description\"]\r\n        }\r\n      },\r\n      \"required\": [\"name\", \"description\"]\r\n    }\r\n  },\r\n  \"required\": [\"name\", \"description\"]\r\n}"