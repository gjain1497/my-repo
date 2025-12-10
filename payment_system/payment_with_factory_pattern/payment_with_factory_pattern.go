package main

import (
	"fmt"
	"log"
)

type Payment struct {
	ID       string
	Amount   float64
	Currency string
	Type     PaymentType
	Status   PaymentStatus
}

type PaymentType string

const (
	CreditCard PaymentType = "CREDIT_CARD"
	UPI        PaymentType = "UPI"
)

type GatewayType string

const (
	Stripe   GatewayType = "STRIPE"
	Razorpay GatewayType = "RAZORPAY"
)

type PaymentStatus string

const (
	Pending PaymentStatus = "PENDING"
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
)

type PaymentGateway interface {
	Charge(payment *Payment) error
}
type StripeGateway struct {
	APIKey string
}

func (s *StripeGateway) Charge(payment *Payment) error {
	log.Printf("💳 [STRIPE] Charging $%.2f\n", payment.Amount)
	payment.Status = Success
	return nil
}

type RazorpayGateway struct {
	APIKey string
}

func (r *RazorpayGateway) Charge(payment *Payment) error {
	log.Printf("💳 [RAZORPAY] Charging ₹%.2f\n", payment.Amount)
	payment.Status = Success
	return nil
}

type PaymentProcessor interface {
	Process(payment *Payment, gateway PaymentGateway) error
}

type CreditCardProcessor struct{}

func (c *CreditCardProcessor) Process(payment *Payment, gateway PaymentGateway) error {
	log.Println("🔒 Validating credit card...")
	return gateway.Charge(payment)
}

type UPIProcessor struct{}

func (u *UPIProcessor) Process(payment *Payment, gateway PaymentGateway) error {
	log.Println("📱 Processing UPI payment...")
	return gateway.Charge(payment)
}

// type PaymentService struct {
// 	Gateway   PaymentGateway   //Direct interface injection
// 	Processor PaymentProcessor //Direct interface injection
// }

type PaymentService struct {
	GatewayFactory   *PaymentGatewayFactory   //Factory injecion instead of Direct interface injection
	ProcessorFactory *PaymentProcessorFactory //Factory injecion instead of Direct interface injection
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		GatewayFactory:   NewPaymentGatewayFactory(),
		ProcessorFactory: NewPaymentProcessorFactory(),
	}
}

type PaymentGatewayFactory struct {
	Gateways map[GatewayType]PaymentGateway
}

func NewPaymentGatewayFactory() *PaymentGatewayFactory {
	return &PaymentGatewayFactory{
		Gateways: map[GatewayType]PaymentGateway{
			Stripe:   &StripeGateway{APIKey: "stp_test_123"},
			Razorpay: &RazorpayGateway{APIKey: "rzp_test_123"},
		},
	}
}

func (f *PaymentGatewayFactory) GetGateway(gatewayType GatewayType) (PaymentGateway, error) {
	gateway, exists := f.Gateways[gatewayType]
	if !exists {
		return nil, fmt.Errorf("gateway %s not found", gatewayType)
	}
	return gateway, nil
}

type PaymentProcessorFactory struct {
	Processors map[PaymentType]PaymentProcessor
}

func NewPaymentProcessorFactory() *PaymentProcessorFactory {
	return &PaymentProcessorFactory{
		Processors: map[PaymentType]PaymentProcessor{
			CreditCard: &CreditCardProcessor{},
			UPI:        &UPIProcessor{},
		},
	}
}

func (f *PaymentProcessorFactory) GetProcessor(paymentType PaymentType) (PaymentProcessor, error) {
	processor, exists := f.Processors[paymentType]
	if !exists {
		return nil, fmt.Errorf("processor %s not found", paymentType)
	}
	return processor, nil
}

func (s *PaymentService) ProcessPayment(payment *Payment) error {
	log.Printf("\n📌 Processing Payment: %s\n", payment.ID)

	//Service decides processor based on payment type
	processor, err := s.ProcessorFactory.GetProcessor(payment.Type)
	if err != nil {
		return err
	}

	//Service decides payment gateway based on currency
	var gatewayType GatewayType
	if payment.Currency == "INR" {
		gatewayType = Razorpay
		log.Println("✅ Currency is INR → Using Gateway Razorpay")
	} else {
		gatewayType = Stripe
		log.Println("✅ Currency is NOT INR → Using Gateway Stripe")
	}
	gateway, err := s.GatewayFactory.GetGateway(gatewayType)

	//Process payment
	payment.Status = Pending
	err = processor.Process(payment, gateway)
	if err != nil {
		payment.Status = Failed
		return err
	}

	log.Printf("✅ Payment %s processed successfully!\n", payment.ID)
	return nil
}

func main() {
	log.Println("Payment System WITH Factory Pattern")

	// ============================================
	// Test 1: Credit Card payment in USD
	// ============================================
	log.Println("=== Test 1: Credit Card + Stripe ===")

	//one service that handled all the cases
	paymentService := NewPaymentService()

	payment1 := &Payment{
		ID:       "pay_001",
		Type:     CreditCard,
		Amount:   100.50,
		Currency: "USD",
		Status:   Pending,
	}
	paymentService.ProcessPayment(payment1)

	// Test 2: UPI payment in INR
	// ============================================

	payment2 := &Payment{
		ID:       "pay_002",
		Type:     UPI,
		Amount:   5000.00,
		Currency: "INR",
		Status:   Pending,
	}

	paymentService.ProcessPayment(payment2)

	// ============================================
	// Test 3: Credit Card payment in INR
	// ============================================
	payment3 := &Payment{
		ID:       "pay_003",
		Type:     CreditCard,
		Amount:   2500.00,
		Currency: "INR",
		Status:   Pending,
	}

	paymentService.ProcessPayment(payment3)

	log.Println("\n✅ All payments processed with ONE service!")
}
