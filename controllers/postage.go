package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"
	"../utils/response"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
)

const getShippingInfoAndPaymentSql = "SELECT company_name, street, city, province_code, country, postal_code, phone, payment_account_id, email FROM companies c INNER JOIN users u on c.head_id = u.id where c.id = $1"
const addLabelsSql = "INSERT INTO labels(company_id, user_id, label, refund_link) VALUES ($1, $2, $3, $4)"
const getShippingPostalCodeSql = "SELECT postal_code FROM companies WHERE id = $1"

// BuyCanadaPostPostageLabel /buy/canadapost returns a pdf of the label just purchased
// request body has name, street, city, province_code, country_code, postal_code, phone, length, width,
// height, weight, service_code
func BuyCanadaPostPostageLabel (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"name", "street", "city", "province_code", "country_code",
		"postal_code", "phone", "length", "width", "height", "weight", "service_code"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}

	var source Shipper
	var dest Address
	var paymentAccount, email string
	var parcel Parcel
	dest.ShipperName = body["name"]
	dest.Street = body["street"]
	dest.City = body["city"]
	dest.ProvinceCode = body["province_code"]
	dest.Country = body["country_code"]
	dest.PostalCode = body["postal_code"]
	dest.Phone, _ = strconv.Atoi(body["phone"])
	parcel.Length, err = strconv.ParseFloat(body["length"], 64)
	parcel.Width, err = strconv.ParseFloat(body["width"], 64)
	parcel.Height, err = strconv.ParseFloat(body["height"], 64)
	parcel.Weight, err = strconv.ParseFloat(body["weight"], 64)

	source.Parcels = []Parcel{parcel}
	getShippingInfoQuery, err := database.DB.Prepare(getShippingInfoAndPaymentSql)
	defer getShippingInfoQuery.Close()
	err = getShippingInfoQuery.QueryRow(tokenClaims.CompanyID).Scan(&source.Address.ShipperName, &source.Address.Street, &source.Address.City, &source.Address.ProvinceCode, &source.Address.Country, &source.Address.PostalCode, &source.Address.Phone, &paymentAccount, &email)
	if err != nil {
		utils.ErrorResponse(w, "Get Shipper Error")
		return
	}

	var rate response.CanadaPostRatesResponse
	priceXml := formatCanadaPostSingleRateBody(source, body["postal_code"], body["weight"], body["service_code"])
	err = canadaPostRequest("POST", "https://ct.soa-gw.canadapost.ca/rs/ship/price", "application/vnd.cpc.ship.rate-v4+xml", "application/vnd.cpc.ship.rate-v4+xml",
		priceXml, &rate)

	rateJSON := rateResponseToJSON(rate)
	total, err := strconv.ParseFloat(rateJSON[0].PriceDetails.Due, 64)
	intAmount := int(total * 100)

	stripe.Key = os.Getenv("STRIPE_SECRET")
	params := &stripe.ChargeParams{
		Amount: stripe.Int64(int64(intAmount)),
		Currency: stripe.String(string(stripe.CurrencyCAD)),
		Description: stripe.String("Label Purchase"),
		Customer: stripe.String(paymentAccount),
	}
	c, _ := charge.New(params)
	if !c.Paid {
		utils.JSONResponse(w, http.StatusConflict, "Payment Error")
		return
	}
	refundParam := &stripe.RefundParams{
		Charge: stripe.String(c.ID),
	}

	xmlBody := formatCanadaPostRequestBody(source, dest, body["weight"], body["service_code"])
	client := &http.Client{
		Timeout: time.Second * 20,
	}

	var canadaPostBody response.CanadaPostPostageResponse

	err = canadaPostRequest("POST", "https://ct.soa-gw.canadapost.ca/rs/"+os.Getenv("CANADA_POST_CUSTNUM")+"/ncshipment", "application/vnd.cpc.ncshipment-v4+xml",
		"application/vnd.cpc.ncshipment-v4+xml",xmlBody, &canadaPostBody)
	if err != nil {
		_, _ = refund.New(refundParam)
		utils.ErrorResponse(w ,"Postage Error")
		return
	}

	var labelLink, refundLink string
	for i := range canadaPostBody.Links {
		if canadaPostBody.Links[i].Name == "label" {
			labelLink = canadaPostBody.Links[i].Link
		}
		if canadaPostBody.Links[i].Name == "refund" {
			refundLink = canadaPostBody.Links[i].Link
		}
	}

	addLabelQuery, err := database.DB.Prepare(addLabelsSql)
	defer addLabelQuery.Close()
	_, err = addLabelQuery.Exec(tokenClaims.CompanyID, tokenClaims.UserID,labelLink, refundLink)
	if err != nil {
		_, _ = refund.New(refundParam)
		utils.ErrorResponse(w ,"Store Postage Error")
		return
	}

	getLabelReq, err := http.NewRequest("GET", labelLink, nil)
	getLabelReq.SetBasicAuth(os.Getenv("CANADA_POST_USER"), os.Getenv("CANADA_POST_PASS"))
	getLabelReq.Header.Add("Accept","application/pdf")
	labelResp, err := client.Do(getLabelReq)
	defer labelResp.Body.Close()
	label, err := ioutil.ReadAll(labelResp.Body)

	receipt := response.CanadaPostLabelPurchase{
		PostalCode: body["postal_code"],
		Total: float64(intAmount)/100.0,
		TrackingNumber: canadaPostBody.Tracking,
	}

	var tmplBuffer bytes.Buffer
	tmpl := template.Must(template.ParseFiles("assets/templates/labelPurchase.html"))
	tmpl.Execute(&tmplBuffer, receipt)

	err = SendEmail("ahmad.ayubi@hotmail.com", "Shipping Label Purchased", tmplBuffer.String(), label)
	if err != nil {
		_, _ = refund.New(refundParam)
		_ = SendEmail("ahmad.ayubi@hotmail.com", "Refund Label Purchase", "Refund Label Purchase", nil)
		utils.ErrorResponse(w ,"Sending Email Error")
		return
	}



	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(label)
}

