type PricingContext struct {
    // Vehicle info
    VehicleType VehicleType
    VehicleId   string
    
    // Booking dates
    StartDate time.Time
    EndDate   time.Time
    Days      int
    
    // Time-based factors
    IsWeekend    bool
    IsPeakSeason bool
    IsHoliday    bool
    
    // Demand factors
    CurrentDemand    float64  // 1.0 = normal, 2.0 = 2x demand
    VehiclesAvailable int     // Fewer available = higher price
    
    // User factors
    UserId           string
    IsLoyaltyMember  bool
    PreviousBookings int
    
    // Discounts
    PromoCode    string
    CorporateId  string
    
    // Location factors
    PickupLocation string
    DropOffLocation string
    IsDifferentLocation bool
    
    // Add-ons
    NeedsInsurance bool
    Equipment      []string
}

// Usage
func (s *VehicleService) CalculatePrice(ctx PricingContext) (float64, error) {
    strategy := s.factory.GetStrategy(ctx.VehicleType)
    baseRate := strategy.GetBaseRate()
    
    price := baseRate * float64(ctx.Days)
    
    // Apply weekend surge
    if ctx.IsWeekend {
        price *= 1.5
    }
    
    // Apply peak season
    if ctx.IsPeakSeason {
        price *= 1.3
    }
    
    // Apply demand multiplier
    price *= ctx.CurrentDemand
    
    // Loyalty discount
    if ctx.IsLoyaltyMember {
        price *= 0.9  // 10% off
    }
    
    return price, nil
}

type RidePricingContext struct {
    // Trip info
    StartLocation  Location
    EndLocation    Location
    Distance       float64  // km
    EstimatedTime  int      // minutes
    
    // Time factors
    CurrentTime    time.Time
    IsRushHour     bool
    IsNightTime    bool     // 10pm-6am
    
    // Demand factors
    SurgeMultiplier float64 // 1.0 to 5.0
    DriversNearby   int
    RidesInArea     int
    
    // Weather
    WeatherCondition string  // "rainy", "snowy", "clear"
    Temperature      int
    
    // User factors
    UserId          string
    UserRating      float64
    TotalRides      int
    IsPremiumMember bool
    
    // Ride type
    ServiceLevel    string  // "UberX", "UberXL", "UberBlack"
    
    // Special requests
    IsScheduled     bool
    StopCount       int     // Additional stops
}

func CalculateRidePrice(ctx RidePricingContext) float64 {
    basePrice := ctx.Distance * 2.0  // $2 per km
    
    // Time-based pricing
    if ctx.IsRushHour {
        basePrice *= 1.5
    }
    if ctx.IsNightTime {
        basePrice *= 1.2
    }
    
    // Surge pricing
    basePrice *= ctx.SurgeMultiplier
    
    // Weather premium
    if ctx.WeatherCondition == "rainy" {
        basePrice *= 1.3
    }
    
    // Service level
    serviceLevelMultiplier := map[string]float64{
        "UberX":     1.0,
        "UberXL":    1.5,
        "UberBlack": 2.0,
    }
    basePrice *= serviceLevelMultiplier[ctx.ServiceLevel]
    
    return basePrice
}



type RoomPricingContext struct {
    // Property info
    PropertyId      string
    RoomType        string
    MaxGuests       int
    
    // Booking dates
    CheckInDate     time.Time
    CheckOutDate    time.Time
    Nights          int
    
    // Demand factors
    OccupancyRate   float64  // 0.0 to 1.0
    LocalEvents     []string // ["concert", "festival"]
    IsPeakSeason    bool
    IsWeekend       bool
    
    // Booking factors
    BookingLeadTime int      // Days in advance
    IsLastMinute    bool     // <48 hours
    
    // User factors
    UserId          string
    IsSuperhost     bool
    GuestRating     float64
    BookingHistory  int
    
    // Pricing factors
    BaseRate        float64
    CleaningFee     float64
    ServiceFee      float64
    
    // Discounts
    WeeklyDiscount  float64  // 7+ nights
    MonthlyDiscount float64  // 28+ nights
    PromoCode       string
    
    // Restrictions
    MinStay         int
    MaxStay         int
    CancellationPolicy string
}


type ProductPricingContext struct {
    // Product info
    ProductId       string
    BasePrice       float64
    Category        string
    
    // Inventory
    StockLevel      int
    IsLowStock      bool    // < 10 items
    
    // Time factors
    CurrentTime     time.Time
    IsFlashSale     bool
    IsCyberMonday   bool
    IsBlackFriday   bool
    
    // User factors
    UserId          string
    UserTier        string  // "bronze", "silver", "gold"
    PurchaseHistory int
    IsNewCustomer   bool
    
    // Location
    Country         string
    Currency        string
    ShippingZone    string
    
    // Discounts
    CouponCode      string
    LoyaltyPoints   int
    ReferralDiscount float64
    
    // Bundling
    IsPartOfBundle  bool
    BundleDiscount  float64
}


type DeliveryPricingContext struct {
    // Order info
    OrderValue      float64
    ItemCount       int
    RestaurantId    string
    
    // Delivery info
    DeliveryDistance float64
    EstimatedTime    int
    
    // Time factors
    CurrentTime      time.Time
    IsMealTime       bool  // Lunch/dinner rush
    IsLateNight      bool
    
    // Demand factors
    DriversAvailable int
    OrdersInQueue    int
    SurgeMultiplier  float64
    
    // Weather
    Weather          string
    
    // User factors
    UserId           string
    SubscriptionTier string  // "DashPass"
    TipAmount        float64
    
    // Restaurant factors
    RestaurantBusy   bool
    PreparationTime  int
}


🎯 BUILDER PATTERN FOR COMPLEX CONTEXTS
When context has many fields, use Builder Pattern:


// Builder for complex contexts
type PricingContextBuilder struct{
	ctx PricingContext
}


func NewPricingContext() *PricingContextBuilder{
	return &PricingContextBuilder{
		// Set sensible defaults
		CurrentDemand: 1.0,
		IsWeekend:     false,
		IsPeakSeason:  false,
	}
}

func(s *PricingContextBuilder) WithVehicleType(vType VehicleType) *PricingContextBuilder{
	b.ctx.VehicleType = vType
}


func (b *PricingContextBuilder) WithDates(start, end time.Time) *PricingContextBuilder {
    b.ctx.StartDate = start
    b.ctx.EndDate = end
    b.ctx.Days = int(end.Sub(start).Hours() / 24)
    b.ctx.IsWeekend = isWeekend(start, end)
    return b
}

func (b *PricingContextBuilder) WithDemand(demand float64) *PricingContextBuilder {
    b.ctx.CurrentDemand = demand
    return b
}

func (b *PricingContextBuilder) WithPromoCode(code string) *PricingContextBuilder {
    b.ctx.PromoCode = code
    return b
}

func (b *PricingContextBuilder) Build() PricingContext {
    return b.ctx
}

// Usage with builder
ctx := NewPricingContext().
    WithVehicleType(Car).
    WithDates(start, end).
    WithDemand(1.5).
    WithPromoCode("SUMMER20").
    Build()

price := CalculatePrice(ctx)