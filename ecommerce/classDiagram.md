┌─────────────────────────────────────────────────────────────────────────────────┐
│                           E-COMMERCE SYSTEM CLASS DIAGRAM                        │
└─────────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════════
                                 DOMAIN MODELS
═══════════════════════════════════════════════════════════════════════════════════

┌──────────────────────┐
│       User           │
├──────────────────────┤
│ - Id: string         │
│ - Name: string       │
│ - Email: string      │
│ - Phone: string      │
│ - Addresses: []Address│
└──────────────────────┘
         │ 1
         │ has
         │ *
         ▼
┌──────────────────────┐
│      Address         │
├──────────────────────┤
│ - Id: string         │
│ - UserId: string     │────┐
│ - Street: string     │    │ belongs to
│ - City: string       │    │
│ - State: string      │    │
│ - ZipCode: string    │    │
│ - Country: string    │    │
│ - IsDefault: bool    │    │
└──────────────────────┘    │
                            │
         ┌──────────────────┘
         │
         │
┌──────────────────────┐
│      Product         │
├──────────────────────┤
│ - Id: string         │
│ - Name: string       │
│ - Price: float64     │
│ - Description: string│
│ - Category: string   │
└──────────────────────┘
         △
         │ refers to
         │
┌──────────────────────┐        ┌──────────────────────┐
│     Inventory        │        │       Cart           │
├──────────────────────┤        ├──────────────────────┤
│ - Stock: map[string]int│      │ - Id: string         │
│ - mu: sync.RWMutex   │        │ - UserId: string     │────┐
└──────────────────────┘        │ - Items: map[string]int│  │ belongs to
         △                      └──────────────────────┘    │
         │                               △                  │
         │ manages                       │                  │
         │                               │ has              │
         │                               │                  │
┌──────────────────────┐                 │                  │
│    OrderItem         │                 │                  │
├──────────────────────┤                 │                  │
│ - ProductId: string  │─────────────────┘                  │
│ - ProductName: string│                                    │
│ - Quantity: int      │                                    │
│ - PriceAtOrder: float64│                                  │
└──────────────────────┘                                    │
         △                                                  │
         │ *                                                │
         │ contains                                         │
         │                                                  │
         │ 1                                                │
┌──────────────────────┐                                    │
│       Order          │                                    │
├──────────────────────┤                                    │
│ - Id: string         │                                    │
│ - Status: OrderStatus│                                    │
│ - UserId: string     │────────────────────────────────────┘
│ - OrderItems: []OrderItem│
│ - TotalAmount: float64│
│ - PaymentId: string  │──────┐
│ - CreatedAt: time.Time│     │
│ - ShippingAddressId: string│ │
└──────────────────────┘     │
                             │ references
                             │
                             ▼
┌──────────────────────┐
│      Payment         │
├──────────────────────┤
│ - Id: string         │
│ - OrderId: string    │
│ - Type: PaymentType  │
│ - Status: PaymentStatus│
│ - Amount: float64    │
│ - TransactionId: string│
│ - Timestamp: time.Time│
└──────────────────────┘

┌──────────────────────┐
│      Shipment        │
├──────────────────────┤
│ - Id: string         │
│ - OrderId: string    │──────┐
│ - Provider: ShippingProviderType│
│ - Status: ShipmentStatus│    │ references
│ - EstimatedDelivery: time.Time│
│ - TrackingNumber: string│    │
│ - ActualDelivery: *time.Time│
└──────────────────────┘      │
                              │
         ┌────────────────────┘
         │
         │
         ▼
    (links to Order)

═══════════════════════════════════════════════════════════════════════════════════
                                    ENUMS
═══════════════════════════════════════════════════════════════════════════════════

┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────────┐
│   OrderStatus        │  │   PaymentStatus      │  │   PaymentType        │
├──────────────────────┤  ├──────────────────────┤  ├──────────────────────┤
│ + PENDING            │  │ + SUCCESS            │  │ + CREDITCARD         │
│ + CONFIRMED          │  │ + FAILED             │  │ + DEBITCARD          │
│ + SHIPPED            │  └──────────────────────┘  │ + UPI                │
│ + DELIVERED          │                            └──────────────────────┘
│ + CANCELLED          │
└──────────────────────┘

┌──────────────────────┐  ┌──────────────────────┐
│  ShipmentStatus      │  │ShippingProviderType  │
├──────────────────────┤  ├──────────────────────┤
│ + PENDING            │  │ + FEDEX              │
│ + IN_TRANSIT         │  │ + UPS                │
│ + DELIVERED          │  │ + DHL                │
└──────────────────────┘  └──────────────────────┘

═══════════════════════════════════════════════════════════════════════════════════
                                  INTERFACES
═══════════════════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────┐
│              <<interface>>                               │
│           PaymentGateway                                 │
├─────────────────────────────────────────────────────────┤
│ + ProcessPayment(orderId, amount, type): (*Payment, error)│
│ + RefundPayment(paymentId, amount): error               │
└─────────────────────────────────────────────────────────┘
                        △
                        │ implements
         ┌──────────────┼──────────────┐
         │              │              │
┌────────────────┐ ┌────────────────┐ ┌────────────────┐
│ StripeGateway  │ │RazorpayGateway │ │MockPaymentGateway│
├────────────────┤ ├────────────────┤ ├────────────────┤
│ - apiKey: string│ │-keyId: string  │ │                │
└────────────────┘ │-keySecret:string│ └────────────────┘
                   └────────────────┘

┌─────────────────────────────────────────────────────────┐
│              <<interface>>                               │
│          ShippingProvider                                │
├─────────────────────────────────────────────────────────┤
│ + ShipOrder(orderId, address): (*Shipment, error)       │
│ + TrackShipment(trackingNumber): (ShipmentStatus, error)│
└─────────────────────────────────────────────────────────┘
                        △
                        │ implements
                 ┌──────┴──────┐
                 │             │
        ┌────────────────┐ ┌────────────────┐
        │     FedX       │ │   Delhivery    │
        ├────────────────┤ ├────────────────┤
        │                │ │                │
        └────────────────┘ └────────────────┘

═══════════════════════════════════════════════════════════════════════════════════
                                   SERVICES
═══════════════════════════════════════════════════════════════════════════════════

┌──────────────────────────────────────────────────────┐
│              UserService                              │
├──────────────────────────────────────────────────────┤
│ - Users: map[string]*User                            │
├──────────────────────────────────────────────────────┤
│ + AddUser()                                          │
│ + UpdateUser()                                       │
│ + DeleteUser()                                       │
│ + AddAddress()                                       │
│ + UpdateAddress()                                    │
│ + GetAddress(userId, addressId): (*Address, error)  │
└──────────────────────────────────────────────────────┘
                        │ manages
                        ▼
                    (User model)

┌──────────────────────────────────────────────────────┐
│            ProductService                             │
├──────────────────────────────────────────────────────┤
│ - products: map[string]*Product                      │
├──────────────────────────────────────────────────────┤
│ + AddProduct()                                       │
│ + UpdateProductPrice(productId, price): error       │
│ + SearchProducts(keyword): []Product                │
│ + GetProduct(productId): (*Product, error)          │
└──────────────────────────────────────────────────────┘
                        │ manages
                        ▼
                  (Product model)

┌──────────────────────────────────────────────────────┐
│          InventoryService                             │
├──────────────────────────────────────────────────────┤
│ - inventory: *Inventory                              │
├──────────────────────────────────────────────────────┤
│ + AddStock(productId, quantity): error              │
│ + RemoveStock(productId, quantity): error           │
│ + CheckStock(productId): (int, error)               │
└──────────────────────────────────────────────────────┘
                        │ manages
                        ▼
                  (Inventory model)

