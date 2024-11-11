package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type telnetConfig struct {
	timeout time.Duration
	host    string
	port    int
}

var config telnetConfig

func init() {
	flag.DurationVar(&config.timeout, "timeout", 10*time.Second, "timeout for connection")
}

func main() {
	readConfig()

	telnetClient := NewTelnetClient(fmt.Sprintf("%v:%v", config.host, config.port), config.timeout, os.Stdin, os.Stdout)

	err := telnetClient.Connect()
	defer telnetClient.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case sig := <-sigs:
				fmt.Printf("Received sig: %v", sig)
				close(done)
				return
			default:
				err := telnetClient.Send()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to send: %v\n", err)
					close(done)
					return
				}
				err = telnetClient.Receive()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to receive: %v\n", err)
					close(done)
					return
				}
			}
		}
	}()

	<-done
}

func readConfig() {
	var err error
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Wrong number of arguments: Expected [options] host port\n")
		os.Exit(1)
	}
	config.host = args[0]
	config.port, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Port must be a number\n")
		os.Exit(1)
	}
}
