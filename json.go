package xjson

import (
	"errors"
	"fmt"
	"github.com/crossoverJie/gscript/syntax"
	"reflect"
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
	if len(tokenize) == 0 {
		return nil, errors.New("input is empty")
	}
	reader := NewTokenReader(tokenize)
	return Parse(reader)
}

func Encode(input string, v interface{}) error {
	value := reflect.ValueOf(v)
	t := value.Type()
	if k := t.Kind(); k != reflect.Ptr {
		return fmt.Errorf("invalid type:%v non-point", t)
	}
	elem := t.Elem()
	loopElem(elem, "")
	//for i := 0; i < elem.NumField(); i++ {
	//	ft := elem.Field(i)
	//	fmt.Printf("%s, %v\n", ft.Name,ft.Type.Kind())
	//}

	return nil
}
func loopElem(elem reflect.Type, indent string) {
	for i := 0; i < elem.NumField(); i++ {
		ft := elem.Field(i)
		fmt.Printf("%s %s, %v\n", indent, ft.Name, ft.Type.Kind())
		if ft.Type.Kind() == reflect.String {
			v := reflect.New(ft.Type)
			v.SetString("1")
		}
		if ft.Type.Kind() == reflect.Struct {
			loopElem(ft.Type, indent+"\t")
		} else if ft.Type.Kind() == reflect.Slice {
			t := ft.Type.Elem()
			if t.Kind() != reflect.Interface {
				loopElem(t, indent+"\t")
			}
		}
	}
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

	return getWithMap(root, grammar)

}

func getWithMap(root map[string]interface{}, grammar string) Result {
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
				Token: token,
				//Raw:   v,
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
				Token: token,
				//Raw:   v,
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
			result := getWithMap(root, read.Value)
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

func (r Result) Map() map[string]interface{} {
	return r.object.(map[string]interface{})
}
func (r Result) Array() *[]interface{} {
	return r.object.(*[]interface{})
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
