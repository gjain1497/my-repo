# LLD & SOLID Principles - Where They Apply

## 🎯 Key Insight: LLD = Service Layer Design

When interviewer asks **"Design a system"**, they're asking:

> **"Design the Service Layer"**

---

## 📊 Layer-wise Analysis

### SRP (Single Responsibility Principle) by Layer:

| Layer | SRP Strict? | Why |
|-------|-------------|-----|
| **Handler/Controller** | ⚠️ Relaxed | Can call multiple services (orchestration is its job) |
| **Service** | 🔒 **Strict** | Core business logic, one responsibility per service |
| **Repository** | ⚠️ Relaxed | Can JOIN multiple tables for performance |
| **Database** | N/A | Just storage |

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────┐
│           Handler/Controller                │
│  Focus: Request/Response, Validation        │
│  Principles: Clean Code, Error Handling     │
│  SRP: Relaxed (can call multiple services)  │
└─────────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────────┐
│              Service Layer                  │
│  Focus: Business Logic                      │
│  Principles: SOLID, Design Patterns, LLD   │
│  SRP: STRICT 🔒                             │
│  🔥 THIS IS WHERE LLD LIVES!                │
└─────────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────────┐
│              Repository Layer               │
│  Focus: Data Access, Performance            │
│  Principles: Query Optimization, Caching    │
│  SRP: Relaxed (can JOIN multiple tables)    │
└─────────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────────┐
│              Database                       │
│  Single source of truth                     │
└─────────────────────────────────────────────┘
```

---

## 🎯 Where SOLID Principles Apply

| Principle | Primary Layer | Example |
|-----------|---------------|---------|
| **S** - Single Responsibility | Service | VehicleService ≠ BookingService |
| **O** - Open/Closed | Service | Strategy pattern for pricing |
| **L** - Liskov Substitution | Service | StripeGateway / RazorpayGateway interchangeable |
| **I** - Interface Segregation | Service | Small focused interfaces |
| **D** - Dependency Inversion | Service | Depend on interfaces, not concrete types |

---

## 🤔 Why Repository Can Break SRP?

### Service with Map (❌ BAD):
```go
type VehicleService struct {
    Vehicles map[string]*Vehicle
    Bookings map[string][]*Booking  // ❌ STORING booking data
}
```
- VehicleService **OWNS/MANAGES** booking data
- Two sources of truth (duplication)
- Has to **maintain, update, sync** this data

### Repository with JOIN (✅ OK):
```go
func (r *VehicleRepository) GetAvailableVehicles(...) {
    query := `
        SELECT v.* FROM vehicles v
        WHERE v.id NOT IN (
            SELECT b.vehicle_id FROM bookings b WHERE ...
        )
    `
}
```
- Repository **READS** booking data (doesn't own it)
- Single source of truth (database)
- No maintenance, no sync issues

### Key Difference:

| Aspect | Service with Map | Repository with JOIN |
|--------|------------------|---------------------|
| Data ownership | ❌ Owns/stores data | ✅ Just reads |
| Source of truth | ❌ Multiple (duplication) | ✅ Single (database) |
| Sync needed | ❌ Yes | ✅ No |
| SRP violation | ❌ Yes | ⚠️ Accepted trade-off |

---

## 💡 The Abstraction

```go
// VehicleService doesn't know HOW availability is checked
// It just asks repository for "available vehicles"

func (s *VehicleService) ListAvailableVehicles(...) ([]*Vehicle, error) {
    return s.vehicleRepo.GetAvailableVehicles(locationId, vehicleType, startDate, endDate)
}

// Repository hides the JOIN complexity
func (r *VehicleRepository) GetAvailableVehicles(...) ([]*Vehicle, error) {
    // JOIN happens here, but service doesn't know!
    query := `SELECT v.* FROM vehicles v WHERE ... NOT IN (SELECT from bookings)`
    // Returns only Vehicle objects
}
```

### The Rule:

| Question | Answer |
|----------|--------|
| Does VehicleService access booking data? | ❌ No |
| Does VehicleService know about bookings table? | ❌ No |
| Does VehicleRepository read from bookings table? | ✅ Yes (for filtering) |
| Does VehicleRepository return booking data? | ❌ No (only vehicles) |

---

## 📝 Interview Mental Model

| Term | Means |
|------|-------|
| LLD | Service Layer Design |
| SOLID | Service Layer Principles |
| Design Patterns | Service Layer Patterns |
| "Design X System" | "Design X Service Layer" |

---

## ✅ One-liners for Interview

### On SRP across layers:
> "SRP is **strictly enforced at Service layer** because that's where business logic lives. Handler orchestrates multiple services, Repository optimizes data access - both have relaxed SRP for practical reasons."

### On Repository JOINs:
> "Repository JOINing multiple tables is an accepted trade-off for performance. Repository only **reads** from other tables, doesn't own/manage them. It returns **one entity type**. SRP is strictly maintained at **service layer**. Database remains **single source of truth**."

### On LLD:
> "LLD Interview = Design the Service Layer following SOLID principles"

---

## 🎯 Summary

```
Handler  → Orchestration      (relaxed SRP)
Service  → Business Logic     (STRICT SRP) 🔒  ← LLD LIVES HERE!
Repository → Data Access      (relaxed SRP)
```

**When we talk about LLD/SOLID, we're talking about the Service Layer!**
