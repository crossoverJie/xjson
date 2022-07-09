package xjson

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrammarToken(t *testing.T) {
	grammar := `obj`
	tokenize, err := GrammarTokenize(grammar)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
func TestGrammarToken2(t *testing.T) {
	grammar := `1o1bj.obj1.name.111[99].1c1`
	tokenize, err := GrammarTokenize(grammar)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
	assert.Equal(t, tokenize[0].T, Key)
	assert.Equal(t, tokenize[0].Value, "1o1bj")
	assert.Equal(t, tokenize[1].T, Dot)
	assert.Equal(t, tokenize[1].Value, ".")
	assert.Equal(t, tokenize[2].T, Key)
	assert.Equal(t, tokenize[2].Value, "obj1")
	assert.Equal(t, tokenize[3].T, Dot)
	assert.Equal(t, tokenize[3].Value, ".")
	assert.Equal(t, tokenize[4].T, Key)
	assert.Equal(t, tokenize[4].Value, "name")
	assert.Equal(t, tokenize[5].T, Dot)
	assert.Equal(t, tokenize[5].Value, ".")
	assert.Equal(t, tokenize[6].T, Key)
	assert.Equal(t, tokenize[6].Value, "111")
	assert.Equal(t, tokenize[7].T, BeginArrayIndex)
	assert.Equal(t, tokenize[7].Value, "[")
	assert.Equal(t, tokenize[8].T, ArrayIndex)
	assert.Equal(t, tokenize[8].Value, "99")
	assert.Equal(t, tokenize[9].T, EndArrayIndex)
	assert.Equal(t, tokenize[9].Value, "]")
}

const input = `
{
    "a1": {
        "1b": [
            {
                "1c1": "c"
            },
            {
                "d": "d"
            }
        ]
    }
}`

func TestGrammarToken3(t *testing.T) {
	// a1.1b[1].1c1
	grammar := `a.b.list[10]`
	tokenize, err := GrammarTokenize(grammar)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
	assert.Equal(t, tokenize[0].T, Key)
	assert.Equal(t, tokenize[0].Value, "a")
	assert.Equal(t, tokenize[1].T, Dot)
	assert.Equal(t, tokenize[1].Value, ".")
	assert.Equal(t, tokenize[2].T, Key)
	assert.Equal(t, tokenize[2].Value, "b")
	assert.Equal(t, tokenize[3].T, Dot)
	assert.Equal(t, tokenize[3].Value, ".")
	assert.Equal(t, tokenize[4].T, Key)
	assert.Equal(t, tokenize[4].Value, "list")
	assert.Equal(t, tokenize[5].T, BeginArrayIndex)
	assert.Equal(t, tokenize[5].Value, "[")
	assert.Equal(t, tokenize[6].T, ArrayIndex)
	assert.Equal(t, tokenize[6].Value, "10")
	assert.Equal(t, tokenize[7].T, EndArrayIndex)
	assert.Equal(t, tokenize[7].Value, "]")
}
func TestGrammarTokenReader(t *testing.T) {
	grammar := `a.b.list[10]`
	tokenize, err := GrammarTokenize(grammar)
	assert.Nil(t, err)
	reader := NewGrammarTokenReader(tokenize)
	for reader.HasNext() {
		read := reader.Read()
		fmt.Printf("%s \t %s \n", read.T, read.Value)
	}

	grammar = `l[10.a`
	_, err = GrammarTokenize(grammar)
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestGrammarTokenize(t *testing.T) {
	str := "a\\.b\\."
	tokenize, err := GrammarTokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}

	str = "1a\\.\\["
	tokenize, err = GrammarTokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}

	str = "\\."
	tokenize, err = GrammarTokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
