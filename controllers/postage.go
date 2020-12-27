package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"../utils"
	"../utils/database"
	"../utils/jwtUtil"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
)

const getShippingInfoAndPaymentSql = "SELECT company_name, street, city, province_code, country, postal_code, phone, payment_account_id, email FROM companies c INNER JOIN users u on c.head_id = u.id where c.id = $1"
const addLabelsSql = "INSERT INTO labels(company_id, user_id, label, refund_link) VALUES ($1, $2, $3, $4)"

type CanadaPostXML struct {
	Tracking string `xml:"tracking-pin"`
	Links []struct {
		Name string `xml:"rel,attr"`
		Link string `xml:"href,attr"`
	} `xml:"links>link"`
}

type LabelPurchase struct {
	PostalCode string
	Total string
	TrackingNumber string
}

func BuyCanadaPostPostageLabel (w http.ResponseWriter, r *http.Request){
	tokenClaims := r.Context().Value("claims").(jwtUtil.TokenClaims)

	var body map[string]string
	err := utils.ParseRequestBody(r, &body,[]string{"name", "street", "city", "province_code", "country_code",
		"postal_code", "phone", "length", "width", "height", "weight", "service_code", "label_total"})
	if err != nil{

		utils.ErrorResponse(w, err)
		return
	}

	var source Shipper
	var dest Address
	var paymentAccount, email string
	dest.ShipperName = body["name"]
	dest.Street = body["street"]
	dest.City = body["city"]
	dest.ProvinceCode = body["province_code"]
	dest.Country = body["country_code"]
	dest.PostalCode = body["postal_code"]
	dest.Phone, _ = strconv.Atoi(body["phone"])
	source.Parcels = []Parcel{
		{
			Length: body["length"],
			Width: body["width"],
			Height: body["height"],
		},
	}
	getShippingInfoQuery, err := database.DB.Prepare(getShippingInfoAndPaymentSql)
	defer getShippingInfoQuery.Close()
	err = getShippingInfoQuery.QueryRow(tokenClaims.CompanyID).Scan(&source.Address.ShipperName, &source.Address.Street, &source.Address.City, &source.Address.ProvinceCode, &source.Address.Country, &source.Address.PostalCode, &source.Address.Phone, &paymentAccount, &email)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET")
	intAmount, err := strconv.Atoi(body["label_total"])
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
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://ct.soa-gw.canadapost.ca/rs/"+os.Getenv("CANADA_POST_CUSTNUM")+"/ncshipment", bytes.NewBuffer([]byte(xmlBody)))
	if err != nil {
		_, _ = refund.New(refundParam)
		utils.ErrorResponse(w, err)
		return
	}

	authEncoded := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CANADA_POST_USER")+":"+os.Getenv("CANADA_POST_PASS")))
	req.Header.Add("Accept", "application/vnd.cpc.ncshipment-v4+xml")
	req.Header.Add("Content-Type", "application/vnd.cpc.ncshipment-v4+xml")
	req.Header.Add("Authorization", "Basic "+authEncoded)
	req.Header.Add("Accept-language", "en-CA")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		_, _ = refund.New(refundParam)
		utils.ErrorResponse(w, err)
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	var canadaPostBody CanadaPostXML

	if err = xml.Unmarshal(respBody, &canadaPostBody); err != nil{
		_, _ = refund.New(refundParam)
		utils.ErrorResponse(w, err)
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
		utils.ErrorResponse(w ,err)
		return
	}

	getLabelReq, err := http.NewRequest("GET", labelLink, nil)
	getLabelReq.SetBasicAuth(os.Getenv("CANADA_POST_USER"), os.Getenv("CANADA_POST_PASS"))
	getLabelReq.Header.Add("Accept","application/pdf")
	labelResp, err := client.Do(getLabelReq)
	label, err := ioutil.ReadAll(labelResp.Body)

	receipt := LabelPurchase{
		PostalCode: body["postal_code"],
		Total: body["label_total"],
		TrackingNumber: canadaPostBody.Tracking,
	}

	var tmplBuffer bytes.Buffer
	tmpl := template.Must(template.ParseFiles("assets/templates/labelPurchase.html"))
	tmpl.Execute(&tmplBuffer, receipt)

	err = SendEmail("ahmad.ayubi@hotmail.com", "Shipping Label Purchased", tmplBuffer.String(), label)
	if err != nil {
		_, _ = refund.New(refundParam)
		_ = SendEmail("ahmad.ayubi@hotmail.com", "Shipping Label Purchased", "Refund Label Purchase", nil)
		utils.ErrorResponse(w ,err)
		return
	}



	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(label)
}

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
	xml += `<length>`+source.Parcels[0].Length+`</length>`
	xml += `<width>`+source.Parcels[0].Width+`</width>`
	xml += `<height>`+source.Parcels[0].Height+`</height>`
	xml += `</dimensions>`
	xml += `</parcel-characteristics>`
	xml += `<preferences>`
	xml += `<show-packing-instructions>true</show-packing-instructions>`
	xml += `</preferences>`
	xml += `</delivery-spec>`
	xml += `</non-contract-shipment>`
	return xml
}