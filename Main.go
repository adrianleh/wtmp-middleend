package main

import (
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

func DialSocketAndSend(path string, message []byte) error {
	sock, err := net.Dial("unix", path)
	if err != nil {
		return err
	}
	_, err = sock.Write(message)
	return err
}

func server(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 0)
	for {
		buf := make([]byte, 512)
		nr, err := conn.Read(buf)
		if err != nil {
			break
		}
		data = append(data, buf[0:nr]...)
	}
	frame, err := command.ParseCommandFrame(data)
	err = command.Handle(frame)
	if err != nil {
		log.Fatal(err)
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
