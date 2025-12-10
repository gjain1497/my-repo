# 📚 README: Models vs Services - Where Does Logic Belong?

## 🎯 The Core Question

**When designing systems, where should business logic live?**
- In the **Models** (Rich Domain Models)?
- In the **Services** (Anemic Domain Models)?

This README explains both approaches and provides a consistent strategy for LLD.

---

## 📊 Two Approaches

### **Approach 1: Rich Domain Models (Logic in Models)**

**Philosophy:** Objects should contain both data AND behavior. Models "know how to do their own work."

#### **Example: Java with Inheritance**

```java
// Business logic IN the model
abstract class Piece {
    protected Color color;
    protected Position position;
    
    // ✅ Behavior methods in model
    public abstract List<Position> getValidMoves(Board board);
    public abstract boolean isValidMove(Position to, Board board);
}

// Each piece type is a subclass
class King extends Piece {
    @Override
    public List<Position> getValidMoves(Board board) {
        // King-specific movement logic HERE
        List<Position> moves = new ArrayList<>();
        // ... king movement logic ...
        return moves;
    }
}

class Queen extends Piece {
    @Override
    public List<Position> getValidMoves(Board board) {
        // Queen-specific movement logic HERE
        List<Position> moves = new ArrayList<>();
        // ... queen movement logic ...
        return moves;
    }
}
```

**Usage:**
```java
Piece king = new King(Color.WHITE, position);
List<Position> validMoves = king.getValidMoves(board);  // ✅ Model does the work
```

#### **Characteristics:**
- **Data + Behavior** together in model
- **Inheritance** for polymorphism
- **Encapsulation** - piece "knows" how to move itself

#### **Pros:**
- ✅ Encapsulation - Logic close to data
- ✅ Traditional OOP principles
- ✅ Reads like business rules
- ✅ Domain experts can understand code

#### **Cons:**
- ❌ Models become "fat" (lots of logic)
- ❌ Hard to test models independently
- ❌ Violates Single Responsibility Principle
- ❌ Tight coupling between data and behavior

#### **When Used:**
- Domain-Driven Design (DDD) projects
- Complex business domains (Banking, E-commerce)
- When domain experts are involved
- Financial systems, Booking systems

---

### **Approach 2: Anemic Domain Models (Logic in Services)**

**Philosophy:** Models contain ONLY data. Services contain ALL behavior. Clear separation of concerns.

#### **Example: Java/Go with Services**

**Models (Data Only):**
```java
// Java
enum PieceType {
    KING, QUEEN, ROOK, BISHOP, KNIGHT, PAWN
}

class Piece {
    private PieceType type;  // ✅ Just an enum!
    private Color color;
    private Position position;
    private boolean hasMoved;
    
    // Only getters/setters - NO business logic
}
```

```go
// Go
type PieceType string

const (
    King   PieceType = "KING"
    Queen  PieceType = "QUEEN"
    Rook   PieceType = "ROOK"
    Bishop PieceType = "BISHOP"
    Knight PieceType = "KNIGHT"
    Pawn   PieceType = "PAWN"
)

type Piece struct {
    Type     PieceType  // ✅ Just an enum!
    Color    Color
    HasMoved bool
    Position Position
    // NO methods - just data
}
```

**Services (All Logic):**
```java
// Java
interface PieceMovementService {
    List<Position> getValidMoves(Piece piece, Board board);
    boolean isValidMove(Piece piece, Position to, Board board);
}

class KingMovementService implements PieceMovementService {
    @Override
    public List<Position> getValidMoves(Piece piece, Board board) {
        // King-specific movement logic HERE
        List<Position> moves = new ArrayList<>();
        // ... king movement logic ...
        return moves;
    }
}

class QueenMovementService implements PieceMovementService {
    @Override
    public List<Position> getValidMoves(Piece piece, Board board) {
        // Queen-specific movement logic HERE
        List<Position> moves = new ArrayList<>();
        // ... queen movement logic ...
        return moves;
    }
}
```

