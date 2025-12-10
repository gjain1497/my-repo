package myrepo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// ============================================
// MODELS
// ============================================

type User struct {
	ID             string
	Region         string
	IsInExperiment func(experimentName string) bool
}

type Booking struct {
	ID        string
	UserID    string
	VehicleID string
	StartDate string
	EndDate   string
	Payment   *Payment
}

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
	Cash       PaymentType = "CASH"
)

type PaymentStatus string

const (
	Pending PaymentStatus = "PENDING"
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
)

type CreateBookingRequest struct {
	VehicleID   string
	StartDate   string
	EndDate     string
	PaymentType PaymentType
	Amount      float64
	Currency    string
}

// ============================================
// PAYMENT GATEWAY INTERFACE
// ============================================

type PaymentGateway interface {
	Charge(payment *Payment) error
	Refund(payment *Payment) error
}

type StripeGateway struct{}

func (s *StripeGateway) Charge(payment *Payment) error {
	log.Printf("💳 [STRIPE] Charging $%.2f\n", payment.Amount)
	payment.Status = Success
	return nil
}

func (s *StripeGateway) Refund(payment *Payment) error {
	log.Printf("💳 [STRIPE] Refunding $%.2f\n", payment.Amount)
	return nil
}

type RazorpayGateway struct{}

func (r *RazorpayGateway) Charge(payment *Payment) error {
	log.Printf("💳 [RAZORPAY] Charging ₹%.2f\n", payment.Amount)
	payment.Status = Success
	return nil
}

func (r *RazorpayGateway) Refund(payment *Payment) error {
	log.Printf("💳 [RAZORPAY] Refunding ₹%.2f\n", payment.Amount)
	return nil
}

// ============================================
// PAYMENT GATEWAY FACTORY
// ============================================

type GatewayType string

const (
	Stripe   GatewayType = "STRIPE"
	Razorpay GatewayType = "RAZORPAY"
)

type PaymentGatewayFactory struct {
	Gateways map[GatewayType]PaymentGateway
}

func NewPaymentGatewayFactory() *PaymentGatewayFactory {
	return &PaymentGatewayFactory{
		Gateways: map[GatewayType]PaymentGateway{
			Stripe:   &StripeGateway{},
			Razorpay: &RazorpayGateway{},
		},
	}
}

func (f *PaymentGatewayFactory) GetGateway(gatewayType GatewayType) PaymentGateway {
	return f.Gateways[gatewayType]
}

// ============================================
// PAYMENT SERVICE INTERFACE
// ============================================

type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, payment *Payment) error
	RefundPayment(ctx context.Context, paymentID string) error
}

// ============================================
// PAYMENT SERVICE V1 (Old Flow)
// ============================================

type PaymentServiceV1 struct {
	GatewayFactory *PaymentGatewayFactory // ✅ Service chooses gateway!
}

func (s *PaymentServiceV1) ProcessPayment(ctx context.Context, payment *Payment) error {
	log.Println("📌 Using PaymentServiceV1 (Old Flow)")

	// ✅ V1 logic: Choose gateway based on currency (PAYMENT LOGIC!)
	var gatewayType GatewayType
	if payment.Currency == "INR" {
		gatewayType = Razorpay
		log.Println("✅ V1: INR currency → Using Razorpay")
	} else {
		gatewayType = Stripe
		log.Println("✅ V1: Non-INR currency → Using Stripe")
	}

	// ✅ Service selects and uses gateway
	gateway := s.GatewayFactory.GetGateway(gatewayType)
	return gateway.Charge(payment)
}

func (s *PaymentServiceV1) RefundPayment(ctx context.Context, paymentID string) error {
	log.Println("Refunding payment (V1)")
	return nil
}

// ============================================
// PAYMENT SERVICE V2 (New Flow)
// ============================================

type PaymentServiceV2 struct {
	GatewayFactory *PaymentGatewayFactory // ✅ Service chooses gateway!
}

func (s *PaymentServiceV2) ProcessPayment(ctx context.Context, payment *Payment) error {
	log.Println("📌 Using PaymentServiceV2 (New Flow)")

	// ✅ V2 logic: Different gateway selection rules (NEW PAYMENT LOGIC!)
	var gatewayType GatewayType

	// V2 has more sophisticated rules
	if payment.Currency == "INR" {
		gatewayType = Razorpay
		log.Println("✅ V2: INR currency → Using Razorpay")
	} else if payment.Amount > 10000 {
		// V2 feature: Large amounts always use Stripe
		gatewayType = Stripe
		log.Println("✅ V2: Large amount → Using Stripe")
	} else {
		gatewayType = Stripe
		log.Println("✅ V2: Default → Using Stripe")
	}

	// ✅ Service selects and uses gateway
	gateway := s.GatewayFactory.GetGateway(gatewayType)
	return gateway.Charge(payment)
}

func (s *PaymentServiceV2) RefundPayment(ctx context.Context, paymentID string) error {
	log.Println("Refunding payment (V2)")
	return nil
}

