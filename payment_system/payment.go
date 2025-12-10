package main

// // # 🎯 Payment System with Both Processor + Gateway

// // Here's the complete example with **both** Payment Processor (CreditCard/UPI/Cash) **and** Payment Gateway (Stripe/Razorpay):

// // ---

// // ## 📋 Complete Code

// // ```go
// package main

// import (
//     "errors"
//     "fmt"
//     "strings"
//     "sync"
//     "time"
// )

// // ============================================
// // MODELS (Data Only)
// // ============================================

// type PaymentType string

// const (
//     CreditCard PaymentType = "CREDIT_CARD"
//     UPI        PaymentType = "UPI"
//     Cash       PaymentType = "CASH"
// )

// type GatewayType string

// const (
//     Stripe   GatewayType = "STRIPE"
//     Razorpay GatewayType = "RAZORPAY"
// )

// type PaymentStatus string

// const (
//     Pending PaymentStatus = "PENDING"
//     Success PaymentStatus = "SUCCESS"
//     Failed  PaymentStatus = "FAILED"
// )

// type Payment struct {
//     ID            string
//     Type          PaymentType   // ✅ CreditCard, UPI, Cash
//     Amount        float64
//     Status        PaymentStatus
//     Metadata      map[string]string
//     GatewayType   GatewayType   // ✅ Which gateway was used
//     TransactionID string        // ✅ Gateway transaction ID
//     CreatedAt     time.Time
// }

// type GatewayResponse struct {
//     TransactionID string
//     Status        string
//     Message       string
//     RawResponse   map[string]interface{}
// }

// // ============================================
// // HELPER STRUCTS
// // ============================================

// type Logger struct{}

// func (l *Logger) Info(msg string) {
//     fmt.Printf("[INFO] %s\n", msg)
// }

// type Config struct {
//     MaxAmount float64
// }

// // ============================================
// // BASE PROCESSOR (Code Reuse via Embedding)
// // ============================================

// type BaseProcessor struct {
//     Logger *Logger
//     Config *Config
// }

// func (b *BaseProcessor) Log(msg string) {
//     b.Logger.Info(msg)
// }

// func (b *BaseProcessor) ValidateAmount(amount float64) bool {
//     return amount > 0 && amount < b.Config.MaxAmount
// }

// // ============================================
// // PAYMENT GATEWAY INTERFACE
// // ============================================

// type PaymentGateway interface {
//     Charge(payment *Payment) (*GatewayResponse, error)
//     Refund(transactionID string, amount float64) (*GatewayResponse, error)
// }

// // ============================================
// // GATEWAY IMPLEMENTATIONS
// // ============================================

// type StripeGateway struct {
//     APIKey string
// }

// func (s *StripeGateway) Charge(payment *Payment) (*GatewayResponse, error) {
//     fmt.Printf("💳 [STRIPE] Charging $%.2f for payment %s\n", payment.Amount, payment.ID)

//     // Simulate Stripe API call
//     return &GatewayResponse{
//         TransactionID: "ch_stripe_" + payment.ID,
//         Status:        "succeeded",
//         Message:       "Payment successful via Stripe",
//         RawResponse:   map[string]interface{}{"gateway": "stripe", "amount": payment.Amount},
//     }, nil
// }

// func (s *StripeGateway) Refund(transactionID string, amount float64) (*GatewayResponse, error) {
//     fmt.Printf("💳 [STRIPE] Refunding $%.2f for transaction %s\n", amount, transactionID)

//     return &GatewayResponse{
//         TransactionID: "re_stripe_" + transactionID,
//         Status:        "succeeded",
//         Message:       "Refund successful via Stripe",
//         RawResponse:   map[string]interface{}{"gateway": "stripe"},
//     }, nil
// }

// type RazorpayGateway struct {
//     APIKey    string
//     APISecret string
// }

// func (r *RazorpayGateway) Charge(payment *Payment) (*GatewayResponse, error) {
//     fmt.Printf("💳 [RAZORPAY] Charging ₹%.2f for payment %s\n", payment.Amount, payment.ID)

//     // Simulate Razorpay API call
//     return &GatewayResponse{
//         TransactionID: "pay_razorpay_" + payment.ID,
//         Status:        "captured",
//         Message:       "Payment successful via Razorpay",
//         RawResponse:   map[string]interface{}{"gateway": "razorpay", "amount": payment.Amount},
//     }, nil
// }

