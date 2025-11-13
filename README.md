Low-Level Design (LLD) - Complete Reference Guide
Table of Contents

What is LLD?
LLD vs Production Code
Core Components
When to Use Service Layer
Design Patterns
Step-by-Step Process
Common Mistakes
Quick Reference
Practice Problems


What is LLD?
Low-Level Design (LLD) focuses on designing the core business logic of a system, NOT the infrastructure.
What LLD Tests:

✅ Domain modeling (entities and relationships)
✅ Business logic and algorithms
✅ Data structures and their trade-offs
✅ Design patterns and SOLID principles
✅ Handling edge cases and constraints

What LLD Does NOT Test:

❌ HTTP/REST API design
❌ Database schema or SQL queries
❌ Frontend/UI design
❌ Deployment or infrastructure
❌ Authentication/Authorization (unless specifically asked)

Formula:
LLD = Domain/Models + Service Layer (if needed)
Key Principle:

Design as if everything is stored in memory (maps/arrays). Focus on the BRAIN (business logic), not the BODY (infrastructure).


LLD vs Production Code
The Relationship
┌─────────────────────────────────────────┐
│              LLD Code                   │
├─────────────────────────────────────────┤
│  ✅ Domain/Models                       │
│     - Entities (structs)                │
│     - Simple entity methods             │
│                                         │
│  ✅ Service Layer (if needed)           │
│     - In-memory storage (maps)          │
│     - Business logic                    │
│     - Complex operations                │
└─────────────────────────────────────────┘
              ↓
       (Extends to)
              ↓
┌─────────────────────────────────────────┐
│         Production Code                 │
├─────────────────────────────────────────┤
│  ✅ Domain/Models (SAME AS LLD!)        │
│     - Entities (structs)                │
│     - Simple entity methods             │
│                                         │
│  ✅ Service Layer (SAME LOGIC!)         │
│     - Database repos (not maps)         │
│     - SAME business logic               │
│                                         │
│  🆕 HTTP/API Layer (NEW!)               │
│     - Request parsing                   │
│     - Response formatting               │
│                                         │
│  🆕 Repository Layer (NEW!)             │
│     - Database queries                  │
│     - Data persistence                  │
└─────────────────────────────────────────┘
Comparison Table
AspectLLDProduction CodeStorageIn-memory (maps, arrays)Database (SQL, NoSQL)LayersDomain + ServiceDomain + Service + HTTP + RepositoryFocusBusiness logic & algorithmsFull working applicationDependenciesNone (self-contained)Payment gateways, email, etc.Code Volume200-500 linesThousands of linesGoalShow design thinkingProduction-ready system
Code Comparison
LLD Code:
go// Domain Model - SAME in both
type Booking struct {
    ID     string
    UserID string
    Amount float64
}

func (b *Booking) IsExpired() bool {
    return time.Now().After(b.ExpiresAt)
}

// Service - In-memory storage
type BookingService struct {
    bookings map[string]*Booking  // ← In-memory
}

func (s *BookingService) CreateBooking(booking *Booking) error {
    s.bookings[booking.ID] = booking  // ← Direct map access
    return nil
}
Production Code:
go// Domain Model - IDENTICAL to LLD!
type Booking struct {
    ID     string
    UserID string
    Amount float64
}

func (b *Booking) IsExpired() bool {
    return time.Now().After(b.ExpiresAt)
}

// Service - Database storage
type BookingService struct {
    bookingRepo repositories.BookingRepository  // ← Database
}

func (s *BookingService) CreateBooking(booking *Booking) error {
    return s.bookingRepo.Create(booking)  // ← Database call
}

// NEW: Repository Layer
type BookingRepository interface {
    Create(booking *Booking) error
    GetByID(id string) (*Booking, error)
}

// NEW: HTTP Handler
type BookingHandler struct {
    service *BookingService
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
    // Parse HTTP, call service, return response
}
What Stays the Same? ✅
go✅ Entity definitions (structs)
✅ Entity methods (operate on self)
✅ Business rules and validation
✅ Service logic and algorithms
✅ Error handling patterns
What Changes? 🔄
go🔄 Storage mechanism:
   LLD:        s.bookings[id] = booking
   Production: s.repo.Create(booking)

🆕 Added layers:
   - HTTP handlers (parse requests, return responses)
   - Repositories (SQL queries, DB connections)

