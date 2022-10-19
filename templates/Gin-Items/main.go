package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Func
func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	weekDay := t.Weekday()
	log.Println("Esto es lo que hago: ", fmt.Sprintf("%d%02d/%02d", year, month, day))
	return fmt.Sprintf("%d/%02d/%02d - %v", year, month, day, weekDay)
}

type ToValidate struct {
	Value1 int64 `json:"value1" binding:"required,checkvalue"`
	Value2 int64 `json:"value2" binding:"required,gtfield=Value1"`
	Value3 int64 `json:"value3" binding:"required,gtfield=Value2"`
}

var checkValue validator.Func = func(Value validator.FieldLevel) bool {
	value, _ := Value.Field().Interface().(int64)
	// cHECK if the value is between 10 AND 20
	if value >= 10 && value <= 20 {
		return true
	}
	return false
}

func main() {
	router := gin.Default()
	router.Delims("{{", "}}")
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLGlob("./templates/*")
	// Register validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("checkvalue", checkValue)
	}
	ginItems := router.Group("/ginitems")
	{
		// COOKIE GIN
		ginItems.GET("/cookie", func(c *gin.Context) {
			cookie, err := c.Cookie("session-id")
			if err != nil {
				cookie = "NotSet"
				c.SetCookie("session-id", "1234567", 36, "/", "localhost", false, true)
			}
			fmt.Printf("Cookie value: %s \n", cookie)
		})

		//Custom Template Funcs GIN
		ginItems.GET("/customTemplateFunc", func(c *gin.Context) {
			duration := time.Second * 5
			time.Sleep(duration)
			c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
				"title": "Template Funcs",
				"now":   time.Now(), //Date(2017, 12, 23, 10, 11, 12, 13, time.UTC),
			})
		})

		// Model binding and validation
		ginItems.POST("/validator", func(c *gin.Context) {
			var b ToValidate
			if err := c.ShouldBindJSON(&b); err == nil {
				c.JSON(http.StatusOK, gin.H{"message": "Values are valid!"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})

		//Serving data from reader
		ginItems.GET("/someDataFromReader", func(c *gin.Context) {
			response, err := http.Get("https://www.tutorialspoint.com/go/go_tutorial.pdf")
			//response, err := http.Get("http://localhost:4321/ginitem/static/file.pdf")
			if err != nil || response.StatusCode != http.StatusOK {
				c.Status(http.StatusServiceUnavailable)
				return
			}

			reader := response.Body
			contentLength := response.ContentLength
			contentType := response.Header.Get("Content-Type")
			fmt.Println(response.Header)
			extraHeaders := map[string]string{
				//"Content-Disposition": `attachment; filename="go_tutorial.pdf"`,
				"Content-Disposition": `inline; filename="go_tutorial.pdf"`,
			}

			c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
		})
	}

	s := &http.Server{
		Addr:           ":8085",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   6 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
