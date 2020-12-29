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

	"../../utils"
	"../../utils/database"
	"../../utils/jwtUtil"
	"github.com/google/uuid"
)

const updateShopifyTokenSql = "UPDATE companies SET shopify_store = $1 ,shopify_token = $2, temp_data = null WHERE temp_data = $3"
const updateTempUUIDSql = "UPDATE companies SET temp_data = $1 WHERE id = $2"
const getShopifyTokenSql = "SELECT shopify_token, shopify_store FROM companies WHERE id = $1"

type callbackResponse struct {
	AccessToken string `json:"access_token"`
}
type locationResponse struct {
	Locations []struct {
		ID                    int64       `json:"id"`
		Name                  string      `json:"name"`
	} `json:"locations"`
}
type permAuth struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code string `json:"code"`
}
type unfulfilledResponse struct {
	Orders []struct {
		ID                    int64         `json:"id"`
		Email                 string        `json:"email"`
		CreatedAt             string        `json:"created_at"`
		Number                int           `json:"number"`
		Note                  string        `json:"note"`
		Token                 string        `json:"token"`
		TotalPrice            string        `json:"total_price"`
		SubtotalPrice         string        `json:"subtotal_price"`
		TotalWeight           int           `json:"total_weight"`
		TotalTax              string        `json:"total_tax"`
		TaxesIncluded         bool          `json:"taxes_included"`
		Currency              string        `json:"currency"`
		FinancialStatus       string        `json:"financial_status"`
		Confirmed             bool          `json:"confirmed"`
		TotalDiscounts        string        `json:"total_discounts"`
		TotalLineItemsPrice   string        `json:"total_line_items_price"`
		CartToken             interface{}   `json:"cart_token"`
		BuyerAcceptsMarketing bool          `json:"buyer_accepts_marketing"`
		Name                  string        `json:"name"`
		Phone                 interface{}   `json:"phone"`
		OrderNumber           int           `json:"order_number"`
		TaxLines              []struct {
			Price    string  `json:"price"`
			Rate     float64 `json:"rate"`
			Title    string  `json:"title"`
			PriceSet struct {
				ShopMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"shop_money"`
				PresentmentMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"presentment_money"`
			} `json:"price_set"`
		} `json:"tax_lines"`
		TotalLineItemsPriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"total_line_items_price_set"`
		TotalDiscountsSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"total_discounts_set"`
		TotalShippingPriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"total_shipping_price_set"`
		SubtotalPriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"subtotal_price_set"`
		TotalPriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"total_price_set"`
		TotalTaxSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"total_tax_set"`
		LineItems []struct {
			ID                         int64         `json:"id"`
			VariantID                  int64         `json:"variant_id"`
			Title                      string        `json:"title"`
			Quantity                   int           `json:"quantity"`
			Sku                        string        `json:"sku"`
			ProductID                  int64         `json:"product_id"`
			RequiresShipping           bool          `json:"requires_shipping"`
			Taxable                    bool          `json:"taxable"`
			Name                       string        `json:"name"`
			FulfillableQuantity        int           `json:"fulfillable_quantity"`
			Grams                      int           `json:"grams"`
			Price                      string        `json:"price"`
			TotalDiscount              string        `json:"total_discount"`
			PriceSet                   struct {
				ShopMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"shop_money"`
				PresentmentMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"presentment_money"`
			} `json:"price_set"`
			TotalDiscountSet struct {
				ShopMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"shop_money"`
				PresentmentMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"presentment_money"`
			} `json:"total_discount_set"`
			TaxLines            []struct {
				Title    string  `json:"title"`
				Price    string  `json:"price"`
				Rate     float64 `json:"rate"`
				PriceSet struct {
					ShopMoney struct {
						Amount       string `json:"amount"`
						CurrencyCode string `json:"currency_code"`
					} `json:"shop_money"`
					PresentmentMoney struct {
						Amount       string `json:"amount"`
						CurrencyCode string `json:"currency_code"`
					} `json:"presentment_money"`
				} `json:"price_set"`
			} `json:"tax_lines"`
		} `json:"line_items"`
		TotalTipReceived       string        `json:"total_tip_received"`
		BillingAddress         struct {
			FirstName    string      `json:"first_name"`
			Address1     string      `json:"address1"`
			Phone        string      `json:"phone"`
			City         string      `json:"city"`
			Zip          string      `json:"zip"`
			Province     string      `json:"province"`
			Country      string      `json:"country"`
			LastName     string      `json:"last_name"`
			Address2     string `json:"address2"`
			Company      string `json:"company"`
			Name         string      `json:"name"`
			CountryCode  string      `json:"country_code"`
			ProvinceCode string      `json:"province_code"`
		} `json:"billing_address"`
		ShippingAddress struct {
			FirstName    string  `json:"first_name"`
			Address1     string  `json:"address1"`
			Phone        string  `json:"phone"`
			City         string  `json:"city"`
			Zip          string  `json:"zip"`
			Province     string  `json:"province"`
			Country      string  `json:"country"`
			LastName     string  `json:"last_name"`
			Address2     string  `json:"address2"`
			Company      string  `json:"company"`
			Name         string  `json:"name"`
			CountryCode  string  `json:"country_code"`
			ProvinceCode string  `json:"province_code"`
		} `json:"shipping_address"`
		Customer struct {
			ID                        int64         `json:"id"`
			Email                     string        `json:"email"`
			FirstName                 string        `json:"first_name"`
			LastName                  string        `json:"last_name"`
			OrdersCount               int           `json:"orders_count"`
			Note                      interface{}   `json:"note"`
			Phone                     interface{}   `json:"phone"`
			Tags                      string        `json:"tags"`
			LastOrderName             string        `json:"last_order_name"`
		} `json:"customer"`
	} `json:"orders"`
}


