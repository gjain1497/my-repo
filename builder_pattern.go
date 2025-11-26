FOR PricingContext: USE BUILDER! ✅✅✅
Why:

✅ Many fields (10+)
✅ Some optional (PromoCode, Demand, etc.)
✅ Need validation (dates, demand range)
✅ Have derived fields (IsWeekend, Days)
✅ Will grow (more pricing factors)


// Your PricingContext
type PricingContext struct {
    VehicleType      VehicleType
    StartDate        time.Time
    EndDate          time.Time
    Days             int
    IsWeekend        bool
    IsPeakSeason     bool
    CurrentDemand    float64
    UserId           string
    PromoCode        string
    IsLoyaltyMember  bool
    LocationId       string
}

// ❌ PROBLEM 1: Too many fields to set manually
ctx := PricingContext{
    VehicleType:     Car,
    StartDate:       start,
    EndDate:         end,
    Days:            3,
    IsWeekend:       true,
    IsPeakSeason:    false,
    CurrentDemand:   1.5,
    UserId:          "user123",
    PromoCode:       "SUMMER20",
    IsLoyaltyMember: true,
    LocationId:      "NYC",
}

// ❌ PROBLEM 2: Easy to make mistakes
ctx := PricingContext{
    VehicleType:   Car,
    StartDate:     start,
    EndDate:       end,
    Days:          3,
    IsWeekend:     true,
    // ❌ Forgot to set CurrentDemand! Will be 0.0!
    // ❌ Forgot PromoCode!
}

// ❌ PROBLEM 3: Order matters, easy to mix up
ctx := PricingContext{
    start,           // ❌ Which field is this?
    end,             // ❌ Which field is this?
    Car,             // ❌ Confusing!
    "user123",
    3,
    // ...
}

// ❌ PROBLEM 4: Can't have default values
ctx := PricingContext{
    VehicleType:   Car,
    StartDate:     start,
    EndDate:       end,
    Days:          3,
    CurrentDemand: 1.0,  // ❌ Have to manually set default!
}

// ❌ PROBLEM 5: Hard to create variations
baseCtx := PricingContext{...}  // Set everything
weekendCtx := baseCtx
weekendCtx.IsWeekend = true     // Modify one field

// ❌ PROBLEM 6: No validation during construction
ctx := PricingContext{
    Days: -5,  // ❌ Invalid! But no error!
}


// ❌ Builder avoids “breaking 50 files”
💥 Problem: Adding a new parameter normally breaks code

If you use a constructor:

func NewPricingContext(v VehicleType, start, end time.Time, demand float64, promo string)


And tomorrow you add:

WeatherType string


Now your constructor must change:

func NewPricingContext(v VehicleType, start, end time.Time, demand float64, promo string, weather string)


❌ Every single call in 50 files must now pass this parameter
Even if they don’t care about weather.
Even if they want the default value.

This forces a rewrite across the entire codebase.

💚 Builder solves this EXACT problem

If today you add a new field:

WeatherType string


You do not change the builder usage anywhere.
You just update the builder:

Step 1: Add the field in the struct
type PricingContext struct {
    WeatherType string
}


No code breaks.
Because struct fields are not positional → old initializations still compile.

Step 2: Add optional builder method
func (b *PricingContextBuilder) WithWeatherType(w string) *PricingContextBuilder {
    b.ctx.WeatherType = w
    return b
}


That’s it.

🧠 Now check how usage behaves

This code from 50 files:

ctx, _ := NewPricingContext().
    WithVehicleType(Car).
    WithDates(start, end).
    Build()


Compile?
👉 YES
We didn’t add anything compulsory.

WeatherType defaults to "" (zero value).
No caller is forced to update their code.


🟩 Why this is “non-breaking change”?

Because:

✔ 1. Builder method WithWeatherType() is optional

If the caller doesn’t need it:

they don’t call it

nothing changes

compile continues

✔ 2. You didn’t change function signatures

NewPricingContext() signature stays the same.

Unlike constructor which becomes:

NewPricingContext(... old params ..., newParam)

Builder never forces this.

✔ 3. You didn’t change build function signature

Build() stays the same.

✔ 4. Old builder chains still work unchanged

Because builder methods are not positional arguments.


// Builder struct
type PricingContextBuilder struct {
    ctx PricingContext
}

