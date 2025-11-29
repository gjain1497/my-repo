// 📋 Options for Next:
// 1. Complete VRS (Remaining Requirements)
// From your original requirements, we haven't covered:
// #RequirementStatus6Notifications (pick-up date, due date reminders)❌7Barcode scanning❌9Vehicle log (track all events)❌10Rental insurance❌11Additional equipment (navigation, child seat)❌12Additional services (roadside assistance, wifi)❌

// 2. Design Patterns

// Builder Pattern (PricingContext with many optional fields)
// Observer Pattern (Notifications)
// Decorator Pattern (Add-ons like insurance, equipment)

// package main

// import (
// 	"errors"
// 	"sync"
// 	"time"
// )

// type Person struct {
// 	id   string
// 	name string
// }

// type AdminService struct {
// }

// func (s *AdminService) AddVehicle() {

// }
// func (s *AdminService) RemoveVehicle() {

// }
// func (s *AdminService) UpdateVehicle() {

// }

// type Admin struct {
// 	Person Person
// 	//some additional properties
// 	VehicleService *VehicleService
// }

// type UserService struct {
// 	users          map[string]User
// 	VehicleService *VehicleService
// }

// func (s *UserService) GetAllVehiclesRentedByUser(userId string) []*Vehicle {
// 	return s.VehicleService.ListVehiclesCurrentlyAssignedToAUser(userId)
// }

// type User struct {
// 	Person        Person
// 	LicenseNumber string
// 	Phone         string
// 	Email         string
// 	//some additional properties
// }

// type VehicleService struct {
// 	Vehicles        map[string]*Vehicle   //(id , Vehicle object)
// 	UserVehicles    map[string][]*Vehicle //(user_id, list of vehicles)
// 	StrategyFactory *PricingStrategyFactory
// 	mu              sync.RWMutex
// }

// type PricingStrategyFactory struct {
// }

// func (f *PricingStrategyFactory) GetStrategy(vehicleType VehicleType) (VehicleTypeInterface, error) {
// 	switch vehicleType {
// 	case Car:
// 		return &CarVehicle{}, nil
// 	case Bike:
// 		return &BikeVehicle{}, nil
// 	case Truck:
// 		return &TruckVehicle{}, nil
// 	default:
// 		return nil, errors.New("unknown vehicle type")
// 	}
// }

// type VehicleTypeInterface interface {
// 	CalculatePricePerDay() (float64, error)
// 	CalculateLateFeePerDay() (float64, error)
// }

// type CarVehicle struct {
// }

// func (c *CarVehicle) CalculatePricePerDay() (float64, error) {
// 	baseRate := 20.0
// 	price := baseRate
// 	return 20, nil
// }
// func (c *CarVehicle) CalculateLateFeePerDay() (float64, error) {
// 	return 2, nil
// }

// type BikeVehicle struct {
// }

// func (c *BikeVehicle) CalculatePricePerDay() (float64, error) {
// 	return 15, nil
// }
// func (c *BikeVehicle) CalculateLateFeePerDay() (float64, error) {
// 	return 1.5, nil
// }

// type TruckVehicle struct {
// }

// func (c *TruckVehicle) CalculatePricePerDay() (float64, error) {
// 	return 10, nil
// }
// func (c *TruckVehicle) CalculateLateFeePerDay() (float64, error) {
// 	return 1.7, nil
// }

// func (s *VehicleService) ListVehicles(userId string, locationId string, vehicleType VehicleType, dateStart time.Time, dateEnd time.Time) ([]*Vehicle, error) {
// 	//anyone can check no need of verifying that user is admin only
// 	var vehicles []*Vehicle
// 	for _, vehicle := range s.Vehicles {
// 		if vehicle.Type == vehicleType && vehicle.IsAvailable {
// 			vehicles = append(vehicles, vehicle)
// 		}
// 	}

// 	return vehicles, nil
// }

// func (s *VehicleService) ListVehiclesCurrentlyAssignedToAUser(userId string) []*Vehicle {
// 	return s.UserVehicles[userId]
// }

// func (s *VehicleService) BookVehicle(userId string, vehicleId string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	vehicle, ok := s.Vehicles[vehicleId]
// 	if !ok {
// 		return errors.New("Vehicle not found")
// 	}

// 	vehicle.IsAvailable = false
// 	s.UserVehicles[userId] = append(s.UserVehicles[userId], vehicle)
// 	return nil
// }