Core Components
1. Domain Entities (Models)
Definition: Structs/classes representing real-world objects.
Characteristics:

Hold data/state
Simple methods that operate only on self
No external dependencies
Used across all parts of the system
IDENTICAL in LLD and Production

Example:
gotype Movie struct {
    ID       string
    Name     string
    Duration time.Duration
    Language string
}

type Booking struct {
    ID        string
    UserID    string
    Show      *Show
    Amount    float64
    Status    BookingStatus
    CreatedAt time.Time
}

// Simple entity methods - operate only on THIS booking
func (b *Booking) IsExpired() bool {
    return time.Now().After(b.ExpiresAt)
}

func (b *Booking) CanBeCancelled() bool {
    return b.Status == Pending || b.Status == Reserved
}

func (b *Booking) GetTotalAmount() float64 {
    total := 0.0
    for _, seat := range b.Seats {
        total += seat.Price
    }
    return total
}

2. Service/Manager Layer
Definition: Orchestrates complex operations across entities.
Characteristics:

Holds collections of entities
Contains business logic
Coordinates multiple entities
Provides search/query capabilities
May hold indexes for performance
Logic is SAME in LLD and Production, only storage changes

LLD Example:
gotype BookMyShowService struct {
    // Collections (in-memory)
    theaters map[string]*Theater
    shows    map[string]*Show
    bookings map[string]*Booking
    
    // Indexes for fast queries
    theatersByCity map[string][]*Theater
    showsByMovie   map[string][]*Show
}

func NewBookMyShowService() *BookMyShowService {
    return &BookMyShowService{
        theaters:       make(map[string]*Theater),
        shows:          make(map[string]*Show),
        bookings:       make(map[string]*Booking),
        theatersByCity: make(map[string][]*Theater),
        showsByMovie:   make(map[string][]*Show),
    }
}

func (s *BookMyShowService) BookTickets(userID, showID string, seatIDs []string) (*Booking, error) {
    // 1. Get show from map
    show := s.shows[showID]
    if show == nil {
        return nil, errors.New("show not found")
    }
    
    // 2. Business logic (SAME in production!)
    if len(seatIDs) > 10 {
        return nil, errors.New("max 10 tickets")
    }
    
    // 3. Check availability
    for _, seatID := range seatIDs {
        if show.ShowSeats[seatID].Status != Available {
            return nil, errors.New("seat not available")
        }
    }
    
    // 4. Create booking
    booking := &Booking{
        ID:     generateID(),
        UserID: userID,
        Show:   show,
        Amount: calculateAmount(seatIDs),
    }
    
    // 5. Save to map (in LLD)
    s.bookings[booking.ID] = booking
    
    return booking, nil
}
Production Example:
gotype BookMyShowService struct {
    // Database repositories (not maps!)
    showRepo    repositories.ShowRepository
    bookingRepo repositories.BookingRepository
}

func (s *BookMyShowService) BookTickets(userID, showID string, seatIDs []string) (*Booking, error) {
    // 1. Get show from database (not map!)
    show, err := s.showRepo.GetByID(showID)
    if err != nil {
        return nil, err
    }
    
    // 2. SAME business logic as LLD!
    if len(seatIDs) > 10 {
        return nil, errors.New("max 10 tickets")
    }
    
    // 3. SAME availability check as LLD!
    for _, seatID := range seatIDs {
        if show.ShowSeats[seatID].Status != Available {
            return nil, errors.New("seat not available")
        }
    }
    
    // 4. SAME booking creation as LLD!
    booking := &Booking{
        ID:     generateID(),
        UserID: userID,
        Show:   show,
        Amount: calculateAmount(seatIDs),
    }
    
    // 5. Save to database (not map!)
    err = s.bookingRepo.Create(booking)
    if err != nil {
        return nil, err
    }
    
    return booking, nil
}
```

**Key Insight:** The business logic is IDENTICAL! Only the storage mechanism changes.

---

## When to Use Service Layer

### Decision Framework
```
┌─────────────────────────────────────────────────┐
│  Do you need a separate Service Layer?         │
├─────────────────────────────────────────────────┤
│                                                 │
│  Ask these questions:                           │
│                                                 │
│  1. Are you managing MULTIPLE instances         │
│     of entities?                                │
│     └─ NO  → Probably no service needed        │
│     └─ YES → Continue...                        │
│                                                 │
│  2. Do you need to SEARCH/QUERY across         │
│     these instances?                            │
│     └─ NO  → Maybe no service needed           │
│     └─ YES → Continue...                        │
│                                                 │
│  3. Do operations span MULTIPLE entities?       │
│     └─ NO  → Probably no service needed        │
│     └─ YES → Continue...                        │
│                                                 │
│  4. Is there ONE "boss" entity that naturally   │
│     manages everything?                         │
│     └─ YES → No service needed (rich entity)   │
│     └─ NO  → Service needed                    │
│                                                 │
└─────────────────────────────────────────────────┘
Quick Decision Table
QuestionNo Service (Rich Entity)Service NeededSingle vs MultipleManaging ONE instanceManaging MANY instancesBoss Entity?Clear "boss" entity existsNo single boss, many equalsCross-entity queries?No searching neededNeed to search/filterOperationsNaturally belong to entitySpan multiple entitiesExamplesParkingLot, Match, ATMBookMyShow, Uber, Amazon

