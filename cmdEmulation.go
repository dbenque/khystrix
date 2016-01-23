package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

var random = rand.New(rand.NewSource(0))

func durationJitter(d time.Duration, r *rand.Rand) time.Duration {
	if d == 0 {
		return 0
	}
	return d + time.Duration(r.Int63n(2*int64(d)))
}

func printSleep(baseDurationMs, floorDurationMs int64) error {
	d := durationJitter((time.Duration)(baseDurationMs)*time.Millisecond, random)
	n := d.Nanoseconds() / 1000 / 1000
	if n < floorDurationMs {
		return fmt.Errorf("Boouuuh")
	}
	time.Sleep(d)
	//log.Printf("As was sleeping for %d ms\n", n)
	return nil
}

//LoopOverCmd emulate calls to command
func LoopOverCmd(cmdName string, baseDurationMs, floorDurationMs int) {

	defer log.Printf("!!!!!!Exiting command loop for: %s\n", cmdName)

	for {
		output := make(chan string, 1)
		errors := hystrix.Go(cmdName, func() error {
			if err := printSleep(int64(baseDurationMs), int64(floorDurationMs)); err != nil {
				output <- "Ko"
				return err
			}
			output <- "Ok"
			return nil
		}, func(err error) error {
			output <- "Degraded mode"
			return nil
		})

		select {
		case out := <-output:
			log.Printf("Result: %v\n", out)
		case err := <-errors:
			log.Printf("Hystrix Error: %v\n", err)
		}
	}

}
