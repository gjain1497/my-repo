package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
 * System Requirements:

	1) System will support the renting of different automobiles like cars, trucks, SUVs, vans, and motorcycles.

	2) Each vehicle should be added with a unique barcode and other details, including a parking stall number which helps to locate the vehicle.

	3) System should be able to retrieve information like which member took a particular vehicle or what vehicles have been rented out by a specific member.

	4) System should collect a late-fee for vehicles returned after the due date.

	5) Members should be able to search the vehicle inventory and reserve any available vehicle.

	6) The system should be able to send notifications whenever the reservation is approaching the pick-up date, as well as when the vehicle is nearing the due date or has not been returned within the due date.

	7) The system will be able to read barcodes from vehicles.

	8) Members should be able to cancel their reservations.

	9) The system should maintain a vehicle log to track all events related to the vehicles.

	10) Members can add rental insurance to their reservation.

	11) Members can rent additional equipment, like navigation, child seat, ski rack, etc.

	12) Members can add additional services to their reservation, such as roadside assistance, additional driver, wifi, etc.
 *
 * */

// Add/Remove vehicles
// Search available vehicles by type/date
// Book a vehicle
// Return vehicle (calculate cost)
// Cancel booking
// User registration
// Payment processing

//Flow

// Admin adds/removes vehicle

// User registers ->

// Searches car in (current loc, by date, by type) -> gets list of
// vehicles, along with cost -> books a vehicle -> payment

// return vahicles-> calculate cost (if there is any difference with initial cost)

// Admin/Owner of Vehicle Flow
// 1 Adds Vehicle
// 2 Remove vehicle
// 3 Update vehicle

// User Flow
// 1 register
// 2 search vehicles (location, vehcileType, date_range)
// 3 get vehicle list with esitmated cost
// 4 book/reserve vehicle (for particular dates)
// 5 make payment
// 6 Return vehcile
// 7 Final payment cost with adjustment (lets say due to late fees or damage)
// 8 complete booking/ mark vehicle avaialable again

// Cancel Flow
// 1 Cancel
// 2 Refund

type Person struct {
	id   string
	name string
}

type AdminService struct {
}

func (s *AdminService) AddVehicle() {

}
func (s *AdminService) RemoveVehicle() {

}
func (s *AdminService) UpdateVehicle() {

}

type Admin struct {
	Person Person
	//some additional properties
	VehicleService VehicleServiceInterface
}

type UserService struct {
	users          map[string]User
	VehicleService VehicleServiceInterface
}

func (s *UserService) GetAllVehiclesRentedByUser(userId string) []*Vehicle {
	return s.VehicleService.ListVehiclesCurrentlyAssignedToAUser(userId)
}

type User struct {
	Person        Person
	LicenseNumber string
	Phone         string
	Email         string
	//some additional properties
}

type PricingStrategyFactory struct {
}

func (f *PricingStrategyFactory) GetStrategy(vehicleType VehicleType) (VehicleTypeInterface, error) {
	switch vehicleType {
	case Car:
		return &CarVehicle{}, nil
	case Bike:
		return &BikeVehicle{}, nil
	case Truck:
		return &TruckVehicle{}, nil
	default:
		return nil, errors.New("unknown vehicle type")
	}
}

type VehicleTypeInterface interface {
	CalculatePricePerDay() (float64, error)
	CalculateLateFeePerDay() (float64, error)
}

type CarVehicle struct {
}

func (c *CarVehicle) CalculatePricePerDay() (float64, error) {
	return 20, nil
}
func (c *CarVehicle) CalculateLateFeePerDay() (float64, error) {
	return 2, nil
}

type BikeVehicle struct {
}

func (c *BikeVehicle) CalculatePricePerDay() (float64, error) {
	return 15, nil
}
func (c *BikeVehicle) CalculateLateFeePerDay() (float64, error) {
	return 1.5, nil
}

type TruckVehicle struct {
}

