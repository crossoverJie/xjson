package gjson

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArithmeticTokenize(t *testing.T) {
	str := "list2.obj2[0].obj3.list3[0] * (obj_list[1].age+obj.age) / obj.age"
	tokenize, err := ArithmeticTokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}

	str = "name+age"
	tokenize, err = ArithmeticTokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
}

func TestGetWithArithmetic(t *testing.T) {
	str := `{"name":"bob", "age":10,"magic":10.1, "score":{"math":[1,2]}}`
	result := GetWithArithmetic(str, "(age+age)*age+magic")
	assert.Equal(t, result.Float(), 210.1)
	result = GetWithArithmetic(str, "(age+age)*age")
	assert.Equal(t, result.Int(), 200)

	result = GetWithArithmetic(str, "(age+age) * age + score.math[0]")
	assert.Equal(t, result.Int(), 201)

	result = GetWithArithmetic(str, "(age+age) * age - score.math[0]")
	assert.Equal(t, result.Int(), 199)

	result = GetWithArithmetic(str, "score.math[1] / score.math[0]")
	assert.Equal(t, result.Int(), 2)

	result = GetWithArithmetic(str, "age")
	assert.Equal(t, result.Int(), 10)
	str = `{"people":[{"bob":{"age":10}},{"alice":{"age":10}}]}`
	age := GetWithArithmetic(str, "people[0].bob.age + people[1].alice.age")
	assert.Equal(t, age.Int(), 20)
}

func TestGetWithArithmetic2(t *testing.T) {
	str := `{"name":[}`
	x := GetWithArithmetic(str, "x")
	assert.Equal(t, x.String(), "")
}

func TestGetWithArithmeticErr(t *testing.T) {
	str := `[1,2]`
	x := GetWithArithmetic(str, "x")
	assert.Equal(t, x.String(), "")

	str = `{"age":true}`
	x = GetWithArithmetic(str, "age")
	assert.Equal(t, x.String(), "")
}