func GetLocations (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var encryptedToken, store string
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	getShopifyQuery, err := database.DB.Prepare(getShopifyTokenSql)
	defer getShopifyQuery.Close()
	err =getShopifyQuery.QueryRow(tokenClaims.CompanyID).Scan(&encryptedToken, &store)

	token, err := utils.Decrypt(encryptedToken)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	req, err := http.NewRequest("GET", "https://"+store+"/admin/api/2020-04/locations.json", nil)
	req.Header.Add("X-Shopify-Access-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	var jsonResp locationResponse
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, jsonResp)
}

func GetUnfulfilledOrders (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)
	var encryptedToken, store string
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	getShopifyQuery, err := database.DB.Prepare(getShopifyTokenSql)
	defer getShopifyQuery.Close()
	err =getShopifyQuery.QueryRow(tokenClaims.CompanyID).Scan(&encryptedToken, &store)

	token, err := utils.Decrypt(encryptedToken)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	req, err := http.NewRequest("GET", "https://"+store+"/admin/api/2020-04/orders.json?updated_at_min=2005-07-31T15:57:11-04:00&fulfillment_status=unfulfilled", nil)
	req.Header.Add("X-Shopify-Access-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	var jsonResp unfulfilledResponse
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, formatShopifyOrder(jsonResp))
}

func GenerateAuthURL (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"shop"})
	if err != nil{
		utils.ErrorResponse(w, err)
		return
	}
	uid := uuid.New()
	updateTempQuery, err := database.DB.Prepare(updateTempUUIDSql)
	defer updateTempQuery.Close()
	_, err = updateTempQuery.Exec(uid.String(), tokenClaims.CompanyID)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}
	var url utils.UrlResponse
	url.Url = fmt.Sprintf("https://%s.myshopify.com/admin/oauth/authorize?client_id=%s&scope=read_orders,write_orders,read_customers&redirect_uri=%s/shopify/callback&state=%s&grant_options[]=",
		body["shop"], os.Getenv("SHOP_API_KEY"), os.Getenv("BASE_URL"), uid.String())

	utils.JSONResponse(w, http.StatusOK, url)
}

func Callback(w http.ResponseWriter, r *http.Request){
	decrpyt, err := utils.Decrypt("HMQLxYVIDIDdSlK0hrx3iTrsdXK8K7uiq5nQ20Dt_a0w312hi2Fi4PUscTEYdmKxGXJgaiI-")
	if err == nil {
		utils.JSONResponse(w, http.StatusOK, decrpyt)
		return
	}
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
		reqBody := permAuth{
			ClientId: os.Getenv("SHOP_API_KEY"),
			ClientSecret: os.Getenv("SHOP_API_SKEY"),
			Code: code,
		}
		body, err := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", "https://"+shop+"/admin/oauth/access_token",bytes.NewBuffer(body))
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := client.Do(req)
		respBody, err := ioutil.ReadAll(resp.Body)
		var jsonResp callbackResponse
		err = json.Unmarshal(respBody, &jsonResp)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		if jsonResp.AccessToken == "" {
			utils.ForbiddenResponse(w)
			return
		}

		encryptedToken, err := utils.Encrypt(jsonResp.AccessToken)
		updateShopifyTokenQuery, err := database.DB.Prepare(updateShopifyTokenSql)
		defer updateShopifyTokenQuery.Close()
		_, err = updateShopifyTokenQuery.Exec(shop, encryptedToken, state)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusAccepted, "Successfully Authenticated")
		return
	}
}

func formatShopifyOrder (resp unfulfilledResponse) utils.Orders {
	var orders utils.Orders

	for i := range resp.Orders {
		add := utils.Address{
			Name: resp.Orders[i].ShippingAddress.Name,
			Address1: resp.Orders[i].ShippingAddress.Address1,
			Address2: resp.Orders[i].ShippingAddress.Address2,
			City: resp.Orders[i].ShippingAddress.City,
			Province: resp.Orders[i].ShippingAddress.Province,
			Country: resp.Orders[i].ShippingAddress.Country,
			PostalCode: resp.Orders[i].ShippingAddress.Zip,
		}
		var items []utils.Item
		for j := range resp.Orders[i].LineItems{
			items = append(items, utils.Item{
				Sku:      resp.Orders[i].LineItems[j].Sku,
				Title:    resp.Orders[i].LineItems[j].Title,
				Quantity: strconv.Itoa(resp.Orders[i].LineItems[j].Quantity),
				Price:    resp.Orders[i].LineItems[j].Price,
			})
		}

		ord := utils.Order{
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