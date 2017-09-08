package parsing

import (
	"fmt"
	"io"
)

// Statement represents a returned stament, it can be a select,
// insert, update or delete
type Statement interface{}

// SelectStatement represents a SQL SELECT statement.
type SelectStatement struct {
	Fields    []string
	TableName string
}

// InsertStatement represents a SQL INSERT statement.
type InsertStatement struct {
	Cols      []string
	Values    []string
	TableName string
}

type parseFn func(*Parser) (parseFn, error)

// Parser represents a parser.
type Parser struct {
	s    *Scanner
	stmt Statement
	buf  struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (Statement, error) {
	fn := getSentence
	for fn != nil {
		newFn, err := fn(p)
		if err != nil {
			return nil, err
		}
		fn = newFn
	}
	return p.stmt, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// Func to choose the type of the statement

func getSentence(p *Parser) (parseFn, error) {
	tok, lit := p.scanIgnoreWhitespace()
	switch tok {
	case SELECT:
		return selectSentence, nil
	case INSERT:
		return insertSentence, nil
	default:
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
}

// Funcs for insert

func insertSentence(p *Parser) (parseFn, error) {
	p.stmt = &InsertStatement{}
	return intoKeyword, nil
}

func intoKeyword(p *Parser) (parseFn, error) {
	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != INTO {
		return nil, fmt.Errorf("found %q, expected INTO", lit)
	}
	return getTableNameInsert, nil
}

func getTableNameInsert(p *Parser) (parseFn, error) {
	// Finally we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt := p.stmt.(*InsertStatement)
	stmt.TableName = lit
	return extractColumns, nil
}

func extractColumns(p *Parser) (parseFn, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != PAR_LEFT {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt := p.stmt.(*InsertStatement)
		stmt.Cols = append(stmt.Cols, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != PAR_RIGHT {
		return nil, fmt.Errorf("found %q, expected )", lit)
	}
	return valuesKeyword, nil
}

func valuesKeyword(p *Parser) (parseFn, error) {
	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != VALUES {
		return nil, fmt.Errorf("found %q, expected VALUES", lit)
	}
	return extractValues, nil
}

func extractValues(p *Parser) (parseFn, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != PAR_LEFT {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt := p.stmt.(*InsertStatement)
		stmt.Values = append(stmt.Values, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != PAR_RIGHT {
		return nil, fmt.Errorf("found %q, expected )", lit)
	}
	return nil, nil
}

// Funcs for instructions

func selectSentence(p *Parser) (parseFn, error) {
	p.stmt = &SelectStatement{}
	return extractFields, nil
}

func extractFields(p *Parser) (parseFn, error) {
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt := p.stmt.(*SelectStatement)
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	return fromKeyword, nil
}

func fromKeyword(p *Parser) (parseFn, error) {
	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}
	return getTableNameSelect, nil
}

func getTableNameSelect(p *Parser) (parseFn, error) {
	// Finally we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt := p.stmt.(*SelectStatement)
	stmt.TableName = lit
	return nil, nil
}
