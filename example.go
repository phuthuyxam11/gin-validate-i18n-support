package example

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/phuthuyxam11/gin-validate-i18n-support/http_request"
	"github.com/phuthuyxam11/gin-validate-i18n-support/utils"
)

var rootDir, _ = os.Getwd()

/*
	message_key is used to define message validation with the key located in the language file.
	the message order corresponds to the order of the elements inside the tag binding.
*/

type request struct {
	Field1 string `uri:"field1" url:"field1" form:"field1" json:"field_1" binding:"required" message_key:"filed_1"`
	Field2 string `uri:"field2" url:"field2" form:"field2" json:"field_2" binding:"required,len=2" message_key:"filed_1, field_2.field_2_1"`
}

func HttpSuccessResponse() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, "OK")
	}
}

func main() {
	// register i18nSupport
	i18nBundle, err := utils.I18nInit(rootDir + "/lang_example")
	if err != nil {
		log.Fatalln(err)
	}
	// routes
	r := gin.Default()
	r.GET("/get-example", http_request.GetParameterValidate[request](i18nBundle), HttpSuccessResponse())
	// json body
	r.POST("/get-example", http_request.BodyJsonValidate[request](i18nBundle), HttpSuccessResponse())
	// form body
	r.GET("/get-example", http_request.BodyFormValidate[request](i18nBundle), HttpSuccessResponse())
	// uri validate
	r.GET("/get-example", http_request.UriValidate[request](i18nBundle), HttpSuccessResponse())
}
