package response

import "encoding/xml"

type User struct {
	Email string `json:"email"`
	Name string	`json:"name"`
	CompanyID string `json:"company_id"`
	IsHead bool `json:"is_head"`
	ID string `json:"id"`
}

type Token struct {
	Raw string `json:"token"`
}

type UserDetails struct {
	UserData User `json:"user"`
	CompanyName string `json:"company_name"`
}

type CallbackResponse struct {
	AccessToken string `json:"access_token"`
}

type LocationResponse struct {
	Locations []struct {
		ID                    int64       `json:"id"`
		Name                  string      `json:"name"`
	} `json:"locations"`
}

type PermAuth struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code string `json:"code"`
}

type ShopifyFulfillmentRequest struct {
	Fulfillment ShopifyFulfillmentData `json:"fulfillment"`
	NotifyCustomer bool `json:"notify_customer"`
}

type ShopifyFulfillmentData struct {
	LocationID int `json:"location_id"`
	TrackingNumber string `json:"tracking_number"`
	TrackingCompany string `json:"tracking_company"`
}

type ShopifyFulfillmentResponse struct {
	Fulfillment struct {
		ID              int64       `json:"id"`
		OrderID         int64       `json:"order_id"`
		Status          string      `json:"status"`
		CreatedAt       string      `json:"created_at"`
		Service         string      `json:"service"`
		UpdatedAt       string      `json:"updated_at"`
		TrackingCompany string      `json:"tracking_company"`
		ShipmentStatus  interface{} `json:"shipment_status"`
		LocationID      int64       `json:"location_id"`
		LineItems       []struct {
			ID                         int64         `json:"id"`
			VariantID                  int64         `json:"variant_id"`
			Title                      string        `json:"title"`
			Quantity                   int           `json:"quantity"`
			Sku                        string        `json:"sku"`
			VariantTitle               interface{}   `json:"variant_title"`
			Vendor                     string        `json:"vendor"`
			FulfillmentService         string        `json:"fulfillment_service"`
			ProductID                  int64         `json:"product_id"`
			RequiresShipping           bool          `json:"requires_shipping"`
			Taxable                    bool          `json:"taxable"`
			GiftCard                   bool          `json:"gift_card"`
			Name                       string        `json:"name"`
			VariantInventoryManagement string        `json:"variant_inventory_management"`
			Properties                 []interface{} `json:"properties"`
			ProductExists              bool          `json:"product_exists"`
			FulfillableQuantity        int           `json:"fulfillable_quantity"`
			Grams                      int           `json:"grams"`
			Price                      string        `json:"price"`
			TotalDiscount              string        `json:"total_discount"`
			FulfillmentStatus          string        `json:"fulfillment_status"`
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
			DiscountAllocations []interface{} `json:"discount_allocations"`
			Duties              []interface{} `json:"duties"`
			AdminGraphqlAPIID   string        `json:"admin_graphql_api_id"`
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
		TrackingNumber  string   `json:"tracking_number"`
		TrackingNumbers []string `json:"tracking_numbers"`
		TrackingURL     string   `json:"tracking_url"`
		TrackingUrls    []string `json:"tracking_urls"`
		Receipt         struct {
		} `json:"receipt"`
		Name              string `json:"name"`
		AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	} `json:"fulfillment"`
}

type ShopifyUnfulfilledResponse struct {
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

type EtsyUnfulfilledResponse struct {
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

type CanadaPostPostageResponse struct {
	Tracking string `xml:"tracking-pin"`
	Links []struct {
		Name string `xml:"rel,attr"`
		Link string `xml:"href,attr"`
	} `xml:"links>link"`
}

type CanadaPostLabelPurchase struct {
	PostalCode string
	Total string
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
			AmDelivery           string `xml:"am-delivery"`
			GuaranteedDelivery   string `xml:"guaranteed-delivery"`
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
	AmDelivery           string `json:"am-delivery"`
	GuaranteedDelivery   string `json:"guaranteed-delivery"`
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