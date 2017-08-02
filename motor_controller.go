package main

import (
	"fmt"
	"time"
)

type MotorController struct {
	motor          *Motor
	delay          time.Duration
	tripDuration   time.Duration
	roundTripCount int
}

func NewMotorController() *MotorController {
	motor := NewMotor()

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
	fmt.Println("Going to run the motor with a step duration of: ", motorStepDuration)
	fmt.Printf("The motor would cover [%d] steps in [%s]\n", mc.motor.maxSteps, mc.tripDuration)
	stopMotor := make(chan bool)
	motorDone := mc.motor.Run(stopMotor, motorStepDuration)

	fmt.Println("Waiting for the motor to finish")
	<-time.After(mc.tripDuration)
	stopMotor <- true
	<-motorDone

}

func (mc *MotorController) reset() {
	mc.delay = 0
	mc.tripDuration = 0
	mc.roundTripCount = 0
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
