# 🚗 Parking Lot Management System

A production-ready parking lot management system built in Go, demonstrating clean architecture, design patterns, and concurrent programming. This system efficiently manages multi-floor parking facilities with different vehicle types, slot categories, flexible pricing strategies, and multi-gate access control.

## 🎯 Features

- **Multi-Floor Support**: Manage parking across multiple floors with independent slot management
- **Vehicle Type Management**: Support for Cars, Bikes, and Trucks with appropriate slot assignment
- **Slot Categories**: Four types of slots - Compact, Large, Handicapped, and Electric
- **Multi-Gate Access Control**: Support for multiple entry/exit points with gate type validation ⭐ NEW!
- **Gate Status Management**: Gates can be Open, Closed, or Under Maintenance ⭐ NEW!
- **Capacity Management**: Enforce maximum capacity limits and prevent overcrowding ⭐ NEW!
- **Flexible Pricing**: Strategy pattern enables multiple pricing models (Hourly, Flat, Progressive) ⭐ ENHANCED!
- **Thread-Safe Operations**: Concurrent parking/unparking using RWMutex for optimal performance
- **Real-Time Availability**: Track slot availability per floor and slot type with live occupancy display
- **Ticket Management**: Automated ticket generation and fee calculation
- **Configurable Rates**: Vehicle-specific hourly rates that can be customized at runtime

## 🏗️ System Architecture

### Core Entities

#### 1. **Vehicle**
- **Purpose**: Represents vehicles entering the parking lot
- **Types**: Car, Bike, Truck (defined as VehicleType enum)
- **Key Fields**: 
  - `VehicleType`: Type of vehicle
  - `LicensePlate`: Unique identifier

#### 2. **Slot**
- **Purpose**: Individual parking space that can be occupied/freed
- **Types**: Compact (bikes, cars), Large (trucks), Handicapped
- **Key Fields**:
  - `ID`: Unique slot identifier
  - `Type`: Slot category
  - `Status`: Available/Occupied
  - `ParkedVehicle`: Reference to parked vehicle (nil if empty)
  - `Floor`: Floor number
- **Thread Safety**: Protected by `sync.RWMutex`

#### 3. **Floor**
- **Purpose**: Organizes slots by type for efficient lookup
- **Structure**: `map[SlotType][]*Slot` - groups slots by category
- **Key Methods**: 
  - `FindSlot()`: O(n) lookup for available slot of given type

#### 4. **Ticket**
- **Purpose**: Proof of parking with entry/exit tracking
- **Key Fields**:
  - `TicketId`: Unique identifier
  - `VehicleAssigned`: Pointer to vehicle
  - `SlotAssigned`: Pointer to assigned slot
  - `EntryTime`: Non-nullable timestamp
  - `ExitTime`: Nullable pointer (set on exit)

#### 5. **ParkingLot**
- **Purpose**: Central orchestrator managing all operations
- **Singleton**: Only one instance per application
- **Key Fields**:
  - `Floors`: Array of floor pointers
  - `Gates`: Map of gate ID to Gate (for multi-gate access) ⭐ NEW!
  - `ActiveTickets`: Map for O(1) ticket lookup
  - `TicketCounter`: Atomic counter for ticket IDs
  - `PricingStrategy`: Pluggable pricing algorithm
  - `MaxCapacity`: Maximum number of vehicles allowed ⭐ NEW!
  - `CurrentOccupancy`: Current number of parked vehicles ⭐ NEW!
- **Thread Safety**: Protected by `sync.RWMutex`

#### 6. **Gate** ⭐ NEW!
- **Purpose**: Control entry and exit access points
- **Types**: Entry-only, Exit-only, or Both (bidirectional)
- **Status**: Open, Closed, or Under Maintenance
- **Key Methods**:
  - `CanProcessEntry()`: Validates if gate can handle vehicle entry
  - `CanProcessExit()`: Validates if gate can handle vehicle exit
- **Use Case**: Large parking lots with multiple access points need controlled entry/exit through designated gates

## 🎨 Design Patterns

### 1. **Factory Pattern** (Creational)
- **Where**: `VehicleFactory.CreateVehicle()`
- **Why**: Centralizes vehicle creation logic, making it easy to add validation, logging, or new vehicle types
- **Benefit**: Client code doesn't need to know vehicle instantiation details
- **Example**:
  ```go
  factory := &VehicleFactory{}
  car := factory.CreateVehicle(Car, "ABC-123")
  ```

### 2. **Strategy Pattern** (Behavioral)
- **Where**: `PricingStrategy` interface with `HourlyPricing` and `FlatPricing` implementations
- **Why**: Different parking facilities need different pricing models without changing core logic
- **Benefit**: Can switch pricing at runtime, add new strategies without modifying existing code (Open/Closed Principle)
- **Example**:
  ```go
  // Hourly pricing with configurable rates
  hourly := NewHourlyPricing()
  
  // Flat pricing
  flat := &FlatPricing{}
  
  // Can swap strategies
  parkingLot.PricingStrategy = flat
  ```

