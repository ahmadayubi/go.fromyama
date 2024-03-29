package controllers

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"go.fromyama/utils"
	"go.fromyama/utils/database"
	"go.fromyama/utils/jwtUtil"
	"go.fromyama/utils/response"

	"golang.org/x/crypto/bcrypt"
)

const createUserSql = "INSERT INTO users(name, email, password, company_id, is_approved, is_head) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
const addEmployeeSql = "INSERT INTO employees(company_id, user_id) VALUES ($1, $2)"
const loginUserSql = "SELECT id, company_id, password, is_approved FROM users WHERE email = $1"
const getUserDetailSql = "SELECT c.id, u.email, u.name as name, c.company_name, u.id, u.is_head, u.is_approved FROM (SELECT * FROM users WHERE email = $1) u INNER JOIN companies c on c.id = u.company_id"

type userError struct {
	s string
}
func (e *userError) Error() string{
	return e.s
}

// SignUpUser /signup signs user up and returns jwt token
// request body has email, name, company_id, password
func SignUpUser(w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body, []string{"email", "name", "company_id", "password"})
	if err != nil{
		response.Error(w, "Body Parse Error, " + err.Error())
		return
	}

	_, token, err := SignUpUserAndGenerateToken(w, false,false, body["name"], body["email"], body["password"], body["company_id"])
	if err != nil {
		return
	}

	response.JSON(w, http.StatusCreated, *token)
}

// LoginUser /login checks user email and password and returns jwt token
// request body has email, password
func LoginUser(w http.ResponseWriter, r *http.Request){
	var body map[string]string
	err := utils.ParseRequestBody(r, &body, []string{"email", "password"})
	if err != nil{
		response.Error(w, "Body Parse Error, " + err.Error())
		return
	}

	query, err := database.DB.Prepare(loginUserSql)
	defer query.Close()
	if err != nil{
		response.Error(w, "User Credential Fetch Error")
		return
	}

	row := query.QueryRow(body["email"])
	var id, companyID, hash string
	var approved bool
	err = row.Scan(&id, &companyID, &hash, &approved)
	if err != nil {
		response.Error(w, "User Credential Fetch Error")
		return
	}

	authChannel := make(chan bool)

	go func(hash, pass string, c chan bool) {
		c <- bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
	}(hash, body["password"], authChannel)

	tokenClaims := jwtUtil.TokenClaims{
		Email: body["email"],
		UserID: id,
		CompanyID: companyID,
		Approved: approved,
	}

	token, err := jwtUtil.GenerateToken(tokenClaims)
	if !(<-authChannel) {
		response.Forbidden(w)
	} else if err != nil {
		response.Error(w, "Token Generation Error")
	} else {
		response.JSON(w, http.StatusAccepted, response.Token{
			Raw: token,
		})
	}
}

// GetUserDetails /details returns the users details
func GetUserDetails(w http.ResponseWriter, r *http.Request){
	query, err := database.DB.Prepare(getUserDetailSql)
	defer query.Close()
	claims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	if err != nil{
		response.Error(w, "Body Parse Error, " + err.Error())
		return
	}

	var q response.UserDetails
	if err = query.QueryRow(claims.Email).Scan(&q.UserData.CompanyID, &q.UserData.Email,&q.UserData.Name,&q.CompanyName, &q.UserData.ID, &q.UserData.IsHead, &q.UserData.IsApproved); err != nil {
		response.Error(w, "Company Details Fetch Error")
		return
	}
	response.JSON(w, http.StatusAccepted, q)
}

