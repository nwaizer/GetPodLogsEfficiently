package main

import (
	"log"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		log.Println("Good message")
		time.Sleep(1 * time.Second)
	}
	for {
		log.Println("Bad message")
		time.Sleep(1 * time.Second)
	}
}
