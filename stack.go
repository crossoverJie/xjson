package gjson

import "fmt"

type Stack []*StackValue

// IsEmpty  check if Stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the Stack
func (s *Stack) Push(val *StackValue) {
	*s = append(*s, val) // Simply append the new value to the end of the Stack
}

// Pop Remove and return top element of Stack. Return false if Stack is empty.
func (s *Stack) Pop() *StackValue {
	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.
	*s = (*s)[:index]      // Remove it from the Stack by slicing it off.
	return element

}

// Peek Peek value, not remove element.
func (s *Stack) Peek() *StackValue {
	index := len(*s) - 1
	element := (*s)[index] // Index into the slice and obtain the element.
	return element
}

type StackType string

const (
	Object    StackType = "object"
	ObjectKey           = "objectKey"
	Array               = "array"
)

type StackValue struct {
	stackType StackType
	object    interface{}
}

func NewObjectValue(object interface{}) *StackValue {
	return &StackValue{stackType: Object, object: object}
}

func NewObjectKey(key string) *StackValue {
	return &StackValue{stackType: ObjectKey, object: key}
}

//func NewArray(array []interface{}) *StackValue {
//	return &StackValue{stackType: Array, object: array}
//}
func NewArrayPoint(array *[]interface{}) *StackValue {
	return &StackValue{stackType: Array, object: array}
}

func (s *StackValue) ObjectKeyValue() string {
	return fmt.Sprint(s.object)
}
func (s *StackValue) ObjectValue() map[string]interface{} {
	return s.object.(map[string]interface{})
}
func (s *StackValue) Raw() interface{} {
	return s.object
}

//func (s *StackValue) ArrayValue() []interface{} {
//	return s.object.([]interface{})
//}
func (s *StackValue) ArrayValuePoint() *[]interface{} {
	arr := s.object.(*[]interface{})
	return arr
}

func (s *StackValue) StackType() StackType {
	return s.stackType
}