// func (r *RazorpayGateway) Refund(transactionID string, amount float64) (*GatewayResponse, error) {
//     fmt.Printf("💳 [RAZORPAY] Refunding ₹%.2f for transaction %s\n", amount, transactionID)

//     return &GatewayResponse{
//         TransactionID: "rfnd_razorpay_" + transactionID,
//         Status:        "processed",
//         Message:       "Refund successful via Razorpay",
//         RawResponse:   map[string]interface{}{"gateway": "razorpay"},
//     }, nil
// }

// // ============================================
// // GATEWAY FACTORY
// // ============================================

// type PaymentGatewayFactory struct {
//     Gateways map[GatewayType]PaymentGateway
// }

// func NewPaymentGatewayFactory() *PaymentGatewayFactory {
//     return &PaymentGatewayFactory{
//         Gateways: map[GatewayType]PaymentGateway{
//             Stripe: &StripeGateway{
//                 APIKey: "sk_test_stripe_123",
//             },
//             Razorpay: &RazorpayGateway{
//                 APIKey:    "rzp_test_123",
//                 APISecret: "secret_123",
//             },
//         },
//     }
// }

// func (f *PaymentGatewayFactory) GetGateway(gatewayType GatewayType) (PaymentGateway, error) {
//     gateway, exists := f.Gateways[gatewayType]
//     if !exists {
//         return nil, fmt.Errorf("gateway %s not found", gatewayType)
//     }
//     return gateway, nil
// }

// // ============================================
// // PAYMENT PROCESSOR INTERFACE
// // ============================================

// type PaymentProcessor interface {
//     Process(payment *Payment, gateway PaymentGateway) error
//     Refund(payment *Payment, gateway PaymentGateway) error
//     Validate(payment *Payment) bool
// }

// // ============================================
// // PROCESSOR IMPLEMENTATIONS
// // ============================================

// type CreditCardProcessor struct {
//     BaseProcessor // ✅ Embedding for code reuse
// }

// func (c *CreditCardProcessor) Process(payment *Payment, gateway PaymentGateway) error {
//     c.Log("Processing credit card payment")

//     if !c.ValidateAmount(payment.Amount) {
//         return errors.New("invalid amount")
//     }

//     // Credit card specific logic
//     cardNumber := payment.Metadata["card_number"]
//     fmt.Printf("🔒 Validating card: %s\n", maskCardNumber(cardNumber))

//     // ✅ Use gateway to charge
//     response, err := gateway.Charge(payment)
//     if err != nil {
//         return err
//     }

//     // Update payment with gateway response
//     payment.TransactionID = response.TransactionID
//     payment.Status = Success

//     return nil
// }

// func (c *CreditCardProcessor) Refund(payment *Payment, gateway PaymentGateway) error {
//     c.Log("Refunding credit card payment")

//     // ✅ Use gateway to refund
//     response, err := gateway.Refund(payment.TransactionID, payment.Amount)
//     if err != nil {
//         return err
//     }

//     fmt.Printf("✅ Refund successful: %s\n", response.TransactionID)
//     return nil
// }

// func (c *CreditCardProcessor) Validate(payment *Payment) bool {
//     cardNumber := payment.Metadata["card_number"]
//     cvv := payment.Metadata["cvv"]
//     return len(cardNumber) == 16 && len(cvv) == 3
// }

// type UPIProcessor struct {
//     BaseProcessor // ✅ Embedding for code reuse
// }

// func (u *UPIProcessor) Process(payment *Payment, gateway PaymentGateway) error {
//     u.Log("Processing UPI payment")

//     if !u.ValidateAmount(payment.Amount) {
//         return errors.New("invalid amount")
//     }

//     // UPI specific logic
//     upiID := payment.Metadata["upi_id"]
//     fmt.Printf("📱 Processing UPI ID: %s\n", upiID)

//     // ✅ Use gateway to charge
//     response, err := gateway.Charge(payment)
//     if err != nil {
//         return err
//     }

//     // Update payment with gateway response
//     payment.TransactionID = response.TransactionID
//     payment.Status = Success

//     return nil
// }

// func (u *UPIProcessor) Refund(payment *Payment, gateway PaymentGateway) error {
//     u.Log("Refunding UPI payment")

//     // ✅ Use gateway to refund
//     response, err := gateway.Refund(payment.TransactionID, payment.Amount)
//     if err != nil {
//         return err
//     }

