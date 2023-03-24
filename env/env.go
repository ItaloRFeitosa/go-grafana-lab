package env

import "os"

var (
	CheckoutURL  = os.Getenv("CHECKOUT_URL")
	PaymentURL   = os.Getenv("PAYMENT_URL")
	OrderURL     = os.Getenv("ORDER_URL")
	WarehouseURL = os.Getenv("WAREHOUSE_URL")
)
