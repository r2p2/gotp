package main

import (
	"flag"
	"fmt"
	"github.com/r2p2/rfc868"
)

var daemon = flag.Bool("daemon", false, "run application as daemon. share your time.")
var address = flag.String("address", "localhost:1024", "connect/listen to/on address.")
var dry = flag.Bool("dry", false, "on dry run time won't be updates. atm every run is a dry run.")

func main() {
	flag.Parse()

	if *daemon == true {
		err := rfc868.ServeTime(*address)
		if err != nil {
			fmt.Println("error: " + err.Error())
			return
		}
	} else {
		time, err := rfc868.RequestTime(*address)
		if err != nil {
			fmt.Println("error: " + err.Error())
			return
		}

		if /*dry ==*/ true {
			fmt.Println(time)
		}
	}
}
