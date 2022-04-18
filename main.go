package main

import (
	"fmt"
	"os"
	"setdb/data"
)

func main() {
	// Passing arguments to this is interpreted as Junction
	// This can be used a a lambda to other processses
	d := data.Data{}
	d.Init(10)
	d.Generate()
	ars := data.GetArray(&d, os.Args[1:])
	res := data.Junction(ars)
	stringified := res.Stringify()

	for _, filename := range stringified {
		fmt.Println(filename)
	}
}
