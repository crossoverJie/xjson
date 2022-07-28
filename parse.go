package xjson

import (
	"errors"
	"strconv"
)

type status int

const (
	StatusObjectKey   status = 0x0002
	StatusColon       status = 0x0004 //;
	StatusObjectValue status = 0x0008
	StatusComma       status = 0x0010 //,
	StatusBeginObject status = 0x0020
	StatusEndObject   status = 0x0040
	StatusBeginArray  status = 0x0080
	StatusEndArray    status = 0x0100
	StatusArrayValue  status = 0x0200
)

func Parse(reader *TokenReader) (interface{}, error) {
	s := &Stack{}
	status := StatusBeginObject | StatusBeginArray
	for {
		tokenType := reader.Read()
		switch tokenType.T {
		case BeginObject:
			if !includeTokenStatus(StatusBeginObject, status) {
				return nil, errors.New("invalid '{'")
			}
			root := make(map[string]interface{})
			stackValue := NewObjectValue(root)
			s.Push(stackValue)
			status = StatusObjectKey | StatusBeginObject | StatusEndObject
		case String:
			if includeTokenStatus(StatusObjectKey, status) {
				stackValue := NewObjectKey(tokenType.Value)
				s.Push(stackValue)
				status = StatusColon
				continue
			}
			if includeTokenStatus(StatusObjectValue, status) {
				// ObjectValue 后跟的可能是 , 和 }
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = tokenType.Value
				status = StatusComma | StatusEndObject
				continue
			}
			if includeTokenStatus(StatusArrayValue, status) {
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, tokenType.Value)
				status = StatusComma | StatusEndArray
				continue
			}
			return nil, errors.New("invalid string '" + tokenType.Value + "'")

		case Number:
			// todo crossoverJie 优雅转为整形
			if includeTokenStatus(StatusObjectValue, status) {
				i, _ := strconv.Atoi(tokenType.Value)
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = i
				status = StatusComma | StatusEndObject
				continue
			}
			if includeTokenStatus(StatusArrayValue, status) {
				i, _ := strconv.Atoi(tokenType.Value)
				arrayValue := s.Peek().ArrayValuePoint()
				//arrayValue := s.Pop().ArrayValue()
				*arrayValue = append(*arrayValue, i)
				//s.Push(NewArray(arrayValue))
				status = StatusComma | StatusEndArray
				continue
			}
			return nil, errors.New("invalid int")
		case Float:
			// todo crossoverJie 优雅转为整形
			if includeTokenStatus(StatusObjectValue, status) {
				i, _ := strconv.ParseFloat(tokenType.Value, 64)

				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = i
				status = StatusComma | StatusEndObject
				continue
			}
			if includeTokenStatus(StatusArrayValue, status) {
				i, _ := strconv.ParseFloat(tokenType.Value, 64)

				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, i)
				status = StatusComma | StatusEndArray
				continue
			}
			return nil, errors.New("invalid int")
		case True:
			value := tokenType.Value
			if includeTokenStatus(StatusObjectValue, status) {
				b, _ := strconv.ParseBool(value)
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = b
				status = StatusComma | StatusEndObject
				continue
			}
			if includeTokenStatus(StatusArrayValue, status) {
				b, _ := strconv.ParseBool(value)
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, b)
				status = StatusComma | StatusEndArray
				continue
			}
			return nil, errors.New("invalid bool true")
		case False:
			value := tokenType.Value
			if includeTokenStatus(StatusObjectValue, status) {
				b, _ := strconv.ParseBool(value)
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = b
				status = StatusComma | StatusEndObject
				continue
			}
			if includeTokenStatus(StatusArrayValue, status) {
				b, _ := strconv.ParseBool(value)
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, b)
				status = StatusComma | StatusEndArray
				continue
			}
			return nil, errors.New("invalid bool false")
		case Null:
			if includeTokenStatus(StatusObjectValue, status) {
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = ""
				status = StatusComma | StatusEndObject
				continue
			}
			if includeTokenStatus(StatusArrayValue, status) {
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, "")
				status = StatusComma | StatusEndArray
				continue
			}
			return nil, errors.New("invalid null")

		case SepComma:
			//,
			if !includeTokenStatus(StatusComma, status) {
				return nil, errors.New("missing ':'")
			} else {
				// 逗号之前可能是 '}',下一个状态则是 StatusObjectKey
				if includeTokenStatus(StatusEndObject, status) {
					status = StatusObjectKey
					continue
				}
				// 逗号之前可能是 ']',下一个状态则可能是
				if includeTokenStatus(StatusEndArray, status) {
					status = StatusArrayValue | StatusBeginArray | StatusBeginObject
					continue
				}
			}

		case SepColon:
			//:
			if !includeTokenStatus(StatusColon, status) {
				return nil, errors.New("missing ','")
			} else {
				status = StatusObjectValue | StatusBeginObject | StatusBeginArray
			}
		case BeginArray:
			if includeTokenStatus(StatusBeginArray, status) {
				arr := make([]interface{}, 0)
				stackValue := NewArrayPoint(&arr)
				s.Push(stackValue)
				status = StatusArrayValue | StatusBeginArray | StatusBeginObject | StatusEndArray
				continue
			}
		case EndArray:
			if !includeTokenStatus(StatusEndArray, status) {
				return nil, errors.New("invalid ']'")
			}
			root := s.Pop().ArrayValuePoint()
			if s.IsEmpty() {
				array := NewArrayPoint(root)
				s.Push(array)
				continue
			}

			stackType := s.Peek().StackType()
			// 如果栈顶还有一个 key
			if stackType == ObjectKey {
				keyValue := s.Pop().ObjectKeyValue()
				objectValue := s.Peek().ObjectValue()
				objectValue[keyValue] = root
				status = StatusComma | StatusEndObject
				continue
			}

			// 如果栈顶还有一个 [
			if stackType == Array {
				arrayValuePoint := s.Peek().ArrayValuePoint()
				*arrayValuePoint = append(*arrayValuePoint, root)
				status = StatusComma | StatusEndArray
				continue
			}

			return nil, errors.New("miss ']'")

		case EndObject:
			if !includeTokenStatus(StatusEndObject, status) {
				return nil, errors.New("invalid '}'")
			}
			root := s.Pop().ObjectValue()
			if s.IsEmpty() {
				value := NewObjectValue(root)
				s.Push(value)
				// 此时栈已经读完，表名所有 token 解析完毕
				continue
			}

			stackType := s.Peek().StackType()
			// 如果栈顶还有一个 key
			// {"name":"cj", "age":{"a":"a","b":"b","c":true}, "e":"e"}
			if stackType == ObjectKey {
				keyValue := s.Pop().ObjectKeyValue()
				objectValue := s.Peek().ObjectValue()
				objectValue[keyValue] = root
				// 下一个可能是 ',' 和 '}'
				status = StatusComma | StatusEndObject
				continue
			}
			// 如果栈顶还有一个 [
			if stackType == Array {
				arrayValuePoint := s.Peek().ArrayValuePoint()
				*arrayValuePoint = append(*arrayValuePoint, root)
				status = StatusComma | StatusEndArray
				continue
			}

		case EndJson:
			// token 读不到数据了，EOF
			root := s.Pop().Raw()
			if s.IsEmpty() {
				return root, nil
			} else {
				return nil, errors.New("invalid EOF")
			}

		}
	}
}

func includeTokenStatus(current, target status) bool {
	return (current & target) > 0
}
