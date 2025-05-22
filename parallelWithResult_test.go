package controlFlow

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type A struct {
	V string
}

func ParallelWithResult1() (*A, error) {
	time.Sleep(time.Second * 1)
	return &A{
		V: "ParallelWithResult1 executed successfully",
	}, nil
}

func ParallelWithResult2() (*A, error) {
	time.Sleep(time.Second * 2)
	return nil, errors.New("Error occurred in ParallelWithResult2")
}

func ParallelWithResult3() (*A, error) {
	time.Sleep(time.Second * 1)
	panic("panic 3")
	return &A{
		V: "ParallelWithResult3 executed successfully",
	}, nil
}

func TestParallelWithResult(t *testing.T) {
	results, err := parallelWithResult[A](ParallelWithResult1, ParallelWithResult2, ParallelWithResult3)
	if err != nil {
		fmt.Println(err)
	}
	for i, v := range results {
		fmt.Printf("Function %d result: %v\n", i+1, v)
	}
	time.Sleep(time.Second * 1111)
}
