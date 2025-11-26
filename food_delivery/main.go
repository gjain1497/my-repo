package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Flow

// User registers -> User login -> Search restaurant/food item (both will show list of restaurants)
// -> Go to specific restaurant -> shows list of food items -> add specific food items with quantity -> adds to cart
// -> places order (selecting specific address) -> payment -> (if success) -> reastaurant accepts order and starts preparing ->
// after some time some delivery partner accepts who is present in that coordinates range
// -> order confirmed -> (tracks order i.e delivery guy updates coordinates) -> order delivered -> rate order(optional)

// NICE to Have (Extra Points):

// ⚠️ Order Tracking - Location updates
// ⚠️ Promo Codes - Discount system
// ⚠️ Delivery Fee Calculation - Based on distance

// UserService
type UserService struct {
	Users map[string]User //list of users in memory (user_id, User object)
}

func (s *UserService) RegisterUser(userId string) {

}
func (s *UserService) GetUserDetails(userId string) {

}
func (s *UserService) UpdateUserDetails(userId string) {

}

type User struct {
	Id        string
	Name      string
	Phone     string
	Email     string
	Addresses []Address
}

// Observer 2: User Notification
type UserNotificationObserver struct {
}

func (u *UserNotificationObserver) OnOrderPlaced(order *Order) error {
	fmt.Printf("📧 [NOTIFICATION] User: %s | Order %s placed successfully!\n",
		order.UserId, order.Id)

	// In production: Send confirmation email/SMS to user

	return nil
}

type Coordinates struct {
	Latitude  string
	Longitude string
}

type Address struct {
	Id          string
	UserId      string
	Street      string
	City        string
	State       string
	ZipCode     string
	Country     string
	Coordinates Coordinates
	IsDefault   bool
}

// Observer 3: Analytics
type AnalyticsObserver struct{}

func (a *AnalyticsObserver) OnOrderPlaced(order *Order) error {
	fmt.Printf("📊 [ANALYTICS] Revenue: ₹%.2f | Restaurant: %s | Time: %s\n",
		order.TotalAmount, order.RestaurantId, order.CreatedAt.Format("15:04:05"))

	// In production: Send to analytics platform (Google Analytics, Mixpanel, etc.)

	return nil
}

// RestaurantService
type RestaurantService struct {
	Restaurants map[string]*Restaurant //list of restaurants (restaurant_id, Restaurant object)
	// OrderService *OrderService          //but its leading to cyclic dependency
	//so we will not add this
	//rather these functions which are
}

// observer handles notification logic
type RestaurantNotificationObserver struct {
	RestaurantService *RestaurantService //to fetch restaurant details
	SMSClient         *TwilioClient
	EmailClient       *SendGridClient
}

func (r *RestaurantNotificationObserver) OnOrderPlaced(order *Order) error {
	// Fetch restaurant details
	restaurant, err := r.RestaurantService.GetResturantDetails(order.RestaurantId)
	if err != nil {
		return err
	}
	// In production, would send actual SMS/Email
	// For now, just log
	fmt.Printf("📢 Notifying restaurant %s (phone: %s) about order %s\n",
		restaurant.Name, restaurant.Address.Street, order.Id)

	// In production:
	r.SMSClient.Send(restaurant.Phone, "New order!")
	r.EmailClient.Send(restaurant.Email, "New order!", "hello new order")

	return nil
}

// External service clients (just stubs for LLD)
type TwilioClient struct {
	apiKey string
}

func (t *TwilioClient) Send(phone string, message string) error {
	// In production: actual API call to Twilio
	fmt.Printf("📱 SMS to %s: %s\n", phone, message)
	return nil
}

type SendGridClient struct {
	apiKey string
}

func (s *SendGridClient) Send(email string, subject string, body string) error {
	// In production: actual API call to SendGrid
	fmt.Printf("📧 Email to %s: %s\n", email, subject)
	return nil
}