Design Patterns
Pattern 1: Rich Domain Model (No Service)
When to Use:

Single central entity (ATM, ParkingLot, Match, VendingMachine)
Operations naturally belong to that entity
No complex cross-entity queries
Entity can manage its own state

Structure:
gotype CentralEntity struct {
    // State
    field1 Type1
    field2 Type2
    collection map[string]*SubEntity
}

// Entity manages everything
func (e *CentralEntity) Operation1() { ... }
func (e *CentralEntity) Operation2() { ... }
func (e *CentralEntity) ComplexOperation() { ... }
Example 1: Parking Lot
gotype ParkingLot struct {
    floors map[int]*Floor
    spots  map[string]*ParkingSpot
}

// Operations on ParkingLot itself
func (p *ParkingLot) Park(vehicle *Vehicle) (*Ticket, error) {
    spot := p.findAvailableSpot(vehicle.Type)
    if spot == nil {
        return nil, errors.New("no spot available")
    }
    spot.Park(vehicle)
    return &Ticket{Spot: spot, Vehicle: vehicle}, nil
}

func (p *ParkingLot) Unpark(ticket *Ticket) error {
    return ticket.Spot.Unpark()
}

func (p *ParkingLot) GetAvailableSpots() []*ParkingSpot {
    // Search within THIS parking lot
}

func main() {
    parkingLot := &ParkingLot{...}
    
    // Everything operates on THIS parking lot
    ticket, _ := parkingLot.Park(vehicle)
    parkingLot.Unpark(ticket)
    
    // No need for a service!
}
Example 2: Cricket Match
gotype Match struct {
    MatchID     string
    Teams       [2]Team
    TeamBatting *Team
    Status      MatchStatus
    Observers   []Observer
}

// Match manages itself
func (m *Match) startMatch() {
    m.Status = InProgress
    m.TeamBatting = &m.Teams[0]
}

func (m *Match) addRuns(runs int) {
    m.TeamBatting.Score.Runs += runs
    m.NotifyAll()
}

func (m *Match) addWicket() {
    m.TeamBatting.Score.Wickets++
    if m.TeamBatting.Score.Wickets >= 10 {
        m.switchInnings()
    }
    m.NotifyAll()
}

func main() {
    match := &Match{...}
    
    // Everything operates on THIS match
    match.startMatch()
    match.addRuns(4)
    match.addWicket()
    
    // No service needed!
}
Why No Service?

✅ Focus is on ONE parking lot / ONE match
✅ No need to search across multiple parking lots / matches
✅ Entity naturally controls all operations
✅ Self-contained scope


Pattern 2: Service/Manager (Separate Service)
When to Use:

Multiple entities at same level (no single "boss")
Need cross-entity queries (search, filter, aggregate)
Complex workflows spanning multiple entities
Need indexes for performance

Structure:
go// Simple entities
type Entity1 struct {
    ID   string
    Name string
}

type Entity2 struct {
    ID       string
    Entity1  *Entity1
}

// Service manages collections
type SystemService struct {
    entity1s map[string]*Entity1
    entity2s map[string]*Entity2
    
    // Indexes
    entity1sByAttribute map[string][]*Entity1
}

