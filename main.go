package main

import (
	"fmt"

	"github.com/abhra303/qDNS/listener"
)

func main() {
	fmt.Println("starting server")
	listener.PortListener(53)
}
