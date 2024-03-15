package assembler

type TokenID byte
type Token struct {
	ID       TokenID
	Contents string
}

func (t Token) IsInt() bool {
	return t.ID == TOK_HEX_INT || t.ID == TOK_BIN_INT || t.ID == TOK_DEC_INT
}

const (
	TOK_NEWLINE TokenID = iota // '\n'
	TOK_IDENT                  // a name
	TOK_LPAREN                 // '('
	TOK_RPAREN                 // ')'
	TOK_COMMA                  // ','
	TOK_COLON                  // ':'
	TOK_DOT                    // '.'
	TOK_INS                    // LDA, JSR, etc
	TOK_DEC_INT                //  0b<num>
	TOK_HEX_INT                // -11, 23, 44, etc
	TOK_BIN_INT                //0b<num>
	TOK_STR                    // '<stuff>'
	TOK_LBRACK                 // '['
	TOK_RBRACK                 // ']'
	TOK_PERCENT                // '%'
	TOK_MINUS                  //'-'
	TOK_DIR                    // set, pad, fill, etc
	TOK_DEF
	TOK_SUBDEF
)
