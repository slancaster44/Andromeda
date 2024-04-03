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
		result := Tokenize(v, "testing.123")

		for j, v2 := range expected[i] {
			if result[j].ID != v2.ID {
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

	results := Tokenize(testString, "testing.123")
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

	results := Tokenize(testString, "testing.123")
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Fatalf("%v", v)
		}
		t.Fatalf("Failed to tokenize identifiers")
	}
}

func TestString(t *testing.T) {
	testString := "(\"hello\")."
	expected := []Token{
		{ID: TOK_LPAREN, Contents: "("},
		{ID: TOK_STR, Contents: "hello"},
		{ID: TOK_RPAREN, Contents: ")"},
		{ID: TOK_DOT, Contents: "."},
	}

	results := Tokenize(testString, "testing.123")
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Logf("%v", v)
		}
		t.Fatalf("Failed to tokenize strings")
	}
}

func TestInt(t *testing.T) {
	testString := "32 (0xf55) 0b10114"
	expected := []Token{
		{TOK_DEC_INT, "32", "testing.123", 1},
		{TOK_LPAREN, "(", "testing.123", 1},
		{TOK_HEX_INT, "f55", "testing.123", 1},
		{TOK_RPAREN, ")", "testing.123", 1},
		{TOK_BIN_INT, "1011", "testing.123", 1},
		{TOK_DEC_INT, "4", "testing.123", 1},
	}

	results := Tokenize(testString, "testing.123")
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
		{TOK_NEWLINE, "\n", "testing.123", 1},
	}

	results := Tokenize(testString, "testing.123")
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
		if v.ID != b[i].ID || v.Contents != b[i].Contents {
			return false
		}
	}

	return true
}
