package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type ctime struct {
	data     []byte
	epoch868 time.Time

	mtx sync.RWMutex
}

func NewCTime() (*ctime, error) {
	tstart, err := time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")
	if err != nil {
		return nil, err
	}

	return &ctime{
		make([]byte, 4),
		tstart,
		sync.RWMutex{},
	}, nil
}

func (ct *ctime) send(udpconn *net.UDPConn, caddr *net.UDPAddr) error {
	ct.mtx.RLock()

	_, err := udpconn.WriteToUDP(ct.data, caddr)
	if err != nil {
		ct.mtx.RUnlock()
		return err
	}

	ct.mtx.RUnlock()
	return nil
}

func (ct *ctime) update() {
	ct.mtx.Lock()

	ct.mtx.Unlock()
}

func ServeTime(addr string) error {
	tx := make([]byte, 4)
	timer, err := NewCTime()
	if err != nil {
		return err
	}

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

		err = timer.send(udpconn, caddr)
		if err != nil {
			fmt.Println("error: " + err.Error())
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
