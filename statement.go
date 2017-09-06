package main

import (
	"errors"
	"fmt"
	"strings"
)

// statementType type
type statementType int

const (
	STATEMENT_INSERT statementType = iota
	STATEMENT_SELECT
)

type statement struct {
	stmt     string
	stmtType statementType
}

func prepareStatement(stmt string) (statement, error) {
	if strings.HasPrefix(stmt, "insert") {
		return statement{stmt: stmt, stmtType: STATEMENT_INSERT}, nil
	}
	if strings.HasPrefix(stmt, "select") {
		return statement{stmt: stmt, stmtType: STATEMENT_SELECT}, nil
	}
	return statement{}, errors.New("Invalid statement")
}

func (s statement) executeStatement() {
	switch s.stmtType {
	case STATEMENT_INSERT:
		var table, col1, col2, name, email string
		found, err := fmt.Sscanf(s.stmt, "insert into %s %s %s values %s %s", &table, &col1, &col2, &name, &email)
		if err != nil || found < 5 {
			fmt.Printf("Syntax error, %d\n", found)
			return
		}
		tableReal, err := backend.getTable(table)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		id, err := tableReal.insert(map[string]string{
			col1: name,
			col2: email,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Created %s with id %d\n", table, id)
	case STATEMENT_SELECT:
		fmt.Println("Selecting...")
	}
}
