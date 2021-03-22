package main

import (
	"fmt"

	"github.com/kazuki0924/go-echo/src/router"
)

func main() {
	fmt.Println("Welcome to the webserver")
	e := router.New()

	e.Start(":8000")
}
