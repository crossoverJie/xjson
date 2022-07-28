package xjson

import "errors"

type GrammarToken string

const (
	GrammarInit     GrammarToken = "Init"
	Key             GrammarToken = "Key"
	Dot             GrammarToken = "Dot"
	GrammarEscape   GrammarToken = "GrammarEscape"
	BeginArrayIndex GrammarToken = "BeginArrayIndex"
	EndArrayIndex   GrammarToken = "EndArrayIndex"
	ArrayIndex      GrammarToken = "ArrayIndex"
	EOF             GrammarToken = "EOF"
)

type GrammarStatus int

const (
	KeyStatus             GrammarStatus = 0x0004
	DotStatus             GrammarStatus = 0x0008
	BeginArrayIndexStatus GrammarStatus = 0x0010
	EndArrayIndexStatus   GrammarStatus = 0x0020
	ArrayIndexStatus      GrammarStatus = 0x0040
)

type GrammarTokenType struct {
	T     GrammarToken
	Value string
}

func GrammarTokenize(str string) ([]*GrammarTokenType, error) {
	//bytes := []byte(str)
	var result []*GrammarTokenType
	var values []byte
	status := GrammarInit
	for i := 0; i < len(str); i++ {
		b := str[i]

		switch status {
		case GrammarInit:
			status, values = InitGrammarStatus(b, values)
			break
		case Key:
			if b == '\\' {
				status = GrammarEscape
				break
			} else if isGrammarLetter(b) || isDigit(b) {
				values = append(values, b)
			} else if b == '[' {
				t := &GrammarTokenType{
					T:     Key,
					Value: string(values),
				}
				result = append(result, t)
				values = nil

				values = append(values, b)
				t = &GrammarTokenType{
					T:     BeginArrayIndex,
					Value: string(values),
				}
				result = append(result, t)
				status = ArrayIndex
				values = nil
			} else {
				t := &GrammarTokenType{
					T:     Key,
					Value: string(values),
				}
				result = append(result, t)
				values = nil
				status, values = InitGrammarStatus(b, values)
				break
			}
		case GrammarEscape:
			if i < len(str) {
				for ; i < len(str); i++ {
					if str[i] == '\\' {
						i++ //skip escape
						if i < len(str) {
							values = append(values, str[i])
						}
					} else if str[i] == '[' {
						t := &GrammarTokenType{
							T:     Key,
							Value: string(values),
						}
						result = append(result, t)
						values = nil

						values = append(values, str[i])
						t = &GrammarTokenType{
							T:     BeginArrayIndex,
							Value: string(values),
						}
						result = append(result, t)
						status = ArrayIndex
						values = nil
						break
					} else {
						values = append(values, str[i])
					}
				}
			}
			if len(values) == 0 {
				// compatible a\.[0]
				break
			}
			t := &GrammarTokenType{
				T:     Key,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status = GrammarInit
			break

		case Dot:
			t := &GrammarTokenType{
				T:     Dot,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitGrammarStatus(b, values)
			break

		case ArrayIndex:
			if isDigit(b) {
				values = append(values, b)
			} else if b == ']' {
				t := &GrammarTokenType{
					T:     ArrayIndex,
					Value: string(values),
				}
				result = append(result, t)
				values = make([]byte, 0)

				values = append(values, b)
				status = EndArrayIndex
			} else {
				return nil, errors.New("invalid array index ]")
			}

		case EndArrayIndex:
			t := &GrammarTokenType{
				T:     EndArrayIndex,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitGrammarStatus(b, values)
			break
		}
	}

	if len(values) > 0 {
		t := &GrammarTokenType{
			T:     status,
			Value: string(values),
		}
		result = append(result, t)
	}

	return result, nil
}

func InitGrammarStatus(b byte, values []byte) (GrammarToken, []byte) {
	if b == '\\' {
		return GrammarEscape, values
	}
	if isGrammarLetter(b) || isDigit(b) {
		values = append(values, b)
		return Key, values
	}
	if b == '.' {
		values = append(values, b)
		return Dot, values
	}

	return GrammarInit, values
}

func isGrammarLetter(b byte) bool {
	if b == '[' || b == ']' {
		return false
	}
	return b >= 65 && b <= 122
}

type GrammarTokenReader struct {
	tokens []*GrammarTokenType
	pos    uint64
}

func NewGrammarTokenReader(tokens []*GrammarTokenType) *GrammarTokenReader {
	return &GrammarTokenReader{tokens: tokens, pos: 0}
}

func (t *GrammarTokenReader) Read() *GrammarTokenType {
	if int(t.pos) >= len(t.tokens) {
		return &GrammarTokenType{
			T: EOF,
		}
	}
	tokenType := t.tokens[t.pos]
	t.pos += 1
	return tokenType
}

func (t *GrammarTokenReader) HasNext() bool {
	return t.pos < uint64(len(t.tokens))
}
