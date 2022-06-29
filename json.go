package gjson

func Decode(input string) (interface{}, error) {
	tokenize, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	reader := NewTokenReader(tokenize)
	return Parse(reader)
}

func GetString(json, grammar string) string {
	_, err := Decode(json)
	if err != nil {
		return ""
	}
	return ""
}
