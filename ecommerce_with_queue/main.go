package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// UserService

type UserService struct {
	//list of all the user so that we can perform the crud operations
	//Users []User not very optimised, so instead take map
	Users map[string]*User //(user_id->User) //so now it becomes easy to do crud ops
}

// admin level function, can check whether he is admin through JWT
func (s *UserService) AddUser() {

}

func (s *UserService) UpdateUser() {

}

func (s *UserService) DeleteUser() {

}

func (s *UserService) AddAddress() {

}

func (s *UserService) UpdateAddress() {

}

func (s *UserService) GetAddress(userId string, addressId string) (*Address, error) {
	user := s.Users[userId]
	for _, address := range user.Addresses {
		if address.Id == addressId {
			return &address, nil
		}
	}
	return nil, errors.New("address not found")
}

//bla bla

type User struct {
	Id        string
	Name      string
	Email     string
	Phone     string
	Addresses []Address
}

type Address struct {
	Id string
	//UserId string not needed because user already has list of addresses UserId string
	//no just kidding needed because lets say we have to verifyu that this address belongs to this user for eg.

	/*
			func (s *UserService) UpdateAddress(userId, addressId string, updates Address) error {
		    address := getAddress(addressId)

		    // ❓ Security check: Does this address belong to this user?
		    if address.UserId != userId {  // ✅ Need UserId field!
		        return errors.New("cannot modify another user's address")
		    }

		    // Update address...
		}
	*/
	UserId    string
	Street    string
	City      string
	State     string
	ZipCode   string
	Country   string
	IsDefault bool // Is this the default shipping address?
}

//ProductService
//Why didn't include product in Inventory Service
// - **Product** = Catalog information (name, price, description, category)
// - **Inventory** = Stock tracking (how many available)

// These are **two different concerns**:
// - Product catalog can exist even with 0 stock
// - Inventory tracks quantities for existing products

type ProductService struct {
	products map[string]*Product //(product_id, Product)
}

// add product to catalog with details
// func (s *ProductService) AddProduct() //confused here as we are already taking care of this in inventory

// Update product price
func (s *ProductService) UpdateProductPrice(productId string, newPrice float64) error {
	// Update price in catalog
	return nil
}

// Search products
func (s *ProductService) SearchProducts(keyword string) []Product {
	// Search by name, category
	return nil
}

// getProduct
func (s *ProductService) GetProduct(productId string) (*Product, error) {
	product := s.products[productId] // Has Price field
	return product, nil
}

type Product struct {
	Id          string
	Name        string
	Price       float64
	Description string
	Category    string
}

/*
### The Difference:

| ProductService | InventoryService |
|----------------|------------------|
| Manages product DETAILS | Manages product QUANTITIES |
| Add new product to catalog | Add stock for existing product |
| Update price, description | Update quantity |
| Search products | Check availability |
| Product info (what) | Stock levels (how many) |

**Example Flow:**
```
1. Admin: AddProduct("Laptop", $1000, "Gaming laptop", "Electronics")
   → ProductService creates product in catalog

2. Warehouse: AddStock("laptop-id", 50)
   → InventoryService adds 50 units to inventory

3. Customer: View product page
   → ProductService shows details ($1000, description)
   → InventoryService shows "50 in stock"

4. Customer: Place order (quantity: 2)
   → InventoryService removes 2 from stock (now 48)
   → Product details unchanged
*/

// InventoryService

// type InventoryService struct {
// 	inventory *Inventory //single inventory not list of inventories here,
// 	// Claude hope my understanding is correct here
// }

// func (s *InventoryService) AddStock(productid string, quantity int) error {
// 	s.inventory.mu.Lock()
// 	defer s.inventory.mu.Unlock()
// 	s.inventory.Stock[productid] += quantity
// 	return nil
// }

// func (s *InventoryService) RemoveStock(productid string, quantity int) error {
// 	s.inventory.mu.Lock()
// 	defer s.inventory.mu.Unlock()
// 	currQuantity := s.inventory.Stock[productid]
// 	if currQuantity > 0 {
// 		s.inventory.Stock[productid] -= quantity
// 		return nil
// 	}
// 	return errors.New("Quantity is not enough")
// }

