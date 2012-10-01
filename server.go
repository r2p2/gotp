package main

import (
	"fmt"
	"net"
	"time"
)

func ServeTime(addr string) error {
	tx := make([]byte, 4)
	var txs int

	tstart, err := time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")
	if err != nil {
		return err
	}

	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpconn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return err
	}

	fps := fpsCounter()
	for {
		_, caddr, err := udpconn.ReadFromUDP(tx)
		if err != nil {
			fmt.Println("error: " + err.Error())
			continue
		}

		dnow := time.Since(tstart)
		now := uint(dnow.Seconds() + dnow.Minutes()*60 + dnow.Hours()*60*60)
		to_byte(now, &tx)

		txs, err = udpconn.WriteToUDP(tx, caddr)
		if err != nil {
			fmt.Println("error: " + err.Error())
			continue
		}

		if txs != 4 {
			fmt.Println("error: just sent", txs, "of 4 byte.")
			continue
		}

		fmt.Println("rps", fps())
	}

	err = udpconn.Close()
	if err != nil {
		return err
	}

	return nil
}