func (s *SystemService) CreateEntity1(...) { ... }
func (s *SystemService) SearchEntity1(...) { ... }
func (s *SystemService) ComplexOperation(...) { ... }
Example 1: BookMyShow
go// Simple entities
type Theater struct {
    ID       string
    Location Location
}

type Show struct {
    ID      string
    Movie   *Movie
    Theater *Theater
}

// Service coordinates everything
type BookMyShowService struct {
    theaters       map[string]*Theater
    shows          map[string]*Show
    bookings       map[string]*Booking
    
    // Indexes for queries
    theatersByCity map[string][]*Theater
    showsByMovie   map[string][]*Show
}

func (s *BookMyShowService) SearchMoviesInCity(city string) []*Movie {
    // Query across multiple theaters
    theaters := s.theatersByCity[city]
    
    movies := make(map[string]*Movie)
    for _, theater := range theaters {
        shows := s.showsByTheater[theater.ID]
        for _, show := range shows {
            movies[show.Movie.ID] = show.Movie
        }
    }
    return mapToSlice(movies)
}

func (s *BookMyShowService) BookTickets(...) (*Booking, error) {
    // Complex workflow across multiple entities
}

func main() {
    service := NewBookMyShowService()
    
    // Create multiple entities
    service.AddTheater(theater1)
    service.AddTheater(theater2)
    service.AddShow(show1)
    service.AddShow(show2)
    
    // Cross-entity operations
    movies := service.SearchMoviesInCity("Mumbai")
    booking, _ := service.BookTickets(...)
}
Example 2: Uber/Ride Sharing
gotype Rider struct {
    ID       string
    Location Location
}

type Driver struct {
    ID        string
    Location  Location
    Available bool
}

type Ride struct {
    ID     string
    Rider  *Rider
    Driver *Driver
    Status RideStatus
}

type UberService struct {
    riders  map[string]*Rider
    drivers map[string]*Driver
    rides   map[string]*Ride
    
    // Indexes
    availableDriversByArea map[string][]*Driver
}

func (s *UberService) RequestRide(riderID string, from, to Location) (*Ride, error) {
    // 1. Find rider
    rider := s.riders[riderID]
    
    // 2. Find nearby available drivers
    drivers := s.findNearbyDrivers(from)
    if len(drivers) == 0 {
        return nil, errors.New("no drivers available")
    }
    
    // 3. Match driver
    driver := drivers[0]
    driver.Available = false
    
    // 4. Create ride
    ride := &Ride{
        ID:     generateID(),
        Rider:  rider,
        Driver: driver,
        From:   from,
        To:     to,
        Status: RideRequested,
    }
    
    s.rides[ride.ID] = ride
    return ride, nil
}

func (s *UberService) GetRiderHistory(riderID string) []*Ride {
    // Search across all rides
}

func main() {
    service := NewUberService()
    
    // Multiple riders and drivers
    service.AddRider(rider1)
    service.AddRider(rider2)
    service.AddDriver(driver1)
    service.AddDriver(driver2)
    
    // Cross-entity operations
    ride, _ := service.RequestRide("rider1", locationA, locationB)
    history := service.GetRiderHistory("rider1")
}
Why Service Needed?

✅ Managing MANY theaters/shows or riders/drivers
✅ Need to search: "movies in Mumbai", "nearby drivers"
✅ Operations span entities: matching rider to driver
✅ No single "boss" entity


Comparison: When Would Each Pattern Need a Service?
Parking Lot: When Service Becomes Needed
go// Current: Single parking lot (no service)
parkingLot.Park(vehicle)

// Extended: Multiple parking lots (need service!)
type ParkingService struct {
    parkingLots map[string]*ParkingLot
}

func (s *ParkingService) FindAvailableParking(location string, vehicleType VehicleType) (*ParkingLot, error) {
    // Search across multiple parking lots
}

func (s *ParkingService) GetAllParkingLots() []*ParkingLot {
    // Query all lots
}
Cricket Match: When Service Becomes Needed
go// Current: Single match (no service)
match.addRuns(4)

// Extended: Multiple matches (need service!)
type CricbuzzService struct {
    matches       map[string]*Match
    matchesByTeam map[string][]*Match
}

func (s *CricbuzzService) GetLiveMatches() []*Match {
    // Query across all matches
}

func (s *CricbuzzService) GetMatchesByTeam(teamName string) []*Match {
    // Search by team
}

