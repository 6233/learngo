package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)
	people := [2]string{"Hippo", "koala"}

	for _, person := range people {
		go isGood(person, c)
	}
	fmt.Println(<-c)
	fmt.Println(<-c)
	// fmt.Println(<-c) // deadlock
}

func isGood(person string, c chan bool) {
	time.Sleep(time.Second * 5)
	fmt.Println(person)
	c <- true
}