package main

import "fmt"

func main() {

	test := "testsestset"
	fmt.Printf("%+v\n", test)

	server := NewServer("127.0.0.1", 9000)
	server.Start()
}
