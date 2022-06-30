package gjson

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
