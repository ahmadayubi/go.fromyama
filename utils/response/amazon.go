package response

import "encoding/xml"

type AmazonUnfulfilledResponse struct {
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

type AmazonOrderDetailResponse struct {
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