func (s *CricbuzzService) GetTeamStatistics(teamName string) *Stats {
    // Aggregate across matches
}
```

---

## Step-by-Step LLD Process

### Step 1: Understand Requirements (5 minutes)

**Questions to Ask:**
1. What are the main use cases?
2. What are the constraints?
3. What operations need to be fast?
4. Are we managing ONE instance or MANY?

**Example: BookMyShow**
- Use cases: Search movies, book tickets, cancel booking
- Constraints: Max 10 tickets per booking
- Optimize: Fast movie search by city
- **Scope: MANY theaters/shows → Need Service**

**Example: Cricket Match**
- Use cases: Track score, notify observers, switch innings
- Constraints: 10 wickets per innings
- Optimize: Real-time notifications
- **Scope: ONE match → No Service needed**

---

### Step 2: Identify Entities (10 minutes)

**Technique:** List all "nouns" from requirements

**Example: BookMyShow**
```
Entities:
- User
- Movie
- Theater
- Screen (Theater has many)
- Show (Movie + Screen + Time)
- Seat (physical seat)
- ShowSeat (seat for a specific show) ← Important!
- Booking
- Payment

Relationships:
- Theater (1) → (M) Screen
- Show (M) → (1) Movie
- Show (M) → (1) Screen
- Booking (1) → (M) ShowSeat
```

---

### Step 3: Decide: Service or Rich Entity? (5 minutes)

**Ask:**

1. **Am I managing ONE instance or MANY?**
   - ONE → Likely no service
   - MANY → Likely need service

2. **Is there a clear "boss" entity?**
   - YES (ParkingLot, Match, ATM) → No service
   - NO (many equal entities) → Need service

3. **Do I need cross-entity queries?**
   - NO → Probably no service
   - YES → Need service

**Example: ParkingLot**
```
1. ONE parking lot ✓
2. ParkingLot is the boss ✓
3. No cross-lot queries ✓

Decision: No service needed!
```

**Example: BookMyShow**
```
1. MANY theaters/shows ✓
2. No single boss entity ✓
3. Need to search across theaters ✓

Decision: Service needed!

Step 4: Design Entities (15 minutes)
gotype Entity struct {
    ID     string
    Field1 Type1
    Field2 Type2
}

// Simple methods (operate on self only)
func (e *Entity) SimpleMethod() bool {
    // Logic using only e's fields
    return true
}
Common Pitfall: Wrong Granularity
go// ❌ BAD: Status on Seat (global)
type Seat struct {
    Status SeatStatus  // Same for all shows!
}

// ✅ GOOD: Status on ShowSeat (per-show)
type Seat struct {
    Number string
    Type   SeatType  // Physical property
}

type ShowSeat struct {
    Seat   *Seat
    Show   *Show
    Status SeatStatus  // Status for THIS show
    Price  float64     // Price for THIS show
}

Step 5: Design Service (if needed) (20 minutes)
gotype ServiceName struct {
    // Collections
    entities map[string]*Entity
    
    // Indexes
    entitiesByAttribute map[string][]*Entity
}

func NewServiceName() *ServiceName {
    return &ServiceName{
        entities:            make(map[string]*Entity),
        entitiesByAttribute: make(map[string][]*Entity),
    }
}

// CRUD
func (s *ServiceName) AddEntity(entity *Entity) {
    s.entities[entity.ID] = entity
    // Update indexes
}

// Queries
func (s *ServiceName) SearchEntities(attr string) []*Entity {
    return s.entitiesByAttribute[attr]
}

// Complex operations
func (s *ServiceName) ComplexOperation(...) error {
    // Business logic
}

Step 6: Handle Edge Cases (10 minutes)
Common Edge Cases:

Race Conditions

go// Two users booking same seat
func (s *Service) BookTickets(...) {
    // Check if seat already being reserved
    if seat.Status == Reserved && seat.ReservedUntil.After(time.Now()) {
        return errors.New("seat temporarily reserved")
    }
}

Expiry/Cleanup

gofunc (s *Service) CleanupExpired() {
    for _, booking := range s.bookings {
        if booking.IsExpired() {
            s.CancelBooking(booking.ID)
        }
    }
}

Consistency

