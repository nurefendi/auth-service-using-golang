package structhelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

func StructToQueryParams(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	var jsonObj map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
		return "", err
	}

	values := url.Values{}
	for key, value := range jsonObj {
		switch v := value.(type) {
		case int, int8, int16, int32, int64:
			values.Set(key, strconv.Itoa(v.(int)))
		case float64:
			values.Set(key, strconv.FormatFloat(v, 'f', -1, 64))

		case string:
			values.Set(key, v)
		}
	}

	return values.Encode(), nil
}

func Map(origin interface{}, target interface{}) error {
	temp := structs.Map(origin)
	err := convertMapToStruct(temp, target)
	return err
}

func convertMapToStruct(m map[string]interface{}, s interface{}) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		if value, ok := m[field.Name]; ok {
			if value != nil {
				stValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func ValidateStruct(s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field '%s' failed validation: %s", err.Field(), err.Tag())
		}
	}
	return nil
}