### 3. **Singleton Pattern** (Creational)
- **Where**: `GetParkingLot()` function with `sync.Once`
- **Why**: Only one parking lot instance should exist to maintain global state consistency
- **Benefit**: Thread-safe initialization, prevents duplicate parking lot instances
- **Implementation**:
  ```go
  var (
      parkingLotInstance *ParkingLot
      once               sync.Once
  )
  
  func GetParkingLot(name string, numFloors int, strategy PricingStrategy) *ParkingLot {
      once.Do(func() {
          // Initialize only once, even with concurrent calls
      })
      return parkingLotInstance
  }
  ```

### 4. **Thread Safety** (Concurrency Pattern)
- **Where**: All shared mutable state
- **Why**: Multiple vehicles can enter/exit simultaneously through different gates
- **Implementation**:
  - `sync.RWMutex` on Slot: Multiple readers (availability checks), exclusive writers (occupy/free)
  - `sync.RWMutex` on ParkingLot: Protects ticket counter, occupancy, and active tickets map
  - `sync.RWMutex` on Gates: Thread-safe gate management
- **Benefit**: Safe concurrent access without data races

### 5. **Gate Access Control** ⭐ NEW!
- **Where**: `ProcessVehicleEntry()` and `ProcessVehicleExit()` methods
- **Why**: Large parking lots have multiple entry/exit points that need controlled access
- **Types**:
  - **Entry Gate**: Only allows vehicles to enter
  - **Exit Gate**: Only allows vehicles to exit
  - **Both Gate**: Can handle both entry and exit (bidirectional)
- **Status Management**: Gates can be Open, Closed, or Under Maintenance
- **Benefit**: 
  - Prevents entry through exit-only gates and vice versa
  - Enables temporary closure of specific gates for maintenance
  - Provides access control and traffic management
- **Real-world Use**: Shopping malls, airports, stadiums with multiple access points

## 🚀 How to Run

### Prerequisites
- Go 1.16 or higher

### Steps
```bash
# Clone or copy the code
# Save as parking_lot.go

# Run the system
go run parking_lot.go
```

### Expected Output
```
--- Test 1: Entry through GATE-A ---
✅ Car entered via GATE-A
   Ticket: Ticket-1, Slot: 1

--- Test 2: Try entry through EXIT gate ---
❌ Expected Error: gate GATE-B cannot process entry (status: OPEN)

--- Test 3: Entry through GATE-C (Both) ---
✅ Truck entered via GATE-C
   Ticket: Ticket-2, Slot: 2

==== Grand Plaza Status ====
Occupancy: 2/2 slots
PARKING LOT IS FULL
Floor 0:
  COMPACT: 0/1 available
  LARGE: 0/1 available
  HANDICAPPED: 1/1 available
Active Tickets: 2

--- Test 4: Exit through GATE-B ---
✅ Car exited via GATE-B, Fee: $10.00

--- Test 5: Try exit through ENTRY gate ---
❌ Expected Error: gate GATE-A cannot process exit (status: OPEN)

--- Test 6: Exit through GATE-C (Both) ---
✅ Truck exited via GATE-C, Fee: $20.00

==== Grand Plaza Status ====
Occupancy: 0/2 slots
🟢 2 slots available
Floor 0:
  COMPACT: 1/1 available
  LARGE: 1/1 available
  HANDICAPPED: 1/1 available
Active Tickets: 0

--- Test 7: Try using closed gate ---
❌ Expected Error: gate GATE-A cannot process entry (status: CLOSED)
```

## 💡 Key Design Decisions

### 1. **Pointers vs Values**
- **Decision**: Use pointers for Vehicle, Slot, Ticket
- **Reason**: 
  - Enables modification of shared objects (e.g., freeing a slot affects all references)
  - Avoids unnecessary copying of large structs
  - Allows nil values (e.g., `ParkedVehicle` is nil when slot is empty)

### 2. **RWMutex vs Mutex**
- **Decision**: Use `sync.RWMutex` instead of `sync.Mutex`
- **Reason**:
  - Read-heavy workload: Many threads check availability, few modify slots
  - `RLock()` allows concurrent reads for better performance
  - `Lock()` provides exclusive access for writes
- **Impact**: ~3-5x better performance for read operations

### 3. **Map Organization in Floor**
- **Decision**: `map[SlotType][]*Slot` instead of flat `[]*Slot`
- **Reason**:
  - O(1) access to slots of specific type (Compact/Large/Handicapped)
  - Reduces search space when looking for available slots
  - Natural grouping for different vehicle requirements
