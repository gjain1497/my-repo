We are designing LLD (Low-Level Design) systems in Go, with a strong focus on:

SOLID principles

Idiomatic Go design principles

Clean architecture & separation of concerns

Extensibility, maintainability, and testability

Production-grade patterns (interfaces, composition, concurrency, error handling, etc.)

All solutions should reflect high-quality Go engineering practices, demonstrate scalable software architecture, and follow LLD best practices.

(V.Imp) 

Again bro learning strategy would be I start you review I do again you review and so on. Don't ever give me all the code at once, remember this


(V.Imp) From every question we will try to extract out patterns which are common across multiple systems and can be reused. This is a crucial step to build solid understanding

This is the method we will follow for our LLD (Vehicle Rental is just an example and the way I had done LLD for VRS) :-


Start with the overall flow of system (understand requirements) -> maybe draw a class diagram if required. Again this I will do, If I am unable to list down/missing something then we can align later

After this is done then we will follow below steps so that it becomes relatively easy to execute the below steps

Step 1: Define ALL models first ✅

Step 2: Identify Services & Group Models ⭐
        
        2a) List all services needed
        2b) Map which models belong to which service
        2c) Identify dependencies between services (who uses whom)

        ┌─────────────────────────────────────────────────────────────┐
        │ 2b) Models Grouping Example (VRS):                          │
        ├─────────────────┬───────────────────────────────────────────┤
        │ Service         │ Models Owned                              │
        ├─────────────────┼───────────────────────────────────────────┤
        │ VehicleService  │ Vehicle, VehicleType, Location            │
        │ BookingService  │ Booking, BookingStatus                    │
        │ PaymentService  │ Payment, PaymentType, PaymentStatus       │
        │ UserService     │ User, Person                              │
        │ AdminService    │ Admin                                     │
        └─────────────────┴───────────────────────────────────────────┘

        // Same in code form:
        
        // --- VehicleService ---
        type VehicleService struct{}
        type Vehicle struct{}
        type VehicleType string
        type Location struct{}
        
        // --- BookingService ---
        type BookingService struct{}
        type Booking struct{}
        type BookingStatus string
        
        // --- PaymentService ---
        type PaymentService struct{}
        type Payment struct{}
        type PaymentType string
        type PaymentStatus string
        
        // --- UserService ---
        type UserService struct{}
        type User struct{}
        type Person struct{}
        
        // --- AdminService ---
        type AdminService struct{}
        type Admin struct{}

        ┌─────────────────────────────────────────────────────────────┐
        │ 2c) Service Dependencies Example (VRS):                    │
        ├─────────────────┬───────────────────────────────────────────┤
        │ Service         │ Depends On (INTERFACES, not concrete!)   │
        ├─────────────────┼───────────────────────────────────────────┤
        │ VehicleService  │ -                                         │
        │ BookingService  │ PaymentServiceInterface,                  │
        │                 │ VehicleServiceInterface                   │
        │ PaymentService  │ PaymentGateway (external interface)       │
        │ UserService     │ VehicleServiceInterface                   │
        │ AdminService    │ VehicleServiceInterface                   │
        └─────────────────┴───────────────────────────────────────────┘

        // Same in code form:
        
        // --- VehicleService (no dependencies) ---
        type VehicleService struct {
            Vehicles map[string]*Vehicle
            mu       sync.RWMutex
        }
        
        // --- BookingService (depends on 2 interfaces) ---
        type BookingService struct {
            PaymentService PaymentServiceInterface  // ✅ Interface!
            VehicleService VehicleServiceInterface  // ✅ Interface!
            Bookings       map[string]*Booking
            mu             sync.RWMutex
        }
        
        // --- PaymentService (depends on external gateway) ---
        type PaymentService struct {
            PaymentGateway PaymentGateway  // ✅ Interface!
            Payments       map[string]*Payment
            mu             sync.RWMutex
        }
        
        // --- UserService ---
        type UserService struct {
            VehicleService VehicleServiceInterface  // ✅ Interface!
            Users          map[string]*User
        }
        
        // --- AdminService ---
        type AdminService struct {
            VehicleService VehicleServiceInterface  // ✅ Interface!
        }

