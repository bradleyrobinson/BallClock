package ballclock

import (
	"encoding/json"
	"errors"
	"reflect"
)

type ballQueue interface {
	checklimit() int
}

// BallTimeQueue stores the information for each section of balls
type BallTimeQueue struct {
	Balls     []int
	numBalls  int
	maxLength int
	nextQueue *BallTimeQueue
}

// MarshalJSON converts a queue into bytes
func (b *BallTimeQueue) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Balls[0:b.numBalls])
}

func createBallQueue(length int) *BallTimeQueue {
	balls := make([]int, length, length)
	newQueue := BallTimeQueue{Balls: balls, maxLength: length}
	return &newQueue
}

func (b *BallTimeQueue) pop() int {
	b.numBalls--
	x := b.Balls[0]
	b.Balls = append(b.Balls, 0)
	b.Balls = b.Balls[1:len(b.Balls)]
	return x
}

// adds a value at the first 0 value
func (b *BallTimeQueue) push(x int) (rvalues []int, err error) {
	if b.numBalls >= b.maxLength {
		start := 0
		// This indicates that this is the minute queue
		if b.Balls[0] == -1 {
			start = 1
		}
		rvalues := make([]int, start+b.numBalls)
		slicedValues := b.Balls[start:b.numBalls]
		copy(rvalues, slicedValues)
		reverse(rvalues)
		if b.nextQueue == nil {
			return nil, errors.New("no queue to pass on")
		}
		b.numBalls = 0
		extraValues, err := b.nextQueue.push(x)
		if err != nil {
			return nil, err
		}
		rvalues = append(rvalues, extraValues...)
		for i := range b.Balls[start:] {
			b.Balls[i] = 0
		}
		return rvalues, nil
	}
	b.Balls[b.numBalls] = x
	b.numBalls++
	return nil, nil
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// BallClock contains the structs and mechanisms for the queue
type BallClock struct {
	Min     *BallTimeQueue `json:"Min"`
	FiveMin *BallTimeQueue `json:"FiveMin"`
	Hour    *BallTimeQueue `json:"Hour"`
	Main    *BallTimeQueue `json:"Main"`
}

// CreateBallClock creates a ball clock
func CreateBallClock(balls int) *BallClock {
	bc := BallClock{
		Min:     createBallQueue(4),
		FiveMin: createBallQueue(11),
		Hour:    createBallQueue(12),
		Main:    createBallQueue(balls),
	}
	bc.Min.nextQueue = bc.FiveMin
	bc.FiveMin.nextQueue = bc.Hour
	bc.Hour.nextQueue = bc.Main
	// The first ball is reserved for the hour
	for i := 1; i <= balls; i++ {
		bc.Main.push(i)
	}
	bc.Hour.Balls[0] = -1
	return &bc
}

// Tick adds one minute to the minute
func (bc *BallClock) Tick() error {
	// first pop the value from the main queue
	val := bc.Main.pop()
	// Place that value in the minute queue
	revValues, err := bc.Min.push(val)
	if err != nil {
		return err
	}
	if revValues != nil {
		for i := range revValues {
			_, err = bc.Main.push(revValues[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *BallClock) Cycled(starter []int) bool {
	return bc.Min.numBalls == 0 && bc.FiveMin.numBalls == 0 &&
		bc.Hour.numBalls == 0 && reflect.DeepEqual(bc.Main.Balls, starter)
}
