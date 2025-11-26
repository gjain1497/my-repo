package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type User struct {
	Id   string
	Name string
}

type Theater struct { //theatre has many screens
	ID       string
	Screens  map[string]*Screen //(screen_id -> screen)
	Location Location           //Theater has location, not movie, to find movies in a city, find all theaters in that city
	//find all shows in that theatre, see which movies are in those shows
}

type Movie struct {
	ID       string
	Name     string
	Duration time.Duration
	Language string
	Genre    string
}

//Q) why movies should not be part of the screen
//Why this is correct:
// Screen is a physical room with seats
// Movies are shown at specific times (Shows)
// A screen doesn't "own" movies permanently

// Real world analogy:

// Your laptop screen doesn't "own" YouTube videos
// It just displays them when you play them
// Same here: Screen displays movies through Shows

type Screen struct { //screen has many movies and many seats
	ScreenId string
	// Movies   map[string]*Movie //(movie_id -> movie)
	Seats   map[string]*Seat //(seat_id -> seat)
	Theater *Theater
}

type Location struct {
	City string
}

type Show struct { //show is movie + screen + time //references movie, references screen, references starttime
	ID      string
	Time    time.Time
	Movie   *Movie
	Screen  *Screen
	Theater *Theater
	Seats   map[string]*ShowSeat //(show_id -> {showseat})
}

func (s *Show) GetSeats() map[string]*ShowSeat {
	return s.Seats
}

func (s *Show) getShowById() {

}

type SeatStatus string

const (
	Booked    SeatStatus = "BOOKED"
	Avaialble SeatStatus = "AVAILBLE"
	Waiting   SeatStatus = "WAITING"
)

type SeatType string

const (
	Deluxe SeatType = "DELUXE"
	Gold   SeatType = "GOLD"
	Silver SeatType = "SILVER"
)

// Removed Status from Seat - CORRECT! ✅
// Why this is correct:

// Seat A1 is always there (physical object)
// But its availability changes per show
// 6 PM show: A1 is booked
// 9 PM show: A1 is available
// Status belongs on ShowSeat, not Seat

type Seat struct {
	Number string
	Type   SeatType
	Screen *Screen
}

type ShowSeat struct {
	Seat   *Seat
	Status SeatStatus
	Price  float64
	mutex  *sync.Mutex
}

func (s *ShowSeat) BookSeat() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Status == Avaialble {
		s.Status = Booked
	}
	if s.Status == Waiting {

	}
}

type BookingStatus string

const (
	Pending   BookingStatus = "PENDING"
	Reserved  BookingStatus = "RESERVED"
	Confirmed BookingStatus = "CONFIRMED"
)

type Booking struct {
	UserID    string
	BookingId string
	Show      *Show
	// Seats     map[string]*Seat //(seat_id -> Seat)
	Seats     map[string]*ShowSeat //(seat_id -> Seat)
	Status    BookingStatus
	CreatedAt time.Time
	ExpiresAt *time.Time
}

//find all movies in Mumbai

type BookMyShowService struct {
	// Raw data
	theaters map[string]*Theater
	shows    map[string]*Show //(show_id -> show)

	// Indexes for fast lookup
	theatersByCity map[string][]*Theater // city → [theaters]
	showsByMovie   map[string][]*Show    // movie_id → [shows]
	showsByTheater map[string][]*Show    // theater_id → [shows]
}

// Find movies in Mumbai
func (s *BookMyShowService) GetMoviesInCity(city string) []*Movie {
	theaters := s.theatersByCity[city] // O(1) lookup

	movies := make(map[string]*Movie)
	for _, theater := range theaters {
		shows := s.showsByTheater[theater.ID] // O(1) lookup
		for _, show := range shows {
			movies[show.Movie.ID] = show.Movie
		}
	}

	return getValues(movies)
}

func (s *BookMyShowService) AddTheatre(Theater) {

}

func (s *BookMyShowService) AddShow() {

}

func (s *BookMyShowService) GetShowByID(showID string) (*Show, error) {
	if _, ok := s.shows[showID]; ok {
		return s.shows[showID], nil
	}
	return nil, fmt.Errorf("Show not found")
}

func (s *BookMyShowService) BookTickets(userId string, showID string, seatIDs []string) (*Booking, error) {
	//get the show
	show, err := s.GetShowByID(showID)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("BookTickets returned error"))
	}

	//show has show seats
	//getSeats of this show
	seats := show.GetSeats()

	//get those particular seats using seat IDS
	for _, seat := range seats {

	}

	//try to book seats now
	for _, seat := range seats {
		seat.BookSeat()
	}
}

func (s *BookMyShowService) CancelBooking() {

}

func getValues(m map[string]*Movie) []*Movie {
	result := make([]*Movie, 0, len(m))
	for _, value := range m {
		result = append(result, value)
	}
	return result
}
