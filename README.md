# go.fromyama - REST API ![Go](https://github.com/FromYama/go.fromyama/workflows/Go/badge.svg)

## Endpoints:

- ### /company
 
    - `POST /register`
    - `DELETE /unregister`
    - `GET /platforms`
    - `GET /details`
    - `GET /employee/all`
    - `PUT /employee/approve/{companyID}`
    - `GET /employee/approved/{companyID}`
    - `POST /add/payment/method`
    - `POST /add/payment/charge`
    - `POST /add/parcel`
    - `GET /shipper`
  
- ### /user
    
  - `POST /signup`
  - `POST /login`
  - `GET /refresh`
  - `GET /details`
  
- ### /postage

  - `POST /buy/canadapost`
  - `GET /rates/canadapost`
  
- ### /shopify

  - `POST /generate-link`
  - `GET /callback`
  - `GET /locations`
  - `GET /orders/all`
  - `POST /fulfill`
  
- ### /etsy
  
  - `POST /generate-link`
  - `GET /callback`
  - `GET /orders/all`
  - `POST /fulfill`
  
- ### /amazon
  
  - `POST /authorize`
  - `GET /orders/all`
  
