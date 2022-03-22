package main

import (
	"log"
	"time"
)

func main() {
	for {
		for i := 0; i < 10; i++ {
			log.Println("Good message")
			time.Sleep(1 * time.Second)
		}
		for i := 0; i < 10; i++ {
			log.Println("Bad message")
			time.Sleep(1 * time.Second)
		}
	}
}
