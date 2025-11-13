# 🏏 Cricbuzz - Live Cricket Score System

A real-time cricket score tracking system built in Go, demonstrating the **Observer Pattern** for event-driven architecture. This system enables multiple displays/observers to automatically receive updates when match scores change, without tight coupling between the match and its observers.

## 🎯 Project Overview

**Learning Goal:** Master the **Observer Pattern** through a practical cricket scoring system

**Key Concept:** When a match score changes (runs, wickets, overs), ALL subscribed observers (scoreboards, apps, commentary) are automatically notified and updated without the match needing to know the specific implementation details of each observer.

---

## 📋 Features Implemented

### Core Functionality
- ✅ **Live Match Tracking**: Track runs, wickets, and overs in real-time
- ✅ **Two-Team Matches**: Support for matches between two teams
- ✅ **Innings Management**: Automatic innings switching after 10 wickets
- ✅ **Event-Based Updates**: Observers notified on scoring events (runs/wickets)
- ✅ **Over Tracking**: Proper over calculation on scoring events
- ✅ **Match Lifecycle**: Start, in-progress, and completion states
- ✅ **Winner Determination**: Automatically determines match winner

**Note:** This is an *event-driven* system that updates on runs and wickets, not true ball-by-ball tracking (which would include dot balls). This design choice reduces notifications while still providing real-time score updates.

### Observer Pattern Implementation
- ✅ **Multiple Observers**: Support for unlimited number of observers
- ✅ **Auto-Notifications**: All observers notified automatically on score changes
- ✅ **Loose Coupling**: Match doesn't know about specific observer implementations
- ✅ **Dynamic Subscription**: Observers can subscribe/unsubscribe at runtime
- ✅ **Different Display Styles**: Each observer can present data differently

---

## 🏗️ System Architecture

### Core Entities

#### 1. **Match** (Subject/Observable)
The central entity that maintains match state and notifies observers.

**Key Fields:**
- `MatchID`: Unique identifier
- `Teams [2]Team`: Two teams playing
- `TeamBatting *Team`: Pointer to currently batting team
- `Status MatchStatus`: NOT_STARTED, IN_PROGRESS, COMPLETED
- `Observers []Observer`: List of subscribed observers

**Key Methods:**
- `Subscribe(observer Observer)`: Add new observer
- `NotifyAll()`: Notify all observers of state change
- `addRuns(runs int)`: Add runs and notify observers
- `addWicket()`: Add wicket and notify observers
- `startMatch()`: Initialize match
- `switchInnings()`: Change batting team
- `endMatch()`: Complete match and determine winner

#### 2. **Team**
Represents a cricket team with players and score.

**Key Fields:**
- `Name string`: Team name (e.g., "India", "Australia")
- `Players []Player`: List of team players
- `Score Score`: Current score (runs, wickets, overs)

#### 3. **Score**
Tracks the scoring metrics for a team.

**Key Fields:**
- `Runs int`: Total runs scored
- `Wickets int`: Total wickets fallen
- `Overs float64`: Overs bowled (e.g., 15.3 = 15 overs + 3 balls)

#### 4. **Player**
Represents a cricket player.

**Key Fields:**
- `Name string`: Player name
- `Role string`: "Batsman", "Bowler", "All-rounder"

#### 5. **Observer** (Interface)
Contract that all observers must implement.

```go
type Observer interface {
    Update(match *Match)
}
```

#### 6. **Concrete Observers**

**Scoreboard** - Stadium display style
```
[Stadium Display] India: 150/3 (15.2 overs) 🏏
```

**CommentaryBox** - Commentary style
```
[Commentary] 📢 That's the score! India currently at 150 runs with 3 wickets down!
```

**MobileApp** - Mobile notification style
```
[Mobile App] 📱 Score Update: India 150/3 (15.2)
```

---

## 🎨 Design Pattern: Observer Pattern

### Problem It Solves

**Without Observer Pattern:**
```go
match.addRuns(4)
scoreboard1.display()    // ❌ Manual call
scoreboard2.display()    // ❌ Manual call
mobileApp.display()      // ❌ Manual call
commentaryBox.display()  // ❌ Manual call
```

**Problems:**
- Match must know about every display
- Adding new display requires changing Match code
- Tight coupling between Match and displays
- Easy to forget to update a display

