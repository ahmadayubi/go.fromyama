package amazon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"strings"
	"time"

	"go.fromyama/utils"
	"go.fromyama/utils/database"
	"go.fromyama/utils/jwtUtil"
	"go.fromyama/utils/response"
)

const updateAmazonAuthSql = "UPDATE companies SET amazon_seller_id = $1, amazon_auth_token = $2, amazon_marketplace = $3 WHERE id = $4"
const getAmazonTokensSql = "SELECT amazon_auth_token, amazon_seller_id, amazon_marketplace FROM companies WHERE id = $1"

// GetUnfulfilledOrders /orders/all returns array of orders from amazon
// TODO: create endpoint for unfufilled orders and for fulfilled orders for amazon
func GetUnfulfilledOrders (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var encryptedToken, sellerID, marketplace string
	getTokenQuery, err := database.DB.Prepare(getAmazonTokensSql)
	defer getTokenQuery.Close()
	err = getTokenQuery.QueryRow(tokenClaims.CompanyID).Scan(&encryptedToken, &sellerID, &marketplace)
	token, err := utils.Decrypt(encryptedToken)
	if err != nil{
		utils.ErrorResponse(w, "Hash Error")
		return
	}

	orderURL := formatAmazonURL("POST", "/Orders/2013-09-01", marketplace, map[string]string{
		"AWSAccessKeyId":os.Getenv("AMAZON_AWS_ACCESS_KEY"),
		"Action":"ListOrders",
		"CreatedAfter":"2020-07-01T04:00:00Z",
		"MWSAuthToken": token,
		"MarketplaceId.Id.1":marketplace,
		"OrderStatus.Status.1":"Unshipped",
		"SellerId":sellerID,
		"SignatureMethod":"HmacSHA256",
		"SignatureVersion":"2",
		"Timestamp": time.Now().UTC().Format(time.RFC3339),
		"Version":"2013-09-01",
	})

	var amazonOrders response.AmazonUnfulfilledResponse
	req, _ := http.NewRequest("POST", orderURL, nil)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err = xml.Unmarshal(respBody, &amazonOrders); err != nil{
		utils.ErrorResponse(w, "Unmarshal Error")
		return
	}
	resp.Body.Close()

	var orderItems [][]utils.Item
	if amazonOrders.ListOrdersResult.Orders[0].Order.AmazonOrderId != "" {
		for i := range amazonOrders.ListOrdersResult.Orders {
			detailURL := formatAmazonURL("POST", "/Orders/2013-09-01", marketplace, map[string]string{
				"AWSAccessKeyId":   os.Getenv("AMAZON_AWS_ACCESS_KEY"),
				"Action":           "ListOrderItems",
				"AmazonOrderId":    amazonOrders.ListOrdersResult.Orders[i].Order.AmazonOrderId,
				"MWSAuthToken":     token,
				"SellerId":         sellerID,
				"SignatureMethod":  "HmacSHA256",
				"SignatureVersion": "2",
				"Timestamp":        time.Now().UTC().Format(time.RFC3339),
				"Version":          "2013-09-01",
			})
			var orderDetails response.AmazonOrderDetailResponse
			reqD, _ := http.NewRequest("POST", detailURL, nil)
			respD, _ := client.Do(reqD)
			respDBody, _ := ioutil.ReadAll(respD.Body)
			if err = xml.Unmarshal(respDBody, &orderDetails); err != nil {
				utils.ErrorResponse(w, "Unmarshal Error")
				return
			}
			respD.Body.Close()
			orderItems = append(orderItems, formatAmazonItem(orderDetails))
		}
	}
	utils.JSONResponse(w, http.StatusOK, formatAmazonOrder(amazonOrders,orderItems))
}

// Authorize /authorize saves amazon tokens to company account
// request body has seller_id, auth_token, marketplace
func Authorize (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"seller_id", "auth_token", "marketplace"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}

	encryptedToken, err := utils.Encrypt(body["auth_token"])
	updateAmazonQuery, err := database.DB.Prepare(updateAmazonAuthSql)
	defer updateAmazonQuery.Close()
	_, err = updateAmazonQuery.Exec(body["seller_id"], encryptedToken, body["marketplace"], tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, "Update Token Error")
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, response.BasicMessage{Message: "Amazon Store Connected."})
}

// formatAmazonURL generates amazon url and signature for requests
func formatAmazonURL (method, path, host string, params map[string]string) string {
	hosts := map[string]string{
		"A2EUQ1WTGCTBG2":"mws.amazonservices.ca",
		"ATVPDKIKX0DER":"mws.amazonservices.ca",
		"A1AM78C64UM0Y8":"mws.amazonservices.com.mx",
		"A1F83G8C2ARO7P":"mws-eu.amazonservices.com",
	}
	signatureBase := method + "\n"
	signatureBase += hosts[host] + "\n"
	signatureBase += path + "\n"

	vals := url.Values{}

	for key, val := range params {
		vals.Add(key,val)
	}
	signatureBase += vals.Encode()

	mac := hmac.New(sha256.New, []byte(os.Getenv("AMAZON_SECRET")))
	mac.Write([]byte(signatureBase))
	hash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	url := "https://"+hosts[host]+"/Orders/2013-09-01?"+strings.Split(signatureBase,"\n")[3]+"&Signature="+url.QueryEscape(hash)
	return url
}

// formatAmazonOrder formats amazon orders response and returns fy orders
func formatAmazonOrder (resp response.AmazonUnfulfilledResponse, items [][]utils.Item) utils.Orders{
	var orders utils.Orders
	if len(items) > 0 {
		for i := range resp.ListOrdersResult.Orders {
			ord := utils.Order{
				Type:      "Amazon",
				OrderID:   resp.ListOrdersResult.Orders[i].Order.AmazonOrderId,
				CreatedAt: resp.ListOrdersResult.Orders[i].Order.PurchaseDate,
				Total:     resp.ListOrdersResult.Orders[i].Order.OrderTotal.Amount,
				Subtotal:  resp.ListOrdersResult.Orders[i].Order.OrderTotal.Amount,
				Tax:       "0",
				Currency:  resp.ListOrdersResult.Orders[i].Order.OrderTotal.CurrencyCode,
				Name:      resp.ListOrdersResult.Orders[i].Order.AmazonOrderId,
				WasPaid:   true,
				Items:     items[i],
				ShippingAddress: utils.Address{
					City:       resp.ListOrdersResult.Orders[i].Order.ShippingAddress.City,
					PostalCode: resp.ListOrdersResult.Orders[i].Order.ShippingAddress.PostalCode,
					Province:   resp.ListOrdersResult.Orders[i].Order.ShippingAddress.StateOrRegion,
					Country:    resp.ListOrdersResult.Orders[i].Order.ShippingAddress.CountryCode,
				},
			}
			orders.Orders = append(orders.Orders, ord)
		}
	}
	return orders
}

// formatAmazonItem formats amazon response for item details and returns fy items
func formatAmazonItem (resp response.AmazonOrderDetailResponse) []utils.Item {
	var items []utils.Item
	for i := range resp.ListOrderItemsResult.OrderItems {
		var item utils.Item
		item.Sku = resp.ListOrderItemsResult.OrderItems[i].OrderItem.SellerSKU
		item.Price = resp.ListOrderItemsResult.OrderItems[i].OrderItem.ItemPrice.Amount
		item.Quantity = resp.ListOrderItemsResult.OrderItems[i].OrderItem.QuantityOrdered
		item.Title = resp.ListOrderItemsResult.OrderItems[i].OrderItem.Title
		items = append(items, item)
	}
	return items
}