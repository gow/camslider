package main

import (
	"fmt"
	"time"
)

const (
	redLEDPin   = 10
	greenLEDPin = 7
)

var (
//greenLED *Blinker
//redLED   *Blinker
)

type MotorController struct {
	motor                *Motor
	delay                time.Duration
	tripDuration         time.Duration
	tripEndPauseDuration time.Duration
	roundTripCount       int
}

func NewMotorController() *MotorController {
	redLED := NewBlinker(redLEDPin, BlinkFast)
	redLED.Blink()
	motor := NewMotor()
	redLED.Stop()

	redLED = NewBlinker(redLEDPin, BlinkSlow)
	redLED.Blink()
	return &MotorController{
		motor: motor,
	}
}

func (mc *MotorController) Run(options ...ControllerOption) {
	mc.reset()
	for _, option := range options {
		option(mc)
	}

	if mc.delay != 0 {
		fmt.Println("Delaying the start of the motor by ", mc.delay)
		<-time.After(mc.delay)
	}

	motorStepDuration := mc.tripDuration / time.Duration(mc.motor.maxSteps)
	for trip := 0; trip <= mc.roundTripCount*2; trip++ {
		fmt.Printf("###################################################################\n")
		fmt.Printf("[Trip-%d]: Going to run the motor with a step duration of [%s]\n", trip/2, motorStepDuration)
		fmt.Printf("[Trip-%d]: The motor would cover [%d] steps in [%s]\n", trip/2, mc.motor.maxSteps, mc.tripDuration)
		stopMotor := make(chan bool)
		motorDone := mc.motor.Run(stopMotor, motorStepDuration)

		fmt.Printf("[Trip-%d]: Waiting [%s] for the motor to finish\n", trip/2, mc.tripDuration)
		<-time.After(mc.tripDuration)
		close(stopMotor)
		<-motorDone
		mc.motor.ToggleDirection()

		if mc.tripEndPauseDuration != 0 {
			fmt.Printf("[Trip-%d]: Trip end reached. Delaying the start of the motor by [%s]\n", trip/2, mc.tripEndPauseDuration)
			<-time.After(mc.delay)
		}
	}

}

func (mc *MotorController) reset() {
	mc.delay = 0
	mc.tripDuration = 0
	mc.roundTripCount = 0
	mc.tripEndPauseDuration = 0
}

type ControllerOption func(mc *MotorController)

func WithDelay(delay time.Duration) ControllerOption {
	return func(mc *MotorController) {
		mc.delay = delay
	}
}

func WithTripDuration(tripDur time.Duration) ControllerOption {
	return func(mc *MotorController) {
		mc.tripDuration = tripDur
	}
}

func WithRoundTripCount(count int) ControllerOption {
	return func(mc *MotorController) {
		mc.roundTripCount = count
	}
}

func WithTripEndPause(pauseDur time.Duration) ControllerOption {
	return func(mc *MotorController) {
		mc.tripEndPauseDuration = pauseDur
	}
}
