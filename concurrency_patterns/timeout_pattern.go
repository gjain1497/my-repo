package main

// import (
// 	"errors"
// 	"fmt"
// 	"time"
// )

// var ErrTimeout = errors.New("operation timed out")

// func DoWorkWithTimeout(work func() string, timeout time.Duration) (string, error) {
// 	//create a channel to receive the result
// 	resultChan := make(chan string, 1) //buffered channel of size 1

// 	//run work in a seperate go routine
// 	go func() {
// 		result := work()     //call the work function
// 		resultChan <- result //send result to channel
// 	}()

// 	//wait for either result or timeout
// 	select {
// 	case result := <-resultChan:
// 		return result, nil
// 	case <-time.After(timeout):
// 		return "", ErrTimeout
// 	}
// }

// func main() {
// 	result, err := DoWorkWithTimeout(func() string {
// 		fmt.Println("Working for 1 second")
// 		time.Sleep(1 * time.Second)
// 		return "Success"
// 	}, 2*time.Second)

// 	fmt.Printf("Result: %s, Error: %v\n\n", result, err)

// 	fmt.Println("=== Test 2: Slow work (times out) ===")
// 	result, err = DoWorkWithTimeout(func() string {
// 		fmt.Println("Working for 3 seconds...")
// 		time.Sleep(3 * time.Second)
// 		return "Too slow!"
// 	}, 1*time.Second) // Timeout is 1 second

// 	fmt.Printf("Result: %s, Error: %v\n", result, err)
// }
