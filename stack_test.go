package gjson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack_Peek(t *testing.T) {
	s := &Stack{}
	s.Push(NewObjectKey("name"))
	s.Push(NewObjectKey("age"))
	s.Push(NewObjectKey("sex"))

	peek := s.Peek()
	assert.Equal(t, peek.ObjectKeyValue(), "sex")

	pop := s.Pop()
	assert.Equal(t, pop.ObjectKeyValue(), "sex")
	pop = s.Pop()
	assert.Equal(t, pop.ObjectKeyValue(), "age")
	pop = s.Pop()
	assert.Equal(t, pop.ObjectKeyValue(), "name")

	empty := s.IsEmpty()
	assert.Equal(t, empty, true)
}
