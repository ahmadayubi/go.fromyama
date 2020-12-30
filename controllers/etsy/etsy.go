package etsy

import (
	"encoding/json"
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

const addEtsyTempSql = "UPDATE companies SET etsy_store = $1,temp_data = $2, etsy_token = $3 WHERE id = $4"
const findCompanyByTempTokenSql = "SELECT id, temp_data FROM companies WHERE etsy_token = $1"
const updateEtsyTokenSql = "UPDATE companies SET temp_data = '', etsy_token=$1, etsy_token_secret=$2 WHERE id = $3"
const getEtsyTokenSql = "SELECT etsy_store, etsy_token, etsy_token_secret FROM companies WHERE id = $1"

type unfulfilledResponse struct {
	Count   int `json:"count"`
	Results []struct {
		ShopID                         int         `json:"shop_id"`
		ShopName                       string      `json:"shop_name"`
		UserID                         int         `json:"user_id"`
		CreationTsz                    int         `json:"creation_tsz"`
		Title                          interface{} `json:"title"`
		Announcement                   interface{} `json:"announcement"`
		CurrencyCode                   string      `json:"currency_code"`
		IsVacation                     bool        `json:"is_vacation"`
		VacationMessage                interface{} `json:"vacation_message"`
		SaleMessage                    interface{} `json:"sale_message"`
		DigitalSaleMessage             interface{} `json:"digital_sale_message"`
		LastUpdatedTsz                 int         `json:"last_updated_tsz"`
		ListingActiveCount             int         `json:"listing_active_count"`
		DigitalListingCount            int         `json:"digital_listing_count"`
		LoginName                      string      `json:"login_name"`
		URL                            string      `json:"url"`
		ImageURL760X100                interface{} `json:"image_url_760x100"`
		NumFavorers                    int         `json:"num_favorers"`
		Languages                      []string    `json:"languages"`
		UpcomingLocalEventID           interface{} `json:"upcoming_local_event_id"`
		IconURLFullxfull               interface{} `json:"icon_url_fullxfull"`
		IsUsingStructuredPolicies      bool        `json:"is_using_structured_policies"`
		HasOnboardedStructuredPolicies bool        `json:"has_onboarded_structured_policies"`
		HasUnstructuredPolicies        bool        `json:"has_unstructured_policies"`
		IncludeDisputeFormLink         bool        `json:"include_dispute_form_link"`
		IsDirectCheckoutOnboarded      bool        `json:"is_direct_checkout_onboarded"`
		IsCalculatedEligible           bool        `json:"is_calculated_eligible"`
		IsOptedInToBuyerPromise        bool        `json:"is_opted_in_to_buyer_promise"`
		IsShopUsBased                  bool        `json:"is_shop_us_based"`
	} `json:"results"`
	Params struct {
		ShopID string `json:"shop_id"`
	} `json:"params"`
	Type       string `json:"type"`
	Pagination struct {
	} `json:"pagination"`
}

func GetUnfulfilledOrders (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var store, encryptedToken, encryptedSecret string
	getTokenQuery, err := database.DB.Prepare(getEtsyTokenSql)
	defer getTokenQuery.Close()
	err = getTokenQuery.QueryRow(tokenClaims.CompanyID).Scan(&store, &encryptedToken, &encryptedSecret)

	token, err := utils.Decrypt(encryptedToken)
	tokenSecret, err := utils.Decrypt(encryptedSecret)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// https://openapi.etsy.com/v2/shops/teststorefy/receipts/open

	authHeader, _ := utils.GenerateOAuthHeader("GET", "https://openapi.etsy.com/v2/shops/"+store+"/receipts/open",
		os.Getenv("ETSY_API_KEY"), os.Getenv("ETSY_SECRET_KEY"), token,tokenSecret,nil)

	req, _ := http.NewRequest("GET", "https://openapi.etsy.com/v2/shops/"+store+"/receipts?was_shopped=false&api_key="+os.Getenv("ETSY_API_KEY"), nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse unfulfilledResponse
	err = json.Unmarshal(respBody, &jsonResponse)

	jsonOrders := formatEtsyOrder(jsonResponse)

	utils.JSONResponse(w, http.StatusOK, jsonOrders)
}

func GenerateAuthURL (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"shop"})
	if err != nil{
		utils.ErrorResponse(w, err)
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
		utils.ErrorResponse(w, err)
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

	req, _ := http.NewRequest("POST", "https://openapi.etsy.com/v2/oauth/access_token?oauth_token="+authToken+"&oauth_verifier="+authVerifier, nil)
	req.Header.Add("Authorization", authHeader)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	values, err := url.ParseQuery(string(respBody))
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	pAuthToken, err := utils.Encrypt(values.Get("oauth_token"))
	pAuthSecret, err := utils.Encrypt(values.Get("oauth_token_secret"))

	setTokenQuery, err := database.DB.Prepare(updateEtsyTokenSql)
	defer setTokenQuery.Close()
	_, err = setTokenQuery.Exec(pAuthToken, pAuthSecret, companyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, "Etsy Store Added")
}

func formatEtsyOrder (resp unfulfilledResponse) utils.Orders{
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