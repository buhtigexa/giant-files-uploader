package main

import (
	"fmt"
	"os"
)

func main() {
	//server := NewStreamServer()
	//server.Start()

	f, err := os.Create("pepe")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	finfo, err := os.Stat("pepucho")
	fmt.Printf("%#v\n", finfo)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
		}
		panic(err)
	}
}
