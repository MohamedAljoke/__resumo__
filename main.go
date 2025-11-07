package main

import "fmt"

type Teste struct {
	Name   string
	Age    int
	Mapper map[string]int
}

func main() {

	te := Teste{}

	fmt.Printf("te %+1v", te)
}