// func (s *VehicleService) UnBookVehicle(userId string, vehicleId string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	vehicle, ok := s.Vehicles[vehicleId]
// 	if !ok {
// 		return errors.New("Vehicle not found")
// 	}

// 	vehicle.IsAvailable = true
// 	currentVehicles := s.ListVehiclesCurrentlyAssignedToAUser(userId)
// 	for i, vehicle := range currentVehicles {
// 		if vehicle.Id == vehicleId {
// 			currentVehicles = append(currentVehicles[:i], currentVehicles[i+1:]...)
// 			break
// 		}
// 	}
// 	s.UserVehicles[userId] = currentVehicles
// 	return nil
// }

// type PricingContext struct {
// 	// Vehicle info
// 	VehicleType VehicleType
// 	VehicleId   string

// 	// Booking dates
// 	StartDate time.Time
// 	EndDate   time.Time
// 	Days      int

// 	// Time-based factors (auto-calculated)
// 	IsWeekend    bool
// 	IsPeakSeason bool

// 	// Demand factors
// 	CurrentDemand float64 // 1.0 = normal, 2.0 = high demand

// 	// User factors
// 	UserId          string
// 	IsLoyaltyMember bool

// 	// Discounts
// 	PromoCode string

// 	// Location
// 	LocationId string
// }

// type PricingContextBuilder struct {
// 	ctx PricingContext
// 	err error
// }

// func NewPricingContext() *PricingContextBuilder {
// 	return &PricingContextBuilder{
// 		ctx: PricingContext{
// 			// Set defaults
// 			CurrentDemand:   1.0,
// 			IsWeekend:       false,
// 			IsPeakSeason:    false,
// 			IsLoyaltyMember: false,
// 		},
// 	}
// }

// func (b *PricingContextBuilder) WithVehicleType(vType VehicleType) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	if vType == "" {
// 		b.err = errors.New("vehicle type cannot be empty")
// 		return b
// 	}
// 	b.ctx.VehicleType = vType
// }

// func (b *PricingContextBuilder) WithVehicleId(vehicleId string) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	if vehicleId == "" {
// 		b.err = errors.New("vehicle ID cannot be empty")
// 		return b
// 	}
// 	b.ctx.VehicleId = vehicleId
// 	return b
// }

// func (b *PricingContextBuilder) WithDates(start, end time.Time) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}

// 	if end.Before(start) {
// 		b.err = errors.New("end date cannot be before start date")
// 		return b
// 	}

// 	b.ctx.StartDate = start
// 	b.ctx.EndDate = end

// 	// Auto-calculate days
// 	days := int(end.Sub(start).Hours() / 24)
// 	if days == 0 {
// 		days = 1 // Minimum 1 day
// 	}
// 	b.ctx.Days = days

// 	// Auto-detect weekend
// 	b.ctx.IsWeekend = b.hasWeekend(start, end)

// 	// Auto-detect peak season (June-August)
// 	b.ctx.IsPeakSeason = b.isPeakSeason(start)

// 	return b
// }
// func (b *PricingContextBuilder) WithUserId(userId string) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	if userId == "" {
// 		b.err = errors.New("user ID cannot be empty")
// 		return b
// 	}
// 	b.ctx.UserId = userId
// 	return b
// }

// // Optional fields

// func (b *PricingContextBuilder) WithDemand(demand float64) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	if demand < 0.5 || demand > 5.0 {
// 		b.err = errors.New("demand must be between 0.5 and 5.0")
// 		return b
// 	}
// 	b.ctx.CurrentDemand = demand
// 	return b
// }

// func (b *PricingContextBuilder) WithLoyaltyMember(isLoyalty bool) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	b.ctx.IsLoyaltyMember = isLoyalty
// 	return b
// }

// func (b *PricingContextBuilder) WithPromoCode(code string) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	b.ctx.PromoCode = code
// 	return b
// }

// func (b *PricingContextBuilder) WithLocation(locationId string) *PricingContextBuilder {
// 	if b.err != nil {
// 		return b
// 	}
// 	b.ctx.LocationId = locationId
// 	return b
// }

// // Helper methods

// func (b *PricingContextBuilder) hasWeekend(start, end time.Time) bool {
// 	// Check if any day in range is weekend
// 	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
// 		weekday := d.Weekday()
// 		if weekday == time.Saturday || weekday == time.Sunday {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (b *PricingContextBuilder) isPeakSeason(date time.Time) bool {
// 	month := date.Month()
// 	// Summer: June, July, August
// 	return month >= 6 && month <= 8
// }

