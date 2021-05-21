package main

import (
	"bytes"
	"github.com/adrianleh/WTMP-middleend/command"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const sockPath = "/tmp/wtmp.sock"

func main() {
	listener, err := startServer()
	cleanUpSocketOnExit()
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	accept(listener)
}

func server(conn net.Conn) {
	defer conn.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(conn); err != nil {
		log.Printf("Failed to read: %v", err)
	}
	data := buf.Bytes()
	if err := command.Submit(data); err != nil {
		log.Printf("Command failed: %v", err)
	}
}

func accept(listener net.Listener) {
	for {
		fd, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		go server(fd)
	}

}

func startServer() (net.Listener, error) {
	return net.Listen("unix", sockPath)
}

func cleanUpSocketOnExit() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel
		if err := os.Remove(sockPath); err != nil {
			os.Exit(3)
		}
		os.Exit(0)
	}()
}
