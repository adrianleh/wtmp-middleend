package main

import (
	"encoding/binary"
	"github.com/adrianleh/WTMP-middleend/command"
	"io"
	"io/ioutil"
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
	for {
		headerReader := io.LimitReader(conn, 25)
		cmdFrameHeader, err := ioutil.ReadAll(headerReader)
		if err != nil {
			log.Println(err)
			continue
		}
		sizeRaw := cmdFrameHeader[16+1 : 25]
		size := binary.BigEndian.Uint64(sizeRaw)
		dataReader := io.LimitReader(conn, int64(size))
		data, err := ioutil.ReadAll(dataReader)
		if err != nil {
			log.Println(err)
			continue
		}
		cmdFrame := append(cmdFrameHeader, data...)
		err = command.Submit(cmdFrame)
		if err != nil {
			log.Printf("Command failed: %v", err)
			continue
		}
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
