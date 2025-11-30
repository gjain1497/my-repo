package main

import "fmt"

// State interface
type LightState interface {
	PressButton(light *Light)
}

// Context
type Light struct { //light has lightstate interface and light also implements lightstate interface
	state LightState
}

func NewLight() *Light {
	return &Light{
		state: &OffState{},
	}
}

func (l *Light) PressButton() {
	l.state.PressButton(l)
}

func (l *Light) SetState(state LightState) {
	l.state = state
}

type OffState struct{}

func (s *OffState) PressButton(light *Light) {
	fmt.Println("Light is on")
	light.SetState(&OnState{})
}

type OnState struct{}

func (s *OnState) PressButton(light *Light) {
	fmt.Println("Light is DIMMED")
	light.SetState(&DimmedState{})
}

type DimmedState struct{}

func (s *DimmedState) PressButton(light *Light) {
	fmt.Println("Light is OFF")
	light.SetState(&OffState{})
}

func main() {
	light := NewLight()
	light.PressButton() //OFF -> ON
	light.PressButton() // ON → DIMMED
	light.PressButton() // DIMMED → OFF
	light.PressButton() // OFF → ON
}
