package main

import (
	"fmt"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	maxSpeed     int64 = 1000 // steps per sec
	maxInitSpeed int64 = 500

	//switchInputReadFrequency time.Duration = 100 * time.Millisecond
)

const (
	limitSwitchPinA rpio.Pin = rpio.Pin(5)
	limitSwitchPinB rpio.Pin = rpio.Pin(13)
)

// Define GPIO signals to use
// Physical pins 11,15,16,18
// GPIO17,GPIO22,GPIO23,GPIO24
var StepPins = []uint8{17, 22, 23, 24}

var stepSequence = [][]int{
	{1, 0, 0, 0},
	{1, 1, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 1, 0},
	{0, 0, 1, 0},
	{0, 0, 1, 1},
	{0, 0, 0, 1},
	{1, 0, 0, 1},
}

type Motor struct {
	maxSteps      int64
	isInitialized bool
	stepPins      []rpio.Pin
	stepDirection int
	currentStep   int
	limitSwitchA  *LimitSwitch
	limitSwitchB  *LimitSwitch
}

func NewMotor() *Motor {
	m := &Motor{
		stepDirection: 1,
		limitSwitchA:  NewLimitSwitch(limitSwitchPinA),
		limitSwitchB:  NewLimitSwitch(limitSwitchPinB),
	}
	for _, pinNum := range StepPins {
		pin := rpio.Pin(pinNum)
		pin.Output()
		m.stepPins = append(m.stepPins, pin)
	}
	m.init()
	m.ToggleDirection()
	fmt.Println("Max number of steps: ", m.maxSteps)
	return m
}

func (m *Motor) Step() {
	//fmt.Printf("Current step [%d]: [%#v]\n", m.currentStep, stepSequence[m.currentStep])

	nextStep := m.currentStep + m.stepDirection
	if nextStep < 0 {
		nextStep = len(stepSequence) - 1
	}
	nextStep = nextStep % len(stepSequence)
	//fmt.Printf("Next step [%d]: [%#v]\n", nextStep, stepSequence[nextStep])

	for i, val := range stepSequence[nextStep] {
		if val == 1 {
			m.stepPins[i].High()
		} else {
			m.stepPins[i].Low()
		}
	}

	// Save the current step position.
	m.currentStep = nextStep
}

func (m *Motor) ToggleDirection() {
	m.stepDirection = -1 * m.stepDirection
}

func (m *Motor) Run(stopChan <-chan bool, stepDuration time.Duration) {
	ticker := time.NewTicker(stepDuration)
	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping the motor")
			ticker.Stop()
			return
		case <-ticker.C:
			m.Step()
		}
	}
}

func (m *Motor) init() {
	fmt.Println("Resetting the motor")
	for _, pin := range m.stepPins {
		pin.Low()
	}
	m.Reset()
	m.ToggleDirection()

	switchNotificationA := m.limitSwitchA.Notify()
	switchNotificationB := m.limitSwitchB.Notify()

	stepTicker := time.NewTicker(time.Second / time.Duration(maxInitSpeed))
	m.maxSteps = 0
	for {
		select {
		case <-switchNotificationA:
			stepTicker.Stop()
			fmt.Println("Received notification from the switch-A:")
			return
		case <-switchNotificationB:
			stepTicker.Stop()
			fmt.Println("Received notification from the switch-B:")
			return
		case <-stepTicker.C:
			m.maxSteps++
			m.Step()
		}
	}
}

func (m *Motor) Reset() {
	stepTicker := time.NewTicker(time.Second / time.Duration(maxInitSpeed))
	switchNotificationA := m.limitSwitchA.NotifyAfterRelease()
	switchNotificationB := m.limitSwitchB.NotifyAfterRelease()

	for {
		select {
		case <-switchNotificationA:
			stepTicker.Stop()
			fmt.Println("Received notification from the switch-A:")
			return
		case <-switchNotificationB:
			stepTicker.Stop()
			fmt.Println("Received notification from the switch-B:")
			return
		case <-stepTicker.C:
			m.Step()
		}
	}
}
