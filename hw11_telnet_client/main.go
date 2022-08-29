package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

const (
	minArgsNumber       = 3
	errMsgStderrInvalid = "can't write to stderr"
	usageMsg            = "Usage: go-telnet <host> <port> [--timeout=<timeout>]\n"
)

func init() {
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "timeout of the connection")
}

func main() {
	if len(os.Args) < minArgsNumber {
		log.Fatal(usageMsg)
	}
	flag.Parse()

	args := flag.Args()
	host := args[0]
	port := args[1]

	client := NewTelnetClient(
		net.JoinHostPort(host, port),
		timeout,
		os.Stdin,
		os.Stdout,
	)

	if err := client.Connect(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("%v\n", err)
		}
	}()
	_, err := fmt.Fprintf(os.Stderr, "...Connected to %s:%s\n", host, port)
	if err != nil {
		log.Panic(errMsgStderrInvalid)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go doJob(ctx, cancel, client.Receive)
	go doJob(ctx, cancel, client.Send)

	listenToSignal(ctx, cancel)
}

func doJob(ctx context.Context, cancel context.CancelFunc, job func() error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := job()
			errAndCtxActive := err != nil && ctx.Err() == nil
			if errAndCtxActive {
				_, err := fmt.Fprintf(os.Stderr, "...%v\n", err)
				if err != nil {
					log.Panic(errMsgStderrInvalid)
				}
				cancel()
			}
		}
	}
}

func listenToSignal(ctx context.Context, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	select {
	case <-sigCh:
		cancel()
		signal.Stop(sigCh)
	case <-ctx.Done():
		close(sigCh)
		return
	}
}
