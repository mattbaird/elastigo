package elastigo

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestSanitizeReferenceSentence(t *testing.T) {
	input := "AND there! are? (lots of) char*cters 2 ^escape!"
	expectation := `\A\N\D there\! are\? \(lots of\) char\*cters 2 \^escape\!`

	actual := Sanitize(input)
	assert.Equal(t, expectation, actual)
}