**With Observer Pattern:**
```go
match.addRuns(4)
// ✅ ALL observers automatically notified and updated!
```

**Benefits:**
- Loose coupling: Match doesn't know about specific observers
- Easy to add/remove observers without changing Match
- Each observer can react differently
- Automatic updates - no manual calls needed

---

## 🎯 Observer Pattern Implementation

### 1. Define Observer Interface

```go
type Observer interface {
    Update(match *Match)  // Called when match state changes
}
```

### 2. Subject (Match) Manages Observers

```go
type Match struct {
    // ... other fields
    Observers []Observer
}

func (m *Match) Subscribe(observer Observer) {
    m.Observers = append(m.Observers, observer)
}

func (m *Match) NotifyAll() {
    for _, observer := range m.Observers {
        observer.Update(m)  // Pass match to each observer
    }
}
```

### 3. Concrete Observers Implement Interface

```go
type Scoreboard struct {
    name  string
    match *Match
}

func (s *Scoreboard) Update(match *Match) {
    s.match = match  // Store reference
    s.Display()      // Display in own style
}

func (s *Scoreboard) Display() {
    // Custom display logic
}
```

### 4. Usage

```go
// Create observers
scoreboard := &Scoreboard{name: "Stadium"}
mobileApp := &MobileApp{name: "App"}

// Subscribe to match
match.Subscribe(scoreboard)
match.Subscribe(mobileApp)

// Automatic updates!
match.addRuns(4)  // Both observers notified automatically
```

---

## 📊 UML Diagram

```
┌─────────────────────────────────────┐
│           <<interface>>              │
│            Observer                  │
├─────────────────────────────────────┤
│ + Update(match *Match)              │
└─────────────────────────────────────┘
           ▲         ▲         ▲
           │         │         │
           │         │         │
    ┌──────┴───┐ ┌──┴─────┐ ┌─┴────────┐
    │Scoreboard│ │MobileApp│ │Commentary│
    └──────────┘ └─────────┘ └──────────┘

┌─────────────────────────────────────┐
│             Match                    │
│         (Subject)                    │
├─────────────────────────────────────┤
│ - observers: []Observer             │
│ - teams: [2]Team                    │
│ - teamBatting: *Team                │
├─────────────────────────────────────┤
│ + Subscribe(Observer)               │
│ + NotifyAll()                       │
│ + addRuns(int)                      │
│ + addWicket()                       │
└─────────────────────────────────────┘
          │
          │ notifies
          ▼
    [All Observers]
```

---

## 🚀 How to Run

### Prerequisites
- Go 1.16 or higher

### Steps

```bash
# Save the code as cricbuzz.go
go run cricbuzz.go
```

### Expected Output

```
Match started

[Stadium Display] India: 4/0 (0.1 overs) 🏏
[Commentary] 📢 That's the score! India currently at 4 runs with 0 wickets down!
[Mobile App] 📱 Score Update: India 4/0 (0.1)

[Stadium Display] India: 10/0 (0.2 overs) 🏏
[Commentary] 📢 That's the score! India currently at 10 runs with 0 wickets down!
[Mobile App] 📱 Score Update: India 10/0 (0.2)

[Stadium Display] India: 11/0 (0.3 overs) 🏏
[Commentary] 📢 That's the score! India currently at 11 runs with 0 wickets down!
[Mobile App] 📱 Score Update: India 11/0 (0.3)

India: 11/0 (0.3 overs)

[Stadium Display] India: 11/1 (0.4 overs) 🏏
[Commentary] 📢 That's the score! India currently at 11 runs with 1 wickets down!
[Mobile App] 📱 Score Update: India 11/1 (0.4)

...

Winner: Australia
```

---

## 💡 Key Design Decisions

### 1. **Pass Whole Match vs. Individual Fields**

**Decision:** Pass `*Match` to `Update()` instead of individual values

**Reasoning:**
- Observers can access any match data they need
- Flexible: New observers can use different data without changing interface
- Observers can compare teams, check match status, etc.

**Alternative Considered:**
```go
Update(runs, wickets, overs int)  // ❌ Too rigid
```

### 2. **Store Match Reference in Observers**

**Decision:** Observers store `match *Match` field

**Reasoning:**
- Enables multiple methods to access match data
- Allows observers to implement additional features (run rate calculation, trends, etc.)
- Can display data at any time, not just during Update()

### 3. **Overs as float64**

