// goconf project main.go
package main

import (
	_ "fmt"
	"log"
)

func main() {
	c, _ := LoadConfigFile("test.ini")
	SaveConfigFile(c, "test1.ini")
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
