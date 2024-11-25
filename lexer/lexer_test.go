package lexer

import (
	"testing"

	"github.com/Hyuga-Tsuki/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}
}
