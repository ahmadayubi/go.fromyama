package etsy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"../../utils"
	"../../utils/database"
	"../../utils/jwtUtil"
)

const addEtsyTempSql = "UPDATE companies SET temp_data = $1, etsy_token = $2 WHERE id = $3"
const findCompanyByTempTokenSql = "SELECT id, temp_data FROM companies WHERE etsy_token = $1"

func GenerateAuthURL (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	authHeader, _ := utils.GenerateRequestOAuthHeader("POST", "https://openapi.etsy.com/v2/oauth/request_token?scope=profile_r%20transactions_w%20transactions_r%20feedback_r",
		"https://29700a4b7f5b.ngrok.io/etsy/callback",
		os.Getenv("ETSY_API_KEY"), os.Getenv("ETSY_SECRET_KEY"))

	req, _ := http.NewRequest("POST", "https://openapi.etsy.com/v2/oauth/request_token?scope=profile_r%20transactions_w%20transactions_r%20feedback_r", nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	values, err := url.ParseQuery(string(respBody))
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	requestToken := values.Get("oauth_token")
	requestSecret := values.Get("oauth_token_secret")
	var authUrl utils.UrlResponse
	authUrl.Url, _ = url.QueryUnescape(strings.Split(string(respBody),"login_url=")[1])

	addTempQuery, err := database.DB.Prepare(addEtsyTempSql)
	defer addTempQuery.Close()
	_, err = addTempQuery.Exec(requestSecret, requestToken, tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	utils.JSONResponse(w, 200, authUrl)
}

func Callback(w http.ResponseWriter, r *http.Request){
	authToken := r.URL.Query().Get("oauth_token")
	authVerifier := r.URL.Query().Get("oauth_verifier")
	var companyID, tempSecret string
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	findCompanySql, err := database.DB.Prepare(findCompanyByTempTokenSql)
	defer findCompanySql.Close()
	err = findCompanySql.QueryRow(authToken).Scan(&companyID, &tempSecret)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	authHeader, _ := utils.GenerateOAuthHeader("POST", "https://openapi.etsy.com/v2/oauth/access_token",
		os.Getenv("ETSY_API_KEY"), os.Getenv("ETSY_SECRET_KEY"), authToken, tempSecret, map[string]string{
		"oauth_verifier":authVerifier,
		})

	req, _ := http.NewRequest("GET", "https://openapi.etsy.com/v2/oauth/access_token", nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	utils.JSONResponse(w, 200, respBody)
}