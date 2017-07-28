package main

import (
	"fmt"

	rpio "github.com/stianeikeland/go-rpio"
)

func main() {
	fmt.Println("Hello slider!")
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	blinker := NewBlinker(10)
	blinker.Blink()
	//<-time.After(20 * time.Second)
	//blinker.Stop()
}
