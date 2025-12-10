# Interface Extension Patterns - Extending Without Breaking

## 🤔 The Problem

When you need to add new methods to an interface, existing implementations break:

```go
type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    ListVehicles() ([]*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

// V1 implements all 3 ✅
type VehicleServiceV1 struct{}

// V2 implements all 3 ✅
type VehicleServiceV2 struct{}

// V3 needs a NEW method: GetVehicleHistory()
// If we add to interface → V1, V2 break! ❌
```

---

## 🎯 Solutions

### Option 1: Interface Segregation (Split)

Split into smaller, focused interfaces:

```go
// Base interface (all versions implement)
type VehicleReader interface {
    GetVehicle(id string) (*Vehicle, error)
    ListVehicles() ([]*Vehicle, error)
}

type VehicleWriter interface {
    AddVehicle(req AddVehicleRequest) error
}

// Extended interface (only V3 implements)
type VehicleHistoryProvider interface {
    GetVehicleHistory(id string) ([]*VehicleEvent, error)
}

// V1, V2 implement base only
type VehicleServiceV1 struct{}
func (s *VehicleServiceV1) GetVehicle(...) {}
func (s *VehicleServiceV1) ListVehicles(...) {}
func (s *VehicleServiceV1) AddVehicle(...) {}

// V3 implements base + extended
type VehicleServiceV3 struct{}
func (s *VehicleServiceV3) GetVehicle(...) {}
func (s *VehicleServiceV3) ListVehicles(...) {}
func (s *VehicleServiceV3) AddVehicle(...) {}
func (s *VehicleServiceV3) GetVehicleHistory(...) {}  // Extra!
```

---

### Option 2: Interface Composition (Embed)

Compose interfaces by embedding:

```go
// Base
type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    ListVehicles() ([]*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

// Extended (embeds base)
type VehicleServiceV3Interface interface {
    VehicleServiceInterface  // All base methods
    GetVehicleHistory(id string) ([]*VehicleEvent, error)  // Extra
}
```

**Usage:**
```go
// Function that needs only base functionality
func ProcessVehicles(svc VehicleServiceInterface) {
    // Works with V1, V2, V3
}

// Function that needs history
func AuditVehicles(svc VehicleServiceV3Interface) {
    // Works only with V3
}
```

---

### Option 3: Type Assertion (Runtime Check)

Check capability at runtime:

```go
type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    ListVehicles() ([]*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

type VehicleHistoryProvider interface {
    GetVehicleHistory(id string) ([]*VehicleEvent, error)
}

// Usage
func DoSomething(svc VehicleServiceInterface) {
    // Base functionality works
    svc.GetVehicle("123")
    
    // Check if this version supports history
    if historyProvider, ok := svc.(VehicleHistoryProvider); ok {
        historyProvider.GetVehicleHistory("123")  // Only if supported
    }
}
```

---

## 🔑 The REAL Difference: Consumer Side

The difference is NOT in implementation, but in **USAGE/CONSUMER side**!

### Option 1: Split (No Composition)

```go
type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

type VehicleHistoryProvider interface {
    GetVehicleHistory(id string) ([]*VehicleEvent, error)
}

// V3 implements BOTH
type VehicleServiceV3 struct{}
func (s *VehicleServiceV3) GetVehicle(...) {}
func (s *VehicleServiceV3) AddVehicle(...) {}
func (s *VehicleServiceV3) GetVehicleHistory(...) {}
```

**Consumer side:**
```go
// Need to accept TWO separate interfaces
func DoSomething(svc VehicleServiceInterface, history VehicleHistoryProvider) {
    svc.GetVehicle("123")
    history.GetVehicleHistory("123")
}

// OR pass same object twice - shit way
DoSomething(v3, v3)
```

### Option 2: Compose

```go
type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

type VehicleServiceV3Interface interface {
    VehicleServiceInterface  // Embedded
    GetVehicleHistory(id string) ([]*VehicleEvent, error)
}

// V3 implements composed interface
type VehicleServiceV3 struct{}
```

