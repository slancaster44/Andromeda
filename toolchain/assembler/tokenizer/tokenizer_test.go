package tokenizer

import (
	"testing"
)

func TestSingleChar(t *testing.T) {
	testStrings := []string{
		"\n(.",
		"\n", "(",
	}

	expected := [][]Token{
		{Token{ID: TOK_NEWLINE, Contents: "\n"}, Token{ID: TOK_LPAREN, Contents: "("}, Token{ID: TOK_DOT, Contents: "."}},
		{Token{ID: TOK_NEWLINE, Contents: "\n"}},
		{Token{ID: TOK_LPAREN, Contents: "("}},
	}

	for i, v := range testStrings {
		result := Tokenize(v)

		for j, v2 := range expected[i] {
			if result[j] != v2 {
				t.Fatalf("Expected '%v' got '%v'", v2, result[j])
			}
		}
	}
}

func TestIdent(t *testing.T) {
	testString := "(hello)"
	expected := []Token{
		{ID: TOK_LPAREN, Contents: "("},
		{ID: TOK_IDENT, Contents: "hello"},
		{ID: TOK_RPAREN, Contents: ")"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Fatalf("%v", v)
		}
		t.Fatalf("Failed to tokenize identifiers")
	}
}

func TestKeyword(t *testing.T) {
	testString := "(equ lda)"
	expected := []Token{
		{ID: TOK_LPAREN, Contents: "("},
		{ID: TOK_DIR, Contents: "equ"},
		{ID: TOK_INS, Contents: "lda"},
		{ID: TOK_RPAREN, Contents: ")"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Fatalf("%v", v)
		}
		t.Fatalf("Failed to tokenize identifiers")
	}
}

func TestString(t *testing.T) {
	testString := "(\"Hello\")."
	expected := []Token{
		{ID: TOK_LPAREN, Contents: "("},
		{ID: TOK_STR, Contents: "Hello"},
		{ID: TOK_RPAREN, Contents: ")"},
		{ID: TOK_DOT, Contents: "."},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Fatalf("%v", v)
		}
		t.Fatalf("Failed to tokenize strings")
	}
}

func TestInt(t *testing.T) {
	testString := "32 (0xf55) 0b10114"
	expected := []Token{
		{TOK_DEC_INT, "32"},
		{TOK_LPAREN, "("},
		{TOK_HEX_INT, "f55"},
		{TOK_RPAREN, ")"},
		{TOK_BIN_INT, "1011"},
		{TOK_DEC_INT, "4"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Fatalf("%v", v)
		}
		t.Fatalf("Failed to tokenize strings")
	}

	t.Log(results)
}

func TestComment(t *testing.T) {
	testString := ";Some Comment\n"
	expected := []Token{
		{TOK_NEWLINE, "\n"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Fatalf("%v", v)
		}
		t.Fatalf("Failed to tokenize strings")
	}
}

func tokListMatch(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
