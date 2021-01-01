package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"
	"../utils/response"
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
const addParcelSql = "INSERT INTO parcel_options (company_id, length, width, height, name, weight) VALUES ($1, $2, $3, $4, $5, $6)"
const getParcelSql = "SELECT id, name, length, width, height, weight FROM parcel_options WHERE company_id = $1"
const getShippingInfoSql = "SELECT company_name, street, city, province_code, country, postal_code, phone FROM companies where id = $1"
const isEmployeeHeadSql = "SELECT is_head FROM users WHERE id = $1"
const unregisterCompanySql = "DELETE FROM companies WHERE id = $1"
const getPlatformsSql = "SELECT amazon_auth_token, etsy_token_secret, shopify_token FROM companies WHERE id = $1"

type Employee struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Approved bool `json:"is_approved"`
	ID string `json:"id"`
}

type Address struct {
	ShipperName string `json:"shipper"`
	Street string `json:"street"`
	City string `json:"city"`
	ProvinceCode string `json:"province_code"`
	Country string `json:"country"`
	PostalCode string `json:"postal_code"`
	Phone int `json:"phone"`
}

type Company struct {
	CompanyName string `json:"company_name"`
	HeadID string `json:"head_id"`
	TotalDue float64 `json:"total_due"`
	Address Address `json:"address"`
}