func (c *TruckVehicle) CalculatePricePerDay() (float64, error) {
	return 10, nil
}
func (c *TruckVehicle) CalculateLateFeePerDay() (float64, error) {
	return 1.7, nil
}

type VehicleServiceInterface interface {
	ListVehicles(userId string, locationId string, vehicleType VehicleType, dateStart time.Time, dateEnd time.Time) ([]*Vehicle, error)
	ListVehiclesCurrentlyAssignedToAUser(userId string) []*Vehicle
	BookVehicle(userId string, vehicleId string) error
	UnBookVehicle(userId string, vehicleId string) error
	CalculateEstimatedPrice(userId string, vehicleId string, dateStart time.Time, dateEnd time.Time) (float64, error)
	CalculateNetPriceAfteReturn(userId string, vehicleId string, dateStart time.Time, dateEnd time.Time, actualDateEnd time.Time) (float64, error)
	GetVehicle(vehicleId string) (*Vehicle, error)
	AddVehicle(personId string, req AddVehicleRequest) error // ✅ Takes request struct
	RemoveVehicle(personId string, vehicleId string) error
	UpdateVehicle(personId string, req UpdateVehicleRequest) error // ✅ Takes request struct
}

type VehicleService struct {
	Vehicles        map[string]*Vehicle   //(id , Vehicle object)
	UserVehicles    map[string][]*Vehicle //(user_id, list of vehicles)
	Locations       map[string]*Location  // ✅ Store locations
	StrategyFactory *PricingStrategyFactory
	mu              sync.RWMutex
}

func (s *VehicleService) ListVehicles(userId string, locationId string, vehicleType VehicleType, dateStart time.Time, dateEnd time.Time) ([]*Vehicle, error) {
	//anyone can check no need of verifying that user is admin only
	s.mu.RLock()
	defer s.mu.RUnlock()

	var vehicles []*Vehicle
	for _, vehicle := range s.Vehicles {
		if vehicle.Type == vehicleType && vehicle.IsAvailable {
			vehicles = append(vehicles, vehicle)
		}
	}

	return vehicles, nil
}

func (s *VehicleService) ListVehiclesCurrentlyAssignedToAUser(userId string) []*Vehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.UserVehicles[userId]
}

func (s *VehicleService) BookVehicle(userId string, vehicleId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle, ok := s.Vehicles[vehicleId]
	if !ok {
		return errors.New("Vehicle not found")
	}

	vehicle.IsAvailable = false
	s.UserVehicles[userId] = append(s.UserVehicles[userId], vehicle)
	return nil
}

func (s *VehicleService) UnBookVehicle(userId string, vehicleId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle, ok := s.Vehicles[vehicleId]
	if !ok {
		return errors.New("Vehicle not found")
	}

	// ✅ VehicleService's responsibility: manage vehicle availability
	vehicle.IsAvailable = true

	// ✅ VehicleService's responsibility: manage user-vehicle mapping
	// This maintains data consistency within VehicleService's domain
	currentVehicles := s.ListVehiclesCurrentlyAssignedToAUser(userId)
	for i, vehicle := range currentVehicles {
		if vehicle.Id == vehicleId {
			currentVehicles = append(currentVehicles[:i], currentVehicles[i+1:]...)
			break
		}
	}
	s.UserVehicles[userId] = currentVehicles
	return nil
}