// Constructor with sensible defaults
func NewPricingContext() *PricingContextBuilder {
    return &PricingContextBuilder{
        ctx: PricingContext{
            CurrentDemand: 1.0,      // ✅ Default!
            IsWeekend:     false,    // ✅ Default!
            IsPeakSeason:  false,    // ✅ Default!
        },
    }
}

// Fluent methods (chainable)
func (b *PricingContextBuilder) WithVehicleType(vType VehicleType) *PricingContextBuilder {
    b.ctx.VehicleType = vType
    return b  // ✅ Return self for chaining!
}

func (b *PricingContextBuilder) WithDates(start, end time.Time) *PricingContextBuilder {
    b.ctx.StartDate = start
    b.ctx.EndDate = end
    b.ctx.Days = int(end.Sub(start).Hours() / 24)
    
    // ✅ Auto-calculate derived fields!
    b.ctx.IsWeekend = isWeekend(start, end)
    
    return b
}

func (b *PricingContextBuilder) WithUser(userId string, isLoyalty bool) *PricingContextBuilder {
    b.ctx.UserId = userId
    b.ctx.IsLoyaltyMember = isLoyalty
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

func (b *PricingContextBuilder) WithPeakSeason(isPeak bool) *PricingContextBuilder {
    b.ctx.IsPeakSeason = isPeak
    return b
}

// Build with validation
func (b *PricingContextBuilder) Build() (PricingContext, error) {
    // ✅ Validate before returning
    if b.ctx.Days <= 0 {
        return PricingContext{}, errors.New("invalid days")
    }
    if b.ctx.EndDate.Before(b.ctx.StartDate) {
        return PricingContext{}, errors.New("end date before start date")
    }
    if b.ctx.VehicleType == "" {
        return PricingContext{}, errors.New("vehicle type required")
    }
    
    return b.ctx, nil
}

// ✅ USAGE: Clean, readable, self-documenting!
ctx, err := NewPricingContext().
    WithVehicleType(Car).
    WithDates(start, end).
    WithUser("user123", true).
    WithPromoCode("SUMMER20").
    WithDemand(1.5).
    Build()

if err != nil {
    return err
}

// Use context
price := CalculatePrice(ctx)


🎯 WHEN TO USE BUILDER PATTERN:
✅ Use Builder When:

1. Many parameters (>4)

// ❌ Without builder: Too many params
func CreateBooking(vehicleId, userId string, start, end time.Time, 
	insurance bool, equipment []string, promoCode string,
	isPeakSeason bool, demand float64) Booking

// ✅ With builder: Clean
func CreateBooking(ctx PricingContext) Booking


2. Optional parameters

// ✅ Some bookings have promo, some don't
ctx1, _ := NewPricingContext().
    WithVehicleType(Car).
    WithDates(start, end).
    WithUserId("user1").
    Build()  // ✅ No promo

ctx2, _ := NewPricingContext().
    WithVehicleType(Car).
    WithDates(start, end).
    WithUserId("user2").
    WithPromoCode("SAVE20").  // ✅ With promo
    Build()


3. Need validation during construction
ctx, err := NewPricingContext().
    WithDemand(10.0).  // ❌ Invalid! Must be 0.5-5.0
    Build()

if err != nil {
    fmt.Println("Error:", err)
    // Output: Error: demand must be between 0.5 and 5.0
}


4. Derived fields need calculation
ctx, _ := NewPricingContext().
    WithDates(start, end).  // ✅ Auto-calculates Days and IsWeekend!
    Build()

// You don't have to manually calculate:
// - Days
// - IsWeekend
// - Other derived fields


5. Want immutable objects

// Builder creates new context, original unchanged
ctx1, _ := NewPricingContext().WithVehicleType(Car).Build()
ctx2, _ := NewPricingContext().WithVehicleType(Bike).Build()
// ctx1 and ctx2 are independent


❌ Dont Use Builder When:

Simple objects (1-3 fields)

go// ❌ Overkill
type Point struct {
    X, Y int
}

// ✅ Just use struct literal
p := Point{X: 10, Y: 20}


All fields are required

go// If ALL fields must be set, constructor with params is fine
func NewUser(id, name, email string) User {
    return User{Id: id, Name: name, Email: email}
}


Performance critical code

go// Builder has slight overhead (chaining, validation)
// For hot paths, direct struct creation might be better
