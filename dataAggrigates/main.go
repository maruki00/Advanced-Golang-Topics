package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	now := time.Now()

	wg := &sync.WaitGroup{}

	age := make(chan int, 0)
	name := make(chan string, 0)
	go getName(name, wg)
	go getAge(age, wg)
	wg.Wait()
	// close(age)
	// close(name)
	fmt.Printf("%s, has %d years old.\n", <-name, <-age)
	fmt.Println("tooks : ", time.Since(now))
}

func getName(name chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	name <- "abdellah"
	wg.Done()
}

func getAge(age chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	age <- 20
	wg.Done()
}
