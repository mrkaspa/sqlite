package main

import (
	"bufio"
	"errors"
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

// StatementType type
type StatementType int

const (
	STATEMENT_INSERT StatementType = iota
	STATEMENT_SELECT
)

type statement struct {
	stmtType StatementType
}

func prepareStatement(stmt string) (statement, error) {
	if strings.HasPrefix(stmt, "insert") {
		return statement{stmtType: STATEMENT_INSERT}, nil
	}
	if strings.HasPrefix(stmt, "select") {
		return statement{stmtType: STATEMENT_SELECT}, nil
	}
	return statement{}, errors.New("Invalid statement")
}

func (s statement) executeStatement() {
	switch s.stmtType {
	case STATEMENT_INSERT:
		fmt.Println("Inserting...")
	case STATEMENT_SELECT:
		fmt.Println("Selecting...")
	}
}