type Parcel struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Length float64 `json:"length"`
	Width float64 `json:"width"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
}

type Shipper struct {
	Address Address `json:"address"`
	Parcels []Parcel `json:"parcels"`
}

type approvedResponse struct {
	Approved bool `json:"approved"`
}

type platformResponse struct {
	AmazonConnected bool `json:"amazon_connected"`
	EtsyConnected bool `json:"etsy_connected"`
	ShopifyConnected bool `json:"shopify_connected"`
}

func GetConnectedPlatforms (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	getPlatformQuery, err := database.DB.Prepare(getPlatformsSql)
	defer getPlatformQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Get Platforms Error")
		return
	}
	var amazon, etsy, shopify string
	err = getPlatformQuery.QueryRow(tokenClaims.CompanyID).Scan(&amazon, &etsy, &shopify)
	if err != nil {
		utils.ErrorResponse(w, "Get Platforms Error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, platformResponse{
		AmazonConnected: amazon != "",
		EtsyConnected: etsy != "",
		ShopifyConnected: shopify != "",
	})
}

// GetAllEmployeeDetails /employee/all returns array of all employees registered to company and the details
func GetAllEmployeeDetails (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	allEmployeeQuery, err := database.DB.Prepare(getAllEmployeesSql)
	defer allEmployeeQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Get Employee Data Error")
		return
	}
	rows, err := allEmployeeQuery.Query(tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, "Get Employee Data Error")
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

// RegisterCompany /register registers a company and user
// request body has email, name, company_name, password, street, city, province_code, country, postal_code, phone
func RegisterCompany (w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"email", "name", "company_name", "password", "street",
		"city", "province_code", "country", "postal_code", "phone"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, "+err.Error())
		return
	}

	takenQuery, err := database.DB.Prepare("SELECT exists(SELECT 1 from users where email=$1)")
	defer takenQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Check Existing Email Error")
		return
	}
	var taken bool
	err = takenQuery.QueryRow(body["email"]).Scan(&taken)
	if !taken {
		registerCompanyQuery, err := database.DB.Prepare(registerCompanySql)
		defer registerCompanyQuery.Close()
		if err != nil {
			utils.ErrorResponse(w, "Register Company Error")
			return
		}
		var companyID string
		err = registerCompanyQuery.QueryRow(body["company_name"], body["street"], body["city"], body["province_code"], body["country"], body["postal_code"], body["phone"]).Scan(&companyID)

		userID, token, err := SignUpUserAndGenerateToken(w, true, true, body["name"], body["email"], body["password"], companyID)

		setHeadQuery, err := database.DB.Prepare(setCompanyHeadSql)
		defer setHeadQuery.Close()
		_, err = setHeadQuery.Exec(userID, companyID)
		if err != nil {
			utils.ErrorResponse(w, "Register Company Error")
			return
		}

		utils.JSONResponse(w, http.StatusCreated, token)
		return
	}
	utils.JSONResponse(w, http.StatusAlreadyReported, response.BasicMessage{Message: "Email Already In Use"})
}

// ApproveEmployee /employee/approve/{employeeID} approves employee
// request url has employeeID
func ApproveEmployee (w http.ResponseWriter, r *http.Request){
	employeeID := chi.URLParam(r, "employeeID")
	approveEmployeeQuery, err := database.DB.Prepare(approveEmployeeSql)
	defer approveEmployeeQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Approve Employee Error")
		return
	}
	_, err = approveEmployeeQuery.Exec(employeeID)
	if err != nil {
		utils.ErrorResponse(w, "Approve Employee Error")
		return
	}
	utils.JSONResponse(w, http.StatusAccepted, response.BasicMessage{Message: "Employee Approved"})
}

// IsEmployeeApproved /employee/approved/{employeeID} returns if the employee is approved
// request url has employeeID
func IsEmployeeApproved (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	employeeID := chi.URLParam(r, "employeeID")
	isApprovedQuery, err := database.DB.Prepare(isEmployeeApprovedSql)
	defer isApprovedQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Employee Approve Check Error")
		return
	}
	var approved bool
	var companyID string
	err = isApprovedQuery.QueryRow(employeeID).Scan(&approved, &companyID)
	if err != nil {
		utils.ErrorResponse(w, "Employee Approve Check Error")
		return
	}
	if tokenClaims.CompanyID != companyID {
		utils.JSONResponse(w, http.StatusForbidden, response.BasicMessage{Message:"Employee Not Registered To Your Company"})
		return
	}
	utils.JSONResponse(w, http.StatusOK, approvedResponse{Approved: approved})
}

// GetCompanyDetails /details returns the company details, a Company struct
func GetCompanyDetails (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	companyDetailsQuery, err := database.DB.Prepare(getCompanyDetailSql)
	defer companyDetailsQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Get Company Details Error")
		return
	}

	var c Company
	err = companyDetailsQuery.QueryRow(tokenClaims.CompanyID).Scan(&c.CompanyName, &c.HeadID,
		&c.TotalDue, &c.Address.Street, &c.Address.City, &c.Address.ProvinceCode,
		&c.Address.Country, &c.Address.PostalCode, &c.Address.Phone)
	if err != nil {
		utils.ErrorResponse(w, "Get Company Details Error")
		return
	}
	utils.JSONResponse(w, http.StatusOK, c)
}

// AddPaymentMethod /add/payment/method adds Stripe payment method to
// company account
// request body has payment_token,
func AddPaymentMethod (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body, []string{"payment_token"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}

	pToken := body["payment_token"]
	source := stripe.SourceParams{
		Token: &pToken,
	}

	stripe.Key = os.Getenv("STRIPE_SECRET")
	params := &stripe.CustomerParams{
		Source: &source,
		Description: stripe.String(tokenClaims.UserID),
		Email: stripe.String(tokenClaims.Email),
	}

	c, _ := customer.New(params)
	addPaymentQuery, err := database.DB.Prepare(addPaymentMethodSql)
	defer addPaymentQuery.Close()
	if err != nil {
		utils.ErrorResponse(w, "Add Payment Error")
		return
	}
	_, err = addPaymentQuery.Exec(c.ID, tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, "Add Payment Error")
		return
	}
	utils.JSONResponse(w, http.StatusAccepted, response.BasicMessage{Message:"Payment Account Added"})
}

// ChargePaymentAccount /add/payment/charge charges the company account
// request body has amount
func ChargePaymentAccount (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body, []string{"amount"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}
	if body["amount"] == "" {
		utils.JSONResponse(w, http.StatusBadRequest, response.BasicMessage{Message:"Missing Amount "})
		return
	}
	stripe.Key = os.Getenv("STRIPE_SECRET")

	var paymentID string
	getPaymentMethodQuery, err := database.DB.Prepare(getPaymentMethodSql)
	defer getPaymentMethodQuery.Close()
	err = getPaymentMethodQuery.QueryRow(tokenClaims.CompanyID).Scan(&paymentID)
	amount, err := strconv.ParseInt(body["amount"],0, 64)
	if err != nil {
		utils.ErrorResponse(w, "Get Payment Account Error")
		return
	}

	params := &stripe.ChargeParams{
		Amount: stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyCAD)),
		Description: stripe.String("Company Charge"),
		Customer: stripe.String(paymentID),
	}
	c, _ := charge.New(params)
	if !c.Paid {
		utils.JSONResponse(w, http.StatusConflict, response.BasicMessage{Message:"Payment Error"})
	}
	utils.JSONResponse(w, http.StatusAccepted, response.BasicMessage{Message:"Payment Successful"})
}

// AddParcel /add/parcel adds parcel option to company parcels
// request body has length, width, height, name
func AddParcel (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body, []string{"length", "width", "height", "name", "weight"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}
	addParcelQuery, err := database.DB.Prepare(addParcelSql)
	defer addParcelQuery.Close()
	_, err = addParcelQuery.Exec(tokenClaims.CompanyID, body["length"], body["width"], body["height"], body["name"], body["weight"])
	if err != nil {
		utils.ErrorResponse(w, "Add Parcel Error")
		return
	}
	utils.JSONResponse(w, http.StatusCreated, response.BasicMessage{
		Message: "Parcel Added",
	})
}

// GetShipper /shipper returns the seller address and parcel options
func GetShipper (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	getParcelAndShippingQuery, err := database.DB.Prepare(getParcelSql)
	defer getParcelAndShippingQuery.Close()
	rows, err := getParcelAndShippingQuery.Query(tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, "Get Parcel Error")
		return
	}
	var allParcel []Parcel
	for rows.Next(){
		var p Parcel
		err = rows.Scan(&p.ID, &p.Name, &p.Length, &p.Width, &p.Height, &p.Weight)
		allParcel = append(allParcel, p)
	}

	var a Address
	getShippingInfoQuery, err := database.DB.Prepare(getShippingInfoSql)
	defer getShippingInfoQuery.Close()
	err = getShippingInfoQuery.QueryRow(tokenClaims.CompanyID).Scan(&a.ShipperName, &a.Street, &a.City, &a.ProvinceCode, &a.Country, &a.PostalCode, &a.Phone)
	if err != nil {
		utils.ErrorResponse(w, "Get Shipper Error")
		return
	}
	shipper := Shipper{
		Address: a,
		Parcels: allParcel,
	}
	utils.JSONResponse(w, http.StatusOK, shipper)
}

// UnregisterCompany /unregister removes company from database
func UnregisterCompany (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var isHead bool
	isEmployeeHeadQuery, err := database.DB.Prepare(isEmployeeHeadSql)
	defer isEmployeeHeadQuery.Close()
	err = isEmployeeHeadQuery.QueryRow(tokenClaims.UserID).Scan(&isHead)
	if err != nil {
		utils.ErrorResponse(w, "Is Employee Head Error")
		return
	}
	if !isHead{
		utils.JSONResponse(w, http.StatusUnauthorized, response.BasicMessage{Message: "Need To Be Head To Unregister"})
		return
	}

	unregisterCompanyQuery, err := database.DB.Prepare(unregisterCompanySql)
	defer unregisterCompanyQuery.Close()
	_, err = unregisterCompanyQuery.Exec(tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, "Unregister Company Error")
		return
	}
	utils.JSONResponse(w, http.StatusAccepted, response.BasicMessage{Message: "Company Unregistered"})
}