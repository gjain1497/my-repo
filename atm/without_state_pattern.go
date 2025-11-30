package main

//simplest example: light bulb with remote

//you have a light bulb with a remote control that has one button

//how it works

// if light is off -> turn on
// if light is on -> turn off

// but now new req on -> dim -> off

// Press button:
// - If light is OFF    → Turn ON
// - If light is ON     → Turn DIMMED
// - If light is DIMMED → Turn OFF

//without state patten

// type Light struct {
// 	status string
// }

// func (l *Light) PressButton() {
// 	if l.status == "OFF" {
// 		l.status = "ON"
// 		fmt.Println("💡 Light is ON")

// 	} else if l.status == "ON" {
// 		l.status = "DIMMED"
// 		fmt.Println("🌙 Light is DIMMED")

// 	} else if l.status == "DIMMED" {
// 		l.status = "OFF"
// 		fmt.Println("⚫ Light is OFF")
// 	}
// }

// // //Now new requirement comes to add flashing mode
// // Press button:
// // - OFF → ON → DIMMED → FLASHING → OFF
// // You must change:

// func (l *Light) PressButton() {
// 	if l.status == "OFF" {
// 		l.status = "ON"
// 		fmt.Println("💡 Light is ON")

// 	} else if l.status == "ON" {
// 		l.status = "DIMMED"
// 		fmt.Println("🌙 Light is DIMMED")

// 	} else if l.status == "DIMMED" {
// 		l.status = "FLASHING" // ✅ Added
// 		fmt.Println("⚡ Light is FLASHING")

// 	} else if l.status == "FLASHING" { // ✅ Added
// 		l.status = "OFF"
// 		fmt.Println("⚫ Light is OFF")
// 	}
// }

// // ❌ Had to modify the ONE method
// // ❌ Adding more modes = more if-else

// //With state pattern

// func main() {
// 	light := &Light{status: "OFF"}

// 	light.PressButton() // OFF → ON
// 	light.PressButton() // ON → DIMMED
// 	light.PressButton() // DIMMED → OFF
// 	light.PressButton() // OFF → ON
// }
