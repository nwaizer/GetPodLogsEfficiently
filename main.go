package main

import (
	"log"
	"time"
)

func main() {
	for {
		log.Println("Demo message")
		time.Sleep(1 * time.Second)
	}
}