// // Check stock availability
// func (s *InventoryService) CheckStock(productId string) (int, error) {
// 	s.inventory.mu.Lock()
// 	defer s.inventory.mu.Unlock()
// 	quantity := s.inventory.Stock[productId]
// 	return quantity, nil
// }

type Job struct {
	ID       int
	Task     Task
	Priority int // optional (for later advanced features)
	Attempts int // how many times job was tried
}

type Result struct {
	JobID  int
	Err    error
	Output interface{}
}
type WorkerPool struct {
	NoOfWorkers    int
	JobsChannel    chan Job
	ResultsChannel chan Result
	wg             sync.WaitGroup
}

func NewWorkerPool(numWorkers int, jobQueuesize int) *WorkerPool {
	return &WorkerPool{
		NoOfWorkers:    numWorkers,
		JobsChannel:    make(chan Job, jobQueuesize),
		ResultsChannel: make(chan Result, jobQueuesize),
	}
}

func (w *WorkerPool) Start() {
	for i := 1; i <= w.NoOfWorkers; i++ {
		go w.worker(i) // Each worker runs in its own goroutine
	}
}
func (w *WorkerPool) worker(workerID int) {
	// Loop forever, processing jobs
	for job := range w.JobsChannel { // Blocks until job arrives
		//Execute the task
		output, err := job.Task.Execute()

		//Always send result(even if error occured, so
		//that it does not block other workers)
		w.ResultsChannel <- Result{
			JobID:  job.ID,
			Err:    err,
			Output: output,
		}
		w.wg.Done()

		// Optional: log
		if err != nil {
			fmt.Printf("Worker %d: Job %d failed: %v\n", workerID, job.ID, err)
		} else {
			fmt.Printf("Worker %d: Job %d completed\n", workerID, job.ID)
		}
	}
}

// Submit job to queue
func (w *WorkerPool) SubmitJob(job Job) error {
	//just send job  to channel
	select {
	case w.JobsChannel <- job:
		w.wg.Add(1)
		return nil
	default:
		return fmt.Errorf("job queue is full or pool is shut down")
	}
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}

// Shutdown gracefully
func (wp *WorkerPool) Shutdown() {
	close(wp.JobsChannel)
}

type Task interface {
	Execute() (interface{}, error)
}

type CheckStockTask struct {
	ProductId string
	Inventory *Inventory
}

// but what if the return types of each function would have been different.
// How we would have handled that case?
func (r *CheckStockTask) Execute() (interface{}, error) {
	r.Inventory.mu.Lock()
	defer r.Inventory.mu.Unlock()
	quantity, exists := r.Inventory.Stock[r.ProductId]
	if !exists {
		return 0, nil // Or return error if product doesn't exist
	}
	return quantity, nil // ← Return nil error if successful!
}

type AddStockTask struct {
	ProductId string
	Quantity  int
	Inventory *Inventory
}

func (r *AddStockTask) Execute() (interface{}, error) {
	r.Inventory.mu.Lock()
	defer r.Inventory.mu.Unlock()
	r.Inventory.Stock[r.ProductId] += r.Quantity
	return nil, nil
}

type RemoveStockTask struct {
	ProductId string
	Quantity  int
	Inventory *Inventory
}

func (r *RemoveStockTask) Execute() (interface{}, error) {
	r.Inventory.mu.Lock()
	defer r.Inventory.mu.Unlock()

	currentStock := r.Inventory.Stock[r.ProductId]
	if currentStock < r.Quantity { // ← Check if enough stock exists!
		return nil, fmt.Errorf("insufficient stock: have %d, need %d", currentStock, r.Quantity)
	}

	r.Inventory.Stock[r.ProductId] -= r.Quantity
	return nil, nil
}

