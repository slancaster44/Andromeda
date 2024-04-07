package tokenizer

import (
	"strings"
)

type Tokenizer struct {
	fileContents string
	filename     string
	curLine      int
	curPos       int
	curTokens    []Token
}

func Tokenize(s string, filename string) []Token {
	s = strings.ToLower(s)
	tokenizer := &Tokenizer{
		fileContents: s,
		filename:     filename,
		curLine:      1,
		curPos:       0,
		curTokens:    make([]Token, 0),
	}

	var aux func(tokenizer *Tokenizer) []Token
	aux = func(t *Tokenizer) []Token {
		if t.curPos >= len(s) {
			return t.curTokens
		}

		curChar := s[t.curPos]
		possibleId, ok := singleCharMap[curChar]
		if ok {
			tok := Token{ID: possibleId, Contents: string(curChar), Filename: t.filename, LineNumber: t.curLine}
			if possibleId == TOK_NEWLINE {
				t.curLine++
			}

			t.curTokens = append(t.curTokens, tok)
			t.curPos++
			return aux(t)
		} else if curChar == ' ' || curChar == '\r' || curChar == '\t' {
			t.curPos++
			return aux(t)
		} else if isAlpha(curChar) {
			t.tokenizeIdent()
			return aux(t)
		} else if curChar == '\'' || curChar == '"' {
			t.tokenizeString()
			return aux(t) //skip over quote and most recently added char
		} else if isNum(curChar) {
			t.tokenizeNumber()
			return aux(t)
		} else if curChar == ';' {
			for t.curPos < len(s) && s[t.curPos] != '\n' {
				t.curPos++
			}
			return aux(t)
		}

		panic("Untokenizable character '" + string(curChar) + "'")
	}

	return aux(tokenizer)
}

func isNum(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isIn(c byte, set []byte) bool {
	for _, v := range set {
		if v == c {
			return true
		}
	}

	return false
}

func (t *Tokenizer) tokenizeIdent() {
	curChar := t.fileContents[t.curPos]

	identString := string(curChar)
	for t.curPos < len(t.fileContents)-1 && (isAlpha(curChar) || isNum(curChar) || curChar == '_') {
		t.curPos += 1
		curChar = t.fileContents[t.curPos]
		identString += string(curChar)
	}

	var tok Token //Handle some weird off-by-one errors
	if t.curPos == len(t.fileContents)-1 && (isAlpha(curChar) || isNum(curChar) || curChar == '_') {
		tok = Token{ID: TOK_IDENT, Contents: identString, Filename: t.filename, LineNumber: t.curLine}
		t.curPos++
	} else {
		tok = Token{ID: TOK_IDENT, Contents: identString[:len(identString)-1], Filename: t.filename, LineNumber: t.curLine}
	}

	possibleId, ok := keywordMap[tok.Contents]
	if ok {
		tok.ID = possibleId
	}

	t.curTokens = append(t.curTokens, tok)
}

func (t *Tokenizer) tokenizeString() {
	str := ""
	curChar := t.fileContents[t.curPos] //capture what the opening quote was

	for t.curPos != len(t.fileContents) && t.fileContents[t.curPos+1] != curChar {
		t.curPos++
		str += string(t.fileContents[t.curPos])
	}

	var tok Token
	if curChar == '\'' && len(str) == 1 {
		tok = Token{ID: TOK_CHAR, Contents: str, Filename: t.filename, LineNumber: t.curLine}
	} else {
		tok = Token{ID: TOK_STR, Contents: str, Filename: t.filename, LineNumber: t.curLine}
	}

	t.curTokens = append(t.curTokens, tok)
	t.curPos += 2
}

func (t *Tokenizer) tokenizeNumber() {
	acceptableChars := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	id := TOK_DEC_INT
	if t.curPos != len(t.fileContents)-1 && t.fileContents[t.curPos+1] == 'x' {
		acceptableChars = append(acceptableChars, []byte{'a', 'b', 'c', 'd', 'e', 'f'}...)
		id = TOK_HEX_INT
		t.curPos += 2
	} else if t.curPos != len(t.fileContents)-1 && t.fileContents[t.curPos+1] == 'b' {
		acceptableChars = []byte{'0', '1'}
		id = TOK_BIN_INT
		t.curPos += 2
	}

	out := ""
	for t.curPos != len(t.fileContents) && isIn(t.fileContents[t.curPos], acceptableChars) {
		out += string(t.fileContents[t.curPos])
		t.curPos++
	}

	tok := Token{ID: id, Contents: out, Filename: t.filename, LineNumber: t.curLine}
	t.curTokens = append(t.curTokens, tok)
}

var singleCharMap = map[byte]TokenID{
	'.':  TOK_DOT,
	'\n': TOK_NEWLINE,
	'(':  TOK_LPAREN,
	')':  TOK_RPAREN,
	':':  TOK_COLON,
	',':  TOK_COMMA,
	'$':  TOK_DOLLAR,
	'-':  TOK_MINUS,
	'^':  TOK_CARROT,
}

var keywordMap = map[string]TokenID{
	"org":     TOK_DIR,
	"equ":     TOK_DIR,
	"dw":      TOK_DIR,
	"include": TOK_DIR,
	"str":     TOK_DIR,
	"pad":     TOK_DIR,
	"lda":     TOK_INS,
	"sta":     TOK_INS,
	"add":     TOK_INS,
	"nnd":     TOK_INS,
	"xor":     TOK_INS,
	"sub":     TOK_INS,
	"jsr":     TOK_INS,
	"jmp":     TOK_INS,
	"jnz":     TOK_INS,
	"jns":     TOK_INS,
	"hlt":     TOK_INS,
	"nop":     TOK_INS,
	"inp":     TOK_INS,
	"out":     TOK_INS,
	"def":     TOK_DEF,
	"subdef":  TOK_SUBDEF,
	"imm":     TOK_ADDR_MODE,
	"ind":     TOK_ADDR_MODE,
	"dir":     TOK_ADDR_MODE,
	"inc":     TOK_ADDR_MODE,
	"dec":     TOK_ADDR_MODE,
	"rel":     TOK_ADDR_MODE,
	"off":     TOK_ADDR_MODE,
}
