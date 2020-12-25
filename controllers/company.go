package controllers

import (
	"net/http"
	"strings"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"
)

const getAllEmployeesSql = "SELECT employees FROM companies WHERE id = $1"
const getEmployeeSql = "SELECT email, name, is_approved FROM users WHERE id = $1"

type Employee struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Approved bool `json:"is_approved"`
	ID string `json:"id"`
}

func GetAllEmployeeDetails (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	allEmployeeQuery, err := database.DB.Prepare(getAllEmployeesSql)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	var employees string
	err = allEmployeeQuery.QueryRow(tokenClaims.CompanyID).Scan(&employees)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	employees = employees[:len(employees)-1]
	employees = employees[1:]
	employeesList := strings.Split(employees, ",")
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	allEmployees := make([]Employee, len(employeesList))
	for i := 0; i< len(employeesList); i++{
		employeeQuery, err := database.DB.Prepare(getEmployeeSql)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		allEmployees[i].ID = employeesList[i]
		err = employeeQuery.QueryRow(employeesList[i]).Scan(&allEmployees[i].Email, &allEmployees[i].Name, &allEmployees[i].Approved)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
	}
	utils.JSONResponse(w, http.StatusOK, allEmployees)
}

func RegisterCompany (w http.ResponseWriter, r *http.Request){

}