func (s *RestaurantService) GetResturantDetails(restaurantId string) (*Restaurant, error) {
	restaurant, ok := s.Restaurants[restaurantId]
	if !ok {
		return nil, errors.New("Restaurant does not exist")
	}
	return restaurant, nil
}

func (s *RestaurantService) ListFoodItems() (*[]MenuItem, error) {
	//will implement later
	return nil, nil
}

func (s *RestaurantService) GetMenuItemByRestaurant(restaurantId string, menuItemId string) (*MenuItem, error) {
	restaurant, ok := s.Restaurants[restaurantId]
	if !ok {
		return nil, errors.New("Restaurant does not exist")
	}
	menuItem, ok := restaurant.Menu[menuItemId]
	if !ok {
		return nil, errors.New("Menu Item does not exist")
	}
	return menuItem, nil
}

// Search by multiple filters (all optional)
func (s *RestaurantService) SearchRestaurants(cuisine string, city string, minRating float64) []*Restaurant {
	var restaurants []*Restaurant

	for _, restaurant := range s.Restaurants {
		// Only active restaurants
		if !restaurant.IsActive {
			continue
		}

		// Filter by cuisine (if provided)
		if cuisine != "" && restaurant.Cuisine != cuisine {
			continue
		}

		// Filter by city (if provided)
		if city != "" && restaurant.Address.City != city {
			continue
		}

		// Filter by minimum rating
		if restaurant.Rating < minRating {
			continue
		}

		restaurants = append(restaurants, restaurant)
	}

	return restaurants
}

// Search food items across all restaurants
func (s *RestaurantService) SearchFoodItems(query string) []*MenuItem {
	var menuItems []*MenuItem

	query = strings.ToLower(query) // Case-insensitive

	for _, restaurant := range s.Restaurants {
		if !restaurant.IsActive {
			continue
		}

		for _, item := range restaurant.Menu {
			if !item.IsAvailable {
				continue
			}

			// Search in name, description, cuisine
			if strings.Contains(strings.ToLower(item.Name), query) ||
				strings.Contains(strings.ToLower(item.Description), query) ||
				strings.Contains(strings.ToLower(item.Cuisine), query) {
				menuItems = append(menuItems, item)
			}
		}
	}

	return menuItems
}

// Get restaurants near location
func (s *RestaurantService) GetNearbyRestaurants(location Coordinates, radiusKm float64) []*Restaurant {
	var nearbyRestaurants []*Restaurant

	for _, restaurant := range s.Restaurants {
		if !restaurant.IsActive {
			continue
		}

		distance := calculateDistance(location, restaurant.Address.Coordinates)

		if distance <= radiusKm {
			nearbyRestaurants = append(nearbyRestaurants, restaurant)
		}
	}

	return nearbyRestaurants
}

//Claude I have a doubt here should it be part of RestaurantService or
//I should have a new Service maybe InventoryService for Seperation of Concerns where
/*
### The Difference:

| RestaurantService | InventoryService |
|----------------|------------------|
| Manages product DETAILS | Manages product QUANTITIES |
| Add new product to catalog | Add stock for existing product |
| Update price, description | Update quantity |
| Search products | Check availability |
| Product info (what) | Stock levels (how many) |

**Example Flow:**
```
1. Admin: AddProduct("Laptop", $1000, "Gaming laptop", "Electronics") //I mean food here obviously
   → RestaurantService creates product in catalog

2. Warehouse: AddStock("laptop-id", 50)
   → InventoryService adds 50 units to inventory

3. Customer: View product page
   → RestaurantService shows details ($1000, description)
   → InventoryService shows "50 in stock"

4. Customer: Place order (quantity: 2)
   → InventoryService removes 2 from stock (now 48)
   → Product details unchanged
*/

func (s *RestaurantService) RemoveMenuItems(restaurantId string, menuItemId string) error {
	_, ok := s.Restaurants[restaurantId]
	if !ok {
		return errors.New("Restaurant does not exist")
	}
	stock := s.Restaurants[restaurantId].Stock[menuItemId]
	if stock > 0 {
		s.Restaurants[restaurantId].Stock[menuItemId]--
	}
	return nil
}

