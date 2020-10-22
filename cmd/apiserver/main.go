package main

import (
	"flag"

	"github.com/SergioBravo/http-rest-api/cmd"
)

func main() {
	flag.Parse()
	cmd.Execute()
}
