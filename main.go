package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abhra303/qDNS/listener"
)

func main() {
	var port int
	var err error
	arguments := os.Args

	if len(arguments) == 2 {
		port, err = strconv.Atoi(arguments[1])

		if err != nil {
			fmt.Println("the given port number is not an integer:", arguments[1])
			return
		}
	} else {
		port = listener.DefaultPort
	}

	fmt.Printf("starting server at port %v ...\n", port)
	listener.PortListener(port)
}
