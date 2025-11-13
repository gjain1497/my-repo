# 🎯 Enum vs Interface - Quick Reference Guide

---

## Decision Framework

### Ask Two Questions:

**Q1: Do we need to STORE which variant was used?**
- Store in database/struct?
- Display to user?
- Query/filter by this value?

**Q2: Do we have MULTIPLE implementations with different behavior?**
- Different external APIs?
- Different algorithms?
- Strategy Pattern needed?

---

## Decision Matrix

| Q1: Store? | Q2: Multiple Impls? | Use |
|------------|---------------------|-----|
| ❌ No | ❌ No | Neither |
| ✅ Yes | ❌ No | **Enum Only** |
| ❌ No | ✅ Yes | **Interface Only** |
| ✅ Yes | ✅ Yes | **Enum + Interface** |

---

## Pattern 1: Interface ONLY

**When:** Multiple implementations, but don't store which one

**Examples:**
- Payment Gateway (don't store which gateway)
- Cache Eviction Policy (LRU, LFU, FIFO)
- Sorting Algorithms
- Compression Methods

```go
// ✅ Interface (for behavior)
type PaymentGateway interface {
    ProcessPayment(...) (*Payment, error)
}

type StripeGateway struct {}
type PayPalGateway struct {}

// ❌ No enum
type Payment struct {
    Id string
    // No "Gateway" field
}
```

**Why:** Gateway is implementation detail, not business data

---

## Pattern 2: Enum + Interface

**When:** Multiple implementations AND need to store which one

**Examples:**
- Shipping Provider (show customer which provider)
- Notification Channel (Email, SMS, Push)
- Database Type (MySQL, Postgres, MongoDB)
- Cloud Storage (S3, GCS, Azure)

```go
// ✅ Enum (for storing/display)
type ShippingProvider string

const (
    FedEx ShippingProvider = "FEDEX"
    UPS   ShippingProvider = "UPS"
)

// ✅ Interface (for behavior)
type ShippingProviderInterface interface {
    CreateShipment(...) (*Shipment, error)
}

type FedExProvider struct {}
func (f *FedExProvider) CreateShipment(...) (*Shipment, error) {
    return &Shipment{Provider: FedEx}, nil  // Set enum
}

// ✅ Store in struct
type Shipment struct {
    Provider ShippingProvider  // Business data
}

type ShippingService struct {
    provider ShippingProviderInterface  // Technical
}
```

**Why:** Customer wants to see "FedEx is delivering" (enum) + Need Strategy Pattern (interface)

---

## Pattern 3: Enum ONLY

**When:** Store which variant, but no complex implementations

**Examples:**
- Order Status (Pending, Shipped, Delivered)
- User Role (Admin, User, Guest)
- Auth Method (simple switch case)
- Payment Type (CreditCard, UPI) ← Just label, not gateway

```go
// ✅ Enum
type OrderStatus string

const (
    Pending   OrderStatus = "PENDING"
    Confirmed OrderStatus = "CONFIRMED"
)

// ❌ No interface needed
type Order struct {
    Status OrderStatus  // Just store value
}

// Simple switch
func UpdateStatus(status OrderStatus) {
    switch status {
    case Pending: // ...
    case Confirmed: // ...
    }
}
```

**Why:** Status is simple label, no complex behavior differences

---

## The Mental Model

### Enum = "WHAT" (Identity/Label)
```
Business asks: "What provider is this?"
Answer: FedEx, UPS, DHL
→ Use Enum
```

### Interface = "HOW" (Behavior/Strategy)
```
Technical asks: "How to create shipment?"
Answer: Call FedEx API, UPS API, DHL API (different)
→ Use Interface
```

### Enum + Interface = "WHAT + HOW"
```
Business: "What provider?" → Enum (store & display)
Technical: "How to call API?" → Interface (Strategy)
→ Use Both
```

---

## Common Use Cases

### ✅ Interface ONLY:
- Payment Gateway (Stripe, PayPal)
- Cache Strategy (LRU, LFU)
- Sorting Algorithm
- Encryption Method

### ✅ Enum + Interface:
- Shipping Provider (FedEx, UPS, DHL)
- Notification Service (Email, SMS, Push)
- Database Driver (MySQL, Postgres)
- Cloud Storage (AWS, GCP, Azure)

### ✅ Enum ONLY:
- Order Status
- User Role
- Payment Type (just label)
- File Format (if no processing)

---

## Quick Rules

```markdown
Customer/Business sees it → Enum
Multiple external APIs → Enum + Interface
Just algorithm variation → Interface
Just status/label → Enum
```

---

## Example Comparison

### Payment (Interface ONLY):

```go
// Why no enum?
// - Customer doesn't see "Paid via Stripe"
// - No need to query "all Stripe payments"
// - TransactionId is enough

type PaymentGateway interface { ... }  // ✅
type StripeGateway struct {}           // ✅

type Payment struct {
    TransactionId string  // ✅
    // No Gateway field    ❌
}
```

### Shipping (Enum + Interface):

```go
// Why enum?
// - Customer sees "FedEx is delivering"
// - Need to query "all FedEx shipments"
// - Analytics by provider

type ShippingProvider string           // ✅ Enum
type ShippingProviderInterface interface { ... }  // ✅ Interface

type Shipment struct {
    Provider ShippingProvider  // ✅ Store
}
```

---

## Interview Tips

**Show your thinking:**

```
"For payment gateway, I'll use interface only because 
the customer doesn't need to know if we used Stripe 
or PayPal - that's an implementation detail."

"For shipping, I'll use enum + interface because 
customers want to see which provider is delivering 
(enum for display), but we also need Strategy Pattern 
for different APIs (interface)."
```

---

## Cheat Sheet

| Pattern | When | Example |
|---------|------|---------|
| **Neither** | Simple value | Price (just float64) |
| **Enum** | Label/status | OrderStatus, UserRole |
| **Interface** | Hidden strategy | CachePolicy, SortAlgorithm |
| **Both** | Visible strategy | ShippingProvider, NotificationChannel |

---

**Print this and keep handy!** 📋✅