Step 3: Go service by service
    For each service:
    
    Step 3a: Define Interface FIRST (ALWAYS - for DIP)
             → Only PUBLIC methods (APIs)
             → Internal helpers stay OUT
    
    Step 3b: Define Service struct with dependencies
    
             ┌─────────────────────────────────────────────────────────────┐
             │ ⭐ FACTORY PATTERN DECISION (NEW!)                         │
             ├─────────────────────────────────────────────────────────────┤
             │ Ask: "Does this service need to choose between MULTIPLE     │
             │       implementations of a SUBCOMPONENT based on            │
             │       REQUEST DATA?"                                        │
             │                                                             │
             │ ✅ YES → Use FACTORY in service                            │
             │    - Subcomponent is SAME domain                           │
             │    - Service has DATA to decide (payment.Type, etc.)       │
             │    - Service understands DOMAIN LOGIC                      │
             │    - Decision changes PER OPERATION                        │
             │                                                             │
             │    Example:                                                 │
             │    type PaymentService struct {                            │
             │        GatewayFactory *PaymentGatewayFactory  // ✅ Factory│
             │        ProcessorFactory *ProcessorFactory     // ✅ Factory│
             │    }                                                        │
             │                                                             │
             │ ❌ NO → Use INTERFACE injection only                       │
             │    - Dependency is DIFFERENT domain                        │
             │    - Service doesn't have relevant data                    │
             │    - Decision at handler/caller layer                      │
             │    - Decision per USER/CONFIG (not per operation)          │
             │                                                             │
             │    Example:                                                 │
             │    type BookingService struct {                            │
             │        PaymentService PaymentServiceInterface  // ✅ Interface│
             │        VehicleService VehicleServiceInterface  // ✅ Interface│
             │    }                                                        │
             └─────────────────────────────────────────────────────────────┘
    
             // ❌ WRONG - Concrete types for cross-domain dependencies
             type BookingService struct {
                 PaymentService *PaymentService
                 VehicleService *VehicleService
             }
             
             // ✅ CORRECT - Interface types for cross-domain dependencies
             type BookingService struct {
                 PaymentService PaymentServiceInterface
                 VehicleService VehicleServiceInterface
             }
             
             // ✅ CORRECT - Factory for same-domain subcomponents
             type PaymentService struct {
                 GatewayFactory   *PaymentGatewayFactory
                 ProcessorFactory *PaymentProcessorFactory
             }
    
    Step 3c: (If Factory needed) Define Strategy Interface & Implementations
    
             ┌─────────────────────────────────────────────────────────────┐
             │ Strategy Pattern Setup:                                     │
             ├─────────────────────────────────────────────────────────────┤
             │ 1. Define strategy interface                                │
             │    type PaymentGateway interface {                          │
             │        Charge(payment *Payment) error                       │
             │    }                                                         │
             │                                                             │
             │ 2. Define implementations                                   │
             │    type StripeGateway struct {}                            │
             │    type RazorpayGateway struct {}                          │
             │                                                             │
             │ 3. Define factory (Map-based - DEFAULT)                    │
             │    type PaymentGatewayFactory struct {                     │
             │        Gateways map[GatewayType]PaymentGateway             │
             │    }                                                        │
             │                                                             │
             │    func NewPaymentGatewayFactory() *PaymentGatewayFactory {│
             │        return &PaymentGatewayFactory{                      │
             │            Gateways: map[GatewayType]PaymentGateway{       │
             │                Stripe:   &StripeGateway{},                 │
             │                Razorpay: &RazorpayGateway{},               │
             │            },                                               │
             │        }                                                    │
             │    }                                                        │
             │                                                             │
             │    func (f *Factory) GetGateway(t GatewayType)             │
             │                         (PaymentGateway, error) {          │
             │        gateway, exists := f.Gateways[t]                    │
             │        if !exists {                                        │
             │            return nil, fmt.Errorf("not found: %s", t)     │
             │        }                                                    │
             │        return gateway, nil                                 │
             │    }                                                        │
             │                                                             │
             │ ⚠️ Only use switch-based factory if implementations        │
             │    have per-request state (rare!)                          │
             └─────────────────────────────────────────────────────────────┘
    
    Step 3d: Define APIs/methods (the interface methods)
             → "What operations does this service expose?"
    
    Step 3e: Implement methods
             → Public methods (in interface)
             → Private helpers (lowercase, NOT in interface)
             → Use factory to select strategies when needed
             
             Example:
             func (s *PaymentService) ProcessPayment(payment *Payment) error {
                 // ✅ Service decides gateway based on domain data
                 var gatewayType GatewayType
                 if payment.Currency == "INR" {
                     gatewayType = Razorpay
                 } else {
                     gatewayType = Stripe
                 }
                 
                 gateway, err := s.GatewayFactory.GetGateway(gatewayType)
                 if err != nil {
                     return err
                 }
                 
                 return gateway.Charge(payment)
             }
    
    Step 3f: (Optional) Alternate implementations
             → MockPaymentService for testing
             → Alternative service versions (PaymentServiceV1, V2)

