package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Design a parking lot system that can:

Park vehicles (different types: Car, Motorcycle, Truck)
Assign parking slots based on availability
Track which slots are occupied/free
Calculate parking fees when vehicle exits
Support multiple floors
Handle different slot sizes (Compact, Large, Handicapped)
*/

/*
📋 Requirements Analysis
Functional Requirements

Vehicle Entry: Assign a slot when vehicle enters
Vehicle Exit: Calculate bill, free the slot
Slot Management: Track availability per floor, per slot type
Multiple Vehicle Types: Car, Bike, Truck (different sizes)
Pricing: Hourly rates (may vary by vehicle type)

Non-Functional Requirements

Thread-safe: Multiple vehicles entering/exiting concurrently
Extensible: Easy to add new vehicle types, pricing strategies
Efficient: Fast slot lookup (O(1) or O(log n))
*/

//questions to think about

//what happens when a vehicle enters??
//find available slot -> assign -> generate ticket

//what happens when a car exists??
//find ticket -> calcualte fee -> free slot

//what if parking is full
//return error //waitlist

//pricing rules
//hourly //flat //different for diff vehicles

//core entities
//vehicle
//slot
//floor
//ticket
//parking lot itself

type VehicleType string

const (
	Truck VehicleType = "TRUCK"
	Car   VehicleType = "CAR"
	Bike  VehicleType = "BIKE"
)

type Vehicle struct {
	VehicleType  VehicleType
	LicensePlate string
}

type SlotType string

const (
	Compact     SlotType = "COMPACT"
	Large       SlotType = "LARGE"
	Handicapped SlotType = "HANDICAPPED"
	Electric    SlotType = "ELECTRIC"
)

type SlotStatus string

const (
	Available SlotStatus = "AVAILABLE"
	Occupied  SlotStatus = "OCCUPIED"
)

type Slot struct {
	ID            string
	Type          SlotType
	Status        SlotStatus
	ParkedVehicle *Vehicle
	Floor         int
	mutex         sync.RWMutex
}

type Floor struct {
	FloorNumber int
	Slots       map[SlotType][]*Slot //grouped by type
}

type Ticket struct {
	TicketId        string
	VehicleAssigned *Vehicle
	SlotAssigned    *Slot
	EntryTime       time.Time  //cannot be nil in any case
	ExitTime        *time.Time //can be nil so a pointer
}

type ParkingLot struct {
	Name             string
	Floors           []*Floor
	Gates            map[string]*Gate //gate_id -> gate
	TicketCounter    int
	PricingStrategy  PricingStrategy
	ActiveTickets    map[string]*Ticket //(ticket_id -> Ticket) //for easy lookup of ticket details on exit to calcuate charge
	MaxCapacity      int                // Total slots across all floors
	CurrentOccupancy int                // How many are occupied
	mutex            sync.RWMutex
}

var (
	parkingLotInstance *ParkingLot
	once               sync.Once
)

// we need singelton pattern for parkingLot
func GetParkingLot(name string, numFloors int, maxCapacity int, startegy PricingStrategy) *ParkingLot {
	once.Do(func() {
		parkingLotInstance = &ParkingLot{
			Name:            name,
			Floors:          make([]*Floor, numFloors),
			Gates:           make(map[string]*Gate),
			ActiveTickets:   map[string]*Ticket{},
			TicketCounter:   0,
			PricingStrategy: startegy,
			MaxCapacity:     maxCapacity,
		}
		// Initialize floors
		for i := 0; i < numFloors; i++ {
			parkingLotInstance.Floors[i] = &Floor{
				FloorNumber: i,
				Slots:       make(map[SlotType][]*Slot),
			}
		}
	})
	return parkingLotInstance
}

func (f *Floor) FindSlot(slotType SlotType) *Slot {
	if slots, ok := f.Slots[slotType]; ok { // ← O(1) map lookup!
		for _, slot := range slots { // ← O(S) where S = slots of this type
			if slot.isAvailable() { // ← O(1)
				return slot
			}
		}
	}
	return nil
}

//add entry exit gates

//multiple entry gates,
// gate type: entry only exit only both
// track which gate was used: store in ticket
// gate status: open/closed/under maintenance

type GateType string

