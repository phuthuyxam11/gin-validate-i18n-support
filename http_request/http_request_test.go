package http_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/go-querystring/query"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type request struct {
	Field1 string `uri:"field1" url:"field1" form:"field1" json:"field_1" binding:"required" message_key:"validate.verifyAcc.code.require"`
	Field2 string `uri:"field2" url:"field2" form:"field2" json:"field_2" binding:"required,len=2" message_key:"validate.verifyAcc.code.require"`
}

var bundle = i18n.NewBundle(language.English)

var tests = []struct {
	name                string
	request             request
	expectingStatusCode int
}{
	{
		name: "should pass validate",
		request: request{
			Field1: "field 1 has value",
			Field2: "ph",
		},
		expectingStatusCode: 200,
	},
	{
		name: "should not pass validate",
		request: request{
			Field1: "field1",
			Field2: "field2",
		},
		expectingStatusCode: 400,
	},
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func HttpSuccessResponse() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, "OK")
	}
}

func Test_GetValidate(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := SetUpRouter()
			r.GET("/", GetParameterValidate[request](bundle), HttpSuccessResponse())
			queryString, _ := query.Values(tc.request)
			req, _ := http.NewRequest("GET", "/?"+queryString.Encode(), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, tc.expectingStatusCode)
		})
	}
}

func Test_PostJsonValidate(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := SetUpRouter()
			r.POST("/", BodyJsonValidate[request](bundle), HttpSuccessResponse())
			jsonValue, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, tc.expectingStatusCode)
		})
	}
}

func Test_PostFormValidate(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := SetUpRouter()
			r.POST("/", BodyFormValidate[request](bundle), HttpSuccessResponse())
			formRequest := url.Values{}
			formRequest.Add("field1", tc.request.Field1)
			formRequest.Add("field2", tc.request.Field2)
			req, err := http.NewRequest("POST", "/", strings.NewReader(formRequest.Encode()))
			if err != nil {
				t.Errorf("got error: %s", err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, tc.expectingStatusCode)
		})
	}
}

func Test_UriValidate(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := SetUpRouter()
			r.GET("/:field1/:field2", UriValidate[request](bundle), HttpSuccessResponse())
			fmt.Println("/" + url.PathEscape(tc.request.Field1) + "/" + url.PathEscape(tc.request.Field2))
			req, err := http.NewRequest("GET", "/"+url.PathEscape(tc.request.Field1)+"/"+url.PathEscape(tc.request.Field2), nil)
			if err != nil {
				t.Errorf("got error: %s", err)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			fmt.Println(w.Code)
			assert.Equal(t, w.Code, tc.expectingStatusCode)
		})
	}
}
