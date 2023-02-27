package main

import (
	"fmt"
	"math/rand"
	"time"
)

func producer(id int, out chan<- int) {
	for {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		out <- id
	}
}

func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	for _, in := range inputs {
		go func(in <-chan int) {
			for {
				out <- <-in
			}
		}(in)
	}
	return out
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go producer(1, ch1)
	go producer(2, ch2)
	go producer(3, ch3)

	for i := range fanIn(ch1, ch2, ch3) {
		fmt.Println(i)
	}
}
