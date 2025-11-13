package main

import (
	"sync"
	"time"
)

// UserService
type User struct {
	Id        string
	Name      string
	Email     string
	Phone     string
	Addresses []Address
}

type Address struct {
	Id        string
	Street    string
	City      string
	State     string
	ZipCode   string
	Country   string
	IsDefault bool // Is this the default shipping address?
}

//ProductService
//Why didn't include product in Inventory Service
// - **Product** = Catalog information (name, price, description, category)
// - **Inventory** = Stock tracking (how many available)

// These are **two different concerns**:
// - Product catalog can exist even with 0 stock
// - Inventory tracks quantities for existing products
type Product struct {
	Id          string
	Name        string
	Price       float64
	Description string
	Category    string
}

// InventoryService
type Inventory struct {
	Stock map[string]int //(product_id, quantity)
	mu    sync.RWMutex
}

// CartService
type Cart struct {
	Id     string
	UserId string
	//Price float64 -> no price should not be stored in the cart because price is dynamic,
	//if we store
	//Price:     1000,  // Price at time of adding
	// 2 hours later (before checkout), price changes to $800 (sale!)
	// But cart still shows $1000!
	// User sees: "Total: $1000"
	// But actual checkout price: $800
	// User confused: "Why did price change?"
	//Total  float64  // ❌ Remove - should be calculated dynamically, not stored
	Items map[string]int //(product_id, quantity), (product-> quantity does not make sense, because we should not
	//store entire object as a key in map)
}

//OrderService

type OrderService struct {
	paymentGateway PaymentGateway //interface for multiple payment types
	inventory      *InventoryService
	cart           *CartService
	shipping       *ShippingService
	orders         map[string]*Order //in real life it will be repo which will store in db
}

func (s *OrderService) PlaceOrder(userId string, cartId string, addressId string, paymentType PaymentType) (*Order, error) {
	//will define this later
	return nil, nil
}

func (s *OrderService) ViewOrderDetails(orderId string) (*Order, error) {
	//will define this later
	return nil, nil
}

func (s *OrderService) CancelOrder(orderId string) error {
	//will define this later
	return nil
}

func (s *OrderService) GetOrderHistory(userId string) ([]Order, error) {
	return nil, nil
}

func (s *OrderService) TrackOrder(orderId string) (*Shipment, error) {
	return nil, nil
}

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Confirmed OrderStatus = "CONFIRMED"
	Shipped   OrderStatus = "SHIPPED"
	Delivered OrderStatus = "DELIVERED"
	Cancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	// OrderItems string
	ProductId    string
	ProductName  string
	Quantity     int
	PriceAtOrder float64
}

type Order struct {
	// CartId string //already discussed why not needed
	Id     string
	Status OrderStatus
	UserId string
	//OrderItems map[string]OrderItems//(order-item-id -> OrderItems) //this will be more optimised
	// why map should not be used, explained in readme
	OrderItems        []OrderItem
	TotalAmount       float64
	PaymentId         string
	CreatedAt         time.Time
	ShippingAddressId string
}

// ShipmentService
type ShippingService struct {
	provider ShippingProviderInterface
}

type ShippingProviderInterface{
	ShipOrder()
}

func(s *FedXService) ShipOrder(){

}


type ShippingProvider string

const (
	FedEx ShippingProvider = "FEDEX"
	UPS   ShippingProvider = "UPS"
	DHL   ShippingProvider = "DHL"
)

type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "PENDING"
	ShipmentInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentDelivered ShipmentStatus = "DELIVERED"
)

type Shipment struct {
	Id                string
	OrderId           string           //if the order is placed otheriwse empty
	Provider          ShippingProvider //(Delhivery, Ekart etc)
	Status            ShipmentStatus
	EstimatedDelivery time.Time
	TrackingNumber    string
	ActualDelivery    *time.Time
}

// PaymentService
type PaymentService struct {
	gateway PaymentGateway //payment service has a paymentgateway
}

func (s *PaymentService) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	//add business logic (validation, logging, etc.)

	//delegate to gateway
	return s.gateway.ProcessPayment(orderId, amount, paymentType) //can be any gateway based on runtime
}

func (s *PaymentService) RefundPayment(paymentId string, amount float64) error {
	//business logic
	//delegate to gateway
	return s.gateway.RefundPayment(paymentId, amount)
}

type PaymentGateway interface {
	ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error)
	RefundPayment(paymentId string, amount float64) error
}

// implementations of payment gateway
type StripeGateway struct {
	apiKey string // ✅ Shows you understand external API auth
}

func (s *StripeGateway) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	return nil, nil
}

func (s *StripeGateway) RefundPayment(paymentId string, amount float64) error {
	return nil
}

type RazorpayGateway struct {
	keyId     string // ✅ Shows you know different services
	keySecret string //    have different auth formats
}

func (r *RazorpayGateway) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	return nil, nil
}
func (r *RazorpayGateway) RefundPayment(paymentId string, amount float64) error {
	return nil
}

type MockPaymentGateway struct {
}

func (m *MockPaymentGateway) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	return nil, nil
}
func (m *MockPaymentGateway) RefundPayment(paymentId string, amount float64) error {
	return nil
}

type PaymentStatus string

const (
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
)

type PaymentType string

const (
	CreditCard PaymentType = "CREDITCARD"
	DebitCard  PaymentType = "DEBITCARD"
	UPI        PaymentType = "UPI"
)

type Payment struct {
	Id            string
	OrderId       string
	Type          PaymentType
	Status        PaymentStatus
	Amount        float64
	TransactionId string
	Timestamp     time.Time
}
