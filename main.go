package main

import (
	"ApiGin/controllers"
	"ApiGin/initializers"
	"ApiGin/middleware"
	"ApiGin/models"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Booking contains binded and validated data.
type Booking struct {
	CheckIn  time.Time `json:"check_in" binding:"required,bookabledate"`
	CheckOut time.Time `json:"check_out" binding:"required,gtfield=CheckIn"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)

	today := time.Now()
	if ok {

		if date.After(today) {
			fmt.Println("Dentro del IF")
			return false
		}
	}
	return true
}

func init() {
	initializers.LoadEnvVars()
	initializers.DBConnection()
}

func main() {
	router := gin.Default()
	router.Delims("{{", "}}")
	router.LoadHTMLGlob("./templates/*")
	//Register CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"},
		AllowHeaders:  []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header", "X-Requested-With", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}
	//Register log with middleware
	file, _ := os.OpenFile(os.Getenv("LOG_FILE"), os.O_RDWR|os.O_CREATE, 0755)
	gin.DefaultWriter = io.MultiWriter(file)
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("ClientIP: %s - [%s] Method: %s  StatusCode:%d  Path:%s  Latency:%s  %s  %s \n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.StatusCode,
			param.Path,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	// GIN ITEMS ROUTES
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
	//AUTH ROUTES
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
	}
	// API ROUTES
	api := router.Group("/api")
	{
		api.GET("/employee", middleware.Auth(false), controllers.GetAllemployees)
		api.GET("/employeewithskills", middleware.Auth(false), controllers.GetAllemployeesWithSkills)
		api.GET("/employee/:employeeId", middleware.Auth(false), controllers.GetAnEmployee)
		api.DELETE("/employee/:employeeId", middleware.Auth(true), controllers.DeleteEmployee)
		api.POST("/employee", middleware.Auth(true), controllers.CreateEmployee)
		api.PUT("/employee/:employeeId", middleware.Auth(true), controllers.UpdateEmployee)
		api.POST("/employee/skill", middleware.Auth(true), controllers.AddSkillemployee)
		api.GET("/employee/skill/:employeeId", middleware.Auth(false), controllers.GetSkillsByEmployeeId)

		api.GET("/employeesbyskill/:skillsId", middleware.Auth(false), controllers.GetEmployeesBySkillID)
		api.GET("/employeesbylevel/:levelId", middleware.Auth(false), controllers.GetEmployeesByLevelId)
		api.DELETE("/employee_skill/:skillID/:employeeID", middleware.Auth(true), controllers.DeleteSkillemployee)

		api.GET("/skill", middleware.Auth(false), controllers.GetSkills)
		api.POST("/skill", middleware.Auth(true), controllers.PostSkill)
		api.DELETE("/skill/:skillId", middleware.Auth(true), controllers.DeleteSkill)

		api.POST("/expertise", middleware.Auth(true), controllers.PostLevelRating)
		api.GET("/expertise", middleware.Auth(false), controllers.GetAllLevelRating)
		api.DELETE("expertise/:levelId", middleware.Auth(true), controllers.DeleteLevelRating)
	}
	router.GET("/bookable", getBookable)
	router.Run()

}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindJSON(&b); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
