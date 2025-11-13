# 🎯 Interfaces & Design Patterns - Quick Reference Guide

---

## 1. When to Use Interfaces

### ✅ Use Interface When:

| Scenario                          | Example |
|----------|---------|
| **Multiple real implementations** | Payment: Stripe, PayPal, Razorpay |
| **One real + mock for testing** | Arkose: Real service + Mock |
| **Future extensibility needed** | "Support multiple shipping providers" |
| **Requirement mentions "different types"** | "Multiple payment methods" |

### ❌ Don't Use Interface When:

| Scenario |                                      Example |
|----------|---------|
| **Single implementation, no mocking needed** | Simple CRUD operations |
| **No variation exists** | Cart operations, User registration |
| **Internal helper functions** | Validation utilities |

---

## 2. Two Interface Patterns

### Pattern A: Multiple Real Implementations (Strategy Pattern)

**Use When:** Multiple different providers/algorithms exist

```go
// Interface
type PaymentGateway interface {
    ProcessPayment(orderId string, amount float64) (*Payment, error)
}

// Multiple REAL implementations
type StripeGateway struct { apiKey string }
type PayPalGateway struct { clientId string }
type RazorpayGateway struct { keyId string }

// Service USES the interface
type PaymentService struct {
    gateway PaymentGateway  // Can be any gateway
}
```

**Example:** E-commerce with Stripe, PayPal, Razorpay

---

### Pattern B: One Real + Mock (Testability Pattern)

**Use When:** One real implementation but need mocking for tests

```go
// Interface
type ArkoseAPI interface {
    IsValidationEnabled() bool
    ValidateArkoseHeader(r *http.Request) error
}

// Real implementation
type ArkoseService struct { cfg Config }
func (a *ArkoseService) IsValidationEnabled() bool { ... }

// Mock implementation (for testing)
type MockArkoseInterface struct { mock.Mock }
func (m *MockArkoseInterface) IsValidationEnabled() bool { ... }

// Handler uses interface (can swap real/mock)
type Handler struct {
    arkose ArkoseService  // Real in prod, mock in tests
}
```

**Example:** Your company's Arkose service

---

## 3. IS-A vs HAS-A Decision

### IS-A (Implementation/Inheritance)

```go
// StripeGateway IS-A PaymentGateway
type StripeGateway struct {}
func (s *StripeGateway) ProcessPayment(...) { ... }
// StripeGateway implements PaymentGateway ✅
```

**Use When:**
- ✅ Class has ONE primary responsibility
- ✅ Class IS the implementation itself
- ✅ No extra business logic
- ✅ "X is a Y" sounds natural

**Examples:**
- `StripeGateway` IS-A `PaymentGateway` ✅
- `FedExProvider` IS-A `ShippingProvider` ✅
- `ArkoseService` IS-A `ArkoseAPI` ✅

---

### HAS-A (Composition)

```go
// PaymentService HAS-A PaymentGateway
type PaymentService struct {
    gateway PaymentGateway  // HAS-A relationship
    cache   Cache
    logger  Logger
}
```

**Use When:**
- ✅ Class has MULTIPLE responsibilities
- ✅ Class orchestrates/coordinates
- ✅ Business logic exists
- ✅ Different layers (business vs integration)
- ✅ "X has a Y" sounds natural

**Examples:**
- `PaymentService` HAS-A `PaymentGateway` ✅
- `OrderService` HAS-A `PaymentService` ✅

---

## 4. When to Use Separate Services vs Service + Gateways

### Option 1: Separate Services (No Shared Logic)

```go
type StripeService struct {}
type PayPalService struct {}
type RazorpayService struct {}

// Each implements PaymentGateway independently
```

**Use When:**
- ❌ NO shared business logic
- ✅ Each provider has unique methods
- ✅ Simple, thin API wrappers

---

### Option 2: Service + Gateways (Shared Logic Exists)

```go
// Thin gateways (only API calls)
type StripeGateway struct {}
type PayPalGateway struct {}

// Service with shared business logic
type PaymentService struct {
    gateway PaymentGateway  // Uses any gateway
    cache   Cache           // Shared logic
}

func (s *PaymentService) ProcessPayment(...) {
    // Validation (shared) ✅
    // Fraud check (shared) ✅
    // Gateway call (delegated) ✅
    // Caching (shared) ✅
}
```

**Use When:**
- ✅ Common validation across providers
- ✅ Shared caching/retry/logging
- ✅ Business rules independent of provider
- ✅ Avoid code duplication (DRY)

---

## 5. Quick Decision Tree