//     fmt.Printf("✅ Refund successful: %s\n", response.TransactionID)
//     return nil
// }

// func (u *UPIProcessor) Validate(payment *Payment) bool {
//     upiID := payment.Metadata["upi_id"]
//     return strings.Contains(upiID, "@")
// }

// type CashProcessor struct {
//     BaseProcessor // ✅ Embedding for code reuse
// }

// func (c *CashProcessor) Process(payment *Payment, gateway PaymentGateway) error {
//     c.Log("Processing cash payment")

//     if !c.ValidateAmount(payment.Amount) {
//         return errors.New("invalid amount")
//     }

//     // Cash doesn't need gateway
//     fmt.Printf("💵 Cash payment received: $%.2f\n", payment.Amount)

//     payment.TransactionID = "cash_" + payment.ID
//     payment.Status = Success

//     return nil
// }

// func (c *CashProcessor) Refund(payment *Payment, gateway PaymentGateway) error {
//     c.Log("Refunding cash payment")

//     // Cash refund doesn't need gateway
//     fmt.Printf("💵 Cash refund: $%.2f\n", payment.Amount)

//     return nil
// }

// func (c *CashProcessor) Validate(payment *Payment) bool {
//     // Cash always valid
//     return true
// }

// // ============================================
// // PROCESSOR FACTORY
// // ============================================

// type PaymentProcessorFactory struct {
//     Logger *Logger
//     Config *Config
// }

// func NewPaymentProcessorFactory(logger *Logger, config *Config) *PaymentProcessorFactory {
//     return &PaymentProcessorFactory{
//         Logger: logger,
//         Config: config,
//     }
// }

// func (f *PaymentProcessorFactory) GetProcessor(paymentType PaymentType) (PaymentProcessor, error) {
//     baseProcessor := BaseProcessor{
//         Logger: f.Logger,
//         Config: f.Config,
//     }

//     switch paymentType {
//     case CreditCard:
//         return &CreditCardProcessor{BaseProcessor: baseProcessor}, nil
//     case UPI:
//         return &UPIProcessor{BaseProcessor: baseProcessor}, nil
//     case Cash:
//         return &CashProcessor{BaseProcessor: baseProcessor}, nil
//     default:
//         return nil, fmt.Errorf("payment type %s not supported", paymentType)
//     }
// }

// // ============================================
// // PAYMENT SERVICE
// // ============================================

// type PaymentService struct {
//     Payments         map[string]*Payment
//     ProcessorFactory *PaymentProcessorFactory
//     GatewayFactory   *PaymentGatewayFactory
//     mu               sync.RWMutex
// }

// func NewPaymentService(processorFactory *PaymentProcessorFactory, gatewayFactory *PaymentGatewayFactory) *PaymentService {
//     return &PaymentService{
//         Payments:         make(map[string]*Payment),
//         ProcessorFactory: processorFactory,
//         GatewayFactory:   gatewayFactory,
//     }
// }

// func (s *PaymentService) ProcessPayment(payment *Payment, gatewayType GatewayType) error {
//     // ✅ Get processor based on payment type
//     processor, err := s.ProcessorFactory.GetProcessor(payment.Type)
//     if err != nil {
//         return err
//     }

//     // ✅ Get gateway based on gateway type
//     gateway, err := s.GatewayFactory.GetGateway(gatewayType)
//     if err != nil {
//         return err
//     }

//     // Validate payment
//     if !processor.Validate(payment) {
//         return errors.New("invalid payment details")
//     }

//     // Process payment
//     payment.Status = Pending
//     payment.GatewayType = gatewayType

//     err = processor.Process(payment, gateway)
//     if err != nil {
//         payment.Status = Failed
//         return err
//     }

//     // Store payment
//     s.mu.Lock()
//     s.Payments[payment.ID] = payment
//     s.mu.Unlock()

//     return nil
// }

// func (s *PaymentService) RefundPayment(paymentID string) error {
//     s.mu.RLock()
//     payment, exists := s.Payments[paymentID]
//     s.mu.RUnlock()

//     if !exists {
//         return errors.New("payment not found")
//     }

//     if payment.Status != Success {
//         return errors.New("can only refund successful payments")
//     }

//     // Get processor
//     processor, err := s.ProcessorFactory.GetProcessor(payment.Type)
//     if err != nil {
//         return err
//     }