// // Build
// func (b *PricingContextBuilder) Build() (PricingContext, error) {
// 	// Check for errors during building
// 	if b.err != nil {
// 		return PricingContext{}, b.err
// 	}
// 	// Validate required fields
// 	if b.ctx.VehicleType == "" {
// 		return PricingContext{}, errors.New("vehicle type is required")
// 	}
// 	if b.ctx.VehicleId == "" {
// 		return PricingContext{}, errors.New("vehicle ID is required")
// 	}
// 	if b.ctx.UserId == "" {
// 		return PricingContext{}, errors.New("user ID is required")
// 	}
// 	if b.ctx.Days <= 0 {
// 		return PricingContext{}, errors.New("days must be positive")
// 	}
// 	return b.ctx, nil
// }

// func (s *VehicleService) CalculateEstimatedPrice(userId string, vehicleId string, dateStart time.Time, dateEnd time.Time) (float64, error) {
// 	//anyone can check no need of verifying that user is admin only
// 	//we know price for one day

// 	//we use some strategy for different prices
// 	//calculateTimeDifference -> pricePerday * (number of days)
// 	s.mu.RLock() //Read lock (allows concurrent reads)
// 	vehicle, ok := s.Vehicles[vehicleId]
// 	s.mu.RUnlock()
// 	if !ok {
// 		return 0, errors.New("Vehicle not found")
// 	}
// 	timeDifference := dateEnd.Sub(dateStart)
// 	days := (timeDifference.Hours() / 24)
// 	if days <= 0 {
// 		return 0, errors.New("invalid date range")
// 	}

// 	//Get strategy using factory
// 	strategy, err := s.StrategyFactory.GetStrategy(vehicle.Type)
// 	if err != nil {
// 		return 0, err
// 	}

// 	ctx, err := NewPricingContext().
// 		WithVehicleId(vehicleId).
// 		WithUserId(userId).
// 		WithDates(dateStart, dateEnd).
// 		Build()

// 	if err != nil {
// 		return 0, err
// 	}

// 	//based on vehicleType different strategies
// 	pricePerDay, err := strategy.CalculatePricePerDay(ctx)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Calculate total
// 	return pricePerDay * float64(days), nil
// }

// func (s *VehicleService) CalculateLateFee(vehicleId string, dateStart time.Time, dateEnd time.Time) (float64, error) {
// 	//anyone can check no need of verifying that user is admin only
// 	//we know price for one day

// 	//we use some strategy for different prices
// 	//calculateTimeDifference -> pricePerday * (number of days)
// 	s.mu.RLock() //Read lock (allows concurrent reads)
// 	vehicle, ok := s.Vehicles[vehicleId]
// 	s.mu.RUnlock()
// 	if !ok {
// 		return 0, errors.New("Vehicle not found")
// 	}
// 	timeDifference := dateEnd.Sub(dateStart)
// 	days := (timeDifference.Hours() / 24)
// 	if days <= 0 {
// 		return 0, errors.New("invalid date range")
// 	}

// 	//Get strategy using factory
// 	strategy, err := s.StrategyFactory.GetStrategy(vehicle.Type)
// 	if err != nil {
// 		return 0, err
// 	}

// 	//based on vehicleType different strategies
// 	lateFeePerDay, err := strategy.CalculateLateFeePerDay()
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Calculate total
// 	return lateFeePerDay * float64(days), nil
// }

// func (s *VehicleService) CalculateNetPriceAfteReturn(userId string, vehicleId string, dateStart time.Time, dateEnd time.Time, actualDateEnd time.Time) (float64, error) {
// 	//anyone can check no need of verifying that user is admin only
// 	//we know price for one day
// 	price, err := s.CalculateEstimatedPrice(userId, vehicleId, dateStart, dateEnd)
// 	if err != nil {
// 		return 0, err
// 	}
// 	var lateFee float64
// 	if actualDateEnd.After(dateEnd) {
// 		lateFee, err = s.CalculateLateFee(vehicleId, dateEnd, actualDateEnd)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	return price + lateFee, nil
// }

// func (s *VehicleService) GetVehicle(vehicleId string) (*Vehicle, error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	vehicle, ok := s.Vehicles[vehicleId]
// 	if !ok {
// 		return nil, errors.New("vehicle not found")
// 	}
// 	return vehicle, nil
// }

// func (v *VehicleService) AddVehicle(personId string, model string, brand string) error {
// 	//will verify that person in admin maybe through the access token/jwt
// 	return nil
// }

