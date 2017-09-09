package backend

import (
	"fmt"
	"strings"

	"github.com/mrkaspa/sqlite/parsing"
)

func ExecuteStatement(text string) error {
	parser := parsing.NewParser(strings.NewReader(text))
	stmt, err := parser.Parse()
	if err != nil {
		return err
	}
	switch s := stmt.(type) {
	case *parsing.InsertStatement:
		tableReal, err := engine.GetTable(s.TableName)
		if err != nil {
			return err
		}
		data := make(map[string]string)
		for i := 0; i < len(s.Cols); i++ {
			data[s.Cols[i]] = s.Values[i]
		}
		id, err := tableReal.Insert(data)
		if err != nil {
			return err
		}
		fmt.Printf("Created %s with id %d\n", s.TableName, id)
	case *parsing.SelectStatement:
		fmt.Println("Selecting...")
	default:

	}
	return nil
}
