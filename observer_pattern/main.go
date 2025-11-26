package main

import "fmt"

// You (YouTuber): "Hey, I uploaded a new video!"
// You manually call each person:
// - Call friend 1: "Hey, new video!"
// - Call friend 2: "Hey, new video!"
// - Call friend 3: "Hey, new video!"
// ```

// **Problem:** You have to remember everyone and call them manually!

// **With Observer Pattern (YouTube Subscription):**
// ```
// People subscribe to your channel ✅
// You upload a video ✅
// YouTube automatically notifies ALL subscribers! ✅

// 🎯 Observer Pattern in Simple Terms:
// 3 Key Players:

// Subject (Observable) - The thing being watched

// Example: YouTube Channel, Cricket Match, Stock Price

// Observer - The thing that wants to be notified

// Example: Subscriber, Scoreboard, Stock Trader

// Notification - When something changes

// Example: New video, Score update, Price change

// Let's start with a Weather Station example:
// Scenario:

// Weather station tracks temperature
// Multiple displays want to show temperature
// When temperature changes, ALL displays should update automatically

// Question: What should ALL observers be able to do?
// Answer: They should be able to update when notified!

// Observer interface - ANY observer/display here must implement this
type Observer interface {
	Update(temperature float64)
}

//This means: "Any observer MUST have an Update() method"

// Step 2: Create the Subject/Observable (Thing Being Observed)

type WeatherStation struct {
	temperature float64
	observers   []Observer //list of subscribers
}

// Subscribe - Add an observer
func (w *WeatherStation) Subscribe(observer Observer) {
	w.observers = append(w.observers, observer)
}

// Notify - Tell all subscribers about the change
func (w *WeatherStation) NotifyAll() {
	for _, observer := range w.observers {
		observer.Update(w.temperature)
	}
}

// Set temperature when the temperature changes
func (w *WeatherStation) SetTemperature(temp float64) {
	w.temperature = temp
	w.NotifyAll()
}

//Now different observers which implement this update method in their own way

type PhoneDisplay struct {
	name        string
	temperature float64
}

func (p *PhoneDisplay) Update(temperature float64) {
	fmt.Printf("Temperature of %s = %f ", p.name, temperature)
}

type TVDisplay struct {
	name        string
	temperature float64
}

func (p *TVDisplay) Update(temperature float64) {
	fmt.Printf("Temperature of %s = %f ", p.name, temperature)
}

type WindowDisplay struct {
	name        string
	temperature float64
}

func (p *WindowDisplay) Update(temperature float64) {
	fmt.Printf("Temperature of %s = %f ", p.name, temperature)
}

func main() {
	weather := &WeatherStation{}

	//create observers
	phone := &PhoneDisplay{name: "MyPhone"}
	tv := &TVDisplay{name: "Living Room TV"}
	window := &WindowDisplay{name: "Window Display"}

	//subscribe observers to weather station
	weather.Subscribe(phone)
	weather.Subscribe(tv)
	weather.Subscribe(window)

	//change temperature of observable/display
	fmt.Println("Setting temperature to 25°C:")
	weather.SetTemperature(25.0)

	fmt.Println("\nSetting temperature to 30°C:")
	weather.SetTemperature(30.0)

	fmt.Println("\nSetting temperature to 18°C:")
	weather.SetTemperature(18.0)
}























