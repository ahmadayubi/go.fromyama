package response

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

type LocationResponse struct {
	Locations []struct {
		ID                    int64       `json:"id"`
		Name                  string      `json:"name"`
	} `json:"locations"`
}