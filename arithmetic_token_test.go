package xjson

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

	str = "(alice.age+bob.age) * tom.age"
	tokenize, err = ArithmeticTokenize(str)
	assert.Nil(t, err)
	for _, tokenType := range tokenize {
		fmt.Printf("%s  %s\n", tokenType.T, tokenType.Value)
	}
	input := `{"alice":{"age":10},"bob":{"age":20},"tom":{"age":20}}`
	arithmetic := GetWithArithmetic(input, str)
	assert.Equal(t, arithmetic.Int(), 600)
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

	grammar := `
if ((people[0].bob.age + people[1].alice.age)==20){
	return 10+10
}else{
	return 20
}
`
	ret := GetWithArithmetic(str, grammar)
	assert.Equal(t, ret.Int(), 20)

	grammar = `
if ((people[0].bob.age + people[1].alice.age)==20){
	return true
}else{
	return 20
}
`
	ret = GetWithArithmetic(str, grammar)
	assert.Equal(t, ret.Bool(), true)

	grammar = `
if ((people[0].bob.age + people[1].alice.age)!=20){
	return true
}else{
	return false
}
`
	ret = GetWithArithmetic(str, grammar)
	assert.Equal(t, ret.Bool(), false)
}

func TestGetWithArithmetic2(t *testing.T) {
	str := `{"name":[}`
	x := GetWithArithmetic(str, "x")
	assert.Equal(t, x.String(), "")
}

func TestGetWithArithmetic4(t *testing.T) {
	str := `{"name":"bob", "iii":10,"magic":10.1, "score":{"math":[1,2]}}`
	result := GetWithArithmetic(str, "(iii+iii)*iii+magic")
	assert.Equal(t, result.Float(), 210.1)

	str = `{"name":"bob", "iii":10,"ee":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+ee")
	assert.Equal(t, result.Float(), 210.1)

	str = `{"name":"bob", "iii":10,"ell":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+ell")
	assert.Equal(t, result.Float(), 210.1)

	str = `{"name":"bob", "iii":10,"elss":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+elss")
	assert.Equal(t, result.Float(), 210.1)

	str = `{"name":"bob", "iii":10,"rr":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+rr")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"ree":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+ree")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"rett":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+rett")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"retuu":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+retuu")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"returr":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+returr")
	assert.Equal(t, result.Float(), 210.1)

	str = `{"name":"bob", "iii":10,"tt":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+tt")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"trr":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+trr")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"truu":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+truu")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"ff":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+ff")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"faa":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+faa")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"fall":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+fall")
	assert.Equal(t, result.Float(), 210.1)
	str = `{"name":"bob", "iii":10,"falss":10.1, "score":{"math":[1,2]}}`
	result = GetWithArithmetic(str, "(iii+iii)*iii+falss")
	assert.Equal(t, result.Float(), 210.1)

}

func TestGetWithArithmeticErr(t *testing.T) {
	str := `[1,2]`
	x := GetWithArithmetic(str, "x")
	assert.Equal(t, x.String(), "")

	str = `{"age":true}`
	x = GetWithArithmetic(str, "age")
	assert.Equal(t, x.String(), "")
}
