package main

import (
	"fmt"
	"modules/mod1"
	"modules/mod2"
)

func main() {
	fmt.Println("Running Program")
	mod1.Test()
	mod2.Test()
	fmt.Println("Program Success")
}