func (s *RestaurantService) AddMenuItems(restaurantId string, menuItemId string) error {
	_, ok := s.Restaurants[restaurantId]
	if !ok {
		return errors.New("Restaurant does not exist")
	}
	s.Restaurants[restaurantId].Stock[menuItemId]++
	return nil
}

type MenuItem struct {
	Id           string
	RestaurantId string
	Name         string
	Price        float64
	Cuisine      string
	Description  string
	IsAvailable  bool
	Category     string
}

type Restaurant struct {
	Id           string
	Name         string
	Phone        string
	Email        string
	Address      Address
	Cuisine      string
	Menu         map[string]*MenuItem //(item_id, MenuItem Object)
	Stock        map[string]int       //(item_id, Quantity)
	IsActive     bool
	Rating       float64 //avg out of 5 stars let's say
	TotalRatings int
}

type CartService struct {
	Carts             map[string]*Cart //(user_id -> cart object) //because one user has one cart
	RestaurantService *RestaurantService
}

func (s *CartService) AddToCart(userId string, restaurantId string, itemId string, quantity int) error {
	cart, ok := s.Carts[userId]
	if !ok {
		cart = &Cart{
			Id:           userId + "_cart",
			UserId:       userId,
			RestaurantId: restaurantId,
			Items:        make(map[string]int),
		}
	} else {
		//cart exist then validate that it should be from the same restaurant
		if cart.RestaurantId != "" && cart.RestaurantId != restaurantId {
			return errors.New("Cannot add items from different restaurants. Place a new order")
		}
	}
	//Validate Item exists and its avaialble
	menuItem, err := s.RestaurantService.GetMenuItemByRestaurant(restaurantId, itemId)
	if err != nil {
		return fmt.Errorf("item not found: %w", err)
	}
	if !menuItem.IsAvailable {
		return fmt.Errorf("item %s is not available", menuItem.Name)
	}
	cart.Items[itemId] += quantity
	s.Carts[userId] = cart
	return nil
}

func (s *CartService) DeleteFromCart(userId string, restaurantId string, itemId string, quantity int) error {
	cart, ok := s.Carts[userId]
	if !ok {
		return errors.New("Cart does not exist")
	}

	// ✅ Check if item exists
	currentQty, exists := cart.Items[itemId]
	if !exists {
		return errors.New("Item not in cart")
	}

	// ✅ Validate quantity
	if quantity > currentQty {
		return fmt.Errorf("Cannot remove %d items. Only %d in cart", quantity, currentQty)
	}

	// Update or remove
	if quantity == currentQty {
		delete(cart.Items, itemId) // ✅ Remove completely
	} else {
		cart.Items[itemId] -= quantity
	}

	// ✅ If cart is now empty, clear restaurant
	if len(cart.Items) == 0 {
		cart.RestaurantId = ""
	}

	return nil
}

func (s *CartService) GetCart(userId string) (*Cart, error) {
	//return cart
	cart, ok := s.Carts[userId]
	if !ok {
		return nil, fmt.Errorf("Cart does not exist for this user")
	}
	return cart, nil
}

