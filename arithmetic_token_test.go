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
}
