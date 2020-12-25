package controllers

import (
	"net/http"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email string `json:"email"`
	Name string	`json:"name"`
	CompanyID string `json:"company_id"`
	IsHead bool `json:"is_head"`
	ID string `json:"id"`
}

type Token struct {
	Raw string `json:"token"`
}

type UserDetails struct {
	UserData User `json:"user"`
	CompanyName string `json:"company_name"`
}


const createUserSql = "INSERT INTO users(name, email, password, company_id, is_approved, is_head) VALUES ($1, $2, $3, $4, $5, $6); INSERT INTO employees(company_id, user_id) VALUES ($4, (SELECT DISTINCT id FROM users WHERE email = $2)) RETURNING user_id;"
const addEmployee = "INSERT INTO employees(company_id, user_id) VALUES ($1, $2)"

const loginUserSql = "SELECT id, company_id, password, is_approved FROM users WHERE email = $1"

const getUserDetailSql = "SELECT c.id, u.email, u.name as name, c.company_name, u.id, u.is_head FROM (SELECT * FROM users WHERE email = $1) u INNER JOIN companies c on c.id = u.company_id"

func SignUpUser(w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}
	_, token, err := SignUpUserAndGenerateToken(w, false,false, body["name"], body["email"], body["password"], body["company_id"])
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, token)
}

func LoginUser(w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}

	query, err := database.DB.Prepare(loginUserSql)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}

	row := query.QueryRow(body["email"])
	var id, companyID, hash string
	var approved bool
	err = row.Scan(&id, &companyID, &hash, &approved)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	authChannel := make(chan bool)
	go ComparePassword(hash, body["password"], authChannel)

	tokenClaims := jwtUtil.TokenClaims{
		Email: body["email"],
		UserID: id,
		CompanyID: companyID,
		Approved: approved,
	}

	token, err := jwtUtil.GenToken(tokenClaims)
	if authorized := <-authChannel;err != nil || !(authorized) {
		if !authorized {
			utils.ForbiddenResponse(w)
		} else {
			utils.ErrorResponse(w, err)
		}
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, token)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request){
	query, err := database.DB.Prepare(getUserDetailSql)
	claims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}

	var q UserDetails
	if err = query.QueryRow(claims.Email).Scan(&q.UserData.CompanyID, &q.UserData.Email,&q.UserData.Name,&q.CompanyName, &q.UserData.ID, &q.UserData.IsHead); err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusAccepted, q)
}

func ComparePassword(hash, password string, c chan bool){
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		c <- false
	}
	c <- true
}

func SignUpUserAndGenerateToken(w http.ResponseWriter, approved, isHead bool,name, email, password, companyID string) (string, string, error) {
	userCreateQuery, err := database.DB.Prepare(createUserSql)
	if err != nil{
		utils.ErrorResponse(w, err)
		return "", "", err
	}
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		utils.ErrorResponse(w, err)
		return "", "", err
	}
	var userIdString string
	err = userCreateQuery.QueryRow(name, email, string(hashByte), companyID, approved, isHead).Scan(&userIdString)
	if err != nil {
		utils.ErrorResponse(w, err)
		return "", "", err
	}
	addEmployeeQuery, err := database.DB.Prepare(addEmployee)
	_, err = addEmployeeQuery.Exec(companyID, userIdString)
	if err != nil {
		utils.ErrorResponse(w, err)
		return "", "", err
	}

	tokenClaims := jwtUtil.TokenClaims{
		Email: email,
		UserID: userIdString,
		CompanyID: companyID,
		Approved: approved,
	}

	token, err := jwtUtil.GenToken(tokenClaims)
	if err != nil {
		utils.ErrorResponse(w, err)
		return "", "", err
	}
	return userIdString, token, nil
}