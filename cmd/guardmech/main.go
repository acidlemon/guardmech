package main

import (
	"fmt"

	"github.com/acidlemon/guardmech"
)

func main() {
	fmt.Println("Hello World!")

	m := guardmech.New()
	m.Run()
}
