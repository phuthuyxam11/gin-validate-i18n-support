package http_request

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
	"reflect"
	"strings"
)

// Name of the struct tag used in examples
const tagName = "message_key"

type fieldError struct {
	err        validator.FieldError
	messageKey string
}
type validateErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Key        string `json:"error_key"`
}

func (q fieldError) toI18nMessage(c *gin.Context, i18nBundle *i18n.Bundle) (string, error) {
	lang := c.Request.FormValue("lang")
	accept := c.Request.Header["Accept-Language"]
	localize := i18n.NewLocalizer(i18nBundle, language.English.String())
	if len(lang) > 0 || len(accept) > 0 {
		localize = i18n.NewLocalizer(i18nBundle, lang, accept[0])
	}

	// get from lang from request
	localizeValidateMessage := i18n.LocalizeConfig{
		MessageID: q.messageKey,
		TemplateData: map[string]string{
			"ErrorTag":             q.err.Tag(),
			"ErrorParam":           q.err.Param(),
			"ErrorField":           q.err.Field(),
			"ErrorStructField":     q.err.StructField(),
			"ErrorActualTag":       q.err.ActualTag(),
			"ValidateError":        q.err.Error(),
			"ErrorNamespace":       q.err.Namespace(),
			"ErrorStructNamespace": q.err.StructNamespace(),
		},
	}

	message, err := localize.Localize(&localizeValidateMessage)
	return message, err
}

func (q fieldError) toString(c *gin.Context, i18nBundle *i18n.Bundle) string {
	var sb strings.Builder
	// load message from message file

	message, err := q.toI18nMessage(c, i18nBundle)

	if len(message) > 0 && err == nil {
		return message
	}

	sb.WriteString("validation failed on field '" + q.err.Field() + "'")
	sb.WriteString(", condition: " + q.err.ActualTag())

	if q.err.Param() != "" {
		sb.WriteString(" { " + q.err.Param() + " }")
	}

	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
	}

	return sb.String()
}

func parseField(inputType reflect.StructField, tagName string) []string {
	// Get the field tag value
	tags := inputType.Tag.Get(tagName)
	tags = strings.Replace(tags, " ", "", -1)
	return strings.Split(tags, ",")
}

func ValidationRender[T any](c *gin.Context, err error, request T, i18nBundle *i18n.Bundle) {
	t := reflect.TypeOf(request)
	mapMessage := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagMessageMap := parseField(field, tagName)
		tagBindingMap := parseField(field, "binding")

		for i := 0; i < len(tagBindingMap); i++ {
			if i < len(tagMessageMap) {
				tag := strings.Split(tagBindingMap[i], "=")
				mapMessage[field.Name+"_"+tag[0]] = tagMessageMap[i]
			}
		}
	}

	for _, fieldErr := range err.(validator.ValidationErrors) {
		mess := fieldError{fieldErr, mapMessage[fieldErr.Field()+"_"+fieldErr.Tag()]}.toString(c, i18nBundle)
		errResponse := validateErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    mess,
			Key:        "validate_err",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errResponse)
		return
	}
}