**Decision:** Use `float64` for overs, not `int`

**Reasoning:**
- Cricket notation: 15.3 overs = 15 overs + 3 balls
- float64 allows: 0.1, 0.2, 0.3, 0.4, 0.5, then 1.0 (next over)
- Simpler than tracking overs and balls separately

**Implementation:**
```go
func (m *Match) addBall() {
    currentOvers := m.TeamBatting.Score.Overs
    balls := int((currentOvers - float64(int(currentOvers))) * 10)
    
    if balls >= 5 {
        // 6 balls = 1 over completed
        m.TeamBatting.Score.Overs = float64(int(currentOvers) + 1)
    } else {
        m.TeamBatting.Score.Overs += 0.1
    }
}
```

### 4. **Pointer to Batting Team**

**Decision:** Use `TeamBatting *Team` (pointer) instead of value

**Reasoning:**
- Modifying batting team's score updates the actual team in the array
- Avoids copying entire Team struct
- Enables easy innings switching by just reassigning pointer

### 5. **Notify on Every Score Change**

**Decision:** Call `NotifyAll()` after every `addRuns()` and `addWicket()`

**Reasoning:**
- Real-time updates (like actual Cricbuzz)
- Observers always have latest data
- Demonstrates Observer pattern's real-time capability

---

## 🔄 How Observer Pattern Works Here

### Event Flow

```
1. User Action
   ↓
   match.addRuns(4)
   
2. Match Updates Internal State
   ↓
   m.TeamBatting.Score.Runs += 4
   m.addBall()
   
3. Match Notifies All Observers
   ↓
   m.NotifyAll()
   
4. Each Observer Reacts
   ↓
   scoreboard.Update(m)  → displays stadium style
   mobileApp.Update(m)   → displays notification style
   commentary.Update(m)  → displays commentary style
```

### Why This Works

1. **Match doesn't know about Scoreboard, MobileApp, or CommentaryBox classes**
2. **Match only knows about Observer interface**
3. **New observers can be added without changing Match code**
4. **Each observer decides HOW to display, not the Match**

---

## 🎓 Learning Outcomes

### Observer Pattern Mastery

After building this system, you should understand:

1. ✅ **When to use Observer Pattern**
   - When one object (subject) needs to notify multiple objects (observers)
   - When objects should be loosely coupled
   - When you want automatic propagation of changes

2. ✅ **How to implement Observer Pattern**
   - Define Observer interface
   - Subject maintains list of observers
   - Subject calls Update() on all observers when state changes
   - Observers implement Update() to react

3. ✅ **Benefits of Observer Pattern**
   - Loose coupling between subject and observers
   - Easy to add/remove observers
   - Supports broadcast communication
   - Observers can react differently to same event

4. ✅ **Trade-offs**
   - More complexity (interfaces, lists of observers)
   - Potential performance impact (notifying many observers)
   - Order of notification is not guaranteed

---

## 🔧 Extensibility

The system can be easily extended:

### Add New Observer Types

```go
type WebsiteDisplay struct {
    name  string
    match *Match
}

func (w *WebsiteDisplay) Update(match *Match) {
    w.match = match
    // Display on website
}

// Add to match
match.Subscribe(&WebsiteDisplay{name: "Website"})
```

### Add Event Types

```go
type EventType string
const (
    EventRuns   EventType = "RUNS"
    EventWicket EventType = "WICKET"
    EventSix    EventType = "SIX"
)

type Observer interface {
    Update(match *Match, event EventType)
}

// Observers can react differently to different events
func (c *CommentaryBox) Update(match *Match, event EventType) {
    switch event {
    case EventSix:
        fmt.Println("🎆 That's a MASSIVE SIX!")
    case EventWicket:
        fmt.Println("💥 WICKET! What a delivery!")
    }
}
```

### Add Unsubscribe

```go
func (m *Match) Unsubscribe(observer Observer) {
    for i, obs := range m.Observers {
        if obs == observer {
            m.Observers = append(m.Observers[:i], m.Observers[i+1:]...)
            break
        }
    }
}
```

### Add Player Statistics

```go
type Player struct {
    Name        string
    Role        string
    RunsScored  int
    BallsFaced  int
    WicketsTaken int
}

// Track which player scored runs
func (m *Match) addRunsForPlayer(player *Player, runs int) {
    player.RunsScored += runs
    m.TeamBatting.Score.Runs += runs
    m.NotifyAll()
}
```

