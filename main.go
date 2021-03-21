package main

import "fmt"

func main() {
	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	fmt.Printf("cccccccccccccccccccc")
	go sum(s[:len(s)/2], c)
	test := "testsestset"
	fmt.Printf("%+v\n", test)
	server := NewServer("127.0.0.1", 9000)
	server.Start()
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}
