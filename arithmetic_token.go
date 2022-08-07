package xjson

type ArithmeticToken string

const (
	ArithmeticInit ArithmeticToken = "Init"
	Plus                           = "Plus"
	Minus                          = "Minus"
	Star                           = "Star"
	Slash                          = "Slash"
	LeftParen                      = "LeftParen"
	RightParen                     = "RightParen"
	Identifier                     = "Identifier"
	If                             = "If"
	If1                            = "If1"
	Else                           = "Else"
	Else1                          = "Else1"
	Else2                          = "Else2"
	Else3                          = "Else3"
	Return                         = "Return"
	Return1                        = "Return1"
	Return2                        = "Return2"
	Return3                        = "Return3"
	Return4                        = "Return4"
	Return5                        = "Return5"
	Equals                         = "Equals"
	Equals1                        = "Equals1"
	Semi                           = ";"
	LeftBra                        = "LeftBra"
	RightBra                       = "RightBra"
	ArithmeticEOF                  = "EOF"
)

type ArithmeticTokenType struct {
	T     ArithmeticToken
	Value string
}

func ArithmeticTokenize(str string) ([]*ArithmeticTokenType, error) {
	bytes := []byte(str)
	var result []*ArithmeticTokenType
	var values []byte
	status := ArithmeticInit
	for _, b := range bytes {
		switch status {
		case ArithmeticInit:
			status, values = InitArithmeticStatus(b, values)
			break
		case Number:
			if isDigit(b) {
				values = append(values, b)
			} else {
				t := &ArithmeticTokenType{
					T:     Number,
					Value: string(values),
				}
				result = append(result, t)
				values = nil
				status, values = InitArithmeticStatus(b, values)
				break
			}
		case If1:
			if b == 'f' {
				values = append(values, b)
				status = If
			} else {
				values = append(values, b)
				status = Identifier
			}
		case If:
			t := &ArithmeticTokenType{
				T:     If,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Else1:
			if b == 'l' {
				values = append(values, b)
				status = Else2
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Else2:
			if b == 's' {
				values = append(values, b)
				status = Else3
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Else3:
			if b == 'e' {
				values = append(values, b)
				status = Else
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Else:
			t := &ArithmeticTokenType{
				T:     Else,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Return1:
			if b == 'e' {
				values = append(values, b)
				status = Return2
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Return2:
			if b == 't' {
				values = append(values, b)
				status = Return3
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Return3:
			if b == 'u' {
				values = append(values, b)
				status = Return4
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Return4:
			if b == 'r' {
				values = append(values, b)
				status = Return5
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Return5:
			if b == 'n' {
				values = append(values, b)
				status = Return
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Return:
			t := &ArithmeticTokenType{
				T:     Return,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Equals1:
			if b == '=' {
				values = append(values, b)
				status = Equals
			} else {
				values = append(values, b)
				status = Identifier
			}
		case Equals:
			t := &ArithmeticTokenType{
				T:     Equals,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Semi:
			t := &ArithmeticTokenType{
				T:     Semi,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case LeftBra:
			t := &ArithmeticTokenType{
				T:     LeftBra,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case RightBra:
			t := &ArithmeticTokenType{
				T:     RightBra,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case True1:
			if b == 'r' {
				values = append(values, b)
				status = True2
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case True2:
			if b == 'u' {
				values = append(values, b)
				status = True3
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case True3:
			if b == 'e' {
				values = append(values, b)
				status = True
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case True:
			t := &ArithmeticTokenType{
				T:     True,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
			break
		case False1:
			if b == 'a' {
				values = append(values, b)
				status = False2
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case False2:
			if b == 'l' {
				values = append(values, b)
				status = False3
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case False3:
			if b == 's' {
				values = append(values, b)
				status = False4
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case False4:
			if b == 'e' {
				values = append(values, b)
				status = False
				break
			} else {
				values = append(values, b)
				status = Identifier
			}
		case False:
			t := &ArithmeticTokenType{
				T:     False,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
			break
		case Identifier:
			if IsArithmetic(b) {
				values = append(values, b)
			} else {
				t := &ArithmeticTokenType{
					T:     Identifier,
					Value: string(values),
				}
				result = append(result, t)
				values = nil
				status, values = InitArithmeticStatus(b, values)
			}
		case Plus:
			t := &ArithmeticTokenType{
				T:     Plus,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Minus:
			t := &ArithmeticTokenType{
				T:     Minus,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Star:
			t := &ArithmeticTokenType{
				T:     Star,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case Slash:
			t := &ArithmeticTokenType{
				T:     Slash,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case LeftParen:
			t := &ArithmeticTokenType{
				T:     LeftParen,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		case RightParen:
			t := &ArithmeticTokenType{
				T:     RightParen,
				Value: string(values),
			}
			result = append(result, t)
			values = nil
			status, values = InitArithmeticStatus(b, values)
		}
	}

	if len(values) > 0 {
		t := &ArithmeticTokenType{
			T:     status,
			Value: string(values),
		}
		result = append(result, t)
	}
	return result, nil
}

func InitArithmeticStatus(b byte, values []byte) (ArithmeticToken, []byte) {
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
	if b == 'i' {
		values = append(values, b)
		return If1, values
	}
	if b == 'e' {
		values = append(values, b)
		return Else1, values
	}

	if b == 'r' {
		values = append(values, b)
		return Return1, values
	}
	if b == ';' {
		values = append(values, b)
		return Semi, values
	}
	if b == '=' {
		values = append(values, b)
		return Equals1, values
	}
	if b == '{' {
		values = append(values, b)
		return LeftBra, values
	}
	if b == '}' {
		values = append(values, b)
		return RightBra, values
	}

	if IsArithmetic(b) {
		values = append(values, b)
		return Identifier, values
	}
	if b == '+' {
		values = append(values, b)
		return Plus, values
	}
	if b == '-' {
		values = append(values, b)
		return Minus, values
	}
	if b == '*' {
		values = append(values, b)
		return Star, values
	}
	if b == '/' {
		values = append(values, b)
		return Slash, values
	}
	if b == '(' {
		values = append(values, b)
		return LeftParen, values
	}
	if b == ')' {
		values = append(values, b)
		return RightParen, values
	}

	return ArithmeticInit, values
}

func IsArithmetic(b byte) bool {
	return IsLetter(b) || isDigit(b) || b == '.'
}

func IsLetter(b byte) bool {
	return b >= 65 && b <= 122
}

type ArithmeticTokenReader struct {
	tokens []*ArithmeticTokenType
	pos    uint64
}

func NewArithmeticTokenReader(tokens []*ArithmeticTokenType) *ArithmeticTokenReader {
	return &ArithmeticTokenReader{tokens: tokens, pos: 0}
}

func (t *ArithmeticTokenReader) Read() *ArithmeticTokenType {
	if int(t.pos) >= len(t.tokens) {
		return &ArithmeticTokenType{
			T: ArithmeticEOF,
		}
	}
	tokenType := t.tokens[t.pos]
	t.pos += 1
	return tokenType
}