```go
// Go
type PieceMovementService interface {
    GetValidMoves(piece *Piece, board *Board) []Position
    IsValidMove(piece *Piece, from, to Position, board *Board) bool
}

type KingMovementService struct{}

func (k *KingMovementService) GetValidMoves(piece *Piece, board *Board) []Position {
    // King-specific movement logic HERE
    validMoves := []Position{}
    // ... king movement logic ...
    return validMoves
}

type QueenMovementService struct{}

func (q *QueenMovementService) GetValidMoves(piece *Piece, board *Board) []Position {
    // Queen-specific movement logic HERE
    validMoves := []Position{}
    // ... queen movement logic ...
    return validMoves
}
```

**Usage:**
```java
// Java
Piece piece = new Piece(PieceType.KING, Color.WHITE, position);
PieceMovementService service = new KingMovementService();
List<Position> validMoves = service.getValidMoves(piece, board);  // ✅ Service does the work
```

```go
// Go
piece := &Piece{Type: King, Color: White, Position: pos}
service := &KingMovementService{}
validMoves := service.GetValidMoves(piece, board)  // ✅ Service does the work
```

#### **Characteristics:**
- **Data** in models (structs/classes)
- **Behavior** in services (interfaces + implementations)
- **Composition** for polymorphism (Go) or Interfaces (Java)

#### **Pros:**
- ✅ **Single Responsibility Principle** - Models = data, Services = logic
- ✅ **Easy to test** - Test services independently
- ✅ **More flexible** - Swap services easily
- ✅ **Clear separation** of concerns
- ✅ **Idiomatic in Go** (composition over inheritance)

#### **Cons:**
- ❌ More service code
- ❌ Services need to know model internals
- ❌ Can feel "anemic" (models are just data bags)

#### **When Used:**
- Clean Architecture projects
- API/Backend services
- CRUD-heavy applications
- Go language projects (idiomatic)
- High testability requirements
- Large teams with junior developers

---

## 🎯 The Key Decision: Inheritance/Composition Location

### **Traditional OOP (Approach 1):**
```
Inheritance/Composition at MODEL layer
├── class King extends Piece
├── class Queen extends Piece
└── class Rook extends Piece
```

### **Clean Architecture (Approach 2):**
```
Composition at SERVICE layer
├── Models: Piece struct with PieceType enum
└── Services:
    ├── KingMovementService
    ├── QueenMovementService
    └── RookMovementService
```

---

## 📊 Comparison Table

| Aspect | Rich Domain Models | Anemic Models + Services |
|--------|-------------------|--------------------------|
| **Data Location** | In models | In models |
| **Logic Location** | In models | In services |
| **Inheritance/Composition** | At model layer | At service layer |
| **Testing** | Test models with logic | Test services independently |
| **SRP** | ❌ Models do 2 things | ✅ Clear separation |
| **OCP** | ✅ Add subclasses | ✅ Add new services |
| **Encapsulation** | ✅ Strong | ⚠️ Weaker |
| **Flexibility** | ⚠️ Behavior fixed | ✅ Can swap services |
| **Go Idiomatic** | ❌ No inheritance | ✅ Composition |
| **Java Common** | ✅ DDD projects | ✅ Spring Boot apps |

---

## 🏢 Industry Usage

### **Rich Domain Models (Approach 1):**
- **Who:** Enterprise DDD teams, Financial systems, E-commerce
- **Examples:** Banking (Account rules), Insurance (Policy rules), Uber (Ride domain)
- **Languages:** Java (with inheritance), C# (with inheritance)

### **Anemic Models + Services (Approach 2):**
- **Who:** Startups, API-focused teams, Microservices
- **Examples:** Google services, Netflix APIs, Most REST APIs
- **Languages:** Go (composition), Java (Spring Boot often), Python, Node.js

