package helper_test

import (
	"github.com/nurefendi/auth-service-using-golang/tools/helper"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

func TestStructToQueryParams(t *testing.T) {
	data := TestStruct{Name: "John", Age: 25, Score: 85.5}
	query, err := helper.StructToQueryParams(data)
	assert.NoError(t, err, "Conversion should not return an error")
	parsedQuery, _ := url.ParseQuery(query)
	assert.Equal(t, "John", parsedQuery.Get("name"), "Name should match")
	assert.Equal(t, "25", parsedQuery.Get("age"), "Age should match")
	assert.Equal(t, "85.5", parsedQuery.Get("score"), "Score should match")
}

func TestMap(t *testing.T) {
	origin := TestStruct{Name: "Alice", Age: 30, Score: 90.0}
	var target TestStruct
	err := helper.Map(origin, &target)
	assert.NoError(t, err, "Mapping should not return an error")
	assert.Equal(t, origin, target, "Mapped struct should match the original")
}

type ValidationStruct struct {
	Username string `validate:"required"`
	Age      int    `validate:"min=18"`
}

func TestValidateStruct(t *testing.T) {
	validData := ValidationStruct{Username: "user1", Age: 20}
	invalidData := ValidationStruct{Username: "", Age: 16}

	assert.NoError(t, helper.ValidateStruct(validData), "Valid data should not return an error")
	err := helper.ValidateStruct(invalidData)
	assert.Error(t, err, "Invalid data should return an error")
}