// func (v *VehicleService) RemoveVehicle(personId string, vehicleId string) error {
// 	//will verify that person in admin maybe through the access token/jwt
// 	return nil
// }

// func (v *VehicleService) UpdateVehicle(personId string, vehicleId string) error {
// 	//will verify that person in admin maybe through the access token/jwt
// 	return nil
// }

// type Vehicle struct {
// 	Id              string
// 	Type            VehicleType
// 	Model           string // "Honda City", "Royal Enfield"
// 	Brand           string // "Honda", "Royal Enfield"
// 	RegistrationNum string // "DL-01-AB-1234"
// 	IsAvailable     bool
// 	Location        Location
// 	PricePerDay     float64
// 	LateFeePerDay   float64 //late fee if vehicle returned after end date
// }

// type VehicleType string

// const (
// 	Car   VehicleType = "CAR"
// 	Truck VehicleType = "TRUCK"
// 	Bike  VehicleType = "BIKE"
// )

// type VehicleStatus string

// const (
// 	Available VehicleStatus = "AVAILABLE"
// 	Booked    VehicleStatus = "BOOKED"
// )

// type Location struct {
// 	City      string
// 	State     string
// 	ZipCode   string
// 	Country   string
// 	IsDefault bool
// }

// type PaymentService struct {
// 	PaymentGateway PaymentGateway
// 	Payments       map[string]*Payment
// 	mu             sync.RWMutex
// }

// func (s *PaymentService) ProcessPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
// 	//do some additional work
// 	if bookingId == "" {
// 		return nil, errors.New("booking ID required")
// 	}
// 	if amount <= 0 {
// 		return nil, errors.New("invalid amount")
// 	}
// 	payment, err := s.PaymentGateway.ProcessPayment(bookingId, amount, paymentType)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s.mu.Lock()
// 	s.Payments[payment.Id] = payment
// 	s.mu.Unlock()

// 	return payment, nil
// }

// func (s *PaymentService) RefundPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
// 	//do some additional work
// 	if bookingId == "" {
// 		return nil, errors.New("booking ID required")
// 	}
// 	if amount <= 0 {
// 		return nil, errors.New("invalid amount")
// 	}
// 	return s.PaymentGateway.RefundPayment(bookingId, amount, paymentType)
// }

// type PaymentGateway interface {
// 	ProcessPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error)
// 	RefundPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error)
// }

// type StripeGateway struct {
// 	apiKey string
// }

// func (s *StripeGateway) ProcessPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
// 	return &Payment{
// 		Id:     "123",
// 		Status: COMPLETE,
// 	}, nil
// }
// func (s *StripeGateway) RefundPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
// 	return &Payment{
// 		Id:     "123",
// 		Status: REFUNDED,
// 	}, nil
// }

// type RazorpayGateway struct {
// 	apiKey string
// }

// func (s *RazorpayGateway) ProcessPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
// 	return &Payment{
// 		Id:     "123",
// 		Status: COMPLETE,
// 	}, nil
// }
// func (s *RazorpayGateway) RefundPayment(bookingId string, amount float64, paymentType PaymentType) (*Payment, error) {
// 	return &Payment{
// 		Id:     "123",
// 		Status: REFUNDED,
// 	}, nil
// }

// type Payment struct {
// 	Id        string
// 	Type      PaymentType
// 	Amount    float64
// 	Status    PaymentStatus
// 	CreatedAt time.Time
// }

// type PaymentType string

// const (
// 	CREDIT_CARD PaymentType = "CREDIT_CARD"
// 	DEBIT_CARD  PaymentType = "DEBIT_CARD"
// 	UPI         PaymentType = "UPI"
// 	CASH        PaymentType = "CASH"
// )

// type PaymentStatus string

// const (
// 	PENDING  PaymentStatus = "PENDING"
// 	COMPLETE PaymentStatus = "COMPLETE"
// 	FAILED   PaymentStatus = "FAILED"
// 	REFUNDED PaymentStatus = "REFUNDED"
// )

// type BookingService struct {
// 	Bookings       map[string]*Booking //[booking_id, Booking object]
// 	PaymentService *PaymentService     //we will initalise paymentService with the specific gateway that we want
// 	VehicleService *VehicleService
// 	mu             sync.RWMutex
// }