// restaurantId also required IMO, because at a time we can order only from one restaurant
// and this will be needed to fetch food details from this restaurant
func (s *CartService) CalculateTotal(userId string) (float64, error) {
	// for each item in cart
	// get detail of food items from restaurant
	cart, ok := s.Carts[userId]
	if !ok {
		return 0, errors.New("Cart not found for user")
	}
	if len(cart.Items) == 0 {
		return 0, nil // Empty cart = 0 total
	}
	if cart.RestaurantId == "" {
		return 0, errors.New("Cart has no restaurant associated")
	}
	restaurantId := cart.RestaurantId
	var total float64 = 0
	for foodItemId, quantity := range cart.Items {
		//get foodItemDetails from this foodItemId
		menuItem, err := s.RestaurantService.GetMenuItemByRestaurant(restaurantId, foodItemId)
		if err != nil {
			return 0, errors.Join(errors.New("CalculateTotal err: "), err) //error wrapping for detailed errors
		}
		// Add this check for availability because it can be possible some other user already orderd it within that time
		// 10:00 AM: User browsing menu, Chicken Biryani is available ✅
		// 10:01 AM: User adds to cart
		// 10:05 AM: User browsing other items, adding more
		// 10:10 AM: Restaurant marks Chicken Biryani as unavailable ❌ (sold out)
		// 10:15 AM: User clicks "Place Order" → CalculateTotal() is called
		if !menuItem.IsAvailable {
			return 0, fmt.Errorf("Item %s is no longer available", menuItem.Name)
		}
		price := menuItem.Price
		total += price * float64(quantity)
	}

	return total, nil
}

// Cart Service
type Cart struct {
	Id string
	//Price -> no, just kidding, price is dynamic we cant store these type of things at cart level
	//as there is only one cart per user
	UserId       string
	RestaurantId string //because at a time we can order only from one restaurant
	//and this will be needed to fetch food details from this restaurant
	Items map[string]int //food_item_id, quantity
}

// OrderObserver - Any Observer must implement this
// type OrderObserver interface {
// 	OnOrderPlaced(order *Order) error
// }

type OrderPlacedObserver interface {
	OnOrderPlaced(order *Order) error
}

type OrderReadyObserver interface {
	OnOrderReady(order *Order) error
}

//Split this into two parts
//OrderPlacedObserver  and OrderReadyObserver

//create the observable (that observes here that is order)

// OrderService
type OrderService struct {
	Orders               map[string]*Order       //[order_id, Order object]
	PaymentService       PaymentGatewayInterface //can take either raozpay gateway or stripe gateway or any new gateway being added
	CartService          *CartService
	RestaurantService    *RestaurantService
	OrderPlacedObservers []OrderPlacedObserver
	OrderReadyObservers  []OrderReadyObserver
}

// subscribe an observer
func (s *OrderService) SubscribeOrderPlaced(observer OrderPlacedObserver) {
	s.OrderPlacedObservers = append(s.OrderPlacedObservers, observer)
}

// subscribe an observer
func (s *OrderService) SubscribeOrderReady(observer OrderReadyObserver) {
	s.OrderReadyObservers = append(s.OrderReadyObservers, observer)
}

// notify - tell all subscribers about the change
// ✅ Async with goroutines
func (s *OrderService) notifyAllOrderPlaced(order *Order) {
	for _, observer := range s.OrderPlacedObservers {
		go func(obs OrderPlacedObserver, o *Order) {
			err := obs.OnOrderPlaced(o)
			if err != nil {
				fmt.Printf("⚠️ Observer Order Placed notification failed: %v\n", err)
			}
		}(observer, order)
	}
}

// // notify - tell all subscribers about the change
// ✅ Async with goroutines
func (s *OrderService) notifyAllOrderReady(order *Order) {
	for _, observer := range s.OrderReadyObservers {
		go func(obs OrderReadyObserver, o *Order) {
			err := obs.OnOrderReady(o)
			if err != nil {
				fmt.Printf("⚠️ Observer Order Ready notification failed: %v\n", err)
			}
		}(observer, order)
	}
}

func (s *OrderService) GetOrderByOrderId(orderId string) (*Order, error) {
	order, ok := s.Orders[orderId]
	if !ok {
		return nil, errors.New("getOrderByOrderId: Order not found")
	}
	return order, nil
}