---

## 🧪 Testing Scenarios

### Test Case 1: Multiple Observers
```go
// Create 5 different observers
observers := []Observer{
    &Scoreboard{name: "Stadium 1"},
    &Scoreboard{name: "Stadium 2"},
    &MobileApp{name: "App 1"},
    &MobileApp{name: "App 2"},
    &CommentaryBox{name: "Radio"},
}

for _, obs := range observers {
    match.Subscribe(obs)
}

match.addRuns(6)
// All 5 observers should display the update
```

### Test Case 2: Dynamic Subscription
```go
// Start with 2 observers
match.Subscribe(scoreboard1)
match.Subscribe(mobileApp)

match.addRuns(4)  // 2 observers notified

// Add new observer mid-match
match.Subscribe(commentary)

match.addWicket()  // 3 observers notified
```

### Test Case 3: Innings Switch
```go
// First innings
for i := 0; i < 10; i++ {
    match.addWicket()  // All observers see each wicket
}
// Automatic innings switch, all observers notified

// Second innings
match.addRuns(10)  // All observers show new batting team
```

---

## 📊 Performance Characteristics

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Subscribe | O(1) | O(n) where n = observers |
| NotifyAll | O(n) where n = observers | O(1) |
| addRuns | O(n) due to notify | O(1) |
| addWicket | O(n) due to notify | O(1) |

**Note:** With many observers (100+), notification time increases linearly. Consider:
- Async notifications (goroutines)
- Observer priorities
- Selective notifications

---

## 🎓 Interview Talking Points

When discussing this system in interviews:

1. **Pattern Choice**: "I used Observer Pattern because multiple displays needed automatic updates when match score changed, without tight coupling."

2. **Design Decision**: "I passed the whole Match object to Update() rather than individual fields, giving observers flexibility to access any data they need."

3. **Real-World Application**: "This is exactly how Cricbuzz works - when a scorer updates the match, all connected devices see the update instantly through push notifications."

4. **Scalability**: "For thousands of observers, I'd use goroutines for async notifications, or implement a pub-sub system with message queues."

5. **Alternative Patterns**: "I could have used a message bus or event-driven architecture, but Observer Pattern is simpler and sufficient for this scale."

---

## 🔗 Related Patterns

- **Pub/Sub Pattern**: Similar to Observer, but with message broker in between
- **Mediator Pattern**: Centralizes communication, whereas Observer is point-to-multipoint
- **Event Sourcing**: Stores all events (every run, wicket) for replay and analytics

---

## 📚 Key Learnings

### What Makes This Observer Pattern?

1. ✅ **One-to-Many Relationship**: One Match, many Observers
2. ✅ **Automatic Notification**: Match calls Update() on all observers
3. ✅ **Loose Coupling**: Match only knows Observer interface
4. ✅ **Dynamic Subscription**: Observers added/removed at runtime
5. ✅ **Push Model**: Match pushes data to observers (vs. observers polling)

### When NOT to Use Observer

- ❌ If only one observer (just call it directly)
- ❌ If observers need to control subject (use two-way communication)
- ❌ If notification order matters (Observer doesn't guarantee order)
- ❌ If complex filtering needed (use Event Bus instead)

---

## 🎉 Congratulations!

You've successfully built a Cricbuzz system demonstrating:
- ✅ Observer Pattern implementation
- ✅ Event-driven architecture
- ✅ Real-time updates
- ✅ Loose coupling
- ✅ Clean code organization

---

## 📝 What's Not Implemented (Future Enhancements)

The following were intentionally left out to focus on Observer Pattern:

- ❌ Squad management (tournament squad vs playing eleven)
- ❌ Player statistics (individual runs, wickets)
- ❌ Ball-by-ball detailed tracking
- ❌ Match types (ODI, Test, T20)
- ❌ Tournament management
- ❌ Historical data and queries
- ❌ Advanced rules (no-balls, wides, DRS, etc.)

These features would add complexity but not teach new design patterns. The focus was on mastering Observer Pattern through a practical, working example.

---

**Next System:** BookMyShow (State Pattern) 🎬

---

## 👨‍💻 Author

Built as part of a comprehensive Low-Level Design learning journey, focusing on design patterns and system design principles.

## 📄 License

Educational project - free to use and modify for learning purposes.