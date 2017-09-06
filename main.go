package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

		stmt, err := prepareStatement(text)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}

		stmt.executeStatement()
	}
}