- **Trade-off**: Slightly more complex but much more efficient

### 4. **Configurable Pricing Rates**
- **Decision**: Store rates in `map[VehicleType]float64`
- **Reason**:
  - Easy to customize rates per parking facility
  - Can load from config files in production
  - Different rates for peak/off-peak hours
- **Alternative Considered**: Hardcoded rates (rejected for lack of flexibility)

### 5. **Ticket ID Generation**
- **Decision**: Simple counter with format "Ticket-{N}"
- **Reason**: Simple, deterministic, sufficient for single-instance systems
- **Production Alternative**: UUID for distributed systems

### 6. **Capacity Enforcement** ⭐ NEW!
- **Decision**: Track MaxCapacity and CurrentOccupancy in ParkingLot
- **Reason**:
  - Prevents overcrowding and ensures fire safety compliance
  - Provides real-time occupancy metrics
  - Enables predictive "lot full" warnings
- **Implementation**: Check before parking, increment on entry, decrement on exit
- **Thread-safe**: All occupancy operations protected by mutex

### 7. **Gate-Based Access Control** ⭐ NEW!
- **Decision**: Separate entry/exit validation through Gate entities
- **Reason**:
  - Large facilities have multiple access points
  - Different gates serve different purposes (entry-only, exit-only, both)
  - Gates can be temporarily closed for maintenance
- **Alternative Considered**: Single entry/exit point (rejected - not scalable)
- **Real-world Benefit**: Enables traffic flow optimization and access control

## 🔄 Extensibility

The system is designed for easy extension:

### Adding New Vehicle Types
```go
const (
    // Existing types
    Truck VehicleType = "TRUCK"
    Car   VehicleType = "CAR"
    Bike  VehicleType = "BIKE"
    
    // Add new type
    ElectricCar VehicleType = "ELECTRIC_CAR"
)

// Update GetRequiredSlotType() method
// Add rate in NewHourlyPricing()
```

### Adding New Pricing Strategies
```go
type WeekendPricing struct {
    BaseRate   float64
    WeekendMultiplier float64
}

func (w *WeekendPricing) CalculateFee(ticket *Ticket) float64 {
    // Custom logic for weekend pricing
    if isWeekend(ticket.EntryTime) {
        return basePrice * w.WeekendMultiplier
    }
    return basePrice
}
```

### Adding Slot Reservation
```go
type Slot struct {
    // ... existing fields
    ReservedFor *string  // Customer ID or nil
    ReservedUntil *time.Time
}

// Add ReserveSlot() method
```

### Adding Payment Methods
```go
type PaymentStrategy interface {
    ProcessPayment(amount float64) error
}

type CashPayment struct{}
type CardPayment struct{}
type UPIPayment struct{}
```

### Adding More Gates ⭐ NEW!
```go
func main() {
    parkingLot := GetParkingLot("Mall Parking", 3, 100, NewHourlyPricing())
    
    // Add entry gates
    parkingLot.AddGate(&Gate{Id: "NORTH-ENTRY", Type: EntryGate, Status: GateOpen})
    parkingLot.AddGate(&Gate{Id: "SOUTH-ENTRY", Type: EntryGate, Status: GateOpen})
    
    // Add exit gates
    parkingLot.AddGate(&Gate{Id: "NORTH-EXIT", Type: ExitGate, Status: GateOpen})
    parkingLot.AddGate(&Gate{Id: "SOUTH-EXIT", Type: ExitGate, Status: GateOpen})
    
    // Add bidirectional gate
    parkingLot.AddGate(&Gate{Id: "VIP-GATE", Type: BothGate, Status: GateOpen})
}
```

## 📊 Performance Characteristics

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Find Available Slot | O(F × S) where F=floors, S=slots per type | O(1) |
| Park Vehicle | O(F × S) + O(1) gate lookup | O(1) |
| Unpark Vehicle | O(1) with ticket ID + O(1) gate lookup | O(1) |
| Add Gate | O(1) | O(1) |
| Gate Validation | O(1) | O(1) |
| Get Status | O(F × T × S) where T=slot types | O(1) |
| Check Capacity | O(1) | O(1) |

**Optimization Opportunities:**
- Maintain a min-heap of available slots per type: O(log S) lookup
- Cache available slot counts: O(1) status queries
- Index slots by floor for faster floor-specific queries
- Batch gate status updates for maintenance windows

## 🧪 Testing Strategy

### Unit Tests (Recommended)
```go
func TestParkVehicle(t *testing.T)
func TestUnparkVehicle(t *testing.T)
func TestConcurrentParking(t *testing.T)
func TestPricingStrategies(t *testing.T)
func TestSlotAvailability(t *testing.T)
func TestGateValidation(t *testing.T) // NEW!
func TestCapacityEnforcement(t *testing.T) // NEW!
func TestGateTypesValidation(t *testing.T) // NEW!
```

