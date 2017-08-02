package main

import (
	"fmt"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

func main() {
	fmt.Println("Hello slider!")
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()

	motorController := NewMotorController()
	motorController.Run(
		WithDelay(5*time.Second),
		WithTripDuration(20*time.Second),
		WithRoundTripCount(3),
	)

	fastestBlinker := NewBlinker(10, BlinkFastest)
	fastestBlinker.Blink()
	<-time.After(5 * time.Second)
	fastestBlinker.Stop()

	fastBlinker := NewBlinker(10, BlinkFast)
	fastBlinker.Blink()
	<-time.After(5 * time.Second)
	fastBlinker.Stop()

	slowBlinker := NewBlinker(10, BlinkSlow)
	slowBlinker.Blink()
	<-time.After(5 * time.Second)
	slowBlinker.Stop()

	slowestBlinker := NewBlinker(10, BlinkSlowest)
	slowestBlinker.Blink()
	<-time.After(5 * time.Second)
	slowestBlinker.Stop()
}
