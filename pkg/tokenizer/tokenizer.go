package tokenizer

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

// Tokenizer interface defines the methods a tokenizer should implement
type Tokenizer interface {
	CountTokens(text string) (int, error)
}

// TikTokenTokenizer implements the Tokenizer interface using the tiktoken library
type TikTokenTokenizer struct {
	encoding string
}

// NewTikTokenTokenizer creates a new TikTokenTokenizer
func NewTikTokenTokenizer(encoding string) *TikTokenTokenizer {
	return &TikTokenTokenizer{encoding: encoding}
}

// CountTokens counts the number of tokens in the given text
func (t *TikTokenTokenizer) CountTokens(text string) (int, error) {
	tke, err := tiktoken.GetEncoding(t.encoding)
	if err != nil {
		return 0, fmt.Errorf("error getting encoding: %v", err)
	}

	tokens := tke.Encode(text, nil, nil)
	return len(tokens), nil
}

// GetTokenizer returns a Tokenizer based on the specified encoding
func GetTokenizer(encoding string) (Tokenizer, error) {
	switch encoding {
	case "cl100k_base", "p50k_base", "r50k_base":
		return NewTikTokenTokenizer(encoding), nil
	case "cl100k":
		return NewTikTokenTokenizer("cl100k_base"), nil
	case "p50k":
		return NewTikTokenTokenizer("p50k_base"), nil
	case "r50k":
		return NewTikTokenTokenizer("r50k_base"), nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}
