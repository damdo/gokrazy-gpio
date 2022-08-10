package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var debouncingDelay = time.Second * 2

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Lookup a pin by its number:
	p := gpioreg.ByName("GPIO2")
	if p == nil {
		log.Fatal("Failed to find GPIO2")
	}

	fmt.Printf("listening for clicks on: %s, current state: %s\n", p, p.Function())

	// Set it as input, with an internal pull down resistor:
	err := p.In(gpio.PullNoChange, gpio.BothEdges)
	// ignore not exported error
	if err != nil && !strings.Contains(err.Error(), "is not exported by sysfs") {
		log.Println(err)
	}

	fmt.Println("entering checking loop")
	// start a timer for keeping track of the delay between clicks
	// to debounce long clicks
	start := time.Now().Add(debouncingDelay)
	// Wait for edges as detected by the hardware, and print the value read:
	previous := gpio.High
	for {
		p.WaitForEdge(-1)
		// Level is normally High.
		l := p.Read()

		now := time.Now()
		elapsed := now.Sub(start)
		// If now level is Low and was previously High,
		// State has changed.
		if l == gpio.Low && previous == gpio.High {
			fmt.Printf("state change High -> Low\n")
			// Only if the last click has happened before the debouncing delay
			// we can consider this as a new click
			if elapsed > debouncingDelay {
				fmt.Println("considering it a click")
				// reset the last click timer
				start = time.Now()
			}
		}
		// Store previous value for comparison.
		previous = l

		// Sleep to pause a little and avoid CPU lock-up.
		time.Sleep(100 * time.Millisecond)
	}
}