func (s *VehicleService) CalculateEstimatedPrice(userId string, vehicleId string, dateStart time.Time, dateEnd time.Time) (float64, error) {
	//anyone can check no need of verifying that user is admin only
	//we know price for one day

	//we use some strategy for different prices
	//calculateTimeDifference -> pricePerday * (number of days)
	s.mu.RLock() //Read lock (allows concurrent reads)
	vehicle, ok := s.Vehicles[vehicleId]
	s.mu.RUnlock()
	if !ok {
		return 0, errors.New("Vehicle not found")
	}
	timeDifference := dateEnd.Sub(dateStart)
	days := (timeDifference.Hours() / 24)
	if days <= 0 {
		return 0, errors.New("invalid date range")
	}

	//Get strategy using factory
	strategy, err := s.StrategyFactory.GetStrategy(vehicle.Type)
	if err != nil {
		return 0, err
	}

	//based on vehicleType different strategies
	pricePerDay, err := strategy.CalculatePricePerDay()
	if err != nil {
		return 0, err
	}

	// Calculate total
	return pricePerDay * float64(days), nil
}

func (s *VehicleService) calculateLateFee(vehicleId string, dateStart time.Time, dateEnd time.Time) (float64, error) {
	//anyone can check no need of verifying that user is admin only
	//we know price for one day

	//we use some strategy for different prices
	//calculateTimeDifference -> pricePerday * (number of days)
	s.mu.RLock() //Read lock (allows concurrent reads)
	vehicle, ok := s.Vehicles[vehicleId]
	s.mu.RUnlock()
	if !ok {
		return 0, errors.New("Vehicle not found")
	}
	timeDifference := dateEnd.Sub(dateStart)
	days := (timeDifference.Hours() / 24)
	if days <= 0 {
		return 0, errors.New("invalid date range")
	}

	//Get strategy using factory
	strategy, err := s.StrategyFactory.GetStrategy(vehicle.Type)
	if err != nil {
		return 0, err
	}

	//based on vehicleType different strategies
	lateFeePerDay, err := strategy.CalculateLateFeePerDay()
	if err != nil {
		return 0, err
	}

	// Calculate total
	return lateFeePerDay * float64(days), nil
}

func (s *VehicleService) CalculateNetPriceAfteReturn(userId string, vehicleId string, dateStart time.Time, dateEnd time.Time, actualDateEnd time.Time) (float64, error) {
	//anyone can check no need of verifying that user is admin only
	//we know price for one day
	price, err := s.CalculateEstimatedPrice(userId, vehicleId, dateStart, dateEnd)
	if err != nil {
		return 0, err
	}
	var lateFee float64
	if actualDateEnd.After(dateEnd) {
		lateFee, err = s.calculateLateFee(vehicleId, dateEnd, actualDateEnd)
		if err != nil {
			return 0, err
		}
	}

	return price + lateFee, nil
}

func (s *VehicleService) GetVehicle(vehicleId string) (*Vehicle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	vehicle, ok := s.Vehicles[vehicleId]
	if !ok {
		return nil, errors.New("vehicle not found")
	}
	return vehicle, nil
}

type AddVehicleRequest struct {
	VehicleType     VehicleType
	Model           string
	Brand           string
	RegistrationNum string
	LocationId      string // ✅ Just ID, not full object
	PricePerDay     float64
	LateFeePerDay   float64
}

// ✅ AddVehicle - Takes request struct, generates ID
func (s *VehicleService) AddVehicle(personId string, req AddVehicleRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Verify personId is admin

	// ✅ System generates vehicle ID
	vehicleId := "112231"

	// ✅ Get location from system
	location, ok := s.Locations[req.LocationId]
	if !ok {
		return errors.New("invalid location ID")
	}

	// ✅ Create vehicle with system-generated values
	vehicle := &Vehicle{
		Id:              vehicleId,
		Type:            req.VehicleType,
		Model:           req.Model,
		Brand:           req.Brand,
		RegistrationNum: req.RegistrationNum,
		IsAvailable:     true, // ✅ Default to available
		Location:        *location,
		PricePerDay:     req.PricePerDay,
		LateFeePerDay:   req.LateFeePerDay,
	}

	s.Vehicles[vehicleId] = vehicle
	return nil
}

func (s *VehicleService) RemoveVehicle(personId string, vehicleId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Verify personId is admin

	_, ok := s.Vehicles[vehicleId]
	if !ok {
		return errors.New("vehicle not found")
	}

	delete(s.Vehicles, vehicleId)
	return nil
}

