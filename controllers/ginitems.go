package controllers

import (
	"ApiGin/initializers"
	"ApiGin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllemployeesXML(cgin *gin.Context) {
	var employees []models.Employee
	initializers.DB.Find(&employees)
	cgin.XML(http.StatusOK, employees[0])
}
