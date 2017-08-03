package main

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	switchInputReadFrequency time.Duration = 20 * time.Millisecond
)

type LimitSwitch struct {
	pin rpio.Pin
}

func NewLimitSwitch(pin rpio.Pin) *LimitSwitch {
	s := &LimitSwitch{
		pin: pin,
	}
	return s
}

func (s *LimitSwitch) Notify() <-chan bool {
	notifyChan := make(chan bool, 1)
	s.pin.PullUp()

	go func() {
		defer s.pin.PullOff()
		switchInputTicker := time.NewTicker(switchInputReadFrequency)
		isSwitchReleased := false
		for {
			select {
			case <-switchInputTicker.C:
				// First make sure the switch is in released position.
				if !isSwitchReleased {
					if s.pin.Read() == rpio.High {
						isSwitchReleased = true
					}
					continue
				}

				// Now check if the switch was depressed.
				if s.pin.Read() != rpio.Low {
					continue
				}
				switchInputTicker.Stop()
				notifyChan <- true
				close(notifyChan)
				return
			}
		}
	}()

	return notifyChan
}

/*
func (s *LimitSwitch) NotifyAfterRelease() <-chan bool {
	notifyChan := make(chan bool, 1)
	s.pin.PullUp()

	go func() {
		defer s.pin.PullOff()
		switchInputTicker := time.NewTicker(switchInputReadFrequency)
		valueToRead := rpio.Low
		state := ""
		for {
			select {
			case <-switchInputTicker.C:
				val := s.pin.Read()
				if val != valueToRead {
					continue
				}
				valueToRead = rpio.High
				if state == "" {
					state = "pressed"
					continue
				}
				if state == "pressed" {
					switchInputTicker.Stop()
					notifyChan <- true
					//close(notifyChan)
					return
				}
			}
		}
	}()

	return notifyChan
}
*/
