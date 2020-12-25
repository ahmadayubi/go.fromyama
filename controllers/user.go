package controllers

import (
	"encoding/json"
	"io/ioutil"
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


const createUserSql = "INSERT INTO users(name, email, password, company_id) VALUES ($1, $2, $3, $4) RETURNING id;"
const addToCompanySql = "UPDATE companies SET employees = array_append(employees, $1) WHERE id = $2"

const loginUserSql = "SELECT id, company_id, password, is_approved FROM users WHERE email = $1"

const getUserDetailSql = "SELECT c.id, u.email, u.name as name, c.company_name, u.id, u.is_head FROM (SELECT * FROM users WHERE email = $1) u INNER JOIN companies c on c.id = u.company_id"

func SignUpUser(w http.ResponseWriter, r *http.Request){
	var body map[string]string

	reqBody, err := ioutil.ReadAll(r.Body)

	if err = json.Unmarshal(reqBody,&body);err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Database requests
	userCreateQuery, err := database.DB.Prepare(createUserSql)
	if err != nil{
		 utils.ErrorResponse(w, err)
		 return
	}
	hashByte, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	var userIdString string
	err = userCreateQuery.QueryRow(body["name"], body["email"], string(hashByte), body["company_id"]).Scan(&userIdString)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	addToCompanyQuery, err := database.DB.Prepare(addToCompanySql)
	_, err = addToCompanyQuery.Exec(userIdString, body["company_id"])
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	tokenClaims := jwtUtil.TokenClaims{
		Email: body["email"],
		UserID: userIdString,
		CompanyID: body["company_id"],
		Approved: false,
	}

	token, err := jwtUtil.GenToken(tokenClaims)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, Token{Raw: token})
}

func LoginUser(w http.ResponseWriter, r *http.Request){
	var body map[string]string
	reqBody, err := ioutil.ReadAll(r.Body)

	if err = json.Unmarshal(reqBody,&body);err != nil {
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

	utils.JSONResponse(w, http.StatusAccepted, Token{Raw: token})
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