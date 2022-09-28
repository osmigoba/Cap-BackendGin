package main

import (
	"ApiGin/controllers"
	"ApiGin/initializers"
	"ApiGin/middleware"
	"ApiGin/models"
	"fmt"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	log.Println("Esto es lo que hago: ", fmt.Sprintf("%d%02d/%02d", year, month, day))
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func init() {
	initializers.LoadEnvVars()
	initializers.DBConnection()
}

func main() {
	router := gin.Default()
	router.Delims("{{", "}}")
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLGlob("./templates/*")
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"},
		AllowHeaders:  []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header", "X-Requested-With", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
	}

	ginItems := router.Group("/ginitem")
	{

		ginItems.GET("/htmlrender", func(c *gin.Context) {
			title := c.Query("title")
			descr := c.Query("description")
			imageuri := c.Query("imageuri")
			log.Println("Titulo: ", title)
			c.HTML(http.StatusOK, "home.html", gin.H{
				"title": title,
				"desc":  descr,
				"image": imageuri,
			})
		})
		ginItems.GET("/raw", func(c *gin.Context) {
			c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
				"now": time.Now(),
			})
		})

		ginItems.Static("/static", "./staticfiles")

		ginItems.GET("/redirect", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "https://www.capgemini.com/us-en/")
		})

		ginItems.GET("/employeexml", controllers.GetAllemployeesXML)
		ginItems.GET("/internalredirect", func(c *gin.Context) {
			c.Request.URL.Path = "/api/employee"
			router.HandleContext(c)
		})

		ginItems.GET("/datafromfile", func(c *gin.Context) {

			contents, err := os.ReadFile("test.txt")
			if err != nil {
				fmt.Println("File reading error", err)
				return
			}
			fmt.Println("Contents of file:", string(contents))

			reader := string(contents)
			reader2 := strings.NewReader(reader)
			contentLength := int64(len(contents))
			contentType := ".txt"

			extraHeaders := map[string]string{
				"fileName": "file.pdf",
			}

			c.DataFromReader(http.StatusOK, contentLength, contentType, reader2, extraHeaders)
		})

		ginItems.POST("/form_employee", func(c *gin.Context) {
			var employeeAPI models.Employee
			employeeAPI.FirstName = c.PostForm("firstName")
			employeeAPI.LastName = c.PostForm("lastName")

			employeeAPI.DOJ = time.Now()
			employeeAPI.Designation = c.PostForm("designation")
			employeeAPI.Email = c.PostForm("email")
			employeeAdded := initializers.DB.Create(&employeeAPI)

			if employeeAdded.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": employeeAdded.Error.Error()})

				return
			}

			c.JSON(http.StatusCreated, employeeAPI)
		})

		ginItems.POST("uploadfile", func(c *gin.Context) {
			root := "C:\\Users\\oscar\\OneDrive\\Documents\\Capgemini\\backendgin\\uploadedfiles\\"
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, root+filename); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			message := fmt.Sprintf("File  %s Uploaded Successfully in path %s", filename, root)
			c.JSON(http.StatusOK, gin.H{"message": message})
		})

	}
	api := router.Group("/api")
	{
		api.GET("/employee", middleware.Auth, controllers.GetAllemployees)
		api.GET("/employee/:employeeId", middleware.Auth, controllers.GetAnEmployee)
		api.DELETE("/employee/:employeeId", middleware.Auth, controllers.DeleteEmployee)
		api.POST("/employee", middleware.Auth, controllers.CreateEmployee)
		api.PUT("/employee/:employeeId", middleware.Auth, controllers.UpdateEmployee)
		api.POST("/employee/skill", middleware.Auth, controllers.AddSkillemployee)
		api.GET("/employee/skill/:employeeId", middleware.Auth, controllers.GetSkillsByEmployeeId)

		api.GET("/employeesbyskill/:skillsId", middleware.Auth, controllers.GetEmployeesBySkillID)
		api.GET("/employeesbylevel/:levelId", middleware.Auth, controllers.GetEmployeesByLevelId)
		api.DELETE("/employee_skill/:skillID/:employeeID", middleware.Auth, controllers.DeleteSkillemployee)

		api.GET("/skill", middleware.Auth, controllers.GetSkills)
		api.POST("/skill", middleware.Auth, controllers.PostSkill)
		api.DELETE("/skill/:skillId", middleware.Auth, controllers.DeleteSkill)

		api.POST("/expertise", middleware.Auth, controllers.PostLevelRating)
		api.GET("/expertise", middleware.Auth, controllers.GetAllLevelRating)
		api.DELETE("expertise/:levelId", middleware.Auth, controllers.DeleteLevelRating)
	}

	router.Run()
}