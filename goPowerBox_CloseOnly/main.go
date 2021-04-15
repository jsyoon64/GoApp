package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {

	l, err := net.Listen("tcp", ":3000")
	if nil != err {
		log.Println(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}
		defer conn.Close()
		go ConnHandler(conn)
	}
}

func ConnHandler(conn net.Conn) {
	recvBuf := make([]byte, 1024)
	n, err := conn.Read(recvBuf)
	if nil != err {
		/*
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
		*/
		return
	}

	if 0 < n {
		data := recvBuf[:n]
		id := dataProcessing(data)
		conn.Close()
		fmt.Println("ID:", id, "Force to close")
	}
}

func dataProcessing(data []byte) int {
	id := int(binary.LittleEndian.Uint32(data[4:8]))
	return id
}
