package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mx sync.Mutex	
	Val int
}

func (s *Counter) Increment() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Val++
}

func main() {
	counter := Counter{Val: 0}

	var wg sync.WaitGroup
	job := make(chan int)
	const workers = 5

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(&wg, &counter, job)
	}

	go prod(job)

	wg.Wait()

	fmt.Println(counter.Val)
}

func worker(wg *sync.WaitGroup, counter *Counter, job <-chan int) {
	defer wg.Done()

	for num := range job {
		time.Sleep(time.Duration(num) * time.Millisecond)
		// some work
		counter.Increment()
	}
}

func prod(job chan<- int) {
	// генерация работы
	defer close(job)
	for i := 0; i < 100; i++ {
		job<-10
	}
}