type InventoryManager struct {
	inventory    *Inventory
	WorkerPool   *WorkerPool
	jobIDCounter int
	mu           sync.Mutex
}

func NewInventoryManager(numWorkers, queueSize int) *InventoryManager {
	//initalise worker pool
	workerPool := NewWorkerPool(numWorkers, queueSize)
	workerPool.Start()
	// create inventory manager
	return &InventoryManager{
		inventory: &Inventory{
			Stock: make(map[string]int),
		},
		WorkerPool: workerPool,
	}
}
func (im *InventoryManager) generateJobID() int {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.jobIDCounter++
	return im.jobIDCounter
}

func (im *InventoryManager) AddStock(productId string, quantity int) error {
	//submit job
	job := Job{
		ID: im.generateJobID(),
		Task: &AddStockTask{
			ProductId: productId,
			Quantity:  quantity,
			Inventory: im.inventory,
		},
	}

	return im.WorkerPool.SubmitJob(job)
}

func (im *InventoryManager) RemoveStock(productId string, quantity int) error {
	// submit job
	job := Job{
		ID: im.generateJobID(),
		Task: &RemoveStockTask{
			ProductId: productId,
			Quantity:  quantity,
			Inventory: im.inventory,
		},
	}

	return im.WorkerPool.SubmitJob(job)
}

// check stock has to return the current quantity of a productId
func (im *InventoryManager) CheckStock(productId string) (int, error) {
	// submit job
	job := Job{
		ID: im.generateJobID(),
		Task: &CheckStockTask{
			ProductId: productId,
			Inventory: im.inventory,
		},
	}

	err := im.WorkerPool.SubmitJob(job)
	if err != nil {
		return 0, err
	}

	//Wait for result because we want to return the quantity here itslef
	result := <-im.WorkerPool.ResultsChannel
	if result.Err != nil {
		return 0, result.Err
	}

	return result.Output.(int), nil
}

type Inventory struct {
	Stock map[string]int //(product_id, quantity)
	mu    sync.RWMutex
}

// CartService
type CartService struct {
	carts            map[string]*Cart //maintains list of carts for different users (user_id, Cart object)
	productService   *ProductService
	inventoryService *InventoryManager
}

func (s *CartService) AddToCart(userId, itemId string, quantity int) error {
	cart, exists := s.carts[userId]
	if !exists {
		cart = &Cart{
			Id:     "12232",
			UserId: userId,
			Items:  make(map[string]int),
		}
		s.carts[userId] = cart
	}

	cart.Items[itemId] += quantity
	return nil
}

func (s *CartService) GetCart(userId string) (*Cart, error) {
	cart, ok := s.carts[userId]
	if !ok {
		return nil, fmt.Errorf("Cart does not exist")
	}
	return cart, nil
}

func (s *CartService) CalculateTotal(userId string) (float64, error) {
	cart, exist := s.carts[userId]
	var total float64 = 0
	if !exist {
		return 0, fmt.Errorf("Cart does not exisit")
	}
	for itemId, quantity := range cart.Items {
		//get realtime price of this itemId/productId from catalog we can get
		product, err := s.productService.GetProduct(itemId)
		if err != nil {
			return 0, errors.Join(errors.New("CalculateTotal Error"), err)
		}
		total += product.Price * float64(quantity)
	}
	return total, nil
}

type Cart struct {
	Id     string
	UserId string
	//Price float64 -> no price should not be stored in the cart because price is dynamic,
	//if we store
	//Price:     1000,  // Price at time of adding
	// 2 hours later (before checkout), price changes to $800 (sale!)
	// But cart still shows $1000!
	// User sees: "Total: $1000"
	// But actual checkout price: $800
	// User confused: "Why did price change?"
	//Total  float64  // ❌ Remove - should be calculated dynamically, not stored
	Items map[string]int //(product_id, quantity), (product-> quantity does not make sense, because we should not
	//store entire object as a key in map)
}

//OrderService