┌──────────────────────────────────────────────────────┐
│             CartService                               │
├──────────────────────────────────────────────────────┤
│ - carts: map[string]*Cart                            │
│ - productService: *ProductService        ────────────┤───┐
│ - inventoryService: *InventoryService    ────────────┤───┼─┐
├──────────────────────────────────────────────────────┤   │ │
│ + AddToCart(userId, itemId, qty): error             │   │ │
│ + GetCart(userId): (*Cart, error)                   │   │ │
│ + CalculateTotal(userId): (float64, error)          │   │ │
└──────────────────────────────────────────────────────┘   │ │
                                                           │ │
                      ┌────────────────────────────────────┘ │
                      │                ┌─────────────────────┘
                      ▼                ▼
              (uses ProductService, InventoryService)

┌──────────────────────────────────────────────────────┐
│            PaymentService                             │
├──────────────────────────────────────────────────────┤
│ - gateway: PaymentGateway                ────────────┤───┐
├──────────────────────────────────────────────────────┤   │
│ + ProcessPayment(orderId, amount, type): (*Payment, error)│
│ + RefundPayment(paymentId, amount): error           │   │
└──────────────────────────────────────────────────────┘   │
                                                           │
                                  ┌────────────────────────┘
                                  │ HAS-A
                                  ▼
                        (PaymentGateway interface)

┌──────────────────────────────────────────────────────┐
│           ShippingService                             │
├──────────────────────────────────────────────────────┤
│ - provider: ShippingProvider             ────────────┤───┐
├──────────────────────────────────────────────────────┤   │
│ + ShipOrder(orderId, address): (*Shipment, error)   │   │
└──────────────────────────────────────────────────────┘   │
                                                           │
                                  ┌────────────────────────┘
                                  │ HAS-A
                                  ▼
                        (ShippingProvider interface)

┌────────────────────────────────────────────────────────────────┐
│                     OrderService                                │
├────────────────────────────────────────────────────────────────┤
│ - productService: *ProductService           ───────────────────┤──┐
│ - userService: *UserService                 ───────────────────┤──┼─┐
│ - paymentService: PaymentGateway            ───────────────────┤──┼─┼─┐
│ - inventoryService: *InventoryService       ───────────────────┤──┼─┼─┼─┐
│ - cartService: *CartService                 ───────────────────┤──┼─┼─┼─┼─┐
│ - shippingService: *ShippingService         ───────────────────┤──┼─┼─┼─┼─┼─┐
│ - orders: map[string]*Order                                    │  │ │ │ │ │ │
├────────────────────────────────────────────────────────────────┤  │ │ │ │ │ │
│ + PlaceOrder(userId, cartId, addressId, type): (*Order, error)│  │ │ │ │ │ │
│ + ViewOrderDetails(orderId): (*Order, error)                  │  │ │ │ │ │ │
│ + CancelOrder(orderId): error                                 │  │ │ │ │ │ │
│ + GetOrderHistory(userId): ([]Order, error)                   │  │ │ │ │ │ │
│ + TrackOrder(orderId): (*Shipment, error)                     │  │ │ │ │ │ │
└────────────────────────────────────────────────────────────────┘  │ │ │ │ │ │
                                                                    │ │ │ │ │ │
        ┌───────────────────────────────────────────────────────────┘ │ │ │ │ │
        │           ┌───────────────────────────────────────────────────┘ │ │ │ │
        │           │         ┌─────────────────────────────────────────────┘ │ │ │
        │           │         │       ┌───────────────────────────────────────────┘ │
        │           │         │       │         ┌─────────────────────────────────────┘
        │           │         │       │         │       ┌───────────────────────────────┘
        ▼           ▼         ▼       ▼         ▼       ▼
    (Uses all services to orchestrate order placement)

═══════════════════════════════════════════════════════════════════════════════════
                              RELATIONSHIPS LEGEND
═══════════════════════════════════════════════════════════════════════════════════

    │
    │  Association (uses/has)
    ▼

    △
    │  Inheritance/Implementation
    │

   ────  Dependency

    *   One-to-Many
    1   One-to-One