package main

import "log"

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

type PaymentService struct {
	Gateway   PaymentGateway   //Direct interface injection
	Processor PaymentProcessor //Direct interface injection
}

func (s *PaymentService) ProcessPayment(payment *Payment) error {
	log.Printf("\n📌 Processing Payment: %s\n", payment.ID)

	payment.Status = Pending
	err := s.Processor.Process(payment, s.Gateway)
	if err != nil {
		payment.Status = Failed
		return err
	}

	log.Printf("✅ Payment %s processed successfully!\n", payment.ID)
	return nil
}

func main() {
	log.Println("Payment System WITHOUT Factory Pattern")

	// ============================================
	// Test 1: Credit Card payment with Stripe
	// ============================================
	log.Println("=== Test 1: Credit Card + Stripe ===")

	// ❌ Problem: Must create service for each combination!
	paymentService1 := &PaymentService{
		Gateway:   &StripeGateway{APIKey: "sk_test_123"},
		Processor: &CreditCardProcessor{},
	}

	payment1 := &Payment{
		ID:       "pay_001",
		Type:     CreditCard,
		Amount:   100.50,
		Currency: "USD",
		Status:   Pending,
	}
	paymentService1.ProcessPayment(payment1)

	// Test 2: UPI payment with Razorpay
	// ============================================
	log.Println("\n=== Test 2: UPI + Razorpay ===")

	// ❌ Problem: Must create ANOTHER service instance!
	paymentService2 := &PaymentService{
		Gateway:   &RazorpayGateway{APIKey: "rzp_test_456"},
		Processor: &UPIProcessor{},
	}
	payment2 := &Payment{
		ID:       "pay_002",
		Type:     UPI,
		Amount:   5000.00,
		Currency: "INR",
		Status:   Pending,
	}

	paymentService2.ProcessPayment(payment2)

	// ============================================
	// Test 3: Credit Card payment with Razorpay
	// ============================================
	log.Println("\n=== Test 3: Credit Card + Razorpay ===")

	// ❌ Problem: Yet ANOTHER service instance!
	paymentService3 := &PaymentService{
		Gateway:   &RazorpayGateway{APIKey: "rzp_test_456"},
		Processor: &CreditCardProcessor{},
	}

	payment3 := &Payment{
		ID:       "pay_003",
		Type:     CreditCard,
		Amount:   2500.00,
		Currency: "INR",
		Status:   Pending,
	}

	paymentService3.ProcessPayment(payment3)

	log.Println("\n✅ All payments processed!")

}

// ❌ Problems with Direct Injection (No Factory)
// Problem 1: Service Can't Decide Dynamically
// type PaymentService struct {
//     Gateway   PaymentGateway   // ❌ Fixed at creation time
//     Processor PaymentProcessor // ❌ Fixed at creation time
// }

// ❌ Can't change gateway based on payment.Currency
// -> let's say bookingservice is calling the payment service then booking servoce would have
// //to make a different payment service based on the currency (which it will get from handler when user /book/user_id)
// // which actually booking service should not decide, it should be responsibility of payment service because
// //otherwise  Payment logic leaks out of PaymentService

// ❌ WITHOUT Factory (Bad Design)

// type PaymentService struct {
// 	Gateway PaymentGateway // ❌ Only ONE gateway hardcoded
// }

// // When BookingService calls PaymentService:
// type BookingService struct {
// 	PaymentService PaymentServiceInterface
// }

// func (b *BookingService) CreateBooking(booking *Booking) error {
// 	// ❌ PROBLEM: BookingService has to decide which PaymentService?

// 	// If payment.Currency == "INR", need PaymentService with Razorpay
// 	// If payment.Currency == "USD", need PaymentService with Stripe

// 	// ❌ BookingService would need to do this:
// 	var paymentService PaymentServiceInterface
// 	if booking.Payment.Currency == "INR" {
// 		paymentService = &PaymentServiceWithRazorpay{}
// 	} else {
// 		paymentService = &PaymentServiceWithStripe{}
// 	}

// 	// ❌ Payment logic leaked to BookingService!
// 	// ❌ BookingService now knows about currencies and gateways!

// 	paymentService.ProcessPayment(booking.Payment)
// }

// Problem: BookingService shouldn't know "INR → Razorpay, USD → Stripe"! That's payment domain knowledge!

// // ❌ Can't change processor based on payment.Type -> again booking service would have to create different
// //payemtn services
// //based on what paymenttype it got from the handler
// // ❌ Must create NEW service for each combination

// // ❌ Need separate service for each combination
// paymentService1 := &PaymentService{Gateway: stripe, Processor: creditCard}
// paymentService2 := &PaymentService{Gateway: razorpay, Processor: upi}
// paymentService3 := &PaymentService{Gateway: razorpay, Processor: creditCard}

// // ❌ Can't have ONE service that handles all cases!

// // ❌ Service can't decide: "If INR, use Razorpay; if USD, use Stripe"
// // ❌ Caller must know payment domain logic
// // ❌ Payment logic leaks out of PaymentService
