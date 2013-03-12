// goconf project main.go
package main

import (
	_ "fmt"
	"log"
)

func main() {
	c, _ := LoadConfigFile("test.ini")
	c.SetKeyComments("Demo", "key1", "")
	SaveConfigFile(c, "test1.ini")
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
