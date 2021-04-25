package main

import (
	"log"
	"net"
)

const sockPath = "/tmp/wtmp.sock"

func main() {
	listener, err := startServer()
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	accept(listener)
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
	println("Server got:", string(data))
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
