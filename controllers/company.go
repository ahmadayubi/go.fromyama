package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

const getAllEmployeesSql = "SELECT email, name, is_approved, id FROM employees INNER JOIN users u on employees.user_id = u.id WHERE u.company_id =  $1"
const registerCompanySql = "INSERT INTO companies (company_name, street, city, province_code, country, postal_code, phone) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
const setCompanyHeadSql = "UPDATE companies SET head_id = $1 WHERE id = $2"
const approveEmployeeSql = "UPDATE users SET is_approved = true WHERE id = $1"
const isEmployeeApprovedSql = "SELECT is_approved, company_id FROM users WHERE id = $1"
const getCompanyDetailSql = "SELECT company_name, head_id, total_due, street, city, province_code, country, postal_code, phone FROM companies WHERE id = $1"
const addPaymentMethodSql = "UPDATE companies SET payment_account_id = $1 WHERE id = $2"
const getPaymentMethodSql = "SELECT payment_account_id FROM companies WHERE id = $1"

type Employee struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Approved bool `json:"is_approved"`
	ID string `json:"id"`
}

type Company struct {
	CompanyName string `json:"company_name"`
	HeadID string `json:"head_id"`
	TotalDue float64 `json:"total_due"`
	Street string `json:"street"`
	City string `json:"city"`
	ProvinceCode string `json:"province_code"`
	Country string `json:"country"`
	PostalCode string `json:"postal_code"`
	Phone int `json:"phone"`
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
		registerCompanyQuery, err := database.DB.Prepare(registerCompanySql)
		defer registerCompanyQuery.Close()
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		var companyID string
		err = registerCompanyQuery.QueryRow(body["company_name"], body["street"], body["city"], body["province_code"], body["country"], body["postal_code"], body["phone"]).Scan(&companyID)

		userID, token, err := SignUpUserAndGenerateToken(w, true, true, body["name"], body["email"], body["password"], companyID)

		setHeadQuery, err := database.DB.Prepare(setCompanyHeadSql)
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

func ApproveEmployee (w http.ResponseWriter, r *http.Request){
	employeeID := chi.URLParam(r, "employeeID")
	approveEmployeeQuery, err := database.DB.Prepare(approveEmployeeSql)
	defer approveEmployeeQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	_, err = approveEmployeeQuery.Exec(employeeID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusAccepted, "Employee Approved")
}

func IsEmployeeApproved (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	employeeID := chi.URLParam(r, "employeeID")
	isApprovedQuery, err := database.DB.Prepare(isEmployeeApprovedSql)
	defer isApprovedQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	var approved bool
	var companyID string
	err = isApprovedQuery.QueryRow(employeeID).Scan(&approved, &companyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	if tokenClaims.CompanyID != companyID {
		utils.JSONResponse(w, http.StatusForbidden, "Employee Not Registered To Your Company")
		return
	}
	utils.JSONResponse(w, http.StatusOK, approved)
}

func GetCompanyDetails (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	companyDetailsQuery, err := database.DB.Prepare(getCompanyDetailSql)
	defer companyDetailsQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	var c Company
	err = companyDetailsQuery.QueryRow(tokenClaims.CompanyID).Scan(&c.CompanyName, &c.HeadID, &c.TotalDue, &c.Street, &c.City, &c.ProvinceCode, &c.Country, &c.PostalCode, &c.Phone)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, c)
}

func AddPaymentMethod (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}
	if body["payment_token"] == ""{
		utils.JSONResponse(w, http.StatusBadRequest, "Missing Payment Token")
		return
	}

	pToken := body["payment_token"]
	source := stripe.SourceParams{
		Token: &pToken,
	}

	stripe.Key = os.Getenv("SHOP_API_KEY")
	params := &stripe.CustomerParams{
		Source: &source,
		Description: stripe.String(tokenClaims.UserID),
		Email: stripe.String(tokenClaims.Email),
	}

	c, _ := customer.New(params)
	addPaymentQuery, err := database.DB.Prepare(addPaymentMethodSql)
	defer addPaymentQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	_, err = addPaymentQuery.Exec(c.ID, tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusAccepted, "Payment Account Added")
}

func ChargePaymentAccount (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body)
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}
	if body["amount"] == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "Missing Amount ")
		return
	}
	stripe.Key = os.Getenv("STRIPE_KEY")

	var paymentID string
	getPaymentMethodQuery, err := database.DB.Prepare(getPaymentMethodSql)
	defer getPaymentMethodQuery.Close()
	err = getPaymentMethodQuery.QueryRow(tokenClaims.CompanyID).Scan(&paymentID)
	amount, err := strconv.ParseInt(body["amount"],0, 64)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	params := &stripe.ChargeParams{
		Amount: stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Label Purchase"),
		Customer: stripe.String(paymentID),
	}
	c, _ := charge.New(params)
	if !c.Paid {
		utils.JSONResponse(w, http.StatusConflict, "Payment Error")
	}
	utils.JSONResponse(w, http.StatusAccepted, "Payment Successful")
}