### **Go Community Specifically:**
- **Strongly favors Approach 2** (Services)
- **Why:** No inheritance, composition over inheritance, interface-driven

---

## 🎯 Recommended Approach for LLD

### **Use Approach 2: Anemic Models + Services**

**Reasons:**
1. ✅ **Clear separation** of concerns (easier to explain)
2. ✅ **Shows SOLID** principles clearly
3. ✅ **More testable** (important for interviews)
4. ✅ **Consistent architecture** across all systems
5. ✅ **Idiomatic in Go**
6. ✅ **What interviewers expect** in most cases

---

## 📋 Universal Pattern for All Systems

```
┌─────────────────────────────────────────────────────────┐
│                  CONSISTENT ARCHITECTURE                 │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  MODELS LAYER (Data Only)                              │
│  ├── Structs/Classes with fields                       │
│  ├── Enums (not inheritance)                           │
│  ├── No business logic methods                         │
│  └── Only getters/setters if needed                    │
│                                                          │
│  SERVICES LAYER (All Logic)                            │
│  ├── Interfaces (define behavior)                      │
│  ├── Implementations (composition/polymorphism)        │
│  ├── All business logic here                           │
│  └── Strategy Pattern for different behaviors          │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## 🎯 Examples Across Different Systems

### **Chess:**
```
Models:
├── Piece (Type: enum, Color, Position)
└── Board, Player, Move

Services:
├── KingMovementService
├── QueenMovementService
└── RookMovementService
```

### **Payment System:**
```
Models:
├── Payment (Type: enum, Amount)
└── Transaction

Services:
├── CreditCardProcessor
├── UPIProcessor
└── CashProcessor
```

### **Notification System:**
```
Models:
├── Notification (Type: enum, Message)
└── Recipient

Services:
├── EmailSender
├── SMSSender
└── PushSender
```

### **Vehicle Rental:**
```
Models:
├── Vehicle (Type: enum, Details)
└── Booking

Services:
├── CarPricingStrategy
├── BikePricingStrategy
└── TruckPricingStrategy
```

---

## 💡 Interview Response Template

**If Asked:** "Why don't you use inheritance in your models?"

**Response:**
```
"I keep models as pure data structures and do all behavior/
polymorphism at the service layer through interfaces and 
composition. This approach gives me:

1. Clear separation of concerns (SRP)
2. Better testability - I can test services independently
3. More flexibility - I can swap service implementations easily
4. Consistent architecture across all my systems
5. Alignment with Go's 'composition over inheritance' philosophy

Both approaches are valid - rich domain models (DDD) vs anemic 
models + services. I choose anemic models because it's more 
aligned with Clean Architecture principles and makes the SOLID 
principles more explicit."
```

---

## 🎯 Key Takeaway

**Statement:**
> "I always do composition/polymorphism at the **SERVICE layer**, not the **MODEL layer**."

**This means:**
- ✅ Models = Data only (enums, simple structs)
- ✅ Services = All logic (interfaces, implementations, composition)
- ✅ Consistent across ALL systems
- ✅ Valid architectural choice
- ✅ Defensible in interviews and production

---

## 📚 Related Patterns

- **Strategy Pattern** - Different service implementations
- **Factory Pattern** - Choose the right service
- **Dependency Injection** - Inject services into other services
- **Repository Pattern** - Data access in services
- **Clean Architecture** - Separation of concerns

---

## 🚀 Summary

| Question | Answer |
|----------|--------|
| **Where does data live?** | Models |
| **Where does logic live?** | Services |
| **Where is inheritance/composition?** | Service layer (not model layer) |
| **Is this production-ready?** | Yes (used by Google, Netflix, etc.) |
| **Is this interview-appropriate?** | Yes (shows SOLID clearly) |
| **Is this Go-idiomatic?** | Yes (composition over inheritance) |

---

**End of README**

---

*This document explains the architectural decision of keeping models as pure data structures and implementing all business logic and polymorphism at the service layer through interfaces and composition.*