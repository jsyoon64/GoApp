package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type FixedPart struct {
	index                                     int
	model, site, group, id, DevSt, BrOv, MsgT int
	oper                                      OperType
}

type OperType struct {
	con1, con2, usb1, usb2 uint8
}

// idIndex[id] = 생성 index
//var idIndex map[int]int
var idIndex map[int]*FixedPart

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
	for {
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
			if data[0] < 2 {
				dataProcessing(data)
				_, err = conn.Write(data[0:14])
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				return
			}
		}
	}
}

func dataProcessing(data []byte) {
	model := int(data[0])
	site := binary.LittleEndian.Uint16(data[1:3])
	//site := 0
	group := int(data[3])
	id := binary.LittleEndian.Uint32(data[4:8])
	operation := int(data[8])
	devstatus := int(data[9])
	breakOver := int(data[10])
	msgtype := int(data[11])
	fmt.Printf("ID:%8d M:%x S:%x G:%x O:%x D:%x B:%x Type:%x\n", id, model, site, group, operation, devstatus, breakOver, msgtype)
}
