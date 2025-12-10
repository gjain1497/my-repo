package main


//Flow
Customer places order -> Selects locker location
Order assigned to delivery agent with desiitnation  that user selected
Delivery agent arrives at locker with the package
System alloocates a specific locker at that location -> based on package size




//user selects package size while ordering items

//GetAllLockersNearToMyLocation(location)

//AddthatLockerLocationTOMyAddressList

//So now address contains this locker address as well


//Now while creating the order pass this id instead

//Based on the order details i.e number of itmes in order/size of items,
//system will decide which locker size to assign at that location
//locker will be assigned for a paricutalr time only


//User picks up by going at that location
//If time is more than given time delivery boy will pick it up after given time

//Customer side
BrowseLockerLocations(userLocation) -> list of locationsconst
AddToAddressBook(lockerLocationID)
PlaceOrder(userId, cartId, deliveryAddressId)


//System Side (after order placed)
PackItems(order_id -> checks items for the order id Determines packageSize (S/M/L/XL))
AssignToDeliveryAgent(order, lockerLocation)


//Delivery Agent side
ArrivesAtLockerLocation(packageSize)
System: FindAvaialbleLocker(lockerLocationId, packageSize)
System: AllocateLocker() -> generates Code
System: NotifyCustomer(lockerNumber, code)


//Customer picks up
EnterCode(code)
System Validates -> open locker
Customr picks up pakcage




type User struct {
    ID             string
    Name           string
    Email          string
    Phone          string
    SavedLocations []string  // Address IDs
}



type Location struct{
	Id string
	Street string
	City string
	Pincode string 
}


type LockerLocation struct{
	ID string
	Name string
	Location Location
	Capacity 
}

type LocationService interface{
	GetLockersByLocationID(locationId string) []Location
}

type LocationServiceV1 struct{

}

func(s *LocationServiceV1) GetLockersByLocationID(locationId string) []Locker{

}

type LockerService interface{
	AssignLockerByPackageSize(locationId string, packageSize string)
	AllocateLocker(lockerId string, locationId string) (string, error)
	GetLockersAtLocation(locationId string) []*Location
}

type LockerServiceV1 struct{
	lockers map[string][]*Locker //(location_id, lockers) 
}

func(s *LocationServiceV1) GetLockersAtLocation(locationId string) []*Location{

}

func(s *LockerServiceV1)AssignLockerByPackageSize(lockerLocationId string, packageSize string) (*Locker, error){
	//getAllLockersfor that location from the map
	//asign based on the package size
}

func(s *LockerServiceV1) AllocateLocker(lockerId string, locationId string) (string, error){ //Code, error
//notify customer after allocating the lock
}



//one location has many lockers
type Locker struct{
	ID string
	LocationId string
	Size LockerSize
	Status LockerStatus
}

type LockerSize string
const (
    Small      LockerSize = "SMALL"
    Medium     LockerSize = "MEDIUM"
    Large      LockerSize = "LARGE"
    ExtraLarge LockerSize = "EXTRA_LARGE"
)

type LockerStatus string
const (
    Available     LockerStatus = "AVAILABLE"
    Occupied      LockerStatus = "OCCUPIED"
    OutOfService  LockerStatus = "OUT_OF_SERVICE"
)

type OrderService interface{
	CreateOrder(cartId, userId, addressId, paymentType)
}

type OrderServiceV1 struct{

}

type Order struct{
	ID string
	UserId string
	ShippingAddressId string
	OrderItems []OrderItems
	Status OrderStatus
	PaymentId string
	TotalAmount float64
	CreatedAt time.Time
}


type OrderItem struct {
	// OrderItems string
	ProductId    string
	ProductName  string
	Quantity     int
	PriceAtOrder float64
}

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Confirmed OrderStatus = "CONFIRMED"
	Shipped   OrderStatus = "SHIPPED"
	Delivered OrderStatus = "DELIVERED"
	Cancelled OrderStatus = "CANCELLED"
)


type Package struct{
	ID string
	OrderId string
	Size PackageSize
	Status PackageStatus
	AssignedLockerId string
	PickUpCode string
	ExpiryTime       time.Time
    CreatedAt        time.Time
    PickedUpAt       *time.Time
}

type PackageSize string
const (
    SmallPackage      PackageSize = "SMALL"
    MediumPackage     PackageSize = "MEDIUM"
    LargePackage      PackageSize = "LARGE"
    ExtraLargePackage PackageSize = "EXTRA_LARGE"
)

type PackageStatus string
const (
    InTransit PackageStatus = "IN_TRANSIT"
    Delivered PackageStatus = "DELIVERED"
    PickedUp  PackageStatus = "PICKED_UP"
    Expired   PackageStatus = "EXPIRED"
    Returned  PackageStatus = "RETURNED"
)



type DeliveryAgent struct {
    ID     string
    Name   string
    Phone  string
    Status DeliveryAgentStatus
}

type DeliveryAgentStatus string
const (
    Active     DeliveryAgentStatus = "ACTIVE"
    Inactive   DeliveryAgentStatus = "INACTIVE"
    OnDelivery DeliveryAgentStatus = "ON_DELIVERY"
)
func main() {

}
