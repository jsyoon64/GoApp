package main

import (
	oper "goPowerBox/dataProcess"
	key "goPowerBox/keyProcess"
	"log"
	"net"
)

func main() {
	//oper.Init_dProcess()
	input := make(chan string, 10)

	l, err := net.Listen("tcp", ":3000")
	if nil != err {
		log.Println(err)
	}
	defer l.Close()

	go key.KeyReader(input)

	for {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}
		go ConnHandler(conn, input)
	}
}

func ConnHandler(conn net.Conn, input chan string) {

	defer conn.Close()
	recvBuf := make([]byte, 1024)
	n, err := conn.Read(recvBuf)
	if nil != err {
		return
	}

	if 0 < n {
		data := recvBuf[:n]
		if data[0] < 2 {
			//id := oper.DataProcessing(data, msg)
			id := oper.DataProcessing(data, input)
			data1 := oper.GetconnRespData(data, id)
			conn.Write(data1)
		}
	}
}