type OrderService struct {
	productService   *ProductService
	userService      *UserService
	paymentService   PaymentGateway //interface for multiple payment types
	inventoryService *InventoryManager
	cartService      *CartService
	shippingService  *ShippingService
	orders           map[string]*Order //[order_id, Order object] //in real life it will be repo which will store in db
}

func (s *OrderService) PlaceOrder(userId string, cartId string, addressId string, paymentType PaymentType) (*Order, error) {
	// ==================== WHY THIS DESIGN? ====================
	//
	// PROBLEM: How to handle cart → order conversion?
	//
	// ❌ BAD APPROACH 1: Store cart reference in Order
	//    type Order struct {
	//        CartId string  // ← Just reference the cart
	//    }
	//    Problem: If we keep cart and just reference it:
	//    - Cart prices change over time (sales, discounts)
	//    - Order would show CURRENT prices, not purchase-time prices
	//    - User buys at $100, price becomes $80, order shows $80 (WRONG!)
	//    - Legal/audit issue: What price did customer actually pay?
	//
	// ❌ BAD APPROACH 2: Multiple carts per user (one per order)
	//    type CartService struct {
	//        carts map[string][]Cart  // userId → multiple carts
	//    }
	//    Problems:
	//    - Complexity: How to know which cart is "active"?
	//    - Storage: Old carts pile up, waste memory
	//    - If we store price in cart at creation time:
	//      * Cart created: Product price $100
	//      * 2 hours later (before checkout): Price becomes $80
	//      * Cart still shows $100 (outdated!)
	//      * User confused: "Why is cart showing $100 when product is $80?"
	//
	// ✅ CORRECT APPROACH: Snapshot pattern (what we're doing)
	//    1. Cart stores ONLY product IDs + quantities (no prices)
	//    2. At checkout time, create OrderItems with CURRENT prices
	//    3. These become permanent snapshot of purchase
	//    4. Clear cart after successful order (reuse for next order)
	//
	// Benefits of our approach:
	//    ✅ Cart always shows real-time prices (via CalculateTotal)
	//    ✅ OrderItems capture exact purchase-time prices
	//    ✅ One cart per user (simple, efficient)
	//    ✅ Historical orders unaffected by future price changes
	//
	// Real-world example (Amazon):
	//    - Your cart shows TODAY's prices (dynamic)
	//    - Once ordered, "Your Orders" shows PURCHASED price (snapshot)
	//    - If product goes on sale next week, your old order price unchanged
	//
	// ==================== IMPLEMENTATION ====================

	// STEP 1: Get and validate cart
	cart, err := s.cartService.GetCart(userId)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Prevent ordering with empty cart
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	items := cart.Items
	orderId := "1122133212" // TODO: Use UUID generator in production

	// STEP 2: Create order items (PRICE SNAPSHOT happens here!)
	// This is the CRITICAL step where we capture current prices
	var orderItems []*OrderItem
	for itemId, quantity := range items {
		// Get current product details from catalog
		product, err := s.productService.GetProduct(itemId)
		if err != nil {
			return nil, errors.Join(errors.New("product not found"), err)
		}

		// IMPORTANT: Validate stock BEFORE accepting order
		// Prevents overselling: If 5 in stock, can't order 10
		stock, err := s.inventoryService.CheckStock(itemId)
		if err != nil || stock < quantity {
			return nil, fmt.Errorf("insufficient stock for %s (available: %d, requested: %d)",
				product.Name, stock, quantity)
		}

		// CREATE SNAPSHOT: Store current price in OrderItem
		// Once stored, this price NEVER changes, even if product price changes
		// Example: Buy laptop at $1000 today
		//          Tomorrow laptop price drops to $800
		//          Your order still shows $1000 (what you actually paid)
		orderItem := &OrderItem{
			ProductId:    product.Id,
			ProductName:  product.Name,
			PriceAtOrder: product.Price, // ← SNAPSHOT! Frozen in time
			Quantity:     quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	// STEP 3: Calculate total (using current prices)
	// This matches what user saw at checkout
	amount, err := s.cartService.CalculateTotal(userId)
	if err != nil {
		return nil, errors.Join(errors.New("cart total calculation failed"), err)
	}

	// STEP 4: Process payment
	// CRITICAL: Payment must succeed before we:
	//           - Reduce inventory
	//           - Create shipment
	//           - Finalize order
	payment, err := s.paymentService.ProcessPayment(orderId, amount, paymentType)
	if err != nil || payment.Status != Success {
		// Payment failed - don't proceed!
		// Inventory not touched, cart still has items
		return nil, errors.Join(fmt.Errorf("payment failure"), err)
	}

	// ==================== PAYMENT SUCCEEDED - POINT OF NO RETURN ====================
	// From here on, we're committed to fulfilling this order

	// STEP 5: Reduce inventory
	// Only AFTER payment succeeds, so we don't hold inventory for unpaid orders
	// Example: User has 2 laptops in cart
	//          Payment succeeds
	//          Inventory: 50 → 48 (reduce by 2)
	for _, orderItem := range orderItems {
		// err := s.inventoryService.RemoveStock(orderItem.ProductId, orderItem.Quantity)
		if err != nil {
			// TODO: ROLLBACK needed here!
			// Payment succeeded but inventory update failed
			// Should refund payment (call s.paymentService.RefundPayment)
			return nil, fmt.Errorf("failed to update inventory for %s", orderItem.ProductName)
		}
	}

	// STEP 6: Clear cart
	// User's cart is now empty, ready for next shopping session
	// Why clear?
	// - Items are now in "Your Orders", not in cart
	// - User can start fresh cart for next purchase
	// - One cart per user = simple and efficient
	delete(s.cartService.carts, userId)

	// STEP 7: Get shipping address
	// Validate user owns this address (security check in GetAddress method)
	address, err := s.userService.GetAddress(userId, addressId)
	if err != nil {
		return nil, err
	}

	// STEP 8: Create shipment
	// Note: Shipment only created AFTER payment succeeds
	// Rationale: Don't want to reserve shipping resources for unpaid orders
	_, err = s.shippingService.ShipOrder(orderId, address)
	if err != nil {
		return nil, errors.New("failed to create shipment")
	}

	// STEP 9: Create order with all details
	order := &Order{
		Id:                orderId,
		Status:            Confirmed, // Order confirmed after successful payment
		UserId:            userId,
		OrderItems:        orderItems, // Snapshot with frozen prices
		TotalAmount:       amount,     // Total amount paid
		PaymentId:         payment.Id, // Link to payment record
		CreatedAt:         time.Now(), // Timestamp of order creation
		ShippingAddressId: addressId,  // Address for delivery
	}

	// STEP 10: Store order
	// In production: This would be database INSERT
	// Here: Storing in map (simulating persistence)
	s.orders[orderId] = order

	// SUCCESS! Return the created order
	// User can now:
	// - View order details
	// - Track shipment
	// - See order in history
	return order, nil

	// ==================== SUMMARY OF KEY DECISIONS ====================
	//
	// 1. WHY OrderItems instead of Cart reference?
	//    → Need price snapshot at purchase time
	//
	// 2. WHY clear cart after order?
	//    → One cart per user, reuse for next order
	//    → Items moved to "Your Orders", not needed in cart
	//
	// 3. WHY validate stock before payment?
	//    → Prevent user paying for unavailable items
	//
	// 4. WHY reduce inventory AFTER payment?
	//    → Don't hold inventory for unpaid orders
	//
	// 5. WHY create shipment AFTER payment?
	//    → Don't reserve shipping for unpaid orders
	//
	// 6. WHY OrderStatus = Confirmed (not Pending)?
	//    → Payment succeeded = Order confirmed
	//    → Pending would mean payment not yet processed
	//
	// ==================== ORDER LIFECYCLE ====================
	//
	// 1. User adds items to cart (Cart service)
	// 2. User clicks "Place Order"
	// 3. → PENDING: Validating...
	// 4. → CONFIRMED: Payment successful (we are here!)
	// 5. → SHIPPED: Package dispatched
	// 6. → DELIVERED: Customer received
	// 7. (Optional) CANCELLED: User/admin cancelled
	//
	// ==================== DATA FLOW ====================
	//
	// Cart (dynamic prices)
	//   ↓
	// OrderItems (frozen snapshot)
	//   ↓
	// Payment (charge customer)
	//   ↓
	// Inventory (reduce stock)
	//   ↓
	// Shipment (create shipping)
	//   ↓
	// Order (store permanently)
	//
	// ==================== END ====================
}

func (s *OrderService) ViewOrderDetails(orderId string) (*Order, error) {
	//will define this later
	return nil, nil
}

func (s *OrderService) CancelOrder(orderId string) error {
	//will define this later
	return nil
}

func (s *OrderService) GetOrderHistory(userId string) ([]Order, error) {
	return nil, nil
}

func (s *OrderService) TrackOrder(orderId string) (*Shipment, error) {
	return nil, nil
}

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Confirmed OrderStatus = "CONFIRMED"
	Shipped   OrderStatus = "SHIPPED"
	Delivered OrderStatus = "DELIVERED"
	Cancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	// OrderItems string
	ProductId    string
	ProductName  string
	Quantity     int
	PriceAtOrder float64
}

type Order struct {
	// CartId string //already discussed why not needed
	Id     string
	Status OrderStatus
	UserId string
	//OrderItems map[string]OrderItems//(order-item-id -> OrderItems) //this will be more optimised
	// why map should not be used, explained in readme
	OrderItems        []*OrderItem
	TotalAmount       float64
	PaymentId         string
	CreatedAt         time.Time
	ShippingAddressId string
}

// ShipmentService
type ShippingService struct {
	provider ShippingProvider
}

func (s *ShippingService) ShipOrder(orderId string, address *Address) (*Shipment, error) {
	//some business logic for shipping at service level
	//based on fedX or Delhivery
	s.provider.ShipOrder(orderId, address)
	// will implement later
	return nil, nil
}

type ShippingProvider interface {
	ShipOrder(orderId string, address *Address) (*Shipment, error)
	TrackShipment(trackingNumber string) (ShipmentStatus, error)
}

type FedX struct {
}

func (f *FedX) ShipOrder(orderId string, address *Address) (*Shipment, error) {
	//will implement later
	return nil, nil
}

func (f *FedX) TrackShipment(trackingNumber string) (ShipmentStatus, error) {
	return ShipmentPending, nil
}

type Delhivery struct {
}

func (d *Delhivery) ShipOrder(orderId string, address *Address) (*Shipment, error) {
	// will implement later
	return nil, nil
}

func (d *Delhivery) TrackShipment(trackingNumber string) (ShipmentStatus, error) {
	return ShipmentPending, nil
}

type ShippingProviderType string

const (
	FedEx ShippingProviderType = "FEDEX"
	UPS   ShippingProviderType = "UPS"
	DHL   ShippingProviderType = "DHL"
)

type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "PENDING"
	ShipmentInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentDelivered ShipmentStatus = "DELIVERED"
)

