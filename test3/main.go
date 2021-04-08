package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

type OperType struct {
	con1, con2, usb1, usb2 uint8
}

type FixedPart struct {
	model, site, group, id, DevSt, BrOv, MsgT int
	oper                                      OperType
	ctr                                       interface{}
}

func (o *OperType) conf(i uint8) {
	temp := i

	o.con1 = temp & 1

	temp >>= 1
	o.con2 = i & 1

	temp >>= 1
	o.usb1 = i & 1

	temp >>= 1
	o.usb2 = i & 1
}

var idIndex map[int]*FixedPart

func main() {

	idIndex = make(map[int]*FixedPart)
	input := make(chan string)

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
		go keyReader(input)

		select {
		case i, ok := <-input:
			if ok {
				fmt.Printf("Input : %v\n", i)
			} else {
				fmt.Println("Channel closed")
			}
		default:
			//case <-time.After(100 * time.Millisecond):
			//fmt.Println("Time out!")
		}
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
		if data[0] < 2 {
			id := dataProcessing(data)
			data1 := getconnRespData(data, id)
			_, err = conn.Write(data1)
		}
	}
}

func getconnRespData(data []byte, id int) []byte {
	if fptr, ok := idIndex[id]; ok {
		defer func() { fptr.ctr = nil }()

		switch v := fptr.ctr.(type) {
		case []byte:
			return v
		case string:
			return []byte(v)
		}
	}
	return data[0:14]
}

func makeControlData(id int, s string) []byte {
	ctrmsg := make([]byte, 14) // model:1, site:2, group:1, id:4, op:1, status:1, break:1, msgtype:1,length:2

	binary.LittleEndian.PutUint32(ctrmsg[0:], 0x0001)
	binary.LittleEndian.PutUint32(ctrmsg[4:], uint32(id))
	binary.LittleEndian.PutUint32(ctrmsg[8:], 0x0000)
	binary.LittleEndian.PutUint16(ctrmsg[12:], 0x00)
	ctrmsg[11] = 0xA0
	binary.LittleEndian.PutUint16(ctrmsg[12:], uint16(len(s)))

	//You need to use "..." as suffix in order to append a slice to another slice.
	//a := []byte("hello")
	//s := "world"
	//a = append(a, s...) // use "..." as suffice
	ctrmsg = append(ctrmsg, s...)
	return ctrmsg
}

func dataProcessing(data []byte) int {
	var fpart FixedPart

	id := int(binary.LittleEndian.Uint32(data[4:8]))
	if fptr, ok := idIndex[id]; ok {
		fpart = *fptr
	} else {
		fpart = FixedPart{}
		//fpart.index = index
		idIndex[id] = &fpart

		payload := "{\"CON1\":0}"
		//payload := "{\"CON1\":1}"
		fpart.ctr = makeControlData(id, payload)
		fmt.Println(fpart.ctr)
	}

	fpart.model = int(data[0])
	fpart.site = int(binary.LittleEndian.Uint16(data[1:3]))
	fpart.id = id
	//site := 0
	fpart.group = int(data[3])
	fpart.oper.conf(uint8(data[8]))
	fpart.DevSt = int(data[9])
	fpart.BrOv = int(data[10])
	fpart.MsgT = int(data[11])

	fmt.Println("Devices:", len(idIndex), fpart)
	return id
}

func keyReader(input chan string) {
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	input <- text
}

// for Graphic
/*
func makeDisplay(id int) {
	if fptr, ok := range idIndex[id];ok {
	for _, val := range idIndex {
		// val is pointer of FixedPart structure
	}
}
*/
