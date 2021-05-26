package main

import (
	"os"

	"github.com/Basic-Components/jwttools/serv"
)

func main() {
	serv.Endpoint.Parse(os.Args)
}