// GetCanadaPostRate /rates/canadapost returns array of rates
// request body has postal_code, weight
func GetCanadaPostRate (w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"postal_code","weight", "length", "width", "height"})
	if err != nil{
		utils.ErrorResponse(w, "Body Parse Error, " + err.Error())
		return
	}

	var sourcePostalCode string

	getShippingInfoQuery, err := database.DB.Prepare(getShippingPostalCodeSql)
	defer getShippingInfoQuery.Close()
	err = getShippingInfoQuery.QueryRow(tokenClaims.CompanyID).Scan(&sourcePostalCode)
	if err != nil {
		utils.ErrorResponse(w, "Get Shipper Error")
		return
	}
	xmlBody := formatCanadaPostRateBody(sourcePostalCode, body["postal_code"], body["length"], body["width"], body["height"], body["weight"])

	var rates response.CanadaPostRatesResponse
	err = canadaPostRequest("POST", "https://ct.soa-gw.canadapost.ca/rs/ship/price", "application/vnd.cpc.ship.rate-v4+xml", "application/vnd.cpc.ship.rate-v4+xml",
		xmlBody, &rates)
	if err != nil{
		utils.ErrorResponse(w, "Marshal Error")
		return
	}

	rateJSON := rateResponseToJSON(rates)
	utils.JSONResponse(w, http.StatusOK, rateJSON)
}

// formatCanadaPostRequestBody formats xml body for postage purchase request
func formatCanadaPostRequestBody (source Shipper, dest Address, weight string, serviceCode string) string{
	xml := `<?xml version="1.0" encoding="utf-8"?>`
	xml += `<non-contract-shipment xmlns="http://www.canadapost.ca/ws/ncshipment-v4">`
	xml += `<requested-shipping-point>`+source.Address.PostalCode+`</requested-shipping-point>`
	xml += `<delivery-spec>`
	xml += `<service-code>`+serviceCode+`</service-code>`
	xml += `<sender>`
	xml += `<company>`+source.Address.ShipperName+`</company>`
	xml += `<contact-phone>`+strconv.Itoa(source.Address.Phone)+`</contact-phone>`
	xml += `<address-details>`
	xml += `<address-line-1>`+source.Address.Street+`</address-line-1>`
	xml += `<city>`+source.Address.City+`</city>`
	xml += `<prov-state>`+source.Address.ProvinceCode+`</prov-state>`
	xml += `<postal-zip-code>`+source.Address.PostalCode+`</postal-zip-code>`
	xml += `</address-details>`
	xml += `</sender>`
	xml += `<destination>`
	xml += `<name>`+dest.ShipperName+`</name>`
	xml += `<address-details>`
	xml += `<address-line-1>`+dest.Street+`</address-line-1>`
	xml += `<city>`+dest.City+`</city>`
	xml += `<prov-state>`+dest.ProvinceCode+`</prov-state>`
	xml += `<country-code>`+dest.Country+`</country-code>`
	xml += `<postal-zip-code>`+dest.PostalCode+`</postal-zip-code>`
	xml += `</address-details>`
	xml += `</destination>`
	xml += `<options>`
	xml += `<option>`
	xml += `<option-code>DC</option-code>`
	xml += `</option>`
	xml += `</options>`
	xml += `<parcel-characteristics>`
	xml += `<weight>`+weight+`</weight>`
	xml += `<dimensions>`
	xml += `<length>`+fmt.Sprintf("%f",source.Parcels[0].Length)+`</length>`
	xml += `<width>`+fmt.Sprintf("%f",source.Parcels[0].Width)+`</width>`
	xml += `<height>`+fmt.Sprintf("%f",source.Parcels[0].Height)+`</height>`
	xml += `</dimensions>`
	xml += `</parcel-characteristics>`
	xml += `<preferences>`
	xml += `<show-packing-instructions>true</show-packing-instructions>`
	xml += `</preferences>`
	xml += `</delivery-spec>`
	xml += `</non-contract-shipment>`
	return xml
}

