package shopify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
	"go.fromyama/utils"
	"go.fromyama/utils/database"
	"go.fromyama/utils/jwtUtil"
	"go.fromyama/utils/response"
)

const updateShopifyTokenSql = "UPDATE companies SET shopify_store = $1 ,shopify_token = $2, temp_data = null WHERE temp_data = $3"
const updateTempUUIDSql = "UPDATE companies SET temp_data = $1 WHERE id = $2"
const getShopifyTokenSql = "SELECT shopify_token, shopify_store FROM companies WHERE id = $1"

// FulfillOrder /fulfill fulfills order
// request body has order_id, location_id, notify_customer
func FulfillOrder (w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"location_id", "notify_customer"})
	if err != nil || orderID == ""{
		response.Error(w, "Body Parse Error, " + err.Error())
		return
	}
	var encryptedToken, store string

	getShopifyQuery, err := database.DB.Prepare(getShopifyTokenSql)
	defer getShopifyQuery.Close()
	err = getShopifyQuery.QueryRow(tokenClaims.CompanyID).Scan(&encryptedToken, &store)

	token, err := utils.AESDecrypt(encryptedToken)
	if err != nil {
		response.Error(w, "Hash Error")
		return
	}

	locationID, err := strconv.Atoi(body["location_id"])
	notifyCustomer, err := strconv.ParseBool(body["notify_customer"])
	fulfillmentData := response.ShopifyFulfillmentRequest{
		Fulfillment: response.ShopifyFulfillmentData{
			LocationID: locationID,
			TrackingNumber: body["tracking_number"],
			TrackingCompany: body["tracking_company"],

		},
		NotifyCustomer: notifyCustomer,
	}

	fulfillmentJSON, err := json.Marshal(fulfillmentData)

	respBody, err := shopifyRequest("POST", "https://"+store+"/admin/api/2021-04/orders/"+orderID+"/fulfillments.json", token, fulfillmentJSON)
	if err != nil {
		response.Error(w, "Shopify Request Error")
		return
	}

	var fulfillmentResponse response.ShopifyFulfillmentResponse
	err = json.Unmarshal(respBody, &fulfillmentResponse)
	if err != nil || fulfillmentResponse.Fulfillment.Status == "" {
		response.Error(w, "Marshal Error")
		return
	}

	response.JSON(w, http.StatusOK, response.BasicMessage{Message: "Order Fulfilled"})
}

// GetLocations /locations returns array of locations that shopify uses for fulfilling orders
func GetLocations (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var encryptedToken, store string

	getShopifyQuery, err := database.DB.Prepare(getShopifyTokenSql)
	defer getShopifyQuery.Close()
	err = getShopifyQuery.QueryRow(tokenClaims.CompanyID).Scan(&encryptedToken, &store)

	token, err := utils.AESDecrypt(encryptedToken)
	if err != nil {
		response.Error(w, "Hash Error")
		return
	}

	respBody, err := shopifyRequest("GET", "https://"+store+"/admin/api/2021-04/locations.json", token, nil)
	if err != nil {
		response.Error(w, "Shopify Request Error")
		return
	}

	var jsonResp response.LocationResponse
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		response.Error(w, "Unmarshal Error")
	} else {
		response.JSON(w, http.StatusOK, jsonResp)
	}
}

// GetUnfulfilledOrders /orders/all returns all unfulfilled orders
// TODO: create endpoint for unfulfilled orders and for fulfilled orders for shopify
func GetUnfulfilledOrders (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var encryptedToken, store string

	getShopifyQuery, err := database.DB.Prepare(getShopifyTokenSql)
	defer getShopifyQuery.Close()
	err = getShopifyQuery.QueryRow(tokenClaims.CompanyID).Scan(&encryptedToken, &store)

	token, err := utils.AESDecrypt(encryptedToken)
	if err != nil {
		response.Error(w, "Hash Error")
		return
	}

	respBody, err := shopifyRequest("GET", "https://"+store+"/admin/api/2021-04/orders.json?updated_at_min=2005-07-31T15:57:11-04:00&fulfillment_status=unfulfilled", token, nil)
	if err != nil {
		response.Error(w, "Shopify Request Error")
		return
	}

	var jsonResp response.ShopifyUnfulfilledResponse
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		response.Error(w, "Unmarshal Error")
	} else {
		response.JSON(w, http.StatusOK, formatShopifyOrder(jsonResp))
	}
}

