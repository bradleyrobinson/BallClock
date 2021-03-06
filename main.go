package main

import (
	"encoding/json"
	"fmt"
	"time"

	"ballclock/clock"
)

func main() {
	bc, err := calculateClockState(30, 325)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(bc.Min)
	fmt.Println(bc.FiveMin)
	fmt.Println(bc.Hour)
	fmt.Println(bc.Main)
	cd, dur, err := cycleDays(30)
	fmt.Println(cd, dur, err)
}

// cycleDays is defined by the first problem.
func cycleDays(balls int) (days int, ms time.Duration, err error) {
	start := time.Now()

	bc, err := ballclock.CreateBallClock(balls)
	if err != nil {
		return 0, time.Duration(0), err
	}
	startingMain := make([]int, balls)
	copy(startingMain, bc.Main.Balls)
	minutes := 1
	bc.Tick()
	for !bc.Cycled(startingMain) {
		err = bc.Tick()
		if err != nil {
			return 0, time.Duration(0), err
		}
		minutes++
	}
	days = minutes / 60 / 24
	duration := time.Since(start)

	return days, duration, nil
}

func clockState(balls, minutes int) ([]byte, error) {
	bc, err := calculateClockState(balls, minutes)
	if err != nil {
		return nil, err
	}
	return json.Marshal(bc)
}

func calculateClockState(balls, minutes int) (bc *ballclock.BallClock,
	err error) {
	bc, err = ballclock.CreateBallClock(balls)
	if err != nil {
		return nil, err
	}
	for i := 0; i < minutes; i++ {
		err := bc.Tick()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return bc, nil
}