const (
	EntryGate GateType = "ENTRY"
	ExitGate  GateType = "EXIT"
	BothGate  GateType = "BOTH"
)

type GateStatus string

const (
	GateOpen        GateStatus = "OPEN"
	GateClosed      GateStatus = "CLOSED"
	GateMaintenance GateStatus = "MAINTENANCE"
)

type Gate struct {
	Id     string
	Type   GateType
	Status GateStatus
	Floor  int
}

func (g *Gate) CanProcessEntry() bool {
	return g.Status == GateOpen && (g.Type == EntryGate || g.Type == BothGate)
}

func (g *Gate) CanProcessExit() bool {
	return g.Status == GateOpen && (g.Type == ExitGate || g.Type == BothGate)
}

//old flow -> parkinglot.ParkVechicle(vehicle)
//finds slot -> creates ticket

//old flow -> parkinglot.ProcessVehicleEntry(vehicle, gateID) -> checks if gate can process entry
//new flow -> gate.ParkVehicle(vehicle) -> finds slot -> creates ticket

func (p *ParkingLot) AddGate(gate *Gate) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Gates[gate.Id] = gate
}

func (p *ParkingLot) ProcessVehicleEntry(v *Vehicle, gateID string) (*Ticket, error) {
	//validate gate
	gate, exists := p.Gates[gateID]
	if !exists {
		return nil, fmt.Errorf("gate not found: %s", gateID)
	}
	if !gate.CanProcessEntry() {
		return nil, fmt.Errorf("gate %s cannot process entry (status: %s)", gateID, gate.Status)
	}

	//rest is same as before (parkvehicle)
	return p.ParkVehicle(v)
}

func (p *ParkingLot) ParkVehicle(v *Vehicle) (*Ticket, error) {
	//add check
	p.mutex.RLock()
	isFull := p.CurrentOccupancy >= p.MaxCapacity
	p.mutex.RUnlock()

	if isFull {
		return nil, fmt.Errorf("parking lot is full")
	}

	// Step 1: Get required slot type - O(1)
	slotType := v.GetRequiredSlotType()

	//loop through all floors to get slot for this
	floors := p.Floors
	var slot *Slot = nil

	for _, floor := range floors { // O(F) where F = number of floors
		slot = floor.FindSlot(slotType)
		if slot != nil {
			break
		}
	}

	if slot == nil {
		return nil, fmt.Errorf("slot not found for given vehicle type")
	}

	//slot found so assign the vehicle to this slot and make it occupied
	slot.Occupy(v)

	//Generate ticket
	p.mutex.Lock()
	p.TicketCounter++
	p.CurrentOccupancy++
	ticketID := fmt.Sprintf("Ticket-%d", p.TicketCounter)

	ticket := &Ticket{
		TicketId:        ticketID,
		VehicleAssigned: v,
		SlotAssigned:    slot,
		EntryTime:       time.Now(),
		ExitTime:        nil,
	}

	//assign in active tickets, (ticket_id -> ticket) for easy lookup using ticket_id
	p.ActiveTickets[ticketID] = ticket
	p.mutex.Unlock()

	return ticket, nil
}

func (p *ParkingLot) ProcessVehicleExit(tickeID string, gateID string) (float64, error) {
	//validate gate
	gate, exists := p.Gates[gateID]

	if !exists {
		return 0, fmt.Errorf("gate not found %s", gateID)
	}

	if !gate.CanProcessExit() {
		return 0, fmt.Errorf("gate %s cannot process exit (status: %s)", gateID, gate.Status)
	}

	//rest is same as before (unpark -> calcuate fee, free slot)
	return p.UnparkVehicle(tickeID)
}

func (p *ParkingLot) UnparkVehicle(ticketID string) (float64, error) {
	p.mutex.Lock()
	var ticket *Ticket = nil
	ticket, ok := p.ActiveTickets[ticketID]
	if !ok {
		return 0, fmt.Errorf("Ticket not found")
	}

	//assign exit time
	now := time.Now()
	ticket.ExitTime = &now

	//getSlot
	slot := ticket.SlotAssigned

	//empty this slot (means set status avaialble )
	slot.Free()

	//remove ticket from Activate Tickets
	delete(p.ActiveTickets, ticketID)
	p.CurrentOccupancy--

	p.mutex.Unlock()

	//calculate fee based on startegy
	fee := p.PricingStrategy.CalculateFee(ticket)

	return fee, nil
}

