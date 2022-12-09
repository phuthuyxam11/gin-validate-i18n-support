# **gin-validate-i18n-support**

## **Description**

This is a small library that helps to generate validation middleware that integrates with gin-gonic. Support the creation of multilingual messages with i18n

## Contents

1. ### Register struct validate
```go
    /*
    message_key is used to define message validation with the key located in the language file.
	the message order corresponds to the order of the elements inside the tag binding 
    */
    type Request struct {
	    Field1 string `uri:"field1" url:"field1" form:"field1" json:"field_1" binding:"required" message_key:"filed_1"`
	    Field2 string `uri:"field2" url:"field2" form:"field2" json:"field_2" binding:"required,len=2" message_key:"filed_1, field_2.field_2_1"`
    }
```
2. ### Register i18nSupport
```go
    utils.I18nInit(rootDir + "/lang_example")
```
3. ### Add middleware
```go
    func main() {
        // register i18nSupport
        i18nBundle, err := utils.I18nInit(rootDir + "/lang_example")
        if err != nil {
            log.Fatalln(err)
        }
        // routes
        r := gin.Default()
        r.GET("/get-example", http_request.GetParameterValidate[Request](i18nBundle), HttpSuccessResponse())
        // json body
        r.POST("/get-example", http_request.BodyJsonValidate[Request](i18nBundle), HttpSuccessResponse())
        // form body
        r.GET("/get-example", http_request.BodyFormValidate[Request](i18nBundle), HttpSuccessResponse())
        // uri validate
        r.GET("/get-example", http_request.UriValidate[Request](i18nBundle), HttpSuccessResponse())
    }
```




 
