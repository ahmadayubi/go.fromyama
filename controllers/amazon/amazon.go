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

	"../../utils"
	"../../utils/database"
	"../../utils/jwtUtil"

)

const updateAmazonAuthSql = "UPDATE companies SET amazon_seller_id = $1, amazon_auth_token = $2, amazon_marketplace = $3 WHERE id = $4"
const getAmazonTokensSql = "SELECT amazon_auth_token, amazon_seller_id, amazon_marketplace FROM companies WHERE id = $1"

type unfulfilledResponse struct {
	XMLName          xml.Name `xml:"ListOrdersResponse"`
	Text             string   `xml:",chardata"`
	Xmlns            string   `xml:"xmlns,attr"`
	ListOrdersResult struct {
		Text   string `xml:",chardata"`
		Orders []struct {
			Text  string `xml:",chardata"`
			Order struct {
				Text                   string `xml:",chardata"`
				LatestShipDate         string `xml:"LatestShipDate"`
				OrderType              string `xml:"OrderType"`
				PurchaseDate           string `xml:"PurchaseDate"`
				BuyerEmail             string `xml:"BuyerEmail"`
				AmazonOrderId          string `xml:"AmazonOrderId"`
				LastUpdateDate         string `xml:"LastUpdateDate"`
				IsReplacementOrder     string `xml:"IsReplacementOrder"`
				NumberOfItemsShipped   string `xml:"NumberOfItemsShipped"`
				ShipServiceLevel       string `xml:"ShipServiceLevel"`
				OrderStatus            string `xml:"OrderStatus"`
				SalesChannel           string `xml:"SalesChannel"`
				ShippedByAmazonTFM     string `xml:"ShippedByAmazonTFM"`
				IsBusinessOrder        string `xml:"IsBusinessOrder"`
				NumberOfItemsUnshipped string `xml:"NumberOfItemsUnshipped"`
				LatestDeliveryDate     string `xml:"LatestDeliveryDate"`
				PaymentMethodDetails   struct {
					Text                string `xml:",chardata"`
					PaymentMethodDetail string `xml:"PaymentMethodDetail"`
				} `xml:"PaymentMethodDetails"`
				IsGlobalExpressEnabled string `xml:"IsGlobalExpressEnabled"`
				IsSoldByAB             string `xml:"IsSoldByAB"`
				EarliestDeliveryDate   string `xml:"EarliestDeliveryDate"`
				IsPremiumOrder         string `xml:"IsPremiumOrder"`
				OrderTotal             struct {
					Text         string `xml:",chardata"`
					Amount       string `xml:"Amount"`
					CurrencyCode string `xml:"CurrencyCode"`
				} `xml:"OrderTotal"`
				EarliestShipDate   string `xml:"EarliestShipDate"`
				MarketplaceId      string `xml:"MarketplaceId"`
				FulfillmentChannel string `xml:"FulfillmentChannel"`
				PaymentMethod      string `xml:"PaymentMethod"`
				ShippingAddress    struct {
					Text                         string `xml:",chardata"`
					City                         string `xml:"City"`
					PostalCode                   string `xml:"PostalCode"`
					IsAddressSharingConfidential string `xml:"isAddressSharingConfidential"`
					StateOrRegion                string `xml:"StateOrRegion"`
					CountryCode                  string `xml:"CountryCode"`
				} `xml:"ShippingAddress"`
				IsPrime                      string `xml:"IsPrime"`
				ShipmentServiceLevelCategory string `xml:"ShipmentServiceLevelCategory"`
			} `xml:"Order"`
		} `xml:"Orders"`
		CreatedBefore string `xml:"CreatedBefore"`
	} `xml:"ListOrdersResult"`
	ResponseMetadata struct {
		Text      string `xml:",chardata"`
		RequestId string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

type orderDetailsResponse struct {
	XMLName              xml.Name `xml:"ListOrderItemsResponse"`
	Text                 string   `xml:",chardata"`
	Xmlns                string   `xml:"xmlns,attr"`
	ListOrderItemsResult struct {
		Text          string `xml:",chardata"`
		AmazonOrderId string `xml:"AmazonOrderId"`
		OrderItems    []struct {
			Text      string `xml:",chardata"`
			OrderItem struct {
				Text          string `xml:",chardata"`
				TaxCollection struct {
					Text             string `xml:",chardata"`
					ResponsibleParty string `xml:"ResponsibleParty"`
					Model            string `xml:"Model"`
				} `xml:"TaxCollection"`
				QuantityOrdered   string `xml:"QuantityOrdered"`
				Title             string `xml:"Title"`
				PromotionDiscount struct {
					Text         string `xml:",chardata"`
					Amount       string `xml:"Amount"`
					CurrencyCode string `xml:"CurrencyCode"`
				} `xml:"PromotionDiscount"`
				ConditionId    string `xml:"ConditionId"`
				IsGift         string `xml:"IsGift"`
				ASIN           string `xml:"ASIN"`
				SellerSKU      string `xml:"SellerSKU"`
				OrderItemId    string `xml:"OrderItemId"`
				IsTransparency string `xml:"IsTransparency"`
				ProductInfo    struct {
					Text          string `xml:",chardata"`
					NumberOfItems string `xml:"NumberOfItems"`
				} `xml:"ProductInfo"`
				QuantityShipped    string `xml:"QuantityShipped"`
				ConditionSubtypeId string `xml:"ConditionSubtypeId"`
				ItemPrice          struct {
					Text         string `xml:",chardata"`
					Amount       string `xml:"Amount"`
					CurrencyCode string `xml:"CurrencyCode"`
				} `xml:"ItemPrice"`
				ItemTax struct {
					Text         string `xml:",chardata"`
					Amount       string `xml:"Amount"`
					CurrencyCode string `xml:"CurrencyCode"`
				} `xml:"ItemTax"`
				PromotionDiscountTax struct {
					Text         string `xml:",chardata"`
					Amount       string `xml:"Amount"`
					CurrencyCode string `xml:"CurrencyCode"`
				} `xml:"PromotionDiscountTax"`
			} `xml:"OrderItem"`
		} `xml:"OrderItems"`
	} `xml:"ListOrderItemsResult"`
	ResponseMetadata struct {
		Text      string `xml:",chardata"`
		RequestId string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}



func GetUnfulfilledOrders (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var token, sellerID, marketplace string
	getTokenQuery, err := database.DB.Prepare(getAmazonTokensSql)
	defer getTokenQuery.Close()
	err = getTokenQuery.QueryRow(tokenClaims.CompanyID).Scan(&token, &sellerID, &marketplace)
	if err != nil{
		utils.ErrorResponse(w, err)
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

	var amazonOrders unfulfilledResponse
	req, _ := http.NewRequest("POST", orderURL, nil)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err = xml.Unmarshal(respBody, &amazonOrders); err != nil{
		utils.ErrorResponse(w, err)
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
			var orderDetails orderDetailsResponse
			reqD, _ := http.NewRequest("POST", detailURL, nil)
			respD, _ := client.Do(reqD)
			respDBody, _ := ioutil.ReadAll(respD.Body)
			if err = xml.Unmarshal(respDBody, &orderDetails); err != nil {
				utils.ErrorResponse(w, err)
				return
			}
			respD.Body.Close()
			orderItems = append(orderItems, formatAmazonItem(orderDetails))
		}
	}
	utils.JSONResponse(w, http.StatusOK, formatAmazonOrder(amazonOrders,orderItems))
}

func Authorize (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"seller_id", "auth_token", "marketplace"})
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}

	updateAmazonQuery, err := database.DB.Prepare(updateAmazonAuthSql)
	defer updateAmazonQuery.Close()
	_, err = updateAmazonQuery.Exec(body["seller_id"], body["auth_token"], body["marketplace"], tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusAccepted, "Amazon Store Connected.")
}

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

func formatAmazonOrder (resp unfulfilledResponse, items [][]utils.Item) utils.Orders{
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

func formatAmazonItem (resp orderDetailsResponse) []utils.Item {
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