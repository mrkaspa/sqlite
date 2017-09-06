package main

import (
	"errors"
	"fmt"
)

const pageSize = 100

type db struct {
	tables []*table
}

type table struct {
	idSeq          int
	tableName      string
	currentPageIdx int
	pages          []*page
}

type page struct {
	pageID        int
	currentRowIdx int
	rows          map[int]interface{}
}

var backend = &db{}

func init() {
	tuser := initTable("user")
	backend.tables = append(backend.tables, tuser)
}

func initPage(pageID int) *page {
	p := page{
		pageID:        pageID,
		currentRowIdx: 0,
		rows:          make(map[int]interface{}),
	}
	start := pageID * pageSize
	for i := start; i <= start+pageSize; i++ {
		p.rows[i] = nil
	}
	return &p
}

func initTable(name string) *table {
	return &table{
		idSeq:          0,
		tableName:      name,
		currentPageIdx: 0,
		pages:          []*page{initPage(0)},
	}
}

func (d *db) getTable(name string) (*table, error) {
	for _, t := range d.tables {
		if t.tableName == name {
			return t, nil
		}
	}
	return nil, fmt.Errorf("Table %s not found", name)
}

func (t *table) insert(data map[string]string) (int, error) {
	p := t.pages[t.currentPageIdx]
	if p.currentRowIdx == (pageSize - 1) {
		np := initPage(t.currentPageIdx + 1)
		t.pages = append(t.pages, np)
		p = np
	}
	id, err := p.insert(t.idSeq, data)
	if err != nil {
		return -1, err
	}
	t.idSeq = id + 1
	return id, nil
}

func (p *page) insert(id int, data map[string]string) (int, error) {
	row, ok := p.rows[id]
	if !ok {
		return -1, errors.New("Memory outbounds error")
	}
	if row != nil {
		return -1, errors.New("Row already inserted")
	}
	p.rows[id] = data
	return id, nil
}
