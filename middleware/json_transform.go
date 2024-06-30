package middleware

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
)

func ToSnakeCaseMiddleware(c *fiber.Ctx) error {
	// Call the next handler
	if err := c.Next(); err != nil {
		return err
	}

	// Get the response data
	resData := c.Response().Body()
	var jsonData interface{}

	// Unmarshal the response data into a map
	if err := c.App().Config().JSONDecoder(resData, &jsonData); err != nil {
		return err
	}

	// Transform the keys to snake_case
	snakeCaseData := transformKeysToSnakeCase(jsonData)

	// Marshal the transformed data back to JSON
	newResData, err := c.App().Config().JSONEncoder(snakeCaseData)
	if err != nil {
		return err
	}

	// Set the new response data
	c.Response().SetBody(newResData)
	return nil
}

func transformKeysToSnakeCase(data interface{}) interface{} {
	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.Map:
		snakeCaseMap := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			snakeCaseKey := strcase.ToSnake(key.String())
			snakeCaseMap[snakeCaseKey] = transformKeysToSnakeCase(v.MapIndex(key).Interface())
		}
		return snakeCaseMap
	case reflect.Slice:
		slice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			slice[i] = transformKeysToSnakeCase(v.Index(i).Interface())
		}
		return slice
	default:
		return data
	}
}
