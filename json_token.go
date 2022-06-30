package gjson

import "errors"

type Token string

const (
	Init        Token = "Init"
	BeginObject       = "BeginObject"
	EndObject         = "EndObject"
	BeginArray        = "BeginArray"
	EndArray          = "EndArray"
	Null              = "Null"
	Null1             = "Null1"
	Null2             = "Null2"
	Null3             = "Null3"
	Number            = "Number"
	Float             = "Float"
	BeginString       = "BeginString"
	EndString         = "EndString"
	String            = "String"
	True              = "True"
	True1             = "True1"
	True2             = "True2"
	True3             = "True3"
	False             = "False"
	False1            = "False1"
	False2            = "False2"
	False3            = "False3"
	False4            = "False4"
	// SepColon :
	SepColon = "SepColon"
	// SepComma ,
	SepComma = "SepComma"
	EndJson  = "EndJson"
)

type TokenType struct {
	T     Token
	Value string
}

func Tokenize(str string) ([]*TokenType, error) {
	bytes := []byte(str)
	var result []*TokenType
	var values []byte
	status := Init
	for _, b := range bytes {
		switch status {
		case Init:
			status, values = InitStatus(b, values)
			break
		case BeginObject:
			t := &TokenType{
				T:     BeginObject,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
		case EndObject:
			t := &TokenType{
				T:     EndObject,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
		case BeginString:
			if b == '"' {
				values = append(values, b)
				status = EndString
			} else {
				values = append(values, b)
			}
		case EndString:
			t := &TokenType{
				T:     String,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case Number:
			if b == '.' {
				values = append(values, b)
				status = Float
				break
			}
			if isDigit(b) {
				values = append(values, b)
			} else {
				t := &TokenType{
					T:     Number,
					Value: string(values),
				}
				result = append(result, t)
				values = nil
				status, values = InitStatus(b, values)
				break
			}
		case Float:
			if b == '.' {
				return nil, errors.New("invalid float")
			}
			if isDigit(b) {
				values = append(values, b)
			} else {
				t := &TokenType{
					T:     Float,
					Value: string(values),
				}
				result = append(result, t)
				values = nil
				status, values = InitStatus(b, values)
				break
			}
		case SepColon:
			t := &TokenType{
				T:     SepColon,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case SepComma:
			t := &TokenType{
				T:     SepComma,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case BeginArray:
			t := &TokenType{
				T:     BeginArray,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case EndArray:
			t := &TokenType{
				T:     EndArray,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case True1:
			if b == 'r' {
				values = append(values, b)
				status = True2
				break
			} else {
				return nil, errors.New("invalid bool true")
			}
		case True2:
			if b == 'u' {
				values = append(values, b)
				status = True3
				break
			} else {
				return nil, errors.New("invalid bool true")
			}
		case True3:
			if b == 'e' {
				values = append(values, b)
				status = True
				break
			} else {
				return nil, errors.New("invalid bool true")
			}
		case True:
			t := &TokenType{
				T:     True,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case Null1:
			if b == 'u' {
				values = append(values, b)
				status = Null2
				break
			} else {
				return nil, errors.New("invalid null")
			}
		case Null2:
			if b == 'l' {
				values = append(values, b)
				status = Null3
				break
			} else {
				return nil, errors.New("invalid null")
			}
		case Null3:
			if b == 'l' {
				values = append(values, b)
				status = Null
				break
			} else {
				return nil, errors.New("invalid null")
			}
		case Null:
			t := &TokenType{
				T:     Null,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		case False1:
			if b == 'a' {
				values = append(values, b)
				status = False2
				break
			} else {
				return nil, errors.New("invalid bool false")
			}
		case False2:
			if b == 'l' {
				values = append(values, b)
				status = False3
				break
			} else {
				return nil, errors.New("invalid bool false")
			}
		case False3:
			if b == 's' {
				values = append(values, b)
				status = False4
				break
			} else {
				return nil, errors.New("invalid bool false")
			}
		case False4:
			if b == 'e' {
				values = append(values, b)
				status = False
				break
			} else {
				return nil, errors.New("invalid bool false")
			}
		case False:
			t := &TokenType{
				T:     False,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitStatus(b, values)
			break
		}
	}

	// 解析最后一个
	if len(values) > 0 {
		t := &TokenType{
			T:     status,
			Value: string(values),
		}
		result = append(result, t)
	}

	return result, nil
}

func InitStatus(b byte, values []byte) (Token, []byte) {

	if b == '{' {
		values = append(values, b)
		return BeginObject, values
	}
	if b == '}' {
		values = append(values, b)
		return EndObject, values
	}
	if b == '"' {
		values = append(values, b)
		return BeginString, values
	}
	if b == ':' {
		values = append(values, b)
		return SepColon, values
	}
	if b == ',' {
		values = append(values, b)
		return SepComma, values
	}
	if b == '[' {
		values = append(values, b)
		return BeginArray, values
	}
	if b == ']' {
		values = append(values, b)
		return EndArray, values
	}
	if isDigit(b) {
		values = append(values, b)
		return Number, values
	}
	if b == 't' {
		values = append(values, b)
		return True1, values
	}
	if b == 'f' {
		values = append(values, b)
		return False1, values
	}
	if b == 'n' {
		values = append(values, b)
		return Null1, values
	}

	return Init, values
}

func isDigit(b byte) bool {
	return b >= 48 && b <= 57
}

type TokenReader struct {
	tokens []*TokenType
	pos    uint64
}

func NewTokenReader(tokens []*TokenType) *TokenReader {
	return &TokenReader{tokens: tokens, pos: 0}
}

func (t *TokenReader) Read() *TokenType {
	if int(t.pos) >= len(t.tokens) {
		return &TokenType{
			T: EndJson,
		}
	}
	tokenType := t.tokens[t.pos]
	t.pos += 1
	return tokenType
}