func (p *ParkingLot) DisplayStatus() {
	fmt.Printf("\n==== %v Status ===\n", p.Name)
	fmt.Printf("Occupancy: %d/%d slots\n", p.CurrentOccupancy, p.MaxCapacity)

	if p.CurrentOccupancy >= p.MaxCapacity {
		fmt.Printf("PARKING LOT IS FULL\n")
	} else {
		available := p.MaxCapacity - p.CurrentOccupancy
		fmt.Printf("🟢 %d slots available\n", available) // Add this!
	}
	for _, floor := range p.Floors {
		fmt.Printf("Floor %v:\n ", floor.FloorNumber)
		for slotType, slots := range floor.Slots {
			avaialble := 0
			for _, slot := range slots {
				if slot.isAvailable() {
					avaialble++
				}
			}
			fmt.Printf("%v: %v/%v available\n", slotType, avaialble, len(slots))
		}
	}
	fmt.Printf("Active Tickets: %v\n", len(p.ActiveTickets))
}

func (v *Vehicle) GetRequiredSlotType() SlotType {
	//Based on v.VehicleType return the slotType
	switch v.VehicleType {
	case Truck:
		return Large
	case Car, Bike:
		return Compact
	default:
		return Large
	}
}

func (s *Slot) isAvailable() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Status == Available
}

func (s *Slot) Occupy(v *Vehicle) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Status = Occupied
	s.ParkedVehicle = v
}

func (s *Slot) Free() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Status = Available
	s.ParkedVehicle = nil
}

// startegy pattern for fee calculation
type PricingStrategy interface {
	CalculateFee(ticket *Ticket) float64
}

type HourlyPricing struct {
	RatePerHour map[VehicleType]float64
}

func NewHourlyPricing() *HourlyPricing {
	return &HourlyPricing{
		RatePerHour: map[VehicleType]float64{
			Truck: 20.0,
			Bike:  15.0,
			Car:   10.0,
		},
	}
}

// Impl1
func (h *HourlyPricing) CalculateFee(ticket *Ticket) float64 {
	//we have to return fee based on entry and exit time
	entryTime := ticket.EntryTime
	exitTime := time.Now()
	duration := exitTime.Sub(entryTime)
	hours := duration.Hours()

	if hours < 1 {
		hours = 1 // Minimum 1 hour
	}

	vehicle := ticket.VehicleAssigned
	vehicleType := vehicle.VehicleType

	rate := h.RatePerHour[vehicleType]
	return hours * rate
}

type FlatPricing struct {
}

// Impl2
func (f *FlatPricing) CalculateFee(ticket *Ticket) float64 {
	fee := 10.2
	return fee
}

// progressive hour strategy
type ProgressivHourStrategy struct {
	RatePerHour     map[VehicleType]float64
	FirstHourRate   float64 // $4
	SecondThirdRate float64 // $3.5
	RemainingRate   float64 // $2.5
}

func (p *ProgressivHourStrategy) CalculateFee(ticket *Ticket) float64 {
	entryTime := ticket.EntryTime
	exitTime := time.Now()
	duration := exitTime.Sub(entryTime)
	hours := duration.Hours()

	if hours <= 1 {
		return p.FirstHourRate
	} else if hours <= 3 {
		return p.FirstHourRate + (hours-1)*p.SecondThirdRate
	} else {
		return p.FirstHourRate + 2*p.SecondThirdRate + (hours-3)*p.RemainingRate
	}
}

// Factory Pattern for vehicle creation as creating vehicles manually is repetitive and error-prone.
type VehicleFactory struct{}

func (f *VehicleFactory) CreateVehicle(vType VehicleType, licensePlate string) *Vehicle {
	// TODO: Return a new Vehicle with given type and license plate
	// Create vehicle
	vehicle := &Vehicle{
		VehicleType:  vType,
		LicensePlate: licensePlate,
	}
	return vehicle
}

