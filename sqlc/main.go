package main

import (
	"fmt"
	"time"
)

func helloworld() string {
	return "helloworld"
}


func hellobarman() string {
	return "you are a barman"
}


func main() {

	fmt.Print("hello wrold")

	go func() {
		fmt.Println("go routime runned ....")
		time.Sleep(1 * time.Microsecond)
	}()

 
