package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

var version = "v0.1.0"

const cmdName = "waitport"

type waitportOpts struct {
	Host string
	Port int
}

func main() {
	opts := &waitportOpts{}
	flag.IntVar(&opts.Port, "port", 0, "wait port")
	flag.StringVar(&opts.Host, "host", "127.0.0.1", "hostname")
	flag.Usage = func() {
		fmt.Println("Usage: " + cmdName + " [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := waitport_(opts); err != nil {
		log.Fatal(err)
	}
}

func checkPort(host string, port int) (bool, error) {
	tcp_addr, err := net.ResolveTCPAddr("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		return false, err
	}
	_, err = net.DialTCP("tcp4", nil, tcp_addr)
	if err == nil {
		return true, nil
	}
	return false, err
}

func waitPort(host string, port int, max_wait float64) (bool, error) {
	if max_wait <= 0 {
		max_wait = 10.
	}
	waiter := makeWaiter(max_wait)
	for waiter() {
		res, err := checkPort(host, port)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, errors.New("timeout")
}

func makeWaiter(max_wait float64) func() bool {
	waited := 0.
	sleep := time.Millisecond
	return func() bool {
		if max_wait >= 0 && waited > max_wait {
			return false
		}
		time.Sleep(sleep)
		waited += sleep.Seconds()
		sleep *= 2.
		return true
	}
}

func waitport_(opt *waitportOpts) error {
	if opt.Port == 0 {
		return errors.New("set port option")
	}
	if _, err := waitPort(opt.Host, opt.Port, 10); err == nil {
		fmt.Printf("%s:%d\n", opt.Host, opt.Port)
		return nil
	} else {
		return err
	}
}