// ✅ Request struct for updating vehicle
type UpdateVehicleRequest struct {
	VehicleId       string
	Model           string
	Brand           string
	RegistrationNum string
	LocationId      string
	PricePerDay     float64
	LateFeePerDay   float64
}

func (s *VehicleService) UpdateVehicle(personId string, req UpdateVehicleRequest) error {
	//will verify that person in admin maybe through the access token/jwt
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Verify personId is admin

	vehicle, ok := s.Vehicles[req.VehicleId]
	if !ok {
		return errors.New("vehicle not found")
	}

	// ✅ Get location from system
	location, ok := s.Locations[req.LocationId]
	if !ok {
		return errors.New("invalid location ID")
	}

	// ✅ Update vehicle fields
	vehicle.Model = req.Model
	vehicle.Brand = req.Brand
	vehicle.RegistrationNum = req.RegistrationNum
	vehicle.Location = *location
	vehicle.PricePerDay = req.PricePerDay
	vehicle.LateFeePerDay = req.LateFeePerDay

	s.Vehicles[req.VehicleId] = vehicle
	return nil

}

type Vehicle struct {
	Id              string
	Type            VehicleType
	Model           string // "Honda City", "Royal Enfield"
	Brand           string // "Honda", "Royal Enfield"
	RegistrationNum string // "DL-01-AB-1234"
	IsAvailable     bool
	Location        Location
	PricePerDay     float64
	LateFeePerDay   float64 //late fee if vehicle returned after end date
}

type VehicleType string

const (
	Car   VehicleType = "CAR"
	Truck VehicleType = "TRUCK"
	Bike  VehicleType = "BIKE"
)

type VehicleStatus string

const (
	Available VehicleStatus = "AVAILABLE"
	Booked    VehicleStatus = "BOOKED"
)

type Location struct {
	City      string
	State     string
	ZipCode   string
	Country   string
	IsDefault bool
}

type PaymentServiceInterface interface {
	ProcessPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error)
	RefundPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error)
	GetPaymentHistory(bookingId string) ([]*Payment, error) // ✅ Extra method
}
type PaymentService struct {
	PaymentGateway PaymentGateway
	Payments       map[string]*Payment
	mu             sync.RWMutex
}

func (s *PaymentService) ProcessPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
	// ✅ Business logic: validation
	if bookingId == "" {
		return nil, errors.New("booking ID required")
	}
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	// ✅ Business logic: convert to cents
	amountInCents := int(amount * 100)

	// ✅ Business logic: generate transaction ID
	transactionId := generateTransactionID()

	// Call gateway (different method name!)
	gatewayResp, err := s.PaymentGateway.Charge(transactionId, amountInCents)
	if err != nil {
		return nil, err
	}

	// ✅ Business logic: convert gateway response to Payment domain model
	payment := &Payment{
		Id:        gatewayResp.TransactionId,
		Type:      paymentType,
		Amount:    amount,
		Status:    s.mapGatewayStatus(gatewayResp.Status),
		CreatedAt: time.Now(),
	}

	// ✅ Business logic: store payment
	s.mu.Lock()
	s.Payments[payment.Id] = payment
	s.mu.Unlock()

	return payment, nil
}

func (s *PaymentService) RefundPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
	if bookingId == "" {
		return nil, errors.New("booking ID required")
	}
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	amountInCents := int(amount * 100)
	transactionId := generateTransactionID()

	gatewayResp, err := s.PaymentGateway.Refund(transactionId, amountInCents)
	if err != nil {
		return nil, err
	}

	payment := &Payment{
		Id:        gatewayResp.TransactionId,
		BookingId: bookingId,
		Type:      paymentType,
		Amount:    amount,
		Status:    REFUNDED,
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.Payments[payment.Id] = payment
	s.mu.Unlock()

	return payment, nil
}