// func (s *BookingService) BookVehicle(vehicleId string, userId string, startDate time.Time, endDate time.Time, paymentType PaymentType) error {
// 	//fetch vehicle details from vehicle service
// 	//ShowEstimatedPrice of that vehicle
// 	//payment at deposit or full amount of this estimatedPrice
// 	//if payment successful
// 	//book that vehicle (of course taking care of concurrency , is available = false, )
// 	price, err := s.VehicleService.CalculateEstimatedPrice(userId, vehicleId, startDate, endDate)
// 	if err != nil {
// 		return err
// 	}

// 	bookingId := "12332131axzddasx"
// 	booking := &Booking{
// 		Id:            bookingId,
// 		UserId:        userId,
// 		PaymentId:     "",
// 		VehicleId:     vehicleId,
// 		Status:        "",
// 		CreatedAt:     time.Now(),
// 		StartDate:     startDate,
// 		EndDate:       endDate,
// 		DepositAmount: price * 0.25,
// 		FinalAmount:   price,
// 	}
// 	//lets say for now only pay deposit amount
// 	payment, err := s.PaymentService.ProcessPayment(bookingId, booking.DepositAmount, paymentType)

// 	if err != nil || payment.Status == FAILED {
// 		return err
// 	}

// 	booking.PaymentId = payment.Id
// 	booking.Status = CONFIRMED

// 	err = s.VehicleService.BookVehicle(userId, vehicleId)
// 	if err != nil {
// 		return err
// 	}
// 	s.mu.Lock()
// 	s.Bookings[bookingId] = booking
// 	s.mu.Unlock()
// 	return nil
// }

// func (s *BookingService) GetBookingDetails(bookingId string) (*Booking, error) {
// 	s.mu.RLock()
// 	booking, ok := s.Bookings[bookingId]
// 	s.mu.RUnlock()
// 	if !ok {
// 		return nil, errors.New("Booking not found for this booking id")
// 	}
// 	return booking, nil
// }

// func (s *BookingService) ReturnVehicle(bookingId string, vehicleId string, userId string, actualDateEnd time.Time) error {
// 	booking, err := s.GetBookingDetails(bookingId)
// 	if err != nil {
// 		return err
// 	}
// 	//calcuate price total including lateFee
// 	priceTotal, err := s.VehicleService.CalculateNetPriceAfteReturn(booking.UserId, booking.VehicleId, booking.StartDate, booking.EndDate, actualDateEnd)
// 	if err != nil {
// 		return err
// 	}

// 	//now only pay remaining amount
// 	alreadyPaid := booking.DepositAmount
// 	payment, err := s.PaymentService.ProcessPayment(bookingId, priceTotal-alreadyPaid, CREDIT_CARD)

// 	if err != nil || payment.Status == FAILED {
// 		return err
// 	}

// 	//make again is available = true
// 	err = s.VehicleService.UnBookVehicle(userId, vehicleId)
// 	if err != nil {
// 		return err
// 	}

// 	s.mu.Lock()
// 	booking.Status = COMPLETED
// 	s.Bookings[bookingId] = booking
// 	s.mu.Unlock()

// 	return nil
// }

// func (s *BookingService) CancelBooking(bookingId string, vehicleId string, userId string) error {
// 	booking, err := s.GetBookingDetails(bookingId)
// 	if err != nil {
// 		return err
// 	}
// 	payment, err := s.PaymentService.RefundPayment(bookingId, booking.DepositAmount, CREDIT_CARD)

// 	if err != nil {
// 		return err
// 	}
// 	if payment.Status == REFUNDED {
// 		//make again is available = true
// 		err = s.VehicleService.UnBookVehicle(userId, vehicleId)
// 		if err != nil {
// 			return err
// 		}
// 		s.mu.Lock()
// 		booking.Status = CANCELLED
// 		s.Bookings[bookingId] = booking
// 		s.mu.Unlock()

// 		return nil
// 	}
// 	return errors.New("Payment not refunded")
// }

// type Booking struct {
// 	Id            string
// 	UserId        string
// 	PaymentId     string
// 	VehicleId     string
// 	Status        BookingStatus
// 	CreatedAt     time.Time
// 	StartDate     time.Time
// 	EndDate       time.Time
// 	DepositAmount float64
// 	FinalAmount   float64
// }

// type BookingStatus string

// const (
// 	CONFIRMED BookingStatus = "CONFIRMED" // Booking created
// 	COMPLETED BookingStatus = "COMPLETED" // Vehicle returned
// 	CANCELLED BookingStatus = "CANCELLED" // User cancelled
// )

// func main() {

// }