type Shipment struct {
	Id                string
	OrderId           string               //if the order is placed otheriwse empty
	Provider          ShippingProviderType //(Delhivery, Ekart etc)
	Status            ShipmentStatus
	EstimatedDelivery time.Time
	TrackingNumber    string
	ActualDelivery    *time.Time
}

// PaymentService
type PaymentService struct {
	gateway PaymentGateway //payment service has a paymentgateway
}

func (s *PaymentService) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	//add business logic (validation, logging, etc.)

	//delegate to gateway
	return s.gateway.ProcessPayment(orderId, amount, paymentType) //can be any gateway based on runtime
}

func (s *PaymentService) RefundPayment(paymentId string, amount float64) error {
	//business logic
	//delegate to gateway
	return s.gateway.RefundPayment(paymentId, amount)
}

type PaymentGateway interface {
	ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error)
	RefundPayment(paymentId string, amount float64) error
}

// implementations of payment gateway
type StripeGateway struct {
	apiKey string // ✅ Shows you understand external API auth
}

func (s *StripeGateway) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	return nil, nil
}

func (s *StripeGateway) RefundPayment(paymentId string, amount float64) error {
	return nil
}

type RazorpayGateway struct {
	keyId     string // ✅ Shows you know different services
	keySecret string //    have different auth formats
}