func main() {
	//create parkingLot using singelton design pattern
	parkingLot := GetParkingLot("Grand Plaza", 1, 2, NewHourlyPricing())

	// Create and add gates
	entryGate := &Gate{
		Id:     "GATE-A",
		Type:   EntryGate,
		Status: GateOpen,
		Floor:  0,
	}

	exitGate := &Gate{
		Id:     "GATE-B",
		Type:   ExitGate,
		Status: GateOpen,
		Floor:  0,
	}

	bothGate := &Gate{
		Id:     "GATE-C",
		Type:   BothGate,
		Status: GateOpen,
		Floor:  0,
	}

	parkingLot.AddGate(entryGate)
	parkingLot.AddGate(exitGate)
	parkingLot.AddGate(bothGate)

	//Add slots to floor 0
	floor0 := parkingLot.Floors[0]

	//create slots
	slot0 := &Slot{
		ID:            "1",
		Type:          Compact,
		Status:        Available,
		ParkedVehicle: nil,
		Floor:         0,
	}

	slot1 := &Slot{
		ID:            "2",
		Type:          Large,
		Status:        Available,
		ParkedVehicle: nil,
		Floor:         0,
	}

	slot2 := &Slot{
		ID:            "1",
		Type:          Handicapped,
		Status:        Available,
		ParkedVehicle: nil,
		Floor:         0,
	}

	floor0.Slots[Compact] = append(floor0.Slots[Compact], slot0)
	floor0.Slots[Large] = append(floor0.Slots[Large], slot1)
	floor0.Slots[Handicapped] = append(floor0.Slots[Handicapped], slot2)

	//create vehicle factory
	factory := &VehicleFactory{}
	car := factory.CreateVehicle(Car, "TG07AH4118")
	truck := factory.CreateVehicle(Truck, "PB46AQ4158")

	// Test 1: Enter through GATE-A (Entry gate - should work)
	fmt.Println("\n--- Test 1: Entry through GATE-A ---")
	ticket1, err := parkingLot.ProcessVehicleEntry(car, "GATE-A")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Car entered via GATE-A\n")
		fmt.Printf("   Ticket: %s, Slot: %s\n", ticket1.TicketId, ticket1.SlotAssigned.ID)
	}

	// Test 2: Enter through GATE-B (Exit gate - should FAIL)
	fmt.Println("\n--- Test 2: Try entry through EXIT gate ---")
	_, err = parkingLot.ProcessVehicleEntry(truck, "GATE-B")
	if err != nil {
		fmt.Printf("❌ Expected Error: %v\n", err)
	}

	// Test 3: Enter through GATE-C (Both gate - should work)
	fmt.Println("\n--- Test 3: Entry through GATE-C (Both) ---")
	ticket2, err := parkingLot.ProcessVehicleEntry(truck, "GATE-C")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Truck entered via GATE-C\n")
		fmt.Printf("   Ticket: %s, Slot: %s\n", ticket2.TicketId, ticket2.SlotAssigned.ID)
	}

	parkingLot.DisplayStatus()

	time.Sleep(2 * time.Second)

	// Test 4: Exit through GATE-B (Exit gate - should work)
	fmt.Println("\n--- Test 4: Exit through GATE-B ---")
	fee1, err := parkingLot.ProcessVehicleExit(ticket1.TicketId, "GATE-B")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Car exited via GATE-B, Fee: $%.2f\n", fee1)
	}

	// Test 5: Exit through GATE-A (Entry gate - should FAIL)
	fmt.Println("\n--- Test 5: Try exit through ENTRY gate ---")
	_, err = parkingLot.ProcessVehicleExit(ticket2.TicketId, "GATE-A")
	if err != nil {
		fmt.Printf("❌ Expected Error: %v\n", err)
	}

	// Test 6: Exit through GATE-C (Both gate - should work)
	fmt.Println("\n--- Test 6: Exit through GATE-C (Both) ---")
	fee2, err := parkingLot.ProcessVehicleExit(ticket2.TicketId, "GATE-C")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Truck exited via GATE-C, Fee: $%.2f\n", fee2)
	}

	parkingLot.DisplayStatus()

	// Test 7: Try using closed gate
	fmt.Println("\n--- Test 7: Try using closed gate ---")
	entryGate.Status = GateClosed
	bike := factory.CreateVehicle(Bike, "MH12XY9876")
	_, err = parkingLot.ProcessVehicleEntry(bike, "GATE-A")
	if err != nil {
		fmt.Printf("❌ Expected Error: %v\n", err)
	}

}
