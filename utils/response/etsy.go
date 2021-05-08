package response

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