//     // Get gateway (use the same gateway that was used for payment)
//     gateway, err := s.GatewayFactory.GetGateway(payment.GatewayType)
//     if err != nil {
//         return err
//     }

//     // Refund
//     return processor.Refund(payment, gateway)
// }

// func (s *PaymentService) GetPayment(paymentID string) (*Payment, error) {
//     s.mu.RLock()
//     defer s.mu.RUnlock()

//     payment, exists := s.Payments[paymentID]
//     if !exists {
//         return nil, errors.New("payment not found")
//     }

//     return payment, nil
// }

// // ============================================
// // HELPER FUNCTIONS
// // ============================================

// func maskCardNumber(cardNumber string) string {
//     if len(cardNumber) < 4 {
//         return cardNumber
//     }
//     return "****-****-****-" + cardNumber[len(cardNumber)-4:]
// }

// // ============================================
// // MAIN (Demo)
// // ============================================

// func main() {
//     // Initialize components
//     logger := &Logger{}
//     config := &Config{MaxAmount: 100000.0}

//     processorFactory := NewPaymentProcessorFactory(logger, config)
//     gatewayFactory := NewPaymentGatewayFactory()

//     paymentService := NewPaymentService(processorFactory, gatewayFactory)

//     fmt.Println("🎮 Payment System Demo\n")

//     // ============================================
//     // Example 1: Credit Card payment via Stripe
//     // ============================================
//     fmt.Println("=" * 60)
//     fmt.Println("📌 Example 1: Credit Card via Stripe")
//     fmt.Println("=" * 60)

//     payment1 := &Payment{
//         ID:     "pay_001",
//         Type:   CreditCard,
//         Amount: 1500.00,
//         Status: Pending,
//         Metadata: map[string]string{
//             "card_number": "4532123456789012",
//             "cvv":         "123",
//         },
//         CreatedAt: time.Now(),
//     }

//     err := paymentService.ProcessPayment(payment1, Stripe)
//     if err != nil {
//         fmt.Printf("❌ Payment failed: %v\n", err)
//     } else {
//         fmt.Printf("✅ Payment successful! Transaction ID: %s\n", payment1.TransactionID)
//     }

//     fmt.Println()

//     // ============================================
//     // Example 2: UPI payment via Razorpay
//     // ============================================
//     fmt.Println("=" * 60)
//     fmt.Println("📌 Example 2: UPI via Razorpay")
//     fmt.Println("=" * 60)

//     payment2 := &Payment{
//         ID:     "pay_002",
//         Type:   UPI,
//         Amount: 2500.00,
//         Status: Pending,
//         Metadata: map[string]string{
//             "upi_id": "user@paytm",
//         },
//         CreatedAt: time.Now(),
//     }

//     err = paymentService.ProcessPayment(payment2, Razorpay)
//     if err != nil {
//         fmt.Printf("❌ Payment failed: %v\n", err)
//     } else {
//         fmt.Printf("✅ Payment successful! Transaction ID: %s\n", payment2.TransactionID)
//     }

//     fmt.Println()

//     // ============================================
//     // Example 3: Cash payment (no gateway)
//     // ============================================
//     fmt.Println("=" * 60)
//     fmt.Println("📌 Example 3: Cash Payment")
//     fmt.Println("=" * 60)

//     payment3 := &Payment{
//         ID:        "pay_003",
//         Type:      Cash,
//         Amount:    500.00,
//         Status:    Pending,
//         CreatedAt: time.Now(),
//     }

//     err = paymentService.ProcessPayment(payment3, Stripe) // Gateway ignored for cash
//     if err != nil {
//         fmt.Printf("❌ Payment failed: %v\n", err)
//     } else {
//         fmt.Printf("✅ Payment successful! Transaction ID: %s\n", payment3.TransactionID)
//     }

//     fmt.Println()

//     // ============================================
//     // Example 4: Refund
//     // ============================================
//     fmt.Println("=" * 60)
//     fmt.Println("📌 Example 4: Refund Payment")
//     fmt.Println("=" * 60)

//     err = paymentService.RefundPayment("pay_001")
//     if err != nil {
//         fmt.Printf("❌ Refund failed: %v\n", err)
//     } else {
//         fmt.Println("✅ Refund successful!")
//     }

//     fmt.Println()

//     // ============================================
//     // Example 5: Get payment details
//     // ============================================
//     fmt.Println("=" * 60)
//     fmt.Println("📌 Example 5: Get Payment Details")
//     fmt.Println("=" * 60)

