package http_request

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetParameterValidate[T any](bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T
		if err := c.ShouldBindQuery(&body); err != nil {
			ValidationRender(c, err, body, bundle)
		}
		c.Next()
	}
}

func BodyJsonValidate[T any](bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		b := bytes.NewBuffer(make([]byte, 0))
		_ = io.TeeReader(c.Request.Body, b)

		var body T
		if err := c.ShouldBindJSON(&body); err != nil {
			ValidationRender(c, err, body, bundle)
		}
		c.Request.Body = ioutil.NopCloser(b)
		c.Next()
	}
}

func BodyFormValidate[T any](bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T
		if err := c.ShouldBind(&body); err != nil {
			ValidationRender(c, err, body, bundle)
		}
		c.Next()
	}
}

func UriValidate[T any](bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T
		if err := c.ShouldBindUri(&body); err != nil {
			ValidationRender(c, err, body, bundle)
		}
		c.Next()
	}
}