Step 4: Wire it all up with Dependency Injection

        ┌─────────────────────────────────────────────────────────────┐
        │ Wiring Pattern:                                             │
        ├─────────────────────────────────────────────────────────────┤
        │ 1. Create factories (for strategies within services)        │
        │ 2. Create services (inject factories + other services)      │
        │ 3. Wire services together                                   │
        └─────────────────────────────────────────────────────────────┘

        // At initialization
        
        // 1. Create strategy factories (same-domain decisions)
        gatewayFactory := NewPaymentGatewayFactory()
        processorFactory := NewPaymentProcessorFactory()
        
        // 2. Create services with factories
        paymentService := &PaymentService{
            GatewayFactory:   gatewayFactory,   // ✅ Factory for strategies
            ProcessorFactory: processorFactory, // ✅ Factory for strategies
        }
        
        vehicleService := &VehicleService{
            PricingFactory: pricingFactory,     // ✅ Factory for pricing strategies
        }
        
        // 3. Wire services together (cross-domain dependencies)
        bookingService := &BookingService{
            PaymentService: paymentService,  // ✅ Interface injection
            VehicleService: vehicleService,  // ✅ Interface injection
        }

┌─────────────────────────────────────────────────────────────────────┐
│ ⭐ FACTORY PATTERN SUMMARY (Quick Reference)                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ USE FACTORY when:                                                  │
│ ✅ Same domain subcomponent (PaymentGateway in PaymentService)    │
│ ✅ Service has request data to decide (payment.Currency)           │
│ ✅ Service understands domain logic (currency → gateway)           │
│ ✅ Decision per operation (different per payment)                  │
│                                                                     │
│ USE INTERFACE when:                                                │
│ ✅ Different domain dependency (PaymentService in BookingService)  │
│ ✅ Handler has data to decide (user.IsInExperiment)               │
│ ✅ Handler understands routing logic (A/B testing)                 │
│ ✅ Decision per user/config (not per operation)                    │
│                                                                     │
│ FACTORY IMPLEMENTATION:                                            │
│ ✅ Default: Map-based (pre-create instances at startup)            │
│ ⚠️ Rare: Switch-based (only if per-request state)                 │
│                                                                     │
│ PATTERN:                                                           │
│ Service Layer: Decides OWN strategies (factory)                    │
│ Handler Layer: Decides OTHER services (injection)                  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘



List of questions I found to practice

Parking Lot -> done
Cricbuzz -> done
Book My Show -> done
IMS(Inventory Management System)/Ecommerce -> done
Worker Pool 
IMS with Queue (uses Worker Pool concept)




Amazon Questions:

Tic-Tac-Toe game
 I was asked to list down all the entities, and explain how the game works and how it ends. There were many follow-up questions, and I answered them.
At the end, I was asked to draw the class diagram and write the code in any language (I chose C++) for the method to check if the game is over.


System: Digital Wallet (Handle all aspects related to wallet creation, balance management, fund transfers, and transaction logging)

Amazon has various products and we need to generate an id for each item. And each item may have some formats for their id. Design a LLD classes and methods for this

//Cache
Network Request Cache LLD
Multi-level Cache System

//Logger 2 times
Design a Logger Framework
Design Logger system with filtering functionality. Not much time was there, covered import classes and logic around how to add multiple filters and combine the filters.

LLD of stack overflow


 "Design an event forwarding framework where event generated from a system(s) is required to be consumed by another system(s)."


//Design Parking Lot
System should be able to handle different parking ways
System should be able to different pricing models.

Book My show -> focus on booking method (concurrency and no two seats should get booked)


Design Stock Broker Platform - Zerodha, Groww



Amazon Warehouses Team uses a software to sort all the incoming online orders based on their priority and
stores the information in a database.
This information is made available to the packaging system which continuously picks the order
with the highest priority and packages it for delivery.

The OrderSortingService has the following functional requirements:

Support adding new orders continuously coming in from amazon.com.
Support getting the order with the highest priority.
Store the incoming orders.



The LLD question was to design a ATM system.

I was not able to do this properly or at least how the interviewer wanted it. The expectation was to list down the requirements, list down core entities with attributes and behaviour and then write code for one or two features. Proper data type for attributes is important

I fumbled the requirements and the core entities so could not perform in this round properly.



Problem: Design Meeting and room reservation system
Approach: It went like typical LLD, I was thinking how i schedule Meetings in teams and outlook.
1. started with requirement gathering
2. Then Made rough flow
3. Identified Objects
4. coded all the classes and code.

It was not a working code. It was more like UML diagram just in codepad. 
So i made all the classes and talked about all is and has relation ships and then created all variables and methods used. 
At the end of it i explained how a good flow would go. 
He was conversing with me throughout, and at the end he asked me what design pattern more i could have used which I already wanted to talk about if i had more time.



