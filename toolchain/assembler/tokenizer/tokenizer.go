package tokenizer

import "strings"

func Tokenize(s string) []Token {
	s = strings.ToLower(s) //

	var aux func(s string, curPos int, curTokens []Token) []Token
	aux = func(s string, curPos int, curTokens []Token) []Token {
		if curPos >= len(s) {
			return curTokens
		}

		curChar := s[curPos]
		possibleId, ok := singleCharMap[curChar]
		if ok {
			tok := Token{ID: possibleId, Contents: string(curChar)}
			curTokens = append(curTokens, tok)
			return aux(s, curPos+1, curTokens)
		} else if curChar == ' ' || curChar == '\r' || curChar == '\t' {
			return aux(s, curPos+1, curTokens)
		} else if isAlpha(curChar) {
			tokenizeIdent(&curPos, s, &curTokens)
			return aux(s, curPos, curTokens)
		} else if curChar == '\'' || curChar == '"' {
			tokenizeString(&curPos, s, &curTokens)
			return aux(s, curPos, curTokens) //skip over quote and most recently added char
		} else if isNum(curChar) {
			tokenizeNumber(&curPos, s, &curTokens)
			return aux(s, curPos, curTokens)
		} else if curChar == ';' {
			for curPos < len(s) && s[curPos] != '\n' {
				curPos++
			}
			return aux(s, curPos, curTokens)
		}

		panic("Untokenizable character '" + string(curChar) + "'")
	}

	return aux(s, 0, []Token{})
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

func tokenizeIdent(curPos *int, s string, curTokens *[]Token) {
	curChar := s[*(curPos)]

	identString := string(curChar)
	for *(curPos) < len(s)-1 && (isAlpha(curChar) || isNum(curChar) || curChar == '_') {
		*(curPos) += 1
		curChar = s[*(curPos)]
		identString += string(curChar)
	}

	var tok Token //Handle some weird off-by-one errors
	if *(curPos) == len(s)-1 && (isAlpha(curChar) || isNum(curChar) || curChar == '_') {
		tok = Token{ID: TOK_IDENT, Contents: identString}
		*(curPos)++
	} else {
		tok = Token{ID: TOK_IDENT, Contents: identString[:len(identString)-1]}
	}

	possibleId, ok := keywordMap[tok.Contents]
	if ok {
		tok.ID = possibleId
	}

	*(curTokens) = append(*(curTokens), tok)
}

func tokenizeString(curPos *int, s string, curTokens *[]Token) {
	str := ""
	curChar := s[*(curPos)] //capture what the opening quote was

	for *(curPos) != len(s) && s[*(curPos)+1] != curChar {
		*(curPos)++
		str += string(s[*(curPos)])
	}

	tok := Token{ID: TOK_STR, Contents: str}
	*(curTokens) = append(*(curTokens), tok)
	*(curPos) += 2
}

func tokenizeNumber(curPos *int, s string, curTokens *[]Token) {
	acceptableChars := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	id := TOK_DEC_INT
	if *(curPos) != len(s)-1 && s[*(curPos)+1] == 'x' {
		acceptableChars = append(acceptableChars, []byte{'a', 'b', 'c', 'd', 'e', 'f'}...)
		id = TOK_HEX_INT
		*(curPos) += 2
	} else if *(curPos) != len(s)-1 && s[*(curPos)+1] == 'b' {
		acceptableChars = []byte{'0', '1'}
		id = TOK_BIN_INT
		*(curPos) += 2
	}

	out := ""
	for *(curPos) != len(s) && isIn(s[*(curPos)], acceptableChars) {
		out += string(s[*(curPos)])
		*(curPos)++
	}

	tok := Token{ID: id, Contents: out}
	*(curTokens) = append(*(curTokens), tok)
}

var singleCharMap = map[byte]TokenID{
	'.':  TOK_DOT,
	'\n': TOK_NEWLINE,
	'(':  TOK_LPAREN,
	')':  TOK_RPAREN,
	':':  TOK_COLON,
	',':  TOK_COMMA,
	'$':  TOK_DOLLAR,
	'+':  TOK_PLUS,
	'-':  TOK_MINUS,
}

var keywordMap = map[string]TokenID{
	"org":     TOK_DIR,
	"equ":     TOK_DIR,
	"dw":      TOK_DIR,
	"include": TOK_DIR,
	"str":     TOK_DIR,
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
