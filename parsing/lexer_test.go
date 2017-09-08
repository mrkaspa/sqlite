package parsing_test

import (
	"reflect"
	"strings"
	"testing"

	ps "github.com/mrkaspa/sqlite/parsing"
)

type TokenLiteral = ps.TokenLiteral

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok ps.Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: ps.EOF},
		{s: `#`, tok: ps.ILLEGAL, lit: `#`},
		{s: ` `, tok: ps.WS, lit: " "},
		{s: "\t", tok: ps.WS, lit: "\t"},
		{s: "\n", tok: ps.WS, lit: "\n"},

		// Misc characters
		{s: `*`, tok: ps.ASTERISK, lit: "*"},

		// Identifiers
		{s: `foo`, tok: ps.IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: ps.IDENT, lit: `Zx12_3U_`},

		// Keywords
		{s: `FROM`, tok: ps.FROM, lit: "FROM"},
		{s: `SELECT`, tok: ps.SELECT, lit: "SELECT"},
	}

	for i, tt := range tests {
		s := ps.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}

func TestScanner_ScanText(t *testing.T) {
	var tests = []struct {
		s   string
		res []TokenLiteral
	}{
		{
			s: ``,
			res: []TokenLiteral{
				TokenLiteral{Tok: ps.EOF},
			},
		},
		{
			s: `SELECT * FROM user`,
			res: []TokenLiteral{
				TokenLiteral{Tok: ps.SELECT, Lit: "SELECT"},
				TokenLiteral{Tok: ps.WS, Lit: " "},
				TokenLiteral{Tok: ps.ASTERISK, Lit: "*"},
				TokenLiteral{Tok: ps.WS, Lit: " "},
				TokenLiteral{Tok: ps.FROM, Lit: "FROM"},
				TokenLiteral{Tok: ps.WS, Lit: " "},
				TokenLiteral{Tok: ps.IDENT, Lit: "user"},
				TokenLiteral{Tok: ps.EOF},
			},
		},
	}

	for i, tt := range tests {
		s := ps.NewScanner(strings.NewReader(tt.s))
		tks := s.ScanText()
		if !reflect.DeepEqual(tks, tt.res) {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q", i, tt.s, tt.res, tks)
		}
	}
}