// API called by User
func (s *OrderService) PlaceOrder(cartId string, userId string, addressId string, paymentType PaymentType) (*Order, error) {
	//from cart -> order first
	cart, err := s.CartService.GetCart(userId)
	if err != nil {
		return nil, errors.Join(errors.New("PlaceOrder: "), err)
	}
	restaurantId := cart.RestaurantId //Claude question here should we keep restaurant id at cart level or keep passing in functions
	//for eg we can pass it in PlaceOrder function itself as a parameter

	//generate new order id maybe uuid
	orderId := "23123123421"
	var orderItems []*OrderItem
	//iterate over cart items and put it into order items
	for itemId, quantity := range cart.Items {
		menuItem, err := s.RestaurantService.GetMenuItemByRestaurant(restaurantId, itemId)
		if err != nil {
			return nil, errors.Join(errors.New("PlaceOrder: "), err)
		}
		orderItem := &OrderItem{
			FoodItemId:   itemId,
			FoodItemName: menuItem.Name,
			Quantity:     quantity,
			Price:        menuItem.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	//calcualte total using current prices
	total, err := s.CartService.CalculateTotal(userId)
	if err != nil {
		return nil, errors.Join(errors.New("PlaceOrder: "), err)
	}

	order := &Order{
		Id:                orderId,
		UserId:            userId,
		RestaurantId:      restaurantId,
		DeliveryPartnerId: "to be assigned yet",
		Status:            Pending,
		OrderItems:        orderItems,
		CreatedAt:         time.Now(),
		TotalAmount:       total,
		PaymentId:         "to be assigned yet",
		ShippingAddressId: addressId,
	}

	//process payment with this order details
	payment, err := s.PaymentService.ProcessPayment(orderId, userId, total)

	if err != nil || payment.Status != PaymentSuccess {
		return nil, errors.Join(errors.New("PlaceOrder Payment failure: "), err)
	}
	order.PaymentId = payment.Id

	//if we reached till here that means payment got succeeded
	//restaurant to accept the order now

	//Claude what I am thinking here is, this should not be blocker for this right
	//I mean in actual life AcceptOrder is something which should be triggered by
	//Restaurant admin API right, why are we calling it here, confusion

	//we just kind of notify restaurant

	//then when restaurantAdmin clicks on accept button this service level function
	//gets triggered and that function then notifies us
	//So I dont think it doesn't make any sense to call it here
	//it shoudl be something like this maybe pub sub pattern or observer pattern
	// err = s.NotifyRestaurant(orderId)

	// err = s.AcceptOrder(orderId, restaurantId)
	// if err != nil {
	// 	return nil, errors.Join(errors.New("PlaceOrder Payment failure: "), err)
	// }
	s.Orders[orderId] = order
	s.notifyAllOrderPlaced(order)

	//maybe a polling mechanism which keeps on polling whats order status now
	// if order.Status == ReadyForPickup {
	//here assign delivery partner
	//Claude I will implement it I know but for now llets review until whatever I have done
	// }

	//then at last create order with all the details
	delete(s.CartService.Carts, userId) // Clear the cart!
	return order, nil
}

// APIs called by Restaurant Admin /accept/restaurantId/orderId
func (s *OrderService) AcceptOrder(orderId string, restaurantId string) error {
	//changes order status
	order, ok := s.Orders[orderId]
	if !ok {
		return errors.New("order does not exist")
	}
	if order.RestaurantId != restaurantId {
		return errors.New("Order doesn't belong to this restaurant")
	}
	order.Status = Confirmed
	return nil
}

// APIs called by Restaurant Admin /update_status/restaurantId/orderId
// Called after food gets prepared
func (s *OrderService) UpdateOrderStatus(orderId string, restaurantId string) error {
	//changes order status
	order, ok := s.Orders[orderId]
	if !ok {
		return errors.New("order does not exist")
	}

	// Transition based on current status
	switch order.Status {
	case Confirmed:
		order.Status = Preparing
		fmt.Printf("Order %s is now being prepared\n", orderId)
	case Preparing:
		order.Status = ReadyForPickup

		fmt.Printf("✅ Order %s is ready for pickup\n", orderId)
		for _, item := range order.OrderItems {
			err := s.RestaurantService.RemoveMenuItems(restaurantId, item.FoodItemId)
			if err != nil {
				return err
			}
		}

		// ✅ NOW trigger delivery partner assignment
		s.notifyAllOrderReady(order)
	default:
		return fmt.Errorf("cannot update order in status: %s", order.Status)
	}

	return nil
}

func (s *OrderService) RejectOrder() {
	//something
}

// enum
type OrderStatus string

const (
	Pending        OrderStatus = "PENDING"
	Confirmed      OrderStatus = "CONFIRMED"
	Preparing      OrderStatus = "PREPARING"
	ReadyForPickup OrderStatus = "READY_FOR_PICKUP"
	OutForDelivery OrderStatus = "OUT_FOR_DELIVERY"
	Delivered      OrderStatus = "DELIVERED"
	Cancelled      OrderStatus = "CANCELLED"
)

type OrderItem struct {
	FoodItemId   string
	FoodItemName string
	Quantity     int
	Price        float64
}

type Order struct {
	Id                string
	UserId            string //for this user
	RestaurantId      string //from this restaurant
	DeliveryPartnerId string //delivery by this guy
	Status            OrderStatus
	OrderItems        []*OrderItem
	CreatedAt         time.Time
	TotalAmount       float64
	PaymentId         string
	ShippingAddressId string
}

// Tracking Service or can be part of order service as well
type TrackingService struct {
	Trackings map[string]Tracking //(tracking_id -> Tracking)
}

func (s *TrackingService) TrackOrder(orderId string) {

}

type Tracking struct {
	OrderId           string
	DeliveryPartnerId string
	RestaurantId      string
	CurrentLocation   Coordinates
}

// Payment Service
//Heyyy Claude I have commented this service, as now we we will be using PaymentGatewayInterface directly in let's say
//other services, for eg in order service we used and similarly we can use interface directly in httpHandler as well
//then why is this even required

// type PaymentService struct {
// 	PaymentGateway PaymentGatewayInterface
// }

// func (p *PaymentService) ProcessPayment(orderId string, userId string) (*Payment, error) {
// 	//some business logic
// 	payment, err := p.PaymentGateway.ProcessPayment(orderId, userId)
// 	return payment, err
// }

// func (p *PaymentService) RefundPayment(orderId string, userId string) error {
// 	//some business logic
// 	err := p.PaymentGateway.RefundPayment(orderId, userId)
// 	return err
// }

type PaymentGatewayInterface interface {
	ProcessPayment(orderId string, userId string, amount float64) (*Payment, error)
	RefundPayment(orderId string, userId string, amount float64) error
}

type StripePaymentGateway struct {
	apiKey string
}

func (s *StripePaymentGateway) ProcessPayment(orderId string, userId string, amount float64) (*Payment, error) {
	payment := &Payment{
		Id:            "stripe_" + orderId,
		Status:        PaymentSuccess,
		OrderId:       orderId,
		Type:          CreditCard, // or parse from somewhere
		Amount:        amount,
		TransactionId: "txn_" + orderId,
		Timestamp:     time.Now(),
	}
	return payment, nil
}

func (s *StripePaymentGateway) RefundPayment(orderId string, userId string, amount float64) error {
	fmt.Printf("💰 Refunding ₹%.2f for order %s\n", amount, orderId)
	return nil
}

type RazorpayPaymentGateway struct {
	apiKey string
}

func (r *RazorpayPaymentGateway) ProcessPayment(orderId string, userId string, amount float64) (*Payment, error) {
	payment := &Payment{
		Id:            "razorpay_" + orderId,
		Status:        PaymentSuccess,
		OrderId:       orderId,
		Type:          UPI,
		Amount:        amount,
		TransactionId: "rzp_" + orderId,
		Timestamp:     time.Now(),
	}
	return payment, nil
}

func (r *RazorpayPaymentGateway) RefundPayment(orderId string, userId string, amount float64) error {
	fmt.Printf("💰 Refunding ₹%.2f for order %s via Razorpay\n", amount, orderId)
	return nil
}

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

type PaymentType string

const (
	CreditCard PaymentType = "CREDITCARD"
	DebitCard  PaymentType = "DEBITCARD"
	UPI        PaymentType = "UPI"
)

type Payment struct {
	Id            string
	Status        PaymentStatus
	OrderId       string
	Type          PaymentType
	Amount        float64
	TransactionId string
	Timestamp     time.Time
}

// observer handles notification logic
type DeliveryPartnerNotificationObserver struct {
	DeliveryMatchingService *DeliveryMatchingService //to assign delivery partner
	RestaurantService       *RestaurantService       //to fetch restaurant details
	SMSClient               *TwilioClient
	EmailClient             *SendGridClient
}

func (d *DeliveryPartnerNotificationObserver) OnOrderReady(order *Order) error {
	// Fetch restaurant details
	restaurant, err := d.RestaurantService.GetResturantDetails(order.RestaurantId)
	if err != nil {
		return fmt.Errorf("failed to get restaurant: %w", err)
	}
	deliveryPartner, err := d.DeliveryMatchingService.AssignDeliveryPartner(order, restaurant.Address.Coordinates)
	if err != nil {
		fmt.Printf("No delivery partner available for order %s\n", order.Id)
		return err
	}

	// ✅ ADD: Notify the delivery partner
	if d.SMSClient != nil && deliveryPartner.Phone != "" {
		d.SMSClient.Send(deliveryPartner.Phone,
			fmt.Sprintf("New delivery assigned! Order: %s, Restaurant: %s", order.Id, restaurant.Name))
	}

	if d.EmailClient != nil && deliveryPartner.Email != "" {
		d.EmailClient.Send(deliveryPartner.Email, "New Delivery Assignment",
			fmt.Sprintf("Order %s is ready for pickup at %s", order.Id, restaurant.Name))
	}

	fmt.Printf("🏍️ Order %s assigned to delivery partner %s\n", order.Id, deliveryPartner.Name)

	return nil
}

// DeliveryMatchingService - Responsible for finding delivery partners
type DeliveryMatchingService struct {
	DeliveryPartners map[string]*DeliveryPartner //(id->DeliveryPartner object)
}

func (s *DeliveryMatchingService) AssignDeliveryPartner(order *Order, restaurantLocation Coordinates) (*DeliveryPartner, error) {
	//find available partner near the restaurant coordinates
	partner := s.findNearestAvailablePartner(restaurantLocation)
	if partner == nil {
		return nil, errors.New("no delivery partners available")
	}
	// Assign to order
	order.DeliveryPartnerId = partner.Id
	partner.IsAvailable = false

	return partner, nil
}

func (d *DeliveryMatchingService) findNearestAvailablePartner(location Coordinates) *DeliveryPartner {
	var nearest *DeliveryPartner
	minDistance := math.MaxFloat64

	for _, partner := range d.DeliveryPartners {
		if !partner.IsAvailable {
			continue
		}

		// Calculate distance
		distance := calculateDistance(partner.CurrentCoordinates, location)

		if distance < minDistance {
			minDistance = distance
			nearest = partner
		}
	}

	return nearest
}

// Helper function to calculate distance
func calculateDistance(coord1, coord2 Coordinates) float64 {
	// For simplicity, using Euclidean distance
	// In production, use Haversine formula for GPS coordinates

	// Convert string to float (you'd parse properly)
	// This is simplified
	return 0.0 // Implement based on your coordinate format
}

// user can be user as well as delivery boy
type DeliveryPartner struct {
	User
	VehicleNumber      string
	CurrentCoordinates Coordinates
	IsAvailable        bool
	//some more fields
}

// Rating Service
type Rating struct {
	Id                string
	Message           string
	Stars             int
	DeliveryPartnerId string
	RestaurantId      string
	OrderId           string
}

//Heyy Claude I think we missed one thing what about MatchEngine that mathces the User with DeliveryBoy (like uber)
//oh but no wait in this case it should be matching based on restaurant coordinates rather than user coordiantes

//restaurant accepts order -> match engine starts matching nearby delivery boys -

func main() {

}
