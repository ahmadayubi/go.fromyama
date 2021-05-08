package response

import "encoding/xml"

type CanadaPostPostageResponse struct {
	Tracking string `xml:"tracking-pin"`
	Links []struct {
		Name string `xml:"rel,attr"`
		Link string `xml:"href,attr"`
	} `xml:"links>link"`
}

type CanadaPostLabelPurchase struct {
	PostalCode string
	Total float64
	TrackingNumber string
}

type CanadaPostRatesResponse struct {
	XMLName    xml.Name `xml:"price-quotes"`
	Text       string   `xml:",chardata"`
	Xmlns      string   `xml:"xmlns,attr"`
	PriceQuote []struct {
		Text        string `xml:",chardata"`
		ServiceCode string `xml:"service-code"`
		ServiceLink struct {
			Text      string `xml:",chardata"`
			Rel       string `xml:"rel,attr"`
			Href      string `xml:"href,attr"`
			MediaType string `xml:"media-type,attr"`
		} `xml:"service-link"`
		ServiceName  string `xml:"service-name"`
		PriceDetails struct {
			Text  string `xml:",chardata"`
			Base  string `xml:"base"`
			Taxes struct {
				Text string `xml:",chardata"`
				Gst  string `xml:"gst"`
				Pst  string `xml:"pst"`
				Hst  struct {
					Text    string `xml:",chardata"`
					Percent string `xml:"percent,attr"`
				} `xml:"hst"`
			} `xml:"taxes"`
			Due     string `xml:"due"`
			Options struct {
				Text   string `xml:",chardata"`
				Option struct {
					Text        string `xml:",chardata"`
					OptionCode  string `xml:"option-code"`
					OptionName  string `xml:"option-name"`
					OptionPrice string `xml:"option-price"`
					Qualifier   struct {
						Text     string `xml:",chardata"`
						Included string `xml:"included"`
					} `xml:"qualifier"`
				} `xml:"option"`
			} `xml:"options"`
			Adjustments struct {
				Text       string `xml:",chardata"`
				Adjustment []struct {
					Text           string `xml:",chardata"`
					AdjustmentCode string `xml:"adjustment-code"`
					AdjustmentName string `xml:"adjustment-name"`
					AdjustmentCost string `xml:"adjustment-cost"`
					Qualifier      struct {
						Text    string `xml:",chardata"`
						Percent string `xml:"percent"`
					} `xml:"qualifier"`
				} `xml:"adjustment"`
			} `xml:"adjustments"`
		} `xml:"price-details"`
		WeightDetails   string `xml:"weight-details"`
		ServiceStandard struct {
			Text                 string `xml:",chardata"`
			AmDelivery           bool `xml:"am-delivery"`
			GuaranteedDelivery   bool `xml:"guaranteed-delivery"`
			ExpectedTransitTime  string `xml:"expected-transit-time"`
			ExpectedDeliveryDate string `xml:"expected-delivery-date"`
		} `xml:"service-standard"`
	} `xml:"price-quote"`
}

type CanadaPostRate struct {
	ServiceCode string `json:"service-code"`
	ServiceName  string `json:"service-name"`
	PriceDetails CanadaPostRatePrice `json:"price-details"`
	WeightDetails   string `json:"weight-details"`
	AmDelivery           bool `json:"am-delivery"`
	GuaranteedDelivery   bool `json:"guaranteed-delivery"`
	ExpectedTransitTime  string `json:"expected-transit-time"`
	ExpectedDeliveryDate string `json:"expected-delivery-date"`
}

type CanadaPostRatePrice struct {
	Base  string `json:"base"`
	Gst string `json:"gst"`
	Pst string `json:"pst"`
	Hst string `json:"hst"`
	Due     string `json:"due"`
	Options struct {
		Option struct {
			OptionCode  string `json:"option-code"`
			OptionName  string `json:"option-name"`
			OptionPrice string `json:"option-price"`
		} `json:"option"`
	} `json:"options"`
	Adjustments []CanadaPostAdjustment `json:"adjustments"`
}

type CanadaPostAdjustment struct {
	AdjustmentCode string `json:"adjustment-code"`
	AdjustmentName string `json:"adjustment-name"`
	AdjustmentCost string `json:"adjustment-cost"`
}