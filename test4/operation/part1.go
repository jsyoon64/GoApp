package operation

import (
	"bufio"
	"os"
)

type OperType struct {
	con1, con2, usb1, usb2 uint8
}

func (o *OperType) Conf(i uint8) {
	temp := i

	o.con1 = temp & 1

	temp >>= 1
	o.con2 = i & 1

	temp >>= 1
	o.usb1 = i & 1

	temp >>= 1
	o.usb2 = i & 1
}

func KeyReader(input chan string) {
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	input <- text
}
