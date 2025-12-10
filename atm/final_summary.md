📝 Quick Reference for Interviews
When asked about State Pattern:

"State Pattern allows an object to change behavior when internal state changes. It's typically a frontend concern - managing user flow through different states. Backend services remain stateless. I implemented this in an ATM system where states like IdleState, CardInsertedState, etc. manage the transaction flow, while stateless services handle business logic."

When asked about ATM design:

"ATMs use monolithic architecture primarily for cost savings - about $42M/year for 10,000 ATMs vs separate frontend/backend. The state machine layer acts as the 'frontend' managing user flow, while stateless services act as the 'backend' processing requests. Hardware interfaces provide denomination detection for deposits."

When asked about Go patterns:

"Go uses composition over inheritance. Instead of abstract classes, we use struct embedding for code reuse. For state pattern, I created a BaseATMState with default implementations, then concrete states embed it and override only what's needed. This is Go's idiomatic approach to achieving inheritance-like behavior."