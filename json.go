package gjson

import (
	"fmt"
	"strconv"
)

const (
	JSONObject Token = "JSONObject"
	Bool             = "Bool"
)

type Result struct {
	Token  Token
	Raw    string
	object interface{}
}

func Decode(input string) (interface{}, error) {
	tokenize, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	reader := NewTokenReader(tokenize)
	return Parse(reader)
}

func Get(json, grammar string) Result {
	decode, err := Decode(json)
	if err != nil {
		return buildEmptyResult()
	}

	root, ok := decode.(map[string]interface{})
	if !ok {
		return buildEmptyResult()
	}

	tokenize, err := GrammarTokenize(grammar)
	if err != nil {
		return buildEmptyResult()
	}
	reader := NewGrammarTokenReader(tokenize)
	status := Key
	result := Result{
		Token: JSONObject,
		//Raw:    "",
		object: root,
	}
	for {
		read := reader.Read()
		switch read.T {
		case Key:
			if status != Key {
				return buildEmptyResult()
			}
			m, ok := result.object.(map[string]interface{})
			if !ok {
				return buildEmptyResult()
			}
			v, ok := m[read.Value]
			var token Token
			switch v.(type) {
			case string:
				token = String
			case int:
				token = Number
			case float64:
				token = Float
			case bool:
				token = Bool
			case map[string]interface{}:
				token = JSONObject
			}
			if !ok {
				return buildEmptyResult()
			}
			result = Result{
				Token: token,
				//Raw:   v,
				object: v,
			}
			status = Dot
			break
		case Dot:
			if status != Dot {
				return buildEmptyResult()
			}
			status = Key

		case EOF:
			return result

		}
	}

}

func (r *Result) String() string {
	switch r.Token {
	case String:
		return fmt.Sprint(r.object)
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return ""
	case Number:
		i, _ := strconv.Atoi(fmt.Sprint(r.object))
		return fmt.Sprintf("%d", i)
	case Float:
		i, _ := strconv.ParseFloat(fmt.Sprint(r.object), 64)
		return fmt.Sprintf("%f", i)
	case JSONObject:
		return fmt.Sprint(r.object)
	default:
		return ""
	}
}

func (r Result) Int() int {
	switch r.Token {
	case String:
		v, _ := strconv.Atoi(fmt.Sprint(r.object))
		return v
	case True:
		return 1
	case False:
		return 0
	case Null:
		return 0
	case Number:
		v, _ := strconv.Atoi(fmt.Sprint(r.object))
		return v
	case Float:
		i, _ := strconv.ParseFloat(fmt.Sprint(r.object), 64)
		return int(i)
	default:
		return 0
	}
}

func buildEmptyResult() Result {
	return Result{}
}