func (r *RazorpayGateway) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	return nil, nil
}
func (r *RazorpayGateway) RefundPayment(paymentId string, amount float64) error {
	return nil
}

type MockPaymentGateway struct {
}

func (m *MockPaymentGateway) ProcessPayment(orderId string, amount float64, paymentType PaymentType) (*Payment, error) {
	return nil, nil
}
func (m *MockPaymentGateway) RefundPayment(paymentId string, amount float64) error {
	return nil
}

type PaymentStatus string

const (
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
)

type PaymentType string

const (
	CreditCard PaymentType = "CREDITCARD"
	DebitCard  PaymentType = "DEBITCARD"
	UPI        PaymentType = "UPI"
)

type Payment struct {
	Id            string
	OrderId       string
	Type          PaymentType
	Status        PaymentStatus
	Amount        float64
	TransactionId string
	Timestamp     time.Time
}

// func testHighLoadWithoutWorkerPool() {
// 	//create inventory service
// 	inventoryService := &InventoryService{
// 		inventory: &Inventory{
// 			Stock: make(map[string]int),
// 		},
// 	}

// 	//Initialise some stock
// 	inventoryService.AddStock("laptop", 100000)

// 	// Print memory before
// 	printMemoryUsage("Before spawning goroutines")
// 	start := time.Now() // ← Time it