One standard LLD problem (Amazon asks this frequently; easily found in resources)
Amazon Locker


design notification system with focus on building message content for different channels.


  Design a loyality program system for Amazon Fresh shoppers, that rewards customers for their shopping behaviour, manages point allocation and handles tier based benefits through a points wallet. 
   Vairious tiers like: Silver, Gold, Platform tier. and a redemption system
   Parameters of the exact question I don't remember.


Initial Problem Statement: code a way to find available space to  store incoming packages.

Initial requirement gathering took 10-15 mins to understand what really i wanted to code here. (hard part)
Then coding part and follow up questions were done (easy part)

Approach:
1. Firstly discussion went on how would we can categorize different containers based on size. 
2. I was thinking in terms of volume but as some packages might fit volume and violate height/length/breadth req. 
3. So eventually, we came with pre defined sizes of lockers like (S, M, L, Xl, XXL) and user could choose his package would fit in which category. 
4. Thinking input variable for method was the hardest part.
5. Then I quicly coded 3 methods, find Locker, add locker, update locker(put/remove package).
6. then discusion went around which data structures i have used and why. and what can i improve in time and space complexity.



implement unix find command as an api , the api will support findings;
files that have given certain requirements
files with a certain naming patterns focus on 2 uses cases at first
find al files over 5 MB somewhere under directory

Vending Machine
Follow ups on Change management/ stock management/ refill scenario.

Rate Limiter LLD

Design Splitwise
LLD -> Splitwise design and focus on Currency conversion.



Design an online food delivery system like Swiggy/Zomato.
Design a Delivery Partner Assignment System


Focus: Low-Level Design (LLD) for a Seller Experience application.



Question was not direct. It was a problem statement - Given N people and (person-id, txn amt), min txns to settle among people.

I honestly didn't do this question before. So I proposed a greedy way of using 2 Max Heaps (debters, creditors) - always fetch highest debter and highest creditor and push back to maxheap, if still due is pending. Interviewer pointed out that it doesnt work for all testcases but still asked me to write pseudo code.

I later figured that was an NP-Hard problem and could've used a brute force - backtracking approach. O(N!)

//Asked twice in amazon
Consider there are differnt types of alexa devices available. One with audio, one with screen, one with audio and screen. These devices may have a battery or may not. Battery devices will have battery percentage. Both battery and non battery devices can be put charging. The task is to show the battery percentage. Include a show methond and that method should show the current battery percentage if it has a battery. If not just say, battery not available. You should also say whether its currently charging or not. There will four statements to print show method like Charging and battery percentage, charging and no battery, just battery percent and no battery.

Expectation is to write interface-driven code using appropriate design patterns

Asked three times in Amazon
1) Design a file filtering and search system". Did you tried clarifying the functional requirements with interviewer? That should be the first step. It gives you 3-5 core feature that your design should address. Having this conversation enables you to collaborate with the interviewer.

2) Design a Filtering System which can filter based on File Size and File Extension

3) very vague LLD question related to designing a filesystem with functions to print, modify or delete the child directories. only discussed the filesystem traversal, classes and their relations.


Was asked to design a notification service.

The service had multiple clients. And each client can have multiple subscribers subscribe to get notifications.

There are three levels of notification urgency : High(H) | Medium(M) | Low(L)
Each subscriber can have a different notification strategy for each severity level for each client.

Example :
Subscriber 2 can subscribe to amazon shopping following strategy:
H(phone, msg, email)|M(msg,email)|L(email)
Subscriber 2 can subscribe to AWS with following strategy:
H(Phone,email)|M(email)|L(msg)

As you can see each subscriber can customise their notification strategy.

The expectation was to discuss database schema/ Database system to use and write a code to retrieve strategies and send notification.The solution should be extensible, we should be able to add more severity levels and more endpoints (like paging,insta etc).




PubSub Like Kafka -> Done
Notification System 
Car Rental -> Done
snake ladder -> done
tic tac toe -> done
chess -> done
Spotify
Uber
Elevator
Chess
Amazon Locker System
Vending Machine 
ATM Machine -> Done
Stock Broker
File System
Logger
Job Scheduler
Meeting Scheduler
Splitwise
Rate Limiter
Linkedin
Cache
TrueCaller
Airline Management System


Design an extensible solution to implement a search filter in OOD for a directory, matching files by size or name.
Design a library to read a directory and perform operations such as filtering by file type and size constraints.
Design an inventory management system with queuing for incoming requests.
Create the low-level design of a download manager capable of handling multiple downloads.
