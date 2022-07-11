package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if len(from)*len(to) == 0 {
		fmt.Println("Arguments exception. Provide valid 'from' and 'to' flags.")
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println("failed: ", err)
		return
	}
	fmt.Println("copied")
}