// 	var wg sync.WaitGroup

// 	//spawn 10000 goroutines
// 	for i := 0; i < 10000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			for j := 0; j < 10; j++ { //<- Each goroutine does 10 ops
// 				inventoryService.RemoveStock("laptop", 1)
// 				time.Sleep(1 * time.Millisecond)
// 			}
// 		}()

// 		// go inventoryService.RemoveStock("laptop", 1)
// 	}
// 	// Print memory DURING goroutines
// 	time.Sleep(50 * time.Millisecond) // Let them start
// 	printMemoryUsage("After spawning goroutines")

// 	wg.Wait() //Wait for all goroutines to finish

// 	elapsed := time.Since(start)

// 	// Print memory after
// 	printMemoryUsage("After all goroutines finished")

// 	// Check final stock
// 	stock, _ := inventoryService.CheckStock("laptop")
// 	fmt.Printf("\nFinal stock: %d (Expected: 0)\n", stock)
// 	fmt.Printf("Time taken: %v\n", elapsed)
// }

func testHighLoadWithWorkerPool() {

	//Initialise some stock
	// inventoryService.WorkerPoolAddStock("laptop", 100000)

	// Print memory before
	printMemoryUsage("Before spawning workerpool")
	start := time.Now() // ← Time it

	//start worker pool
	inventoryManager := NewInventoryManager(10, 100)
	err := inventoryManager.AddStock("laptop", 1)
	if err != nil {
		fmt.Printf("err ", err)
	}
	quantity, err := inventoryManager.CheckStock("laptop")
	if err != nil {
		fmt.Printf("err ", err)
	}

	if quantity > 1 {
		err = inventoryManager.RemoveStock("laptop", 1)
		if err != nil {
			fmt.Printf("err ", err)
		}
	}

	elapsed := time.Since(start)

	// Print memory after
	printMemoryUsage("After all worker pools finished")

	// Check final stock
	fmt.Printf("Time taken: %v\n", elapsed)
}

func printMemoryUsage(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: Alloc = %v MB\n", label, m.Alloc/1024/1024)
}

func main() {
	//simulate 100000 inventory operations
	// testHighLoadWithoutWorkerPool()

	//simulate 100000 inventory operations
	testHighLoadWithWorkerPool()
}
