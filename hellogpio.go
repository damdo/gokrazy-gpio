package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

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

	// log.Printf("Toggling GPIO forever")
	// t := time.NewTicker(5 * time.Second)
	// for l := gpio.Low; ; l = !l {
	// 	log.Printf("setting GPIO pin number 18 (signal BCM24) to %v", l)
	// 	// Lookup a pin by its location on the board:
	// 	if err := rpi.P1_18.In(gpio.PullDown, gpio.); err != nil {
	// 		return err
	// 	}
	// 	<-t.C
	// }
	for i := 0; i < 5; i++ {
		fmt.Printf("%s: %s\n", p, p.Function())
		time.Sleep(1 * time.Second)
	}

	// Set it as input, with an internal pull down resistor:
	err := p.In(gpio.PullNoChange, gpio.BothEdges)
	if err != nil {
		log.Println("errrr:", err)
	}

	fmt.Println("entering wait loop")
	// Wait for edges as detected by the hardware, and print the value read:
	for {
		p.WaitForEdge(-1)
		fmt.Printf("-> %s\n", p.Read())
		time.Sleep(100 * time.Millisecond)
		fmt.Println("next")
	}
}
