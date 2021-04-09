package keyProcess

import (
	"bufio"
	"os"
)

func KeyReader(input chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		input <- text
		//fmt.Println("Input is ", text)
	}
}