gofunc (s *Service) ComplexOperation() error {
    // All-or-nothing: if any step fails, rollback all
    entity1 := s.createEntity1()
    entity2 := s.createEntity2()
    
    if err := s.validateBoth(entity1, entity2); err != nil {
        // Don't save either
        return err
    }
    
    // Both valid, save together
    s.entities1[entity1.ID] = entity1
    s.entities2[entity2.ID] = entity2
    return nil
}

Common Mistakes
1. ❌ Wrong Granularity of State
go// ❌ BAD
type Seat struct {
    Status SeatStatus  // Global state across all shows
}

// ✅ GOOD
type ShowSeat struct {
    Seat   *Seat
    Show   *Show
    Status SeatStatus  // Per-show state
}
Rule: If property value depends on context, create a context-specific entity.

2. ❌ Missing Navigation Paths
go// ❌ BAD
type Show struct {
    Movie  *Movie
    Screen *Screen  // Can't get to Theater!
}

// ✅ GOOD
type Show struct {
    Movie   *Movie
    Screen  *Screen
    Theater *Theater  // Easy navigation
}

3. ❌ Entity Knowing About All Entities
go// ❌ BAD
type Theater struct {
    ID          string
    AllTheaters map[string]*Theater  // What?!
}

// ✅ GOOD
type Service struct {
    theaters map[string]*Theater  // Service manages all
}

4. ❌ Complex Logic in Entities
go// ❌ BAD
type Booking struct {
    paymentGateway PaymentGateway  // External dependency
}

func (b *Booking) ProcessPayment() {
    b.paymentGateway.Charge(b.Amount)
}

// ✅ GOOD
type Service struct {
    paymentGateway PaymentGateway
}

func (s *Service) ProcessPayment(booking *Booking) {
    s.paymentGateway.Charge(booking.Amount)
}

5. ❌ Not Using Indexes
go// ❌ BAD: O(n) every time
func (s *Service) GetMoviesInCity(city string) []*Movie {
    for _, theater := range s.theaters {
        if theater.Location.City == city {
            // ...
        }
    }
}

// ✅ GOOD: O(1) with index
type Service struct {
    theaters       map[string]*Theater
    theatersByCity map[string][]*Theater  // Index
}

func (s *Service) GetMoviesInCity(city string) []*Movie {
    theaters := s.theatersByCity[city]  // O(1)
    // ...
}

6. ❌ Forgetting Constructors
go// ❌ BAD
service := &Service{}
service.entities["id"] = entity  // Panic: nil map

// ✅ GOOD
func NewService() *Service {
    return &Service{
        entities: make(map[string]*Entity),
    }
}
```

---

## Quick Reference

### Entity vs Service Decision
```
Operation affects only THIS entity?
└─ YES → Entity method

Operation searches across entities?
└─ YES → Service method

Operation coordinates multiple entities?
└─ YES → Service method

Operation needs external dependencies?
└─ YES → Service method

Is there a clear "boss" entity?
└─ YES → Rich entity (no service)
└─ NO  → Service pattern
LLD Checklist
Before Submitting:
Entities:

 All key entities identified
 Relationships clearly defined
 No wrong granularity (like status on Seat)
 Entity methods are simple
 No collections of other entities in entities

Service (if used):

 Collections (maps) for all entities
 Indexes for frequently queried attributes
 Constructor initializes all maps
 CRUD methods present
 Search/query methods
 Complex workflows implemented

Business Logic:

 All use cases covered
 Business rules enforced
 Edge cases handled
 Error handling
 Validation

Code Quality:

 Clear naming
 Comments for complex logic
 Constants (no hardcoded values)
 Helper functions
 Main function demonstrates usage


LLD Template
Template 1: With Service (Multiple Entities)
gopackage main

import (
    "errors"
    "fmt"
    "time"
)

// ============================================
// DOMAIN ENTITIES
// ============================================

type Entity1 struct {
    ID   string
    Name string
}

func (e *Entity1) SimpleMethod() bool {
    return true
}

type Entity2 struct {
    ID      string
    Entity1 *Entity1
}

// ============================================
// ENUMS / CONSTANTS
// ============================================

type Status string

const (
    StatusActive   Status = "ACTIVE"
    StatusInactive Status = "INACTIVE"
)

// ============================================
// SERVICE
// ============================================

type SystemService struct {
    entity1s map[string]*Entity1
    entity2s map[string]*Entity2
    
    entity1sByAttribute map[string][]*Entity1
}