// ============================================
// VEHICLE SERVICE
// ============================================

type VehicleServiceInterface interface {
	MarkAsBooked(vehicleID string) error
	MarkAsAvailable(vehicleID string) error
}

type VehicleServiceV1 struct{}

func (v *VehicleServiceV1) MarkAsBooked(vehicleID string) error {
	log.Printf("🚗 Vehicle %s marked as booked\n", vehicleID)
	return nil
}

func (v *VehicleServiceV1) MarkAsAvailable(vehicleID string) error {
	log.Printf("🚗 Vehicle %s marked as available\n", vehicleID)
	return nil
}

// ============================================
// BOOKING SERVICE
// ============================================

type BookingService struct {
	PaymentService PaymentServiceInterface // ✅ Just interface, no gateway knowledge
	VehicleService VehicleServiceInterface
}

func (s *BookingService) CreateBooking(ctx context.Context, booking *Booking) (*Booking, error) {
	log.Println("📝 Creating booking...")

	// ✅ BookingService doesn't know about gateways!
	// Just delegates to PaymentService
	err := s.PaymentService.ProcessPayment(ctx, booking.Payment)
	if err != nil {
		return nil, err
	}

	// Mark vehicle as booked
	err = s.VehicleService.MarkAsBooked(booking.VehicleID)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Booking created successfully!")
	return booking, nil
}

// ============================================
// BOOKING SERVICE FACTORY
// ============================================

type BookingServiceFactory struct {
	PaymentServiceV1 PaymentServiceInterface // ✅ Pre-created V1
	PaymentServiceV2 PaymentServiceInterface // ✅ Pre-created V2
	VehicleService   VehicleServiceInterface
}

func NewBookingServiceFactory() *BookingServiceFactory {
	// ✅ Create shared gateway factory
	gatewayFactory := NewPaymentGatewayFactory()

	return &BookingServiceFactory{
		// ✅ V1 with gateway factory (V1 chooses gateway internally)
		PaymentServiceV1: &PaymentServiceV1{
			GatewayFactory: gatewayFactory,
		},
		// ✅ V2 with gateway factory (V2 chooses gateway internally)
		PaymentServiceV2: &PaymentServiceV2{
			GatewayFactory: gatewayFactory,
		},
		VehicleService: &VehicleServiceV1{},
	}
}

// ✅ Factory only chooses which PaymentService VERSION
// Each PaymentService decides gateway internally
func (f *BookingServiceFactory) CreateBookingService(user *User) *BookingService {
	var paymentService PaymentServiceInterface

	// ✅ HANDLER LOGIC: Choose service version based on user experiment
	if user.IsInExperiment("new_payment_flow") {
		paymentService = f.PaymentServiceV2 // ✅ V2 (new flow)
		log.Printf("🔬 User %s in experiment → Using PaymentServiceV2\n", user.ID)
	} else {
		paymentService = f.PaymentServiceV1 // ✅ V1 (old flow)
		log.Printf("📊 User %s not in experiment → Using PaymentServiceV1\n", user.ID)
	}

	return &BookingService{
		PaymentService: paymentService, // ✅ V1 or V2 (both handle gateways internally)
		VehicleService: f.VehicleService,
	}
}

// ============================================
// HANDLER
// ============================================

type BookingHandler struct {
	ServiceFactory *BookingServiceFactory
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Get user from request
	user := GetUserFromContext(ctx)
	log.Printf("👤 Request from user: %s\n", user.ID)

	// 2. ✅ Handler decides which PaymentService VERSION (V1 or V2)
	//    This is ROUTING logic (handler's job)
	bookingService := h.ServiceFactory.CreateBookingService(user)

	// 3. Parse request
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Create booking with payment data
	booking := &Booking{
		UserID:    user.ID,
		VehicleID: req.VehicleID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Payment: &Payment{
			Amount:   req.Amount,
			Currency: req.Currency, // ✅ Handler passes currency (doesn't interpret it!)
			Type:     req.PaymentType,
			Status:   Pending,
		},
	}

	// 5. ✅ PaymentService (V1 or V2) decides gateway based on currency
	//    This is PAYMENT LOGIC (service's job)
	result, err := bookingService.CreateBooking(ctx, booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ============================================
// CONTEXT HELPERS
// ============================================

type contextKey string

const userContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		// Return default user for demo
		return &User{
			ID:     "user_123",
			Region: "US",
			IsInExperiment: func(experimentName string) bool {
				// Simulate A/B test (50% in experiment)
				return user.ID[len(user.ID)-1]%2 == 0
			},
		}
	}
	return user
}

// ============================================
// MAIN
// ============================================

func main() {
	log.Println("🚀 Starting VRS Booking Service...")

	// ✅ Create factory once at startup
	serviceFactory := NewBookingServiceFactory()

	// ✅ Create handler with factory
	bookingHandler := &BookingHandler{
		ServiceFactory: serviceFactory,
	}

	// Register routes
	http.HandleFunc("/bookings", bookingHandler.CreateBooking)

	log.Println("✅ Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
