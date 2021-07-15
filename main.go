package main

import (
	"fmt"
	"time"
)

func main() {
	go Counter("Hippo")
	go Counter("Koala")

	time.Sleep(time.Second * 5)
}

func Counter(person string) {
	for i := 0; i< 10; i++ {
		fmt.Println(person, "is good", i)
		time.Sleep(time.Second)
	}
}