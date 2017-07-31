package main

import (
	"fmt"
	"sync"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const blinkerPin = 10

type BlinkSpeed int

const (
	BlinkFastest BlinkSpeed = iota
	BlinkFast
	BlinkSlow
	BlinkSlowest
)

func (blinkSpeed BlinkSpeed) Duration() time.Duration {
	return map[BlinkSpeed]time.Duration{
		BlinkFastest: 50 * time.Millisecond,
		BlinkFast:    100 * time.Millisecond,
		BlinkSlow:    700 * time.Millisecond,
		BlinkSlowest: 1500 * time.Millisecond,
	}[blinkSpeed]
}

type Blinker struct {
	pin         rpio.Pin
	onDuration  time.Duration
	offDuration time.Duration
	stopChan    chan bool
	toggle      chan bool
	wg          sync.WaitGroup
}

func NewBlinker(pin int, speed BlinkSpeed) *Blinker {
	rpioPin := rpio.Pin(pin)
	rpioPin.Output()
	return &Blinker{
		pin:         rpioPin,
		onDuration:  speed.Duration(),
		offDuration: speed.Duration(),
		toggle:      make(chan bool),
	}
}

func (blinker *Blinker) Blink() {
	blinker.stopChan = make(chan bool)
	// ON
	go func() {
		for {
			select {
			case <-blinker.stopChan:
				fmt.Println("Stop signal received. [ON]")
				return
			case <-blinker.toggle:
				//fmt.Println("ON")
				blinker.pin.High()
				time.Sleep(blinker.onDuration)
				blinker.toggle <- true
			}
		}
	}()

	// OFF
	go func() {
		for {
			select {
			case <-blinker.stopChan:
				fmt.Println("Stop signal received. [OFF]")
				return
			case <-blinker.toggle:
				//fmt.Println("OFF")
				blinker.pin.Low()
				time.Sleep(blinker.offDuration)
				blinker.toggle <- true
			}
		}
	}()

	//kick off the blinking.
	blinker.toggle <- true
}

func (blinker *Blinker) Stop() {
	fmt.Println("Stopping the blinker")
	close(blinker.stopChan)
	<-blinker.toggle
	blinker.pin.Low()
}