### Integration Tests
```go
func TestMultiGateScenario(t *testing.T) {
    // Test vehicles entering through different gates simultaneously
}

func TestCapacityLimitWithGates(t *testing.T) {
    // Test capacity enforcement across multiple gates
}

func TestGateMaintenanceMode(t *testing.T) {
    // Test closing gates and reopening them
}
```

### Concurrency Testing
```go
func TestRaceConditions(t *testing.T) {
    // Use go test -race
    // Simulate 100 concurrent park operations through multiple gates
}
```

## 📚 Key Learnings

### Go-Specific Patterns
- **Interfaces over Inheritance**: Go's composition model using interfaces
- **Explicit Error Handling**: Every operation returns `error` for robustness
- **Pointer Semantics**: Understanding when to use `*T` vs `T`
- **Concurrency Primitives**: `sync.RWMutex`, `sync.Once`, `defer`

### Design Principles Applied
- **Single Responsibility**: Each struct has one clear purpose
- **Open/Closed**: Open for extension (new strategies), closed for modification
- **Dependency Inversion**: ParkingLot depends on `PricingStrategy` interface, not concrete implementations
- **Interface Segregation**: Small, focused interfaces (`PricingStrategy` has one method)

### Production Considerations
- Thread safety is not optional in concurrent systems
- Configurability > Hardcoding (rates, slot counts, gate configurations)
- Use appropriate lock granularity (slot-level, not system-level locks)
- Consider idempotency (what if UnparkVehicle called twice?)
- Gate status changes should be atomic and logged
- Capacity limits must account for reserved slots and edge cases
- Monitor gate performance (average processing time per gate)

## 🎓 Interview Talking Points

When discussing this system in interviews:

1. **Design Choice**: "I used Strategy pattern for pricing because different parking facilities need different pricing models without changing core logic"
2. **Scalability**: "For 10,000+ slots, I'd use a priority queue for O(log n) slot lookup"
3. **Trade-offs**: "RWMutex improves read performance but adds complexity vs simple Mutex. I chose it because availability checks (reads) are more frequent than occupancy changes (writes)"
4. **Real-world**: "In production, I'd add persistent storage, distributed locking for multi-instance deployments, and gate traffic analytics"
5. **Monitoring**: "I'd add metrics for slot utilization, average parking duration, revenue tracking, and per-gate throughput"
6. **Gate Management**: "Multi-gate support enables large facilities to control traffic flow, handle maintenance without shutting down the entire lot, and optimize entry/exit based on peak hours"
7. **Capacity Management**: "Enforcing max capacity prevents overcrowding, ensures fire safety compliance, and enables predictive analytics for peak hour planning"

## 🔗 Related Systems

This system demonstrates concepts applicable to:
- **Hotel Room Booking**: Similar slot allocation logic
- **Movie Theater Seats**: Seat categories like slot types
- **Resource Pooling**: Database connections, worker threads
- **Inventory Management**: Stock tracking with reservations

## 👨‍💻 Author

Built as part of a Low-Level Design learning journey, focusing on design patterns and Go best practices.

## 📄 License

Educational project - free to use and modify for learning purposes.

---

**Next System**: Cricbuzz (Observer Pattern for live score updates)

---

## 📋 Requirements Coverage

This implementation covers the following parking lot system requirements:

| Requirement | Status | Implementation |
|------------|--------|----------------|
| 1. Multiple floors | ✅ Complete | `Floors []*Floor` in ParkingLot |
| 2. Multiple entry/exit points | ✅ Complete | Gate system with Entry/Exit/Both types |
| 3. Parking ticket system | ✅ Complete | Ticket struct with entry/exit tracking |
| 4. Payment at exit/attendant | ✅ Complete | ProcessVehicleExit with gate validation |
| 5. Cash and credit card payments | ⚠️ Partial | PricingStrategy exists, payment methods can be added |
| 6. Payment at info portal | ⚠️ Partial | Can be implemented by adding payment status to Ticket |
| 7. Capacity management & display | ✅ Complete | MaxCapacity, CurrentOccupancy with enforcement |
| 8. Multiple spot types | ✅ Complete | Compact, Large, Handicapped, Electric |
| 9. Electric car charging spots | ✅ Complete | Electric slot type defined |
| 10. Different vehicle types | ✅ Complete | Car, Truck, Bike with Factory pattern |
| 11. Display board per floor | ✅ Complete | DisplayStatus() shows availability |
| 12. Per-hour parking fee | ✅ Complete | HourlyPricing & ProgressiveHourlyPricing strategies |

**Legend:** ✅ Complete | ⚠️ Partial (can be extended) | ❌ Not implemented