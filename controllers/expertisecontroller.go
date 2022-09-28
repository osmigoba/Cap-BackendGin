package controllers

import (
	"ApiGin/initializers"
	"ApiGin/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostLevelRating(cgin *gin.Context) {

	var level models.Level
	if err := cgin.ShouldBindJSON(&level); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	levelAdded := initializers.DB.Create(&level)
	err := levelAdded.Error
	if err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Level-Rating Added to Catalogue")
	log.Println("Skill Added to Catalogue")
	cgin.JSON(http.StatusCreated, gin.H{"message": "Level-Expertise Created Successfully in the catalogue"})
}

func GetAllLevelRating(cgin *gin.Context) {

	var levels []models.Level
	initializers.DB.Find(&levels)
	cgin.JSON(http.StatusOK, levels)
}

func DeleteLevelRating(cgin *gin.Context) {

	levelId := cgin.Param("levelId")
	var level models.Level

	initializers.DB.First(&level, levelId)
	if level.Id == 0 {
		//Register Not Found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Level-Expertise Not Found"})
		return
	}
	initializers.DB.Delete(&level)
	cgin.JSON(http.StatusCreated, gin.H{"message": "Level-Expertise Deleted Successfully from the catalogue"})
}
