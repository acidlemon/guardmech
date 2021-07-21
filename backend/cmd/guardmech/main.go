package main

import (
	"fmt"
	"log"

	guardmech "github.com/acidlemon/guardmech/backend"
)

func main() {
	fmt.Println("Hello World!")

	// setup
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	m := guardmech.New()
	m.Run()
}
