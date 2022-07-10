package xjson

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitStatus(t *testing.T) {
	str := `{"name":"cj", "age":10}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
func TestToken(t *testing.T) {
	str := `{"name":"cj", "object":{"age":10, "sex":"girl"}}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
func TestToken2(t *testing.T) {
	str := `{"name":"cj", "object":{"age":10, "sex":"girl"},"list":[{"1":"a"},{"2":"b"}]}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
func TestToken3(t *testing.T) {
	str := `{"k":"89","name":true,"age":false,"value":12}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}

func TestToken4(t *testing.T) {
	str := `{"k":null,"name":true,"age":false,"value":12, "x": null}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
func TestToken5(t *testing.T) {
	str := `{"k":"89","name":10.2}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}

func TestTokenErr(t *testing.T) {
	str := `{"k":nul}`
	_, err := Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":t}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":tr}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":tru}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":f}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":fa}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":fal}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":fals}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":n}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":nu}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":nul}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	str = `{"k":10.20.2}`
	_, err = Tokenize(str)
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestEscapeToken(t *testing.T) {
	str := `{"a":"a"}`
	tokenize, err := Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
	str = `{"a":"{\"a\":\"a\"}"}`
	fmt.Println(str)
	tokenize, err = Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
	str = `{"a":"{\"a\":[1.2]}"}`
	fmt.Println(str)
	tokenize, err = Tokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}
