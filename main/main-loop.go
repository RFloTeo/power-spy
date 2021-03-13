package main

import (
	"fmt"
	"time"
)

func mainLoop() {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		fmt.Println("Hi")
	}
}