// SignUpUserAndGenerateToken util function that signs up user and generates the jwt token
func SignUpUserAndGenerateToken(w http.ResponseWriter, approved, isHead bool,name, email, password, companyID string) (string, *response.Token, error) {

	takenQuery, err := database.DB.Prepare("SELECT exists(SELECT 1 from users where email=$1)")
	defer takenQuery.Close()
	if err != nil {
		response.Error(w, "Check Existing Email Error")
		return "", nil, err
	}
	var taken bool
	err = takenQuery.QueryRow(email).Scan(&taken)
	if taken{
		response.JSON(w, http.StatusConflict, response.BasicMessage{
			Message: "Email Already In Use",
		})
		return "", nil, &userError{"Email Already In Use"}
	}

	userCreateQuery, err := database.DB.Prepare(createUserSql)
	defer userCreateQuery.Close()
	if err != nil{
		response.JSON(w, http.StatusConflict, response.BasicMessage{
			Message: "Creating User Error",
		})
		return "", nil, err
	}
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		response.Error(w, "Password Hash Error")
		return "", nil, err
	}
	var userIdString string
	err = userCreateQuery.QueryRow(name, email, string(hashByte), companyID, approved, isHead).Scan(&userIdString)
	if err != nil {
		response.Error(w, "Create User Error")
		return "", nil, err
	}
	addEmployeeQuery, err := database.DB.Prepare(addEmployeeSql)
	defer addEmployeeQuery.Close()
	_, err = addEmployeeQuery.Exec(companyID, userIdString)
	if err != nil {
		response.Error(w, "Adding Employee Error")
		return "", nil, err
	}

	tokenClaims := jwtUtil.TokenClaims{
		Email: email,
		UserID: userIdString,
		CompanyID: companyID,
		Approved: approved,
	}

	token, err := jwtUtil.GenerateToken(tokenClaims)
	if err != nil {
		response.Error(w, "Token Generation Error")
		return "", nil, err
	}

	newUserHTML, err := ioutil.ReadFile("assets/templates/newUser.html")
	err = SendEmail(email, "Welcome To FromYama", string(newUserHTML), nil)
	if err != nil {
		response.Error(w, "Sending Email Error")
		return "", nil, err
	}

	return userIdString, &response.Token{
		Raw: token,
	}, nil
}

// RefreshToken /user/refresh refreshes the token without needing to check password
func RefreshToken(w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	token, err := jwtUtil.GenerateToken(tokenClaims)
	if err != nil {
		response.Error(w, "Token Generation Error")
		return
	}
	response.JSON(w, http.StatusAccepted, response.Token{Raw: token})
}

// SendEmail util function that sends email to user
func SendEmail(toEmail, subject, body string, attachment []byte) error {
	fromString := os.Getenv("MAIL_USER")
	from := mail.Address{Name: "FromYama",Address: fromString}
	to := mail.Address{Name: "", Address: toEmail}


	serverName := "smtppro.zoho.com:465"
	host := "smtppro.zoho.com"

	emailAuth := smtp.PlainAuth("", fromString, os.Getenv("MAIL_PASS"), host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", serverName, tlsConfig)
	if err != nil {
		return err
	}

	cli, err := smtp.NewClient(conn, host)
	defer cli.Close()
	err = cli.Auth(emailAuth)
	err = cli.Mail(from.Address)
	err = cli.Rcpt(to.Address)
	w, err := cli.Data()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	msg := ""
	for i, j := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", i ,j)
	}

	msg += "MIME-Version: 1.0\r\n"
	if attachment != nil {
		msg += "Content-type: multipart/mixed; boundary=\"**=54jfsuf3jng3b\"\r\n"
		msg += "\r\n--**=54jfsuf3jng3b\r\n"
	}

	msg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	msg += "Content-Transfer-Encoding: 7bit\r\n"
	msg += body+"\r\n"

	if attachment != nil {
		msg += "\r\n--**=54jfsuf3jng3b\r\n"
		msg += "Content-Type: application/pdf; charset=\"utf-8\"\r\n"
		msg += "Content-Transfer-Encoding: base64\r\n"
		msg += "Content-Disposition: attachment;filename=\"label.pdf\"\r\n"
		msg += "\r\n" + base64.StdEncoding.EncodeToString(attachment)
	}

	_, err = w.Write([]byte(msg))
	err = w.Close()
	err = cli.Quit()
	if err != nil {
		return err
	}

	return nil
}