func NewSystemService() *SystemService {
    return &SystemService{
        entity1s:            make(map[string]*Entity1),
        entity2s:            make(map[string]*Entity2),
        entity1sByAttribute: make(map[string][]*Entity1),
    }
}

func (s *SystemService) AddEntity1(entity *Entity1) {
    s.entity1s[entity.ID] = entity
    attr := entity.SomeAttribute
    s.entity1sByAttribute[attr] = append(s.entity1sByAttribute[attr], entity)
}

func (s *SystemService) SearchEntity1s(attr string) []*Entity1 {
    return s.entity1sByAttribute[attr]
}

func (s *SystemService) ComplexOperation(param1, param2 string) (*Entity2, error) {
    entity1 := s.entity1s[param1]
    if entity1 == nil {
        return nil, errors.New("entity1 not found")
    }
    
    if !entity1.SimpleMethod() {
        return nil, errors.New("validation failed")
    }
    
    entity2 := &Entity2{
        ID:      generateID(),
        Entity1: entity1,
    }
    
    s.entity2s[entity2.ID] = entity2
    return entity2, nil
}

// ============================================
// HELPERS
// ============================================

func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ============================================
// MAIN
// ============================================

func main() {
    service := NewSystemService()
    
    entity1 := &Entity1{ID: "e1", Name: "Example"}
    service.AddEntity1(entity1)
    
    result, err := service.ComplexOperation("e1", "param")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Success: %+v\n", result)
}
Template 2: Without Service (Rich Entity)
gopackage main

import (
    "errors"
    "fmt"
)

// ============================================
// DOMAIN ENTITIES
// ============================================

type CentralEntity struct {
    ID         string
    SubEntities map[string]*SubEntity
}

func NewCentralEntity(id string) *CentralEntity {
    return &CentralEntity{
        ID:         id,
        SubEntities: make(map[string]*SubEntity),
    }
}

// Entity manages itself
func (c *CentralEntity) AddSubEntity(sub *SubEntity) {
    c.SubEntities[sub.ID] = sub
}

func (c *CentralEntity) Operation1() error {
    // Business logic
    return nil
}

func (c *CentralEntity) Operation2(param string) (*Result, error) {
    sub := c.SubEntities[param]
    if sub == nil {
        return nil, errors.New("not found")
    }
    
    result := &Result{
        SubEntity: sub,
    }
    
    return result, nil
}

type SubEntity struct {
    ID   string
    Name string
}

type Result struct {
    SubEntity *SubEntity
}

// ============================================
// MAIN
// ============================================

func main() {
    entity := NewCentralEntity("main")
    
    sub := &SubEntity{ID: "s1", Name: "Sub"}
    entity.AddSubEntity(sub)
    
    result, err := entity.Operation2("s1")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Success: %+v\n", result)
}
```

---

## Practice Problems

### Level 1: Rich Entity (No Service Needed)

1. **Parking Lot**
   - Single parking lot with multiple floors
   - Operations: Park, unpark, check availability
   - **No service needed** - ParkingLot entity manages itself

2. **ATM**
   - Single ATM machine
   - Operations: Insert card, withdraw, check balance, deposit
   - **No service needed** - ATM entity manages itself

3. **Vending Machine**
   - Single machine with products
   - Operations: Select product, insert coin, dispense, return change
   - **No service needed** - VendingMachine entity manages itself

4. **Cricket Match (Cricbuzz)**
   - Single match with teams
   - Operations: Add runs, wickets, notify observers
   - **No service needed** - Match entity manages itself

5. **Chess Game**
   - Single game with board
   - Operations: Move piece, check valid moves, checkmate
   - **No service needed** - Game entity manages itself

### Level 2: Service Needed

6. **Hotel Booking System**
   - Multiple hotels, rooms, reservations
   - Operations: Search rooms, book, cancel
   - **Service needed** - Manage multiple hotels

7. **Uber/Ride Sharing**
   - Multiple riders, drivers, rides
   - Operations: Request ride, match driver, track ride
   - **Service needed** - Match riders to drivers

8. **Food Delivery (Swiggy/Zomato)**
   - Multiple restaurants, orders, delivery partners
   - Operations: Browse, order, assign delivery, track
   - **Service needed** - Coordinate restaurants and deliveries

9. **BookMyShow**
   - Multiple theaters, shows, bookings
   - Operations: Search movies, book tickets, handle payments
   - **Service needed** - Search across theaters

10. **Library Management**
    - Multiple books, members, loans
    - Operations: Check out, return, search, overdue
    - **Service needed** - Manage multiple books/members

### Level 3: Advanced

11. **LinkedIn**
    - Users, connections, posts, jobs
    - Operations: Connect, post, apply, recommend
    - **Service needed** - Complex social graph

12. **Stock Trading**
    - Users, stocks, orders, portfolio
    - Operations: Place order, match orders, execute trades
    - **Service needed** - Match buy/sell orders

13. **Amazon/E-commerce**
    - Products, cart, orders, inventory
    - Operations: Browse, add to cart, checkout, track
    - **Service needed** - Complex workflows

---

## Key Takeaways

### The Golden Rules

1. **LLD = Domain + Service (if needed)**
   - Focus on business logic, not infrastructure
   - Design with in-memory storage (maps)

2. **Domain models are IDENTICAL in LLD and Production**
   - Entity methods don't change
   - Business logic stays the same

3. **Production adds layers around LLD**
   - HTTP layer (input)
   - Repository layer (output)
   - Core logic remains the same

4. **When to use Service:**
   - Managing MULTIPLE instances
   - Need cross-entity queries
   - No single "boss" entity
   - Complex coordination needed

5. **When NOT to use Service:**
   - Single instance (ONE parking lot, ONE match)
   - Clear "boss" entity manages everything
   - No searching across instances
   - Self-contained operations

### Time Management in Interviews
```
Total: 45 minutes

