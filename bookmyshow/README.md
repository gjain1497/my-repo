📋 Functional Requirements:
Core Features:

Theater Management

System should support multiple theaters in different cities
Each theater can have multiple screens/halls
Each screen has a seating arrangement with different seat types


Movie & Show Management

Movies can be shown across multiple theaters
Each screen can have multiple shows for different movies at different times
Shows have specific start times and durations


Seat Booking

Users can search for movies by name, language, or theater
Users can view available shows for a movie
Users can select seats from the available seats
Multiple users should be able to book seats concurrently without conflicts
Selected seats should be temporarily blocked for a user (e.g., 10 minutes)


Booking Lifecycle

When a user selects seats, they should be temporarily reserved
User has limited time (10 minutes) to complete payment
If payment is not completed in time, seats should be released
After successful payment, booking is confirmed
Users can cancel confirmed bookings (with refund policies)


Payment Processing

Users can pay for their bookings
System should handle payment success/failure scenarios




📋 Non-Functional Requirements:

Concurrency Handling

The system must handle multiple users trying to book the same seat simultaneously
No double-booking should occur (race condition prevention)


Extensibility

Easy to add new seat types (Recliner, Couple seats, etc.)
Easy to add new booking policies (cancellation rules, pricing)


State Management

Proper state transitions for bookings (Pending → Reserved → Confirmed → Cancelled)
Each state should have clearly defined behaviors




🎯 Key Constraints:

Seat Locking: Once a user proceeds to payment, seats must be locked for 10 minutes
Auto-Expiry: If payment not done in 10 minutes, booking should automatically expire
No Double Booking: Same seat cannot be booked by two users for the same show
Seat Types: At least 2 types of seats with different pricing (Regular, Premium)


🎨 Design Focus (Patterns Expected):
This problem is designed to test:

⭐ State Pattern - Booking lifecycle with different states
Concurrency Control - Thread-safe seat reservation
Object-Oriented Design - Proper entity modeling
Error Handling - Edge cases and failure scenarios


📊 Example Scenario:
1. User searches for "Avengers" in "Mumbai"
2. System shows available theaters and showtimes
3. User selects: PVR Juhu, 7:00 PM show, Screen 3
4. System shows seat layout with availability
5. User selects seats: C5, C6 (Premium seats)
6. System reserves seats for 10 minutes
7. User proceeds to payment
8. Payment successful
9. Booking confirmed, tickets generated
```

---

## 🚫 Out of Scope (Don't Implement):

- User authentication/registration
- Movie recommendations
- Reviews and ratings
- Food/beverage ordering
- Actual payment gateway integration (just simulate success/failure)
- Mobile app UI/API design

---

## 🎯 Your Task:

**Step 1: Identify Core Entities**

Think about and list:
- What are the main "things" (nouns) in this system?
- What are their relationships?

**Write down 6-8 core entities you think are needed.**

Example format:
```
1. Theater - Represents a cinema
2. ???
3. ???
Take your time, think it through, and share your entity list!
I'll review it and we'll refine before moving to implementation! 💪RetryClaude can make mistakes. Please double-check responses.