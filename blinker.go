package main

import (
	"fmt"
	"sync"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const blinkerPin = 10

type Blinker struct {
	pin         rpio.Pin
	onDuration  time.Duration
	offDuration time.Duration
	stopChan    chan bool
	wg          sync.WaitGroup
}

func NewBlinker(pin int) *Blinker {
	rpioPin := rpio.Pin(pin)
	rpioPin.Output()
	return &Blinker{
		pin:         rpioPin,
		onDuration:  time.Second,
		offDuration: time.Second,
		stopChan:    make(chan bool),
	}
}

func (blinker *Blinker) Blink() {
	timer := time.NewTimer(time.Second * 1)
	// On Goroutine
	blinker.wg.Add(1)
	go func(stopChan <-chan bool) {
		defer blinker.wg.Done()
		for {
			select {
			case q := <-blinker.stopChan:
				timer.Stop()
				fmt.Println("Stop signal received.", q)
				return
			case <-timer.C:
				fmt.Println("ON")
				blinker.pin.High()
				<-time.After(time.Second)

				fmt.Println("OFF")
				blinker.pin.Low()
				<-time.After(time.Second)
				timer = time.NewTimer(time.Second * 1)
			}
		}
	}(blinker.stopChan)
	fmt.Println("Waiting for 15 sec")
	<-time.After(15 * time.Second)
	blinker.Stop()
	fmt.Println("Done stopping the blinker. Returning...")
}

func (blinker *Blinker) Stop() {
	fmt.Println("Stopping the blinker")
	blinker.stopChan <- true
	blinker.wg.Wait()
}
