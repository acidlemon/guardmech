package main

import (
	"fmt"
	"log"

	"github.com/acidlemon/guardmech"
)

func main() {
	fmt.Println("Hello World!")

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	m := guardmech.New()
	m.Run()
}