**Consumer side:**
```go
// Single interface parameter - CLEAN!
func DoSomething(svc VehicleServiceV3Interface) {
    svc.GetVehicle("123")           // Base method ✅
    svc.GetVehicleHistory("123")    // Extended method ✅
}

DoSomething(v3)  // Pass once ✅
```

### The Real Difference:

| Aspect | Option 1 (Split) | Option 2 (Compose) |
|--------|------------------|-------------------|
| Implementation | Same | Same |
| Consumer needs base only | `func(svc VehicleServiceInterface)` ✅ | `func(svc VehicleServiceInterface)` ✅ |
| Consumer needs both | `func(svc Interface1, h Interface2)` 🤮 | `func(svc V3Interface)` ✅ |

### Bottom Line:

| Scenario | Use |
|----------|-----|
| Features are **independent**, consumers use them **separately** | **Split** (Option 1) |
| Features are **related**, consumers use them **together** | **Compose** (Option 2) |

---

## 🎯 When to Use Which?

### Decision Rule:

| Question | If YES → |
|----------|----------|
| Can the new method exist WITHOUT base methods? | **Option 1** (Split) |
| Is the new version a "superset" of base? | **Option 2** (Compose) |
| Will some consumers need ONLY the new method? | **Option 1** (Split) |
| Will consumers always need base + new together? | **Option 2** (Compose) |
| Need runtime flexibility? | **Option 3** (Type Assertion) |

### Simple Rule:

| Scenario | Use |
|----------|-----|
| New method is **independent/optional** feature | **Option 1** (Split) |
| New method is **natural extension** of existing | **Option 2** (Compose) |
| Need to check capability at **runtime** | **Option 3** (Type Assertion) |

---

## 💡 Examples

### Option 1 Example: Split (Independent Feature)

```go
// GetVehicleHistory is OPTIONAL feature
// Not every vehicle service needs history tracking

type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

// Separate interface - independent feature
type VehicleHistoryProvider interface {
    GetVehicleHistory(id string) ([]*VehicleEvent, error)
}

// V1 - basic (no history)
type VehicleServiceV1 struct{}

// V3 - implements BOTH (has history)
type VehicleServiceV3 struct{}
```

**Why Split?** History is optional - not every consumer needs it.

---

### Option 2 Example: Compose (Natural Extension)

```go
// V3 is a "premium/advanced" version
// It does EVERYTHING V1 does + more

type VehicleServiceInterface interface {
    GetVehicle(id string) (*Vehicle, error)
    AddVehicle(req AddVehicleRequest) error
}

// V3 is superset - naturally extends base
type VehicleServiceV3Interface interface {
    VehicleServiceInterface  // All base methods
    GetVehicleHistory(id string) ([]*VehicleEvent, error)
    GetVehicleAnalytics(id string) (*Analytics, error)
}
```

**Why Compose?** V3 is clearly "base + extras", not a separate feature.

---

## ✅ Real-World Example: Go Standard Library

```go
// io package in Go stdlib does this!

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Composed
type ReadWriter interface {
    Reader
    Writer
}

// Extended
type ReadCloser interface {
    Reader
    Close() error
}
```

---

## 📝 Interview One-liners

> "To extend interfaces without breaking existing implementations, we use **Interface Segregation** - split into smaller focused interfaces, or **compose interfaces** by embedding. This follows Open/Closed principle - open for extension, closed for modification."

> **Split** = Independent features that can exist alone
> 
> **Compose** = Extended version that builds on base

---

## 🎯 Summary

| Approach | When | Example |
|----------|------|---------|
| **Split (ISP)** | Independent/optional features | History tracking, Analytics |
| **Compose (Embed)** | Natural superset/premium version | V3 = V1 + extras |
| **Type Assertion** | Runtime capability check | Plugin systems |

**Key Principle:** Open/Closed - Open for extension, Closed for modification!