func (s *PaymentService) GetPaymentHistory(bookingId string) ([]*Payment, error) {
	// ✅ Additional business method (not in gateway)
	s.mu.RLock()
	defer s.mu.RUnlock()

	var payments []*Payment
	for _, payment := range s.Payments {
		// Assuming we store bookingId in Payment (add to Payment struct)
		payments = append(payments, payment)
	}
	return payments, nil
}
func (s *PaymentService) mapGatewayStatus(gatewayStatus string) PaymentStatus {
	switch gatewayStatus {
	case "success":
		return COMPLETE
	case "failed":
		return FAILED
	case "pending":
		return PENDING
	default:
		return FAILED
	}
}

type PaymentGateway interface {
	Charge(transactionId string, amountInCents int) (*GatewayResponse, error)
	Refund(transactionId string, amountInCents int) (*GatewayResponse, error)
}
type GatewayResponse struct {
	TransactionId string
	Status        string
	RawResponse   map[string]interface{}
}

type StripeGateway struct {
	apiKey string
}

func (s *StripeGateway) Charge(transactionId string, amountInCents int) (*GatewayResponse, error) {
	return &GatewayResponse{
		TransactionId: transactionId,
		Status:        "success",
		RawResponse:   map[string]interface{}{"gateway": "stripe"},
	}, nil
}

func (s *StripeGateway) Refund(transactionId string, amountInCents int) (*GatewayResponse, error) {
	return &GatewayResponse{
		TransactionId: "ref_" + transactionId,
		Status:        "success",
		RawResponse:   map[string]interface{}{"gateway": "stripe"},
	}, nil
}

// ✅ RazorpayGateway - Fixed to match interface
type RazorpayGateway struct {
	apiKey string
}

func (s *RazorpayGateway) Charge(transactionId string, amountInCents int) (*GatewayResponse, error) {
	return &GatewayResponse{
		TransactionId: transactionId,
		Status:        "success",
		RawResponse:   map[string]interface{}{"gateway": "razorpay"},
	}, nil
}

func (s *RazorpayGateway) Refund(transactionId string, amountInCents int) (*GatewayResponse, error) {
	return &GatewayResponse{
		TransactionId: "ref_" + transactionId,
		Status:        "success",
		RawResponse:   map[string]interface{}{"gateway": "razorpay"},
	}, nil
}

type Payment struct {
	Id        string
	BookingId string
	Type      PaymentType
	Amount    float64
	Status    PaymentStatus
	CreatedAt time.Time
}

type PaymentType string

const (
	CREDIT_CARD PaymentType = "CREDIT_CARD"
	DEBIT_CARD  PaymentType = "DEBIT_CARD"
	UPI         PaymentType = "UPI"
	CASH        PaymentType = "CASH"
)

type PaymentStatus string

const (
	PENDING  PaymentStatus = "PENDING"
	COMPLETE PaymentStatus = "COMPLETE"
	FAILED   PaymentStatus = "FAILED"
	REFUNDED PaymentStatus = "REFUNDED"
)

type BookingServiceInterface interface {
	BookVehicle(vehicleId string, userId string, startDate time.Time, endDate time.Time, paymentType PaymentType) error
	GetBookingDetails(bookingId string) (*Booking, error)
	ReturnVehicle(bookingId string, vehicleId string, userId string, actualDateEnd time.Time) error
	CancelBooking(bookingId string, vehicleId string, userId string) error
}
type BookingService struct {
	Bookings       map[string]*Booking     //[booking_id, Booking object]
	PaymentService PaymentServiceInterface //we will initalise paymentService with the specific gateway that we want
	VehicleService VehicleServiceInterface
	mu             sync.RWMutex
}

