package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func fanOut(jobs <-chan int, workerCount int) <-chan int {
	results := make(chan int)

	for i := 0; i < workerCount; i++ {
		go worker(i, jobs, results)
	}

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				close(results)
				return
			}
		}
	}()

	return results
}

func main() {
	const jobCount = 10

	jobs := make(chan int, jobCount)
	results := fanOut(jobs, 3)

	for j := 1; j <= jobCount; j++ {
		jobs <- j
	}

	close(jobs)

	for r := range results {
		fmt.Println("result:", r)
	}
}
