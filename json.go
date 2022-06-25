package gjson

func Decode(input string) (interface{}, error) {
	tokenize, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	reader := NewTokenReader(tokenize)
	return Parse(reader)
}