// GenerateAuthURL /generate-link generates authentication link for shopify stores
// request body has shop
func GenerateAuthURL (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"shop"})
	if err != nil{
		response.Error(w, "Body Parse Error, " + err.Error())
		return
	}
	uid := uuid.New().String()
	updateTempQuery, err := database.DB.Prepare(updateTempUUIDSql)
	defer updateTempQuery.Close()
	_, err = updateTempQuery.Exec(uid, tokenClaims.CompanyID)
	if err != nil {
		response.Error(w, "Update UUID Error")
		return
	}
	var url response.UrlResponse
	url.Url = fmt.Sprintf("https://%s.myshopify.com/admin/oauth/authorize?client_id=%s&scope=read_orders,write_orders,read_customers&redirect_uri=%s/shopify/callback&state=%s&grant_options[]=",
		body["shop"], os.Getenv("SHOP_API_KEY"), os.Getenv("BASE_URL"), uid)

	response.JSON(w, http.StatusOK, url)
}

// Callback /callback authenticates company shopify store and encrypts tokens
// request url has code, hmac, timestamp, state, shop
func Callback(w http.ResponseWriter, r *http.Request){
	code := r.URL.Query().Get("code")
	hc := r.URL.Query().Get("hmac")
	timestamp := r.URL.Query().Get("timestamp")
	state := r.URL.Query().Get("state")
	shop := r.URL.Query().Get("shop")
	h := hmac.New(sha256.New, []byte(os.Getenv("SHOP_API_SKEY")))
	hData := "code="+code+"&shop="+shop+"&state="+state+"&timestamp="+timestamp
	h.Write([]byte(hData))
	sha := hex.EncodeToString(h.Sum(nil))

	if hmac.Equal([]byte(sha), []byte(hc)) {
		reqBody := response.PermAuth{
			ClientId: os.Getenv("SHOP_API_KEY"),
			ClientSecret: os.Getenv("SHOP_API_SKEY"),
			Code: code,
		}
		body, err := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", "https://"+shop+"/admin/oauth/access_token",bytes.NewBuffer(body))
		if err != nil {
			response.Error(w, "Token Error")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := client.Do(req)
		respBody, err := ioutil.ReadAll(resp.Body)
		var jsonResp response.CallbackResponse
		err = json.Unmarshal(respBody, &jsonResp)
		if err != nil {
			response.Error(w, "Unmarshal Error")
			return
		} else if jsonResp.AccessToken == "" {
			response.Forbidden(w)
			return
		}

		encryptedToken, err := utils.AESEncrypt(jsonResp.AccessToken)
		updateShopifyTokenQuery, err := database.DB.Prepare(updateShopifyTokenSql)
		defer updateShopifyTokenQuery.Close()
		_, err = updateShopifyTokenQuery.Exec(shop, encryptedToken, state)
		if err != nil {
			response.Error(w, "Update Token Error")
		} else {
			response.JSON(w, http.StatusAccepted, response.BasicMessage{Message: "Successfully Authenticated"})
		}
	}
}

// formatShopifyOrder util function that takes response from shopify and formats to fy order
func formatShopifyOrder (resp response.ShopifyUnfulfilledResponse) response.Orders {
	var orders response.Orders

	for i := range resp.Orders {
		add := response.Address{
			Name: resp.Orders[i].ShippingAddress.Name,
			Address1: resp.Orders[i].ShippingAddress.Address1,
			Address2: resp.Orders[i].ShippingAddress.Address2,
			City: resp.Orders[i].ShippingAddress.City,
			Province: resp.Orders[i].ShippingAddress.Province,
			Country: resp.Orders[i].ShippingAddress.Country,
			PostalCode: resp.Orders[i].ShippingAddress.Zip,
		}
		var items []response.Item
		for j := range resp.Orders[i].LineItems{
			items = append(items, response.Item{
				Sku:      resp.Orders[i].LineItems[j].Sku,
				Title:    resp.Orders[i].LineItems[j].Title,
				Quantity: strconv.Itoa(resp.Orders[i].LineItems[j].Quantity),
				Price:    resp.Orders[i].LineItems[j].Price,
			})
		}

		ord := response.Order{
			Type:    "Shopify",
			OrderID: strconv.FormatInt(resp.Orders[i].ID, 10),
			CreatedAt: resp.Orders[i].CreatedAt,
			Total: resp.Orders[i].TotalPrice,
			Subtotal: resp.Orders[i].SubtotalPrice,
			Tax: resp.Orders[i].TotalTax,
			Currency: resp.Orders[i].Currency,
			Name: resp.Orders[i].Name,
			WasPaid: resp.Orders[i].FinancialStatus == "paid",
			Items: items,
			ShippingAddress: add,
		}
		orders.Orders = append(orders.Orders, ord)
	}
	return orders
}

// shopifyRequest util function to make shopify api requests
func shopifyRequest (method, url, token string, body []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Add("X-Shopify-Access-Token", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "FromYama/2021.05.01")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody, nil
}