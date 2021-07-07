package message

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

const (
	eventPattern = `{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "data_id": {
      "type": "integer",
	  "minimum": 1
    },
    "access_token": {
      "type": "string"
    },
    "data": {
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "event_name": {
              "type": "string",
			  "minLength": 1
            },
            "event": {
              "type": "object",
              "properties": {
                "content": {
                  "type": "string",
			      "minLength": 1
                }
              },
              "required": [
                "content"
              ]
            },
            "target": {
              "type": "string",
			  "minLength": 1
            },
            "dimension": {
              "type": "object"
            },
            "timestamp": {
              "type": "integer"
            }
          },
          "required": [
            "event_name",
            "event",
            "target",
            "dimension",
            "timestamp"
          ]
        }
      ]
    }
  },
  "required": [
    "data_id",
    "access_token",
    "data"
  ]
}`

	timeSeriesPattern = `{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "data_id": {
      "type": "integer",
	  "minimum": 1
    },
    "access_token": {
      "type": "string"
    },
    "data": {
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "metrics": {
              "type": "object",
              "properties": {}
            },
            "target": {
              "type": "string",
			  "minLength": 1
            },
            "dimension": {
              "type": "object",
              "properties": {}
            },
            "timestamp": {
              "type": "integer"
            }
          },
          "required": [
            "metrics",
            "target",
            "dimension",
            "timestamp"
          ]
        }
      ]
    }
  },
  "required": [
    "data_id",
    "access_token",
    "data"
  ]
}`
)

var (
	eventSchema      *gojsonschema.Schema
	timeSeriesSchema *gojsonschema.Schema
)

func init() {
	eventSchema = mustLoadSchema(eventPattern)
	timeSeriesSchema = mustLoadSchema(timeSeriesPattern)
}

func mustLoadSchema(s string) *gojsonschema.Schema {
	loader := gojsonschema.NewStringLoader(s)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		panic(err)
	}

	return schema
}

func validateWithJSONSchema(schema *gojsonschema.Schema, content string) error {
	documentLoader := gojsonschema.NewStringLoader(content)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("schema: failed to decode schema, err: %v", err.Error())
	}

	if result.Valid() {
		return nil
	}

	var errMsg string
	for _, err := range result.Errors() {
		if err != nil {
			errMsg += fmt.Sprintf("%s\n", err.Description())
		}
	}

	if len(errMsg) > 0 {
		return errors.New(errMsg)
	}

	return nil
}

func ValidateSchema(content string) bool {
	return ValidateEventSchema(content) == nil || ValidateTimeSeriesSchema(content) == nil
}

func ValidateEventSchema(content string) error {
	return validateWithJSONSchema(eventSchema, content)
}

func ValidateTimeSeriesSchema(content string) error {
	return validateWithJSONSchema(timeSeriesSchema, content)
}
