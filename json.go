package xjson

import (
	"errors"
	"fmt"
	"github.com/crossoverJie/gscript/syntax"
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
	object interface{}
}

func Decode(input string) (interface{}, error) {
	tokenize, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(tokenize) == 0 {
		return nil, errors.New("input is empty")
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

	return getWithRoot(root, grammar)

}

func getWithRoot(root map[string]interface{}, grammar string) Result {
	tokenize, err := GrammarTokenize(grammar)
	if err != nil {
		return buildEmptyResult()
	}
	reader := NewGrammarTokenReader(tokenize)
	statuses := []GrammarToken{Key}
	result := Result{
		Token:  JSONObject,
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
				v     interface{}
				token Token
				ok    bool
			)
			switch result.object.(type) {
			case map[string]interface{}:
				m = result.object.(map[string]interface{})
			}

			if m != nil && len(m) > 0 {
				v, ok = m[read.Value]
				token = typeOfToken(v)
				if !ok {
					return buildEmptyResult()
				}
			} else {
				return buildEmptyResult()
			}

			result = Result{
				Token:  token,
				object: v,
			}
			statuses = []GrammarToken{Dot, BeginArrayIndex}
			break
		case BeginArrayIndex:
			// unless code
			//if notIncludeGrammarToken(BeginArrayIndex, statuses) {
			//	return buildEmptyResult()
			//}
			statuses = []GrammarToken{Dot, ArrayIndex}
		case ArrayIndex:
			// unless code
			//if notIncludeGrammarToken(ArrayIndex, statuses) {
			//	return buildEmptyResult()
			//}
			a := result.object.(*[]interface{})
			index, _ := strconv.Atoi(read.Value)
			v := (*a)[index]
			token := typeOfToken(v)
			result = Result{
				Token:  token,
				object: v,
			}
			statuses = []GrammarToken{EndArrayIndex}
		case EndArrayIndex:
			// unless code
			//if notIncludeGrammarToken(EndArrayIndex, statuses) {
			//	return buildEmptyResult()
			//}
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

func GetWithArithmetic(json, arithmetic string) Result {
	tokenize, err := ArithmeticTokenize(arithmetic)
	if err != nil {
		return buildEmptyResult()
	}

	decode, err := Decode(json)
	if err != nil {
		return buildEmptyResult()
	}
	root, ok := decode.(map[string]interface{})
	if !ok {
		return buildEmptyResult()
	}

	reader := NewArithmeticTokenReader(tokenize)
	var builder strings.Builder
	for {
		read := reader.Read()
		switch read.T {
		case Identifier:
			result := getWithRoot(root, read.Value)
			switch result.Token {
			case Number:
				builder.WriteString(fmt.Sprintf("%d", result.Int()))
			case Float:
				builder.WriteString(fmt.Sprintf("%f", result.Float()))
			default:
				return buildEmptyResult()
			}
		case ArithmeticEOF:
			builder.WriteString("\n")
			r := syntax.ArithmeticOperators(builder.String())
			switch r.(type) {
			case int:
				return Result{
					Token:  Number,
					object: r,
				}
			case float64:
				return Result{
					Token:  Float,
					object: r,
				}
			}
		default:
			builder.WriteString(read.Value)
		}
	}

	//return Result{}
}

// String return result of string
func (r Result) String() string {
	switch r.Token {
	case String:
		return fmt.Sprint(r.object)
	case Bool:
		return fmt.Sprint(r.object)
	//case Null:
	//	return ""
	case Number:
		i, _ := strconv.Atoi(fmt.Sprint(r.object))
		return fmt.Sprintf("%d", i)
	case Float:
		i, _ := strconv.ParseFloat(fmt.Sprint(r.object), 64)
		return fmt.Sprintf("%f", i)
	case JSONObject:
		return object2JSONString(r.object)
	case ArrayObject:
		return object2JSONString(r.Array())
	default:
		return ""
	}
}

func object2JSONString(object interface{}) string {
	var builder strings.Builder

	switch data := object.(type) {
	case map[string]interface{}:
		builder.WriteString("{")
		m := data
		count := 0
		for s, v := range m {
			builder.WriteString(fmt.Sprintf("\"%s\"", s))
			builder.WriteString(":")

			switch vv := v.(type) {
			case map[string]interface{}:
				value := object2JSONString(vv)
				builder.WriteString(value)
			case *[]interface{}:
				slice := covertSlice(vv)
				value := object2JSONString(slice)
				builder.WriteString(value)
			default:
				builder.WriteString(interface2String(v))
			}

			count++
			if len(m) != count {
				builder.WriteString(",")
			}

		}
		builder.WriteString("}")
	case []interface{}:
		builder.WriteString("[")
		count := 0
		for _, v := range data {
			switch vv := v.(type) {
			case map[string]interface{}:
				value := object2JSONString(vv)
				builder.WriteString(value)
			case *[]interface{}:
				slice := covertSlice(vv)
				value := object2JSONString(slice)
				builder.WriteString(value)
			default:
				builder.WriteString(interface2String(v))
			}
			count++
			if len(data) != count {
				builder.WriteString(",")
			}
		}
		builder.WriteString("]")
	}

	return builder.String()
}

func interface2String(v interface{}) string {
	switch vv := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int:
		return strconv.Itoa(vv)
	case float64:
		return fmt.Sprintf("%f", vv)
	case bool:
		return strconv.FormatBool(vv)
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
	case Bool:
		v, _ := strconv.ParseBool(fmt.Sprint(r.object))
		if v {
			return 1
		} else {
			return 0
		}
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
	case Bool:
		v, _ := strconv.ParseBool(fmt.Sprint(r.object))
		if v {
			return 1.0
		} else {
			return 0.0
		}
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

// Map return map for object
func (r Result) Map() map[string]interface{} {
	return r.object.(map[string]interface{})
}

// Array return array for object
func (r Result) Array() []interface{} {
	arr := r.object.(*[]interface{})
	return covertSlice(arr)
}

func covertSlice(d *[]interface{}) []interface{} {
	var s []interface{}
	for _, e := range *d {
		s = append(s, e)
	}
	return s
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
