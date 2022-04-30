package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

const greetings string = "Hello, OTUS!"

func main() {
	greetingsReversed := stringutil.Reverse(greetings)
	fmt.Println(greetingsReversed)
}
