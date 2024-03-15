package assembler

import "testing"

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
		Token{ID: TOK_LPAREN, Contents: "("},
		Token{ID: TOK_IDENT, Contents: "hello"},
		Token{ID: TOK_RPAREN, Contents: ")"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Logf("%v", v)
		}
		t.Fatalf("Failed to tokenize identifiers")
	}
}

func TestKeyword(t *testing.T) {
	testString := "(set lda)"
	expected := []Token{
		Token{ID: TOK_LPAREN, Contents: "("},
		Token{ID: TOK_DIR, Contents: "set"},
		Token{ID: TOK_INS, Contents: "lda"},
		Token{ID: TOK_RPAREN, Contents: ")"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Logf("%v", v)
		}
		t.Fatalf("Failed to tokenize identifiers")
	}
}

func TestString(t *testing.T) {
	testString := "[\"Hello\"]."
	expected := []Token{
		{ID: TOK_LBRACK, Contents: "["},
		{ID: TOK_STR, Contents: "Hello"},
		{ID: TOK_RBRACK, Contents: "]"},
		{ID: TOK_DOT, Contents: "."},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Logf("%v", v)
		}
		t.Fatalf("Failed to tokenize strings")
	}
}

func TestInt(t *testing.T) {
	testString := "32 0x55f 0b10114"
	expected := []Token{
		{TOK_DEC_INT, "32"},
		{TOK_HEX_INT, "55f"},
		{TOK_BIN_INT, "1011"},
		{TOK_DEC_INT, "4"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Logf("%v", v)
		}
		t.Fatalf("Failed to tokenize strings")
	}
}

func TestComment(t *testing.T) {
	testString := ";Some Comment\n"
	expected := []Token{
		{TOK_NEWLINE, "\n"},
	}

	results := Tokenize(testString)
	if !tokListMatch(results, expected) {
		for _, v := range results {
			t.Logf("%v", v)
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
