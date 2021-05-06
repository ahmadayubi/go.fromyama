package etsy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"go.fromyama/utils"
	"go.fromyama/utils/database"
	"go.fromyama/utils/jwtUtil"
	"go.fromyama/utils/response"
)

const addEtsyTempSql = "UPDATE companies SET etsy_store = $1,temp_data = $2, etsy_token = $3 WHERE id = $4"
const findCompanyByTempTokenSql = "SELECT id, temp_data FROM companies WHERE etsy_token = $1"
const updateEtsyTokenSql = "UPDATE companies SET temp_data = '', etsy_token=$1, etsy_token_secret=$2 WHERE id = $3"
const getEtsyTokenSql = "SELECT etsy_store, etsy_token, etsy_token_secret FROM companies WHERE id = $1"

// FulfillOrder /fulfill fulfills etsy order
// request body has order_id, optionally tracking_number, tracking_company
// TODO:need to update receipt to mark as shipped
func FulfillOrder (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	orderID := chi.URLParam(r, "orderID")
	var body map[string]string
	err := utils.ParseRequestBody(r, &body, nil)
	if err != nil || orderID == ""{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}

	var store, encryptedToken, encryptedSecret string
	getTokenQuery, err := database.DB.Prepare(getEtsyTokenSql)
	defer getTokenQuery.Close()
	err = getTokenQuery.QueryRow(tokenClaims.CompanyID).Scan(&store, &encryptedToken, &encryptedSecret)

	token, err := utils.Decrypt(encryptedToken)
	tokenSecret, err := utils.Decrypt(encryptedSecret)
	if err != nil {
		utils.ErrorResponse(w, "Hash Error")
		return
	}

	if body["tracking_number"] != "" && body["tracking_company"] != "" {
		params := "&tracking_code="+body["tracking_number"]+"&carrier_name="+body["tracking_company"]+"&api_key="+os.Getenv("ETSY_API_KEY")
		addTrackingResp := etsyRequest("POST", "https://openapi.etsy.com/v2/shops/"+store+"/receipts/"+orderID+"/tracking",
			params, token, tokenSecret)

		utils.JSONResponse(w, http.StatusOK, addTrackingResp)
	}



}

// GetUnfulfilledOrders /orders/all returns array of unfulfilled orders from etsy
// TODO: create endpoint for unfulfilled orders and for fulfilled orders for etsy
func GetUnfulfilledOrders (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var store, encryptedToken, encryptedSecret string
	getTokenQuery, err := database.DB.Prepare(getEtsyTokenSql)
	defer getTokenQuery.Close()
	err = getTokenQuery.QueryRow(tokenClaims.CompanyID).Scan(&store, &encryptedToken, &encryptedSecret)

	token, err := utils.Decrypt(encryptedToken)
	tokenSecret, err := utils.Decrypt(encryptedSecret)
	if err != nil {
		utils.ErrorResponse(w, "Hash Error")
		return
	}

	// https://openapi.etsy.com/v2/shops/teststorefy/receipts/open

	params := "api_key="+os.Getenv("ETSY_API_KEY")
	respBody := etsyRequest("GET", "https://openapi.etsy.com/v2/shops/"+store+"/receipts/open", params,token, tokenSecret)

	var jsonResponse response.EtsyUnfulfilledResponse
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		utils.ErrorResponse(w, "Unmarshal Error")
		return
	}

	jsonOrders := formatEtsyOrder(jsonResponse)

	utils.JSONResponse(w, http.StatusOK, jsonOrders)
}

// GenerateAuthURL /generate-link generates authentication link for etsy
// request body has shop
func GenerateAuthURL (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"shop"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	authHeader, _ := utils.GenerateRequestOAuthHeader("POST", "https://openapi.etsy.com/v2/oauth/request_token?scope=profile_r%20transactions_w%20transactions_r%20feedback_r",
		"https://0e1e72401111.ngrok.io/etsy/callback",
		os.Getenv("ETSY_API_KEY"), os.Getenv("ETSY_SECRET_KEY"))

	req, _ := http.NewRequest("POST", "https://openapi.etsy.com/v2/oauth/request_token?scope=profile_r%20transactions_w%20transactions_r%20feedback_r", nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	values, err := url.ParseQuery(string(respBody))
	if err != nil {
		utils.ErrorResponse(w, "Etsy Request Error")
		return
	}
	requestToken := values.Get("oauth_token")
	requestSecret := values.Get("oauth_token_secret")
	var authUrl utils.UrlResponse
	authUrl.Url, _ = url.QueryUnescape(strings.Split(string(respBody),"login_url=")[1])

	addTempQuery, err := database.DB.Prepare(addEtsyTempSql)
	defer addTempQuery.Close()
	_, err = addTempQuery.Exec(body["shop"], requestSecret, requestToken, tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, "Add Token Error")
		return
	}

	utils.JSONResponse(w, 200, authUrl)
}

// Callback /callback used to authenticate user and store
// request url has oauth_token, oauth_verifier
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
		utils.ErrorResponse(w, "Find Company Token")
		return
	}

	authHeader, _ := utils.GenerateOAuthHeader("POST", "https://openapi.etsy.com/v2/oauth/access_token",
		os.Getenv("ETSY_API_KEY"), os.Getenv("ETSY_SECRET_KEY"), authToken, tempSecret, map[string]string{
		"oauth_verifier":authVerifier,
		})

	req, _ := http.NewRequest("POST", "https://openapi.etsy.com/v2/oauth/access_token?oauth_token="+authToken+"&oauth_verifier="+authVerifier, nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	values, err := url.ParseQuery(string(respBody))
	if err != nil {
		utils.ErrorResponse(w, "Etsy Request Error")
		return
	}
	pAuthToken, err := utils.Encrypt(values.Get("oauth_token"))
	pAuthSecret, err := utils.Encrypt(values.Get("oauth_token_secret"))

	setTokenQuery, err := database.DB.Prepare(updateEtsyTokenSql)
	defer setTokenQuery.Close()
	_, err = setTokenQuery.Exec(pAuthToken, pAuthSecret, companyID)
	if err != nil {
		utils.ErrorResponse(w, "Update Token Error")
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, response.BasicMessage{Message: "Etsy Store Added"})
}

// formatEtsyOrder util function converts etsy response to fy orders
func formatEtsyOrder (resp response.EtsyUnfulfilledResponse) utils.Orders{
	var orders utils.Orders

	for i := range resp.Results {
		ord := utils.Order{
			Type: "Etsy",
			OrderID: string(i),
		}
		orders.Orders = append(orders.Orders, ord)
	}
	return orders
}

// etsyRequest util function that makes etsy api requests and returns response
func etsyRequest (method, url, params, token, tokenSecret string) []byte {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	authHeader, _ := utils.GenerateOAuthHeader(method, url,
		os.Getenv("ETSY_API_KEY"), os.Getenv("ETSY_SECRET_KEY"), token, tokenSecret,nil)
	req, _ := http.NewRequest(method, url+params, nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody
}