```
Working on a service →

Q1: Does this have multiple ways to do it?
    YES → Need interface
    NO  → No interface

Q2: How many real implementations?
    ONE  → Service implements interface (like Arkose)
    MANY → Gateways implement interface, Service uses it (like Payment)

Q3: Is there shared business logic?
    NO  → Separate services (StripeService, PayPalService)
    YES → Service + Gateways (PaymentService + gateways)

Q4: What's the relationship?
    One responsibility → IS-A
    Multiple responsibilities → HAS-A
```

---

## 6. Complete Examples

### Example 1: Payment System (Multiple Implementations + Shared Logic)

```go
// Interface
type PaymentGateway interface {
    ProcessPayment(orderId string, amount float64) (*Payment, error)
}

// Gateways (thin, just API calls)
type StripeGateway struct { apiKey string }
func (s *StripeGateway) ProcessPayment(...) (*Payment, error) {
    // Call Stripe API only
}

type PayPalGateway struct { clientId string }
func (p *PayPalGateway) ProcessPayment(...) (*Payment, error) {
    // Call PayPal API only
}

// Service (business logic)
type PaymentService struct {
    gateway PaymentGateway  // HAS-A
    cache   Cache
}

func (s *PaymentService) ProcessPayment(...) (*Payment, error) {
    // Validation ✅
    // Fraud check ✅
    payment, err := s.gateway.ProcessPayment(...)  // Delegate
    // Caching ✅
    return payment, err
}

// Usage
stripe := &StripeGateway{apiKey: "sk_..."}
service := &PaymentService{gateway: stripe}
```

**Why this design:**
- Multiple gateways → Need interface ✅
- Shared validation/caching → PaymentService + Gateways ✅
- PaymentService orchestrates → HAS-A relationship ✅

---

### Example 2: Arkose System (One Implementation + Mock)

```go
// Interface
type ArkoseAPI interface {
    IsValidationEnabled() bool
    ValidateArkoseHeader(r *http.Request) error
}

// Real implementation
type ArkoseService struct { cfg Config }
func (a *ArkoseService) IsValidationEnabled() bool {
    return a.cfg.IsValidationEnabled
}

// Mock implementation
type MockArkoseInterface struct { mock.Mock }
func (m *MockArkoseInterface) IsValidationEnabled() bool {
    return m.Called().Get(0).(bool)
}

// Handler
type Handler struct {
    arkose ArkoseAPI  // Interface allows real or mock
}

// Production
real := &ArkoseService{cfg: prodConfig}
handler := &Handler{arkose: real}

// Testing
mock := &MockArkoseInterface{}
mock.On("IsValidationEnabled").Return(true)
handler := &Handler{arkose: mock}
```

**Why this design:**
- One real implementation → Service implements interface ✅
- Need testing → Mock also implements interface ✅
- Service has state (cfg) → IS-A relationship ✅

---

## 7. Key Takeaways

| Concept | Rule |
|---------|------|
| **Multiple implementations** | Use interface + gateways |
| **One implementation** | Service implements interface directly (if mocking needed) |
| **Shared logic** | Service layer + gateway layer |
| **No shared logic** | Separate services |
| **One responsibility** | IS-A relationship |
| **Multiple responsibilities** | HAS-A relationship |
| **Business logic** | Goes in Service, not Gateway |
| **API calls** | Goes in Gateway, not Service |

---

## 8. Common Mistakes to Avoid

❌ **Mistake 1:** Putting business logic in gateways
```go
// ❌ Wrong
type StripeGateway struct {}
func (s *StripeGateway) ProcessPayment(...) {
    // Validation here ❌
    // Caching here ❌
    // API call
}
```

✅ **Correct:** Business logic in service, API calls in gateway
```go
// ✅ Correct
type PaymentService struct { gateway PaymentGateway }
func (s *PaymentService) ProcessPayment(...) {
    // Validation ✅
    // Caching ✅
    s.gateway.ProcessPayment(...)  // Delegate API call
}
```

---

❌ **Mistake 2:** Using HAS-A when IS-A is better
```go
// ❌ Wrong (if StripeService has no extra logic)
type PaymentService struct {
    stripe StripeGateway  // HAS-A
}
```

✅ **Correct:** Use IS-A if no business logic
```go
// ✅ Correct
type StripeService struct {}
func (s *StripeService) ProcessPayment(...) {
    // Just API call, no business logic
}
```

---

## 9. Cheat Sheet

**When designing a new service:**

```
Step 1: Does it need interface?
  → Multiple implementations OR mocking needed? YES → Interface

Step 2: How many real implementations?
  → One: Service implements interface
  → Many: Gateways implement interface, Service uses them

Step 3: Shared business logic?
  → Yes: Service + Gateways pattern
  → No: Separate services pattern

Step 4: What's the relationship?
  → Service has extra logic: HAS-A
  → Service is just implementation: IS-A
```

---

**Print this and keep handy!** 📋✅