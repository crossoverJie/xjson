package gjson

import (
	"errors"
	"strconv"
	"strings"
)

type status int

const (
	StatusObjectKey status = iota
	StatusColon            //;
	StatusObjectValue
	StatusComma //,
	StatusBeginObject
	StatusEndObject
	StatusBeginArray
	StatusEndArray
	StatusArrayValue
)

func Parse(reader *TokenReader) (interface{}, error) {
	s := &Stack{}
	statuses := []status{StatusBeginObject, StatusBeginArray}
	for {
		tokenType := reader.Read()
		switch tokenType.T {
		case BeginObject:
			if notIncludeStatus(StatusBeginObject, statuses) {
				return nil, errors.New("invalid '{'")
			}
			root := make(map[string]interface{})
			stackValue := NewObjectValue(root)
			s.Push(stackValue)
			statuses = []status{StatusObjectKey, StatusBeginObject, StatusEndObject}
		case String:
			// todo crossoverJie 优雅去掉引号
			value := strings.Trim(tokenType.Value, "\"")
			if includeStatus(StatusObjectKey, statuses) {
				stackValue := NewObjectKey(value)
				s.Push(stackValue)
				statuses = []status{StatusColon}
				continue
			}
			if includeStatus(StatusObjectValue, statuses) {
				// ObjectValue 后跟的可能是 , 和 }
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = value
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			if includeStatus(StatusArrayValue, statuses) {
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, value)
				statuses = []status{StatusComma, StatusEndArray}
				continue
			}
			return nil, errors.New("invalid string '" + value + "'")

		case Number:
			// todo crossoverJie 优雅转为整形
			value := strings.Trim(tokenType.Value, "\"")
			if includeStatus(StatusObjectValue, statuses) {
				i, _ := strconv.Atoi(value)
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = i
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			if includeStatus(StatusArrayValue, statuses) {
				i, _ := strconv.Atoi(value)
				arrayValue := s.Peek().ArrayValuePoint()
				//arrayValue := s.Pop().ArrayValue()
				*arrayValue = append(*arrayValue, i)
				//s.Push(NewArray(arrayValue))
				statuses = []status{StatusComma, StatusEndArray}
				continue
			}
			return nil, errors.New("invalid int")
		case Float:
			// todo crossoverJie 优雅转为整形
			value := strings.Trim(tokenType.Value, "\"")
			if includeStatus(StatusObjectValue, statuses) {
				i, _ := strconv.ParseFloat(value, 64)

				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = i
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			if includeStatus(StatusArrayValue, statuses) {
				i, _ := strconv.ParseFloat(value, 64)

				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, i)
				statuses = []status{StatusComma, StatusEndArray}
				continue
			}
			return nil, errors.New("invalid int")
		case True:
			value := tokenType.Value
			if includeStatus(StatusObjectValue, statuses) {
				b, _ := strconv.ParseBool(value)
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = b
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			if includeStatus(StatusArrayValue, statuses) {
				b, _ := strconv.ParseBool(value)
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, b)
				statuses = []status{StatusComma, StatusEndArray}
				continue
			}
			return nil, errors.New("invalid bool true")
		case False:
			value := tokenType.Value
			if includeStatus(StatusObjectValue, statuses) {
				b, _ := strconv.ParseBool(value)
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = b
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			if includeStatus(StatusArrayValue, statuses) {
				b, _ := strconv.ParseBool(value)
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, b)
				statuses = []status{StatusComma, StatusEndArray}
				continue
			}
			return nil, errors.New("invalid bool false")
		case Null:
			if includeStatus(StatusObjectValue, statuses) {
				objectKey := s.Pop().ObjectKeyValue()
				rootMap := s.Peek().ObjectValue()
				rootMap[objectKey] = ""
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			if includeStatus(StatusArrayValue, statuses) {
				arrayValue := s.Peek().ArrayValuePoint()
				*arrayValue = append(*arrayValue, "")
				statuses = []status{StatusComma, StatusEndArray}
				continue
			}
			return nil, errors.New("invalid null")

		case SepComma:
			//,
			if notIncludeStatus(StatusComma, statuses) {
				return nil, errors.New("missing ':'")
			} else {
				// 逗号之前可能是 '}',下一个状态则是 StatusObjectKey
				if includeStatus(StatusEndObject, statuses) {
					statuses = []status{StatusObjectKey}
					continue
				}
				// 逗号之前可能是 ']',下一个状态则可能是
				if includeStatus(StatusEndArray, statuses) {
					statuses = []status{StatusArrayValue, StatusBeginArray, StatusBeginObject}
					continue
				}
			}

		case SepColon:
			//:
			if notIncludeStatus(StatusColon, statuses) {
				return nil, errors.New("missing ','")
			} else {
				statuses = []status{StatusObjectValue, StatusBeginObject, StatusBeginArray}
			}
		case BeginArray:
			if includeStatus(StatusBeginArray, statuses) {
				arr := make([]interface{}, 0)
				stackValue := NewArrayPoint(&arr)
				s.Push(stackValue)
				statuses = []status{StatusArrayValue, StatusBeginArray, StatusBeginObject, StatusEndArray}
				continue
			}
		case EndArray:
			if notIncludeStatus(StatusEndArray, statuses) {
				return nil, errors.New("invalid ']'")
			}
			if includeStatus(StatusEndArray, statuses) {
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
					statuses = []status{StatusComma, StatusEndObject}
					continue
				}

				// 如果栈顶还有一个 [
				if stackType == Array {
					arrayValuePoint := s.Peek().ArrayValuePoint()
					*arrayValuePoint = append(*arrayValuePoint, root)
					statuses = []status{StatusComma, StatusEndArray}
					continue
				}

			}
			return nil, errors.New("miss ']'")

		case EndObject:
			if notIncludeStatus(StatusEndObject, statuses) {
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
				statuses = []status{StatusComma, StatusEndObject}
				continue
			}
			// 如果栈顶还有一个 [
			if stackType == Array {
				arrayValuePoint := s.Peek().ArrayValuePoint()
				*arrayValuePoint = append(*arrayValuePoint, root)
				statuses = []status{StatusComma, StatusEndArray}
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

func includeStatus(status status, statuses []status) bool {
	for _, s := range statuses {
		if status == s {
			return true
		}
	}
	return false
}

func notIncludeStatus(status status, statuses []status) bool {
	for _, s := range statuses {
		if status == s {
			return false
		}
	}
	return true
}
