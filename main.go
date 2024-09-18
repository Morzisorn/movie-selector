package main

import (
	"fmt"
)

func main() {
	go StartServer()
	StartBot()
	fmt.Println("Server and bot started")
}
