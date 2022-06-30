package gjson

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	JSONObject  Token = "JSONObject"
	Bool              = "Bool"
	ArrayObject       = "ArrayObject"
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
	statuses := []GrammarToken{Key}
	result := Result{
		Token: JSONObject,
		//Raw:    "",
		object: root,
	}
	for {
		read := reader.Read()
		switch read.T {
		case Key:
			if notIncludeGrammarToken(Key, statuses) {
				return buildEmptyResult()
			}
			var (
				m     map[string]interface{}
				a     *[]interface{}
				v     interface{}
				token Token
			)
			switch result.object.(type) {
			case map[string]interface{}:
				m = result.object.(map[string]interface{})
			case *[]interface{}:
				a = result.object.(*[]interface{})
			}

			if m != nil {
				v, ok = m[read.Value]
				token = typeOfToken(v)
				if !ok {
					return buildEmptyResult()
				}
			} else if a != nil {
				v = a
			} else {
				return buildEmptyResult()
			}

			result = Result{
				Token: token,
				//Raw:   v,
				object: v,
			}
			statuses = []GrammarToken{Dot, BeginArrayIndex}
			break
		case BeginArrayIndex:
			if notIncludeGrammarToken(BeginArrayIndex, statuses) {
				return buildEmptyResult()
			}
			statuses = []GrammarToken{Dot, ArrayIndex}
		case ArrayIndex:
			if notIncludeGrammarToken(ArrayIndex, statuses) {
				return buildEmptyResult()
			}
			a := result.object.(*[]interface{})
			index, _ := strconv.Atoi(read.Value)
			v := (*a)[index]
			token := typeOfToken(v)
			result = Result{
				Token: token,
				//Raw:   v,
				object: v,
			}
			statuses = []GrammarToken{EndArrayIndex}
		case EndArrayIndex:
			if notIncludeGrammarToken(EndArrayIndex, statuses) {
				return buildEmptyResult()
			}
			statuses = []GrammarToken{Dot}

		case Dot:
			if notIncludeGrammarToken(Dot, statuses) {
				return buildEmptyResult()
			}
			statuses = []GrammarToken{Key}

		case EOF:
			if includeGrammarToken(Key, statuses) {
				// syntax error
				return buildEmptyResult()
			}
			return result

		}
	}

}

func (r *Result) String() string {
	switch r.Token {
	case String:
		return fmt.Sprint(r.object)
	case Bool:
		return fmt.Sprint(r.object)
	case Null:
		return ""
	case Number:
		i, _ := strconv.Atoi(fmt.Sprint(r.object))
		return fmt.Sprintf("%d", i)
	case Float:
		i, _ := strconv.ParseFloat(fmt.Sprint(r.object), 64)
		return fmt.Sprintf("%f", i)
	case JSONObject:
		// todo crossoverJie toJSONString
		return fmt.Sprint(r.object)
	case ArrayObject:
		// todo crossoverJie toJSONString
		return fmt.Sprint(r.object)
	default:
		return ""
	}
}

func (r Result) Bool() bool {
	switch r.Token {
	case String:
		v, _ := strconv.ParseBool(strings.ToLower(fmt.Sprint(r.object)))
		return v
	case Bool:
		v, _ := strconv.ParseBool(strings.ToLower(fmt.Sprint(r.object)))
		return v
	case Null:
		return false
	case Number:
		v, _ := strconv.Atoi(fmt.Sprint(r.object))
		return v != 0
	default:
		return false
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

func (r Result) Float() float64 {
	switch r.Token {
	case String:
		v, _ := strconv.ParseFloat(fmt.Sprint(r.object), 64)
		return v
	case True:
		return 1
	case False:
		return 0
	case Null:
		return 0
	case Number:
		v, _ := strconv.Atoi(fmt.Sprint(r.object))
		return float64(v)
	case Float:
		v, _ := strconv.ParseFloat(fmt.Sprint(r.object), 64)
		return v
	default:
		return 0
	}
}

func (r Result) Object() interface{} {
	return r.object
}
func (r Result) Map() map[string]interface{} {
	return r.object.(map[string]interface{})
}

func (r Result) Exists() bool {
	return r.object != nil
}

func buildEmptyResult() Result {
	return Result{}
}

func includeGrammarToken(status GrammarToken, statuses []GrammarToken) bool {
	for _, s := range statuses {
		if status == s {
			return true
		}
	}
	return false
}

func notIncludeGrammarToken(status GrammarToken, statuses []GrammarToken) bool {
	for _, s := range statuses {
		if status == s {
			return false
		}
	}
	return true
}

func typeOfToken(v interface{}) (token Token) {
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
	case *[]interface{}:
		token = ArrayObject
	}
	return
}
