package controllers

import (
	"ApiGin/initializers"
	"ApiGin/models"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteSkillemployee(cgin *gin.Context) { //params [{skillId}/{employeeId}]

	employeeId := cgin.Param("employeeID")
	skillId := cgin.Param("skillID")
	var skill_employee models.EmployeeSkills
	initializers.DB.Where("employee_id = ? AND skill_id = ?", employeeId, skillId).First(&skill_employee)

	if skill_employee.EmployeeID == 0 {
		//Register Not Found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found"})
		return
	}

	err := initializers.DB.Delete(&models.EmployeeSkills{}, "employee_id = ? AND skill_id = ?", employeeId, skillId).Error
	if err != nil {
		//Error deleting
		cgin.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting the Register"})
		return
	}
	log.Println("Employee's skill Deleted")
	cgin.JSON(http.StatusOK, gin.H{"message": "Employee's skill Deleted"})
}

func GetSkills(cgin *gin.Context) {
	var skills []models.Skill
	initializers.DB.Find(&skills)
	cgin.JSON(http.StatusOK, skills)
}

func PostSkill(cgin *gin.Context) {

	var skill models.Skill
	skill.CreatedDate = time.Now()
	if err := cgin.ShouldBindJSON(&skill); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	skillAdded := initializers.DB.Create(&skill)
	err := skillAdded.Error
	if err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Skill Added to Catalogue")
	cgin.JSON(http.StatusCreated, gin.H{"message": "Skill Created Successfully in the catalogue"})
}

func DeleteSkill(cgin *gin.Context) {

	skillid := cgin.Param("skillId")
	var skill models.Skill
	initializers.DB.First(&skill, skillid)
	var employeeskills []models.EmployeeSkills
	if skill.Id == 0 {
		//Skill Not found
		cgin.JSON(http.StatusBadRequest, gin.H{"error": "Skill Not Found"})
		return
	}

	// Look for data related
	initializers.DB.Where("skill_id = ?", skillid).Find(&employeeskills)
	if len(employeeskills) != 0 {
		// Delete data related
		initializers.DB.Where("skill_id = ?", skillid).Delete(&models.EmployeeSkills{})
		initializers.DB.Delete(&skill)
		log.Println("Data related Deleted")
		cgin.JSON(http.StatusOK, gin.H{"message": "Skill Deleted and its data related"})
		return
	}
	initializers.DB.Delete(&skill)
	log.Println("Skill Deleted")
	cgin.JSON(http.StatusOK, gin.H{"message": "Skill Deleted"})
}
