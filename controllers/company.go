package controllers

import (
	"net/http"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"
)

const getAllEmployeesSql = "SELECT email, name, is_approved, id FROM employees INNER JOIN users u on employees.user_id = u.id WHERE u.company_id =  $1"
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
	defer allEmployeeQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	rows, err := allEmployeeQuery.Query(tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	var allEmployees []Employee
	for rows.Next(){
		var e Employee
		err = rows.Scan(&e.Email, &e.Name, &e.Approved, &e.ID)
		allEmployees = append(allEmployees, e)
	}

	utils.JSONResponse(w, http.StatusOK, allEmployees)
}

func RegisterCompany (w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}
	if body["email"] == "" || body["name"] == "" || body["company_name"] == "" || body["password"] == "" ||
		body["street"] == "" || body["city"] == "" || body["province_code"] == "" || body["country"] == "" ||
		body["postal_code"] == "" || body["phone"] == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "Missing Parameter(s)")
		return
	}

	takenQuery, err := database.DB.Prepare("SELECT exists(SELECT 1 from users where email=$1)")
	defer takenQuery.Close()
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

		userID, token, err := SignUpUserAndGenerateToken(w, true, true, body["name"], body["email"], body["password"], companyID)

		setHeadQuery, err := database.DB.Prepare(setCompanyHead)
		defer setHeadQuery.Close()
		_, err = setHeadQuery.Exec(userID, companyID)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusCreated, token)
		return
	}
	utils.JSONResponse(w, http.StatusAlreadyReported, "Email Already In Use")
}