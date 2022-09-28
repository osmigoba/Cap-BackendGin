package controllers

import (
	"ApiGin/dtos"
	"ApiGin/initializers"
	"ApiGin/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllemployees(cgin *gin.Context) {

	var employees []models.Employee
	initializers.DB.Order("id asc").Find(&employees)
	cgin.JSON(http.StatusOK, employees)
}

func GetAnEmployee(cgin *gin.Context) {
	id := cgin.Param("employeeId")
	var employee models.Employee
	initializers.DB.First(&employee, id)
	if employee.Id == 0 {
		//Employee Not found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found"})
		return
	}
	cgin.JSON(http.StatusOK, employee)
}

func DeleteEmployee(cgin *gin.Context) {
	id := cgin.Param("employeeId")
	var employee models.Employee
	var employeeskills []models.EmployeeSkills
	initializers.DB.First(&employee, id)
	if employee.Id == 0 {
		//Employee Not found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found"})
		return
	}
	// Look for data related
	initializers.DB.Where("employee_id = ?", id).Find(&employeeskills)
	if len(employeeskills) != 0 {
		// Delete data related
		initializers.DB.Where("employee_id = ?", employee.Id).Delete(&models.EmployeeSkills{})
		initializers.DB.Delete(&employee)
		log.Println("Data related Deleted")
		cgin.JSON(http.StatusOK, gin.H{"message": "Employee Deleted and its data related"})
		return
	}
	initializers.DB.Delete(&employee)
	cgin.JSON(http.StatusOK, gin.H{"message": "Employee Deleted"})
}

func CreateEmployee(cgin *gin.Context) {
	var employee models.Employee
	if err := cgin.ShouldBindJSON(&employee); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employeeAdded := initializers.DB.Create(&employee)
	fmt.Println(&employeeAdded)
	err := employeeAdded.Error
	if err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	cgin.JSON(http.StatusCreated, &employee)
}

func UpdateEmployee(cgin *gin.Context) {
	id := cgin.Param("employeeId")
	var employee models.Employee
	initializers.DB.First(&employee, id)
	if employee.Id == 0 {
		//Register Not Found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found"})
		return
	}
	//Get the body with the parameters to update
	if err := cgin.ShouldBindJSON(&employee); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	initializers.DB.Save(&employee)
	cgin.JSON(http.StatusOK, employee)
}

func AddSkillemployee(cgin *gin.Context) {

	var dto dtos.AddSkillEmployeeDTO
	if err := cgin.ShouldBindJSON(&dto); err != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var employee models.Employee
	initializers.DB.First(&employee, dto.EmployeeID)
	if employee.Id == 0 {
		//Register Not found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Employee Not Found"})
		return
	}

	var skill models.Skill
	initializers.DB.First(&skill, dto.SkillID)
	if skill.Id == 0 {
		//Register Not found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Skill Not Found"})
		return
	}

	var level models.Level
	initializers.DB.First(&level, dto.LevelRatingId)
	if level.Id == 0 {
		//Register Not found
		cgin.JSON(http.StatusNotFound, gin.H{"error": "Level Not Found"})
		return
	}

	var checkExist models.EmployeeSkills
	initializers.DB.Where("employee_id = ? AND skill_id = ?", dto.EmployeeID, dto.SkillID).First(&checkExist)
	if checkExist.EmployeeID != 0 {
		//Skill Already Exist
		cgin.JSON(http.StatusBadRequest, gin.H{"error": "Skill Already Exist"})
		return
	}

	// Every Validation to database is OK
	var data models.EmployeeSkills
	data.EmployeeID = dto.EmployeeID
	data.SkillID = dto.SkillID
	data.LevelID = dto.LevelRatingId
	data.Experience = dto.Experience

	skillAdded := initializers.DB.Create(&data)

	err1 := skillAdded.Error

	if err1 != nil {
		cgin.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}

	cgin.JSON(http.StatusCreated, gin.H{"message": "Skill Created Successfully"})
}

func GetSkillsByEmployeeId(cgin *gin.Context) {

	id := cgin.Param("employeeId")
	var skillByEmployee []dtos.SkillsByEmployeeDTO
	initializers.DB.Table("employees").
		Select("skills.id as skillId, skills.skill, levels.id as levelId, levels.name as level, employee_skills.experience as experience").
		Joins("JOIN employee_skills ON employee_skills.employee_id = employees.id JOIN skills ON skills.id = employee_skills.skill_id JOIN levels ON levels.id = employee_skills.level_id").
		Where("employees.id = ?", id).
		Scan(&skillByEmployee)

	cgin.JSON(http.StatusOK, skillByEmployee)
}

func GetEmployeesBySkillID(cgin *gin.Context) {

	id := cgin.Param("skillsId")
	var employees []models.Employee
	initializers.DB.Table("employees").
		Select("employees.id, employees.first_name, employees.last_name, employees.doj, employees.designation, employees.email").
		Joins("JOIN employee_skills ON employee_skills.employee_id = employees.id").
		Where("employee_skills.skill_id = ?", id).
		Order("id asc").
		Scan(&employees)

	cgin.JSON(http.StatusOK, employees)
}

func GetEmployeesByLevelId(cgin *gin.Context) {

	id := cgin.Param("levelId")
	var employees []dtos.EmployeeLevel
	initializers.DB.Table("employees").
		Select("employees.id, employees.first_name, employees.last_name, employees.doj, employees.designation, employees.email, skills.skill, levels.name as levelName").
		Joins("JOIN employee_skills ON employee_skills.employee_id = employees.id JOIN skills ON skills.id = employee_skills.skill_id JOIN levels ON levels.id = employee_skills.level_id").
		Where("employee_skills.level_id = ?", id).
		Order("id asc").
		Scan(&employees)

	cgin.JSON(http.StatusOK, employees)
}
