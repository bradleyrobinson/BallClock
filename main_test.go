package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBallClock(t *testing.T) {
	t.Run("input: 30, 125", func(t *testing.T) {
		result, err := clockState(30, 325)
		if err != nil {
			fmt.Printf("unsuccessful %v\n", err)
			return
		}
		resultStr := string(result)
		err = assertEqual(resultStr, "{\"Min\":[],\"FiveMin\":[22,13,25,3,7],\"Hour\":[6,12,17,4,15],\"Main\"[11,5,26,18,2,30,19,8,24,10,29,20,16,21,28,1,23,14,27,9]}")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestCycleDays(t *testing.T) {
	// TODO: Make this method simpler, since it's just doing one thing
	t.Run("input: 30", func(t *testing.T) {
		result, duration, err := cycleDays(30)
		if err != nil {
			fmt.Printf("unsuccessful %v\n", err)
			t.Fatal(err)
			return
		}
		err = assertEqual(15, result)
		if err != nil {
			fmt.Printf("unsuccessful %v\n", err)
			t.Fatal(err)
			return
		}
		fmt.Println("duration: ", duration)
	})
	t.Run("input: 45", func(t *testing.T) {
		result, duration, err := cycleDays(45)
		if err != nil {
			fmt.Printf("unsuccessful %v\n", err)
			t.Fatal(err)
			return
		}
		err = assertEqual(378, result)
		if err != nil {
			fmt.Printf("unsuccessful %v\n", err)
			t.Fatal(err)
			return
		}
		fmt.Println("duration: ", duration)
	})
}

func assertEqual(exp, got interface{}) error {
	if !reflect.DeepEqual(exp, got) {
		return fmt.Errorf("wanted %v, got %v", exp, got)
	}
	return nil
}