//     payment, err := paymentService.GetPayment("pay_002")
//     if err != nil {
//         fmt.Printf("❌ Error: %v\n", err)
//     } else {
//         fmt.Printf("Payment ID: %s\n", payment.ID)
//         fmt.Printf("Type: %s\n", payment.Type)
//         fmt.Printf("Amount: $%.2f\n", payment.Amount)
//         fmt.Printf("Status: %s\n", payment.Status)
//         fmt.Printf("Gateway: %s\n", payment.GatewayType)
//         fmt.Printf("Transaction ID: %s\n", payment.TransactionID)
//     }
// }
// ```

// ---

// ## 🎯 Key Architecture Points

// ### **1. Two Levels of Strategy Pattern:**

// ```
// ┌─────────────────────────────────────────────────────────┐
// │                    PaymentService                        │
// ├─────────────────────────────────────────────────────────┤
// │                                                          │
// │  ProcessorFactory ───┬──> CreditCardProcessor          │
// │  (Payment Type)      ├──> UPIProcessor                  │
// │                      └──> CashProcessor                  │
// │                                                          │
// │  GatewayFactory ─────┬──> StripeGateway                │
// │  (Gateway Type)      └──> RazorpayGateway               │
// │                                                          │
// └─────────────────────────────────────────────────────────┘
// ```

// ---

// ### **2. Processor Uses Gateway:**

// ```go
// // Processor interface takes gateway as parameter
// type PaymentProcessor interface {
//     Process(payment *Payment, gateway PaymentGateway) error
//     Refund(payment *Payment, gateway PaymentGateway) error
//     Validate(payment *Payment) bool
// }

// // Implementation
// func (c *CreditCardProcessor) Process(payment *Payment, gateway PaymentGateway) error {
//     // Do card-specific validation
//     // ...

//     // ✅ Then use gateway to actually charge
//     response, err := gateway.Charge(payment)
//     // ...
// }
// ```

// ---

// ### **3. Separation of Concerns:**

// | Component | Responsibility |
// |-----------|---------------|
// | **PaymentProcessor** | Payment method logic (CreditCard/UPI/Cash validation) |
// | **PaymentGateway** | External API integration (Stripe/Razorpay calls) |
// | **PaymentService** | Orchestration (combines processor + gateway) |

// ---

// ## 📊 Example Output

// ```
// 🎮 Payment System Demo

// ============================================================
// 📌 Example 1: Credit Card via Stripe
// ============================================================
// [INFO] Processing credit card payment
// 🔒 Validating card: ****-****-****-9012
// 💳 [STRIPE] Charging $1500.00 for payment pay_001
// ✅ Payment successful! Transaction ID: ch_stripe_pay_001

// ============================================================
// 📌 Example 2: UPI via Razorpay
// ============================================================
// [INFO] Processing UPI payment
// 📱 Processing UPI ID: user@paytm
// 💳 [RAZORPAY] Charging ₹2500.00 for payment pay_002
// ✅ Payment successful! Transaction ID: pay_razorpay_pay_002

// ============================================================
// 📌 Example 3: Cash Payment
// ============================================================
// [INFO] Processing cash payment
// 💵 Cash payment received: $500.00
// ✅ Payment successful! Transaction ID: cash_pay_003

// ============================================================
// 📌 Example 4: Refund Payment
// ============================================================
// [INFO] Refunding credit card payment
// 💳 [STRIPE] Refunding $1500.00 for transaction ch_stripe_pay_001
// ✅ Refund successful: re_stripe_ch_stripe_pay_001
// ✅ Refund successful!

// ============================================================
// 📌 Example 5: Get Payment Details
// ============================================================
// Payment ID: pay_002
// Type: UPI
// Amount: $2500.00
// Status: SUCCESS
// Gateway: RAZORPAY
// Transaction ID: pay_razorpay_pay_002
// ```

// ---

// ## 🎯 Summary

// **This example shows:**
// - ✅ **Two factories** (ProcessorFactory + GatewayFactory)
// - ✅ **Embedding** for code reuse (BaseProcessor)
// - ✅ **Interfaces** for polymorphism (PaymentProcessor, PaymentGateway)
// - ✅ **Strategy Pattern** at two levels (payment type + gateway type)
// - ✅ **Clean separation** (models, processors, gateways, service)

// **Perfect example of combining all the patterns we've discussed!** 🚀
