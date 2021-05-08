package response

type BasicMessage struct {
	Message string `json:"message"`
}

type Message struct {
	Message string `json:"message"`
	Status int32 `json:"status"`
}

type User struct {
	Email string `json:"email"`
	Name string	`json:"name"`
	CompanyID string `json:"company_id"`
	IsHead bool `json:"is_head"`
	ID string `json:"id"`
	IsApproved bool `json:"is_approved"`
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

type PermAuth struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code string `json:"code"`
}

type Employee struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Approved bool `json:"is_approved"`
	ID string `json:"id"`
}

type PostageAddress struct {
	ShipperName string `json:"shipper"`
	Street string `json:"street"`
	City string `json:"city"`
	ProvinceCode string `json:"province_code"`
	Country string `json:"country"`
	PostalCode string `json:"postal_code"`
	Phone int `json:"phone"`
}

type Company struct {
	CompanyName string         `json:"company_name"`
	HeadID      string         `json:"head_id"`
	TotalDue    float64        `json:"total_due"`
	Address     PostageAddress `json:"address"`
}

type Parcel struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Length float64 `json:"length"`
	Width float64 `json:"width"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
}

type Shipper struct {
	Address PostageAddress `json:"address"`
	Parcels []Parcel       `json:"parcels"`
}

type UrlResponse struct {
	Url string `json:"url"`
}

type Orders struct {
	Orders []Order `json:"orders"`
}


type Order struct {
	Type            string         `json:"type"`
	OrderID         string         `json:"order_id"`
	CreatedAt       string         `json:"created_at"`
	Total           string         `json:"total"`
	Subtotal        string         `json:"subtotal"`
	Tax             string         `json:"tax"`
	Currency        string         `json:"currency"`
	Name            string         `json:"name"`
	WasPaid         bool           `json:"was_paid"`
	Items           []Item         `json:"items"`
	ShippingAddress Address `json:"shipping_address"`
}

type Item struct {
	Sku      string `json:"sku"`
	Title    string `json:"title"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

type Address struct {
	Name       string      `json:"name"`
	Address1   string      `json:"address1"`
	Address2   interface{} `json:"address2"`
	City       string      `json:"city"`
	Province   string      `json:"province"`
	Country    string      `json:"country"`
	PostalCode string      `json:"postal_code"`
	Formatted  interface{} `json:"formatted"`
}