package dataProcess

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
)

type OperType struct {
	con1, con2, usb1, usb2 uint8
}

type FixedPart struct {
	isNew                                     bool
	model, site, group, id, DevSt, BrOv, MsgT int
	oper                                      OperType
}

type JsonPart struct {
	c1v, c1i, c2v, c2i, u1v, u1i, u2v, u2i interface{}
}

type RxedMessage struct {
	FixedPart
	JsonPart
	ctr interface{}
}

var IdIndex map[int]*RxedMessage

var JsonData map[string]interface{}

//
// Local
//
func (o *OperType) conf(i uint8) {

	if o.con1 = 0; i&1 != 0 {
		o.con1 = 1
	}
	if o.con2 = 0; i&2 != 0 {
		o.con2 = 1
	}
	if o.usb1 = 0; i&4 != 0 {
		o.usb1 = 1
	}
	if o.usb2 = 0; i&8 != 0 {
		o.usb2 = 1
	}
	/*
		temp := i

		o.con1 = temp & 1

		temp >>= 1
		o.con2 = temp & 1

		temp >>= 1
		o.usb1 = temp & 1

		temp >>= 1
		o.usb2 = temp & 1
	*/
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

func parseJsonData(fptr *RxedMessage, b []byte) {
	if b != nil {
		err := json.Unmarshal(b, &JsonData)
		if err != nil {
			dd := string(b)
			fmt.Println(dd, err)
		} else {
			if con, ok := JsonData["CON1"]; ok {
				fptr.c1v = con.(map[string]interface{})["V"].(float64)
				fptr.c1i = con.(map[string]interface{})["I"].(float64)
			} else {
				fptr.c1v = nil
				fptr.c1i = nil
			}
			if con, ok := JsonData["CON2"]; ok {
				fptr.c2v = con.(map[string]interface{})["V"].(float64)
				fptr.c2i = con.(map[string]interface{})["I"].(float64)
			} else {
				fptr.c2v = nil
				fptr.c2i = nil
			}
			if con, ok := JsonData["USB1"]; ok {
				fptr.u1v = con.(map[string]interface{})["V"].(float64)
				fptr.u1i = con.(map[string]interface{})["I"].(float64)
			} else {
				fptr.u1v = nil
				fptr.u1i = nil
			}
			if con, ok := JsonData["USB2"]; ok {
				fptr.u2v = con.(map[string]interface{})["V"].(float64)
				fptr.u2i = con.(map[string]interface{})["I"].(float64)
			} else {
				fptr.u2v = nil
				fptr.u2i = nil
			}
		}
	} else {
		fptr.c1v = nil
		fptr.c1i = nil
		fptr.c2v = nil
		fptr.c2i = nil
		fptr.u1v = nil
		fptr.u1i = nil
		fptr.u2v = nil
		fptr.u2i = nil
	}
}

func parseFixedData(fptr *RxedMessage, data []byte, id int) {
	if data == nil {
		return
	}
	fptr.model = int(data[0])
	fptr.site = int(binary.LittleEndian.Uint16(data[1:3]))
	fptr.id = id
	fptr.group = int(data[3])
	fptr.oper.conf(uint8(data[8]))
	fptr.DevSt = int(data[9])
	fptr.BrOv = int(data[10])
	fptr.MsgT = int(data[11])
}

//
// Global
//

func DataProcessing(data []byte, command chan string) int {
	var fptr *RxedMessage
	var ok bool
	var payload string

	select {
	case msg := <-command:
		msg = strings.TrimSuffix(msg, "\r\n")
		switch msg {
		case "on", "ON":
			payload = "{\"CON1\":1, \"USB1\":1, \"USB2\":1}"
		case "off", "OFF":
			payload = "{\"CON1\":0, \"USB1\":0, \"USB2\":0}"
		default:
			payload = ""
		}
	default:
	}

	id := int(binary.LittleEndian.Uint32(data[4:8]))
	if fptr, ok = IdIndex[id]; !ok {
		fptr = &RxedMessage{}
		IdIndex[id] = fptr
		fptr.isNew = true
		//fmt.Println(fptr.ctr)
	} else {
		fptr.isNew = false
	}

	if payload != "" {
		fptr.ctr = makeControlData(id, payload)
	} else {
		fptr.ctr = nil
	}

	//fmt.Println("fptr.ctr ", fptr.ctr)

	// parse fixed part data
	parseFixedData(fptr, data, id)

	// parse json data
	if msgLen := int(binary.LittleEndian.Uint16(data[12:14])); msgLen != 0 {
		parseJsonData(fptr, data[14:])
	} else {
		parseJsonData(fptr, nil)
	}

	fmt.Println("Devices:", len(IdIndex), *fptr)
	//fmt.Println(*fptr)
	return id
}

func GetconnRespData(data []byte, id int) []byte {
	if fptr, ok := IdIndex[id]; ok {
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

func init() {
	IdIndex = make(map[int]*RxedMessage)
	JsonData = make(map[string]interface{})
}
