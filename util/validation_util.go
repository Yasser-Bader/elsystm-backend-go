package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationErrorResponse struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func Validation(err interface{}) []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse
	for _, err := range err.(validator.ValidationErrors) {
		var element ValidationErrorResponse
		element.FailedField = strings.ToLower(err.Field())
		element.Tag = err.Tag()
		element.Value = err.Param()
		errors = append(errors, &element)
	}
	return errors
}

func ValidationReturnErrorResponse(c *gin.Context, field string, tag string, value string) {
	var errors []*ValidationErrorResponse
	var element ValidationErrorResponse
	element.FailedField = field
	element.Tag = tag
	element.Value = value
	errors = append(errors, &element)
	c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
}
