package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mrkaspa/sqlite/backend"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("db >> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.HasPrefix(text, ".") {
			switch text {
			case ".exit":
				os.Exit(0)
				return
			default:
				fmt.Println("Unrecognized command")
				continue
			}
		}

		err := backend.ExecuteStatement(text)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
	}
}
