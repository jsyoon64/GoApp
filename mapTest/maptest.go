package main

import "fmt"

//var PBoxLabel = map[string]interface{}{"ID": 0, "Model": 0, "CON1": 0, "CON2": 0, "USB1": 0, "USB2": 0,
var PBoxLabel = [...]string{"ID", "Model", "CON1", "CON2", "USB1", "USB2", "PwrCMD",
	"LED", "LedCMD", "C1V", "C1I", "C2V", "C2I", "U1V", "U1I", "U2V", "U2I"}

var PBoxmap map[string]interface{}

func main() {

	PBoxmap = make(map[string]interface{})

	for index, str := range PBoxLabel {
		fmt.Println("index:", index, "==>", str)
		PBoxmap[str] = 0
	}

	PBoxmap["CON1"] = "x"
	for key, element := range PBoxmap {
		fmt.Println("Key:", key, "==>", element)
	}
}