func (s *BookingService) BookVehicle(vehicleId string, userId string, startDate time.Time, endDate time.Time, paymentType PaymentType) error {
	//fetch vehicle details from vehicle service
	//ShowEstimatedPrice of that vehicle
	//payment at deposit or full amount of this estimatedPrice
	//if payment successful
	//book that vehicle (of course taking care of concurrency , is available = false, )
	price, err := s.VehicleService.CalculateEstimatedPrice(userId, vehicleId, startDate, endDate)
	if err != nil {
		return err
	}

	bookingId := "12332131axzddasx"
	booking := &Booking{
		Id:            bookingId,
		UserId:        userId,
		PaymentId:     "",
		VehicleId:     vehicleId,
		Status:        "",
		CreatedAt:     time.Now(),
		StartDate:     startDate,
		EndDate:       endDate,
		DepositAmount: price * 0.25,
		FinalAmount:   price,
	}
	//lets say for now only pay deposit amount
	payment, err := s.PaymentService.ProcessPayment(bookingId, booking.DepositAmount, paymentType)

	if err != nil || payment.Status == FAILED {
		return err
	}

	booking.PaymentId = payment.Id
	booking.Status = CONFIRMED

	err = s.VehicleService.BookVehicle(userId, vehicleId)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.Bookings[bookingId] = booking
	s.mu.Unlock()
	return nil
}

func (s *BookingService) GetBookingDetails(bookingId string) (*Booking, error) {
	s.mu.RLock()
	booking, ok := s.Bookings[bookingId]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("Booking not found for this booking id")
	}
	return booking, nil
}

func (s *BookingService) ReturnVehicle(bookingId string, vehicleId string, userId string, actualDateEnd time.Time) error {
	booking, err := s.GetBookingDetails(bookingId)
	if err != nil {
		return err
	}
	//calcuate price total including lateFee
	priceTotal, err := s.VehicleService.CalculateNetPriceAfteReturn(booking.UserId, booking.VehicleId, booking.StartDate, booking.EndDate, actualDateEnd)
	if err != nil {
		return err
	}

	//now only pay remaining amount
	alreadyPaid := booking.DepositAmount
	payment, err := s.PaymentService.ProcessPayment(bookingId, priceTotal-alreadyPaid, CREDIT_CARD)

	if err != nil || payment.Status == FAILED {
		return err
	}

	//make again is available = true
	err = s.VehicleService.UnBookVehicle(userId, vehicleId)
	if err != nil {
		return err
	}

	s.mu.Lock()
	booking.Status = COMPLETED
	s.Bookings[bookingId] = booking
	s.mu.Unlock()

	return nil
}

func (s *BookingService) CancelBooking(bookingId string, vehicleId string, userId string) error {
	booking, err := s.GetBookingDetails(bookingId)
	if err != nil {
		return err
	}
	payment, err := s.PaymentService.RefundPayment(bookingId, booking.DepositAmount, CREDIT_CARD)

	if err != nil {
		return err
	}
	if payment.Status == REFUNDED {
		//make again is available = true
		err = s.VehicleService.UnBookVehicle(userId, vehicleId)
		if err != nil {
			return err
		}
		s.mu.Lock()
		booking.Status = CANCELLED
		s.Bookings[bookingId] = booking
		s.mu.Unlock()

		return nil
	}
	return errors.New("Payment not refunded")
}

type Booking struct {
	Id            string
	UserId        string
	PaymentId     string
	VehicleId     string
	Status        BookingStatus
	CreatedAt     time.Time
	StartDate     time.Time
	EndDate       time.Time
	DepositAmount float64
	FinalAmount   float64
}

type BookingStatus string

const (
	CONFIRMED BookingStatus = "CONFIRMED" // Booking created
	COMPLETED BookingStatus = "COMPLETED" // Vehicle returned
	CANCELLED BookingStatus = "CANCELLED" // User cancelled
)

// heper functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateVehicleID() string {
	return "vehicle_" + generateID()
}

func generateTransactionID() string {
	return "txn_" + generateID()
}

func main() {

}