5 min  → Understand requirements
5 min  → Decide: Service or Rich Entity?
10 min → Design entities
15 min → Implement service/entity methods
5 min  → Handle edge cases
5 min  → Test with main function
```

### Communication Tips

1. **Think out loud:** Explain your reasoning
2. **Ask clarifying questions:** Don't assume
3. **Start simple:** Get basic structure right first
4. **Iterate:** Add complexity after foundation is solid
5. **Discuss trade-offs:** Show you understand alternatives

---

## Summary Diagram
```
┌───────────────────────────────────────────────────────┐
│                   LLD APPROACH                        │
├───────────────────────────────────────────────────────┤
│                                                       │
│   Single Instance?                Multiple Entities? │
│   Boss Entity?                    Need Queries?      │
│         ↓                                ↓            │
│                                                       │
│   Rich Entity                      Service Pattern   │
│   ─────────────                    ────────────────  │
│   type Entity {                    type Service {    │
│     collections                      entities map    │
│   }                                   indexes map     │
│                                     }                 │
│   func (e *Entity)                                   │
│     Operations()                   func (s *Service) │
│                                      Operations()     │
│   Examples:                                          │
│   - ParkingLot                     Examples:         │
│   - Match                          - BookMyShow      │
│   - ATM                            - Uber            │
│   - Chess                          - Hotel Booking   │
│                                                       │
└───────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────┐
│         LLD → PRODUCTION TRANSFORMATION               │
├───────────────────────────────────────────────────────┤
│                                                       │
│   LLD Code              Production Code              │
│   ────────              ───────────────              │
│                                                       │
│   Domain ────────────→  Domain (SAME!)               │
│   (entities +           (entities +                  │
│    methods)              methods)                    │
│                                                       │
│   Service ───────────→  Service (SAME LOGIC!)        │
│   (maps,                (repos,                      │
│    business logic)       same logic)                 │
│                                                       │
│   [Nothing] ─────────→  HTTP Layer (NEW!)            │
│                         (handlers)                   │
│                                                       │
│   [Nothing] ─────────→  Repository (NEW!)            │
│                         (database)                   │
│                                                       │
└───────────────────────────────────────────────────────┘

Final Thoughts

LLD is about showing you can think clearly about domain modeling and business logic. The infrastructure (HTTP, DB) is just plumbing that wraps around your core design.


Your LLD domain models and business logic should be so well-designed that they can be directly used in production code with minimal changes.


The decision between Rich Entity and Service isn't about which is "better" - it's about which naturally fits your problem domain.

Good luck with your LLD practice! 🚀

Last Updated: 2024
Version: 2.0RetryClaude can make mistakes. Please double-check responses. Sonnet 4.5