package ballclock

import (
	"encoding/json"
	"errors"
	"reflect"
)

type ballQueue interface {
	pop() int
	push(x int) ([]int, error)
}

// BallTimeQueue stores the information for each section of balls
type BallTimeQueue struct {
	Balls     []int
	name      string
	numBalls  int
	maxLength int
	nextQueue ballQueue
}

// MarshalJSON converts a queue into bytes
func (b *BallTimeQueue) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Balls[0:b.numBalls])
}

func createBallQueue(length int, name string) *BallTimeQueue {
	balls := make([]int, length, length)
	newQueue := BallTimeQueue{Balls: balls, maxLength: length,
		name: name}
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
	start := 0
	// This indicates that this is the hour queue
	if b.Balls[0] == -1 {
		start = 1
	}
	if b.numBalls+start >= b.maxLength {
		rvalues := make([]int, b.numBalls-start)
		slicedValues := b.Balls[start : b.numBalls+start]
		copy(rvalues, slicedValues)
		reverse(rvalues)
		if b.nextQueue == nil {
			return nil, errors.New("no queue to pass on")
		}
		b.numBalls = start
		extraValues, err := b.nextQueue.push(x)
		if err != nil {
			return nil, err
		}
		rvalues = append(rvalues, extraValues...)
		for i := range b.Balls[start:] {
			if i+i < b.maxLength {
				b.Balls[i+1] = 0
			}
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

type BallMainQueue struct {
	Balls     []int
	numBalls  int
	maxLength int
}

func (bm *BallMainQueue) pop() int {
	bm.numBalls--
	x := bm.Balls[0]
	bm.Balls = append(bm.Balls, 0)
	bm.Balls = bm.Balls[1:len(bm.Balls)]
	return x
}

func createBallMainQueue(n int) *BallMainQueue {
	balls := make([]int, n, n)
	newQueue := BallMainQueue{Balls: balls, maxLength: n}
	return &newQueue
}

func (bm *BallMainQueue) push(x int) ([]int, error) {
	if x == 0 {
		return nil, nil
	}
	bm.Balls[bm.numBalls] = x
	bm.numBalls++
	return nil, nil
}

// BallClock contains the structs and mechanisms for the queue
type BallClock struct {
	Min     *BallTimeQueue `json:"Min"`
	FiveMin *BallTimeQueue `json:"FiveMin"`
	Hour    *BallTimeQueue `json:"Hour"`
	Main    *BallMainQueue `json:"Main"`
}

// CreateBallClock creates a ball clock
func CreateBallClock(balls int) (bc *BallClock, err error) {
	if balls < 27 || balls > 127 {
		return nil, errors.New("ball amount out of valid range," +
			"must be a value within 27-127")
	}
	bc = &BallClock{
		Min:     createBallQueue(4, "Min"),
		FiveMin: createBallQueue(11, "FiveMin"),
		Hour:    createBallQueue(12, "Hour"),
		Main:    createBallMainQueue(balls),
	}
	bc.Min.nextQueue = bc.FiveMin
	bc.FiveMin.nextQueue = bc.Hour
	bc.Hour.nextQueue = bc.Main
	// The first ball is reserved for the hour
	for i := 1; i <= balls; i++ {
		bc.Main.push(i)
	}
	bc.Hour.push(-1)
	return bc, nil
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

// Cycled determines whether the ball clock is the same order
func (bc *BallClock) Cycled(starter []int) bool {
	return reflect.DeepEqual(bc.Main.Balls, starter) &&
		bc.Hour.numBalls == 1 && bc.Min.numBalls == 0 &&
		bc.FiveMin.numBalls == 0
}
