package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 20; i++ {
		fmt.Println(rand.Intn(3) + 1)
	}
}

func test1(func(int) error) {
	
}