// formatCanadaPostRateBody formats xml body for rate checking canada post request
func formatCanadaPostRateBody (sourcePostalCode, destPostalCode, length, width, height, weight string) string {
	xml := `<?xml version="1.0" encoding="utf-8"?>`
	xml += `<mailing-scenario xmlns="http://www.canadapost.ca/ws/ship/rate-v4">`
	xml += `<customer-number>`+os.Getenv("CANADA_POST_CUSTNUM")+`</customer-number>`
	xml += `<parcel-characteristics>`
	xml += `<weight>`+weight+`</weight>`
	xml += `<dimensions><length>`+length+"</length>"
	xml += `<width>`+width+"</width>"
	xml += `<height>`+height+"</height></dimensions>"
	xml += `</parcel-characteristics>`
	xml += `<origin-postal-code>`+sourcePostalCode+`</origin-postal-code>`
	xml += `<destination><domestic>`
	xml += `<postal-code>`+destPostalCode+`</postal-code>`
	xml += `</domestic></destination></mailing-scenario>`
	return xml
}

func formatCanadaPostSingleRateBody (source Shipper, destPostalCode, weight, serviceCode string) string {
	xml := `<?xml version="1.0" encoding="utf-8"?>`
	xml += `<mailing-scenario xmlns="http://www.canadapost.ca/ws/ship/rate-v4">`
	xml += `<customer-number>`+os.Getenv("CANADA_POST_CUSTNUM")+`</customer-number>`
	xml += `<parcel-characteristics>`
	xml += `<weight>`+weight+`</weight>`
	xml += `<dimensions><length>`+fmt.Sprintf("%f",source.Parcels[0].Length)+"</length>"
	xml += `<width>`+fmt.Sprintf("%f",source.Parcels[0].Width)+"</width>"
	xml += `<height>`+fmt.Sprintf("%f",source.Parcels[0].Height)+"</height></dimensions>"
	xml += `</parcel-characteristics>`
	xml += `<services><service-code>`+serviceCode+"</service-code></services>"
	xml += `<origin-postal-code>`+source.Address.PostalCode+`</origin-postal-code>`
	xml += `<destination><domestic>`
	xml += `<postal-code>`+destPostalCode+`</postal-code>`
	xml += `</domestic></destination></mailing-scenario>`
	return xml
}

func rateResponseToJSON (resp response.CanadaPostRatesResponse) []response.CanadaPostRate{
	var rateJSON []response.CanadaPostRate
	for i := range resp.PriceQuote{
		var adjustments []response.CanadaPostAdjustment
		for j := range resp.PriceQuote[i].PriceDetails.Adjustments.Adjustment {
			adjustments = append(adjustments, response.CanadaPostAdjustment{
				AdjustmentCode: resp.PriceQuote[i].PriceDetails.Adjustments.Adjustment[j].AdjustmentCode,
				AdjustmentCost: resp.PriceQuote[i].PriceDetails.Adjustments.Adjustment[j].AdjustmentCost,
				AdjustmentName: resp.PriceQuote[i].PriceDetails.Adjustments.Adjustment[j].AdjustmentName,
			})
		}
		due, _ := strconv.ParseFloat(resp.PriceQuote[i].PriceDetails.Due, 64)
		rateJSON = append(rateJSON, response.CanadaPostRate{
			ServiceCode: resp.PriceQuote[i].ServiceCode,
			ServiceName: resp.PriceQuote[i].ServiceName,
			WeightDetails: resp.PriceQuote[i].WeightDetails,
			PriceDetails: response.CanadaPostRatePrice{
				Base: resp.PriceQuote[i].PriceDetails.Base,
				Gst: resp.PriceQuote[i].PriceDetails.Taxes.Gst,
				Hst: resp.PriceQuote[i].PriceDetails.Taxes.Hst.Percent,
				Pst: resp.PriceQuote[i].PriceDetails.Taxes.Pst,
				Due: fmt.Sprintf("%f",due + 1.1),
				Adjustments: adjustments,
			},
			AmDelivery: resp.PriceQuote[i].ServiceStandard.AmDelivery,
			GuaranteedDelivery: resp.PriceQuote[i].ServiceStandard.GuaranteedDelivery,
			ExpectedTransitTime: resp.PriceQuote[i].ServiceStandard.ExpectedTransitTime,
			ExpectedDeliveryDate: resp.PriceQuote[i].ServiceStandard.ExpectedDeliveryDate,
		})
	}
	return rateJSON
}

func canadaPostRequest(method, url, acceptType, contentType, xmlBody string, parseTo interface{}) error{
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(xmlBody)))
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	authEncoded := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CANADA_POST_USER")+":"+os.Getenv("CANADA_POST_PASS")))
	req.Header.Add("Accept", acceptType)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Basic " + authEncoded)
	req.Header.Add("Accept-language", "en-CA")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	if err = xml.Unmarshal(respBody, parseTo); err != nil{
		return err
	}
	return nil
}