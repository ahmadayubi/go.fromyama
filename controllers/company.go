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
const registerCompany = "INSERT INTO companies (company_name, street, city, province_code, country, postal_code, phone) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
const setCompanyHead = "UPDATE companies SET head_id = $1 WHERE id = $2"

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
	if len(employees) > 2 {
		employees = employees[:len(employees)-1]
		employees = employees[1:]
		employeesList := strings.Split(employees, ",")

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
		return
	}
	utils.JSONResponse(w, http.StatusOK, []Employee{})
}

func RegisterCompany (w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}

	takenQuery, err := database.DB.Prepare("SELECT exists(SELECT 1 from users where email=$1)")
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	var taken bool
	err = takenQuery.QueryRow(body["email"]).Scan(&taken)
	if !taken {
		registerCompanyQuery, err := database.DB.Prepare(registerCompany)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		var companyID string
		err = registerCompanyQuery.QueryRow(body["company_name"], body["street"], body["city"], body["province_code"], body["country"], body["postal_code"], body["phone"]).Scan(&companyID)

		userID, token, err := SignUpUserAndGenerateToken(w, true, true,body["name"], body["email"], body["password"], companyID)

		setHeadQuery, err := database.DB.Prepare(setCompanyHead)
		_, err = setHeadQuery.Exec(userID, companyID)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusAccepted, token)
		return
	}
	utils.JSONResponse(w, http.StatusAlreadyReported, "Email Already In Use")
}