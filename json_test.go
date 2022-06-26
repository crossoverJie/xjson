package gjson

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	str := `{"name":"cj"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
}
func TestDecode1(t *testing.T) {
	str := `{"name","cj"}`
	_, err := Decode(str)
	fmt.Println(err)
	assert.NotNil(t, err)
}
func TestDecode2(t *testing.T) {
	str := `{"name":"cj", "age" :"10"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	assert.Equal(t, v["age"], "10")
}
func TestDecode3(t *testing.T) {
	str := `{"name":"cj", "age" :10}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	assert.Equal(t, v["age"], 10)
}
func TestDecode4(t *testing.T) {
	str := `{"name":"cj", "age" :true,"x":false}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	assert.Equal(t, v["age"], true)
	assert.Equal(t, v["x"], false)
}
func TestDecode5(t *testing.T) {
	str := `{"name":"cj", "age" :true,"x":null}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	assert.Equal(t, v["age"], true)
	assert.Equal(t, v["x"], "")
}
func TestDecode6(t *testing.T) {
	str := `{"name":"cj", "age":{"a":"a","b":"b","c":true}, "e":"e"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	m := v["age"].(map[string]interface{})
	assert.Equal(t, m["a"], "a")
	assert.Equal(t, m["b"], "b")
	assert.Equal(t, m["c"], true)
	assert.Equal(t, v["e"], "e")

	indent, err := json.MarshalIndent(decode, "", "\t")
	assert.Nil(t, err)
	fmt.Println(string(indent))
}
func TestDecode66(t *testing.T) {
	str := `{"age":{"a":"a"}, "e":"e"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	m := v["age"].(map[string]interface{})
	assert.Equal(t, m["a"], "a")
	assert.Equal(t, v["e"], "e")
}
func TestDecode7(t *testing.T) {
	str := `{}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
}
func TestDecode8(t *testing.T) {
	str := `{"name":"cj", "age":[1,2,3],"e":"e"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	ints := v["age"].(*[]interface{})
	x := *ints
	assert.Equal(t, x[0], 1)
	assert.Equal(t, x[1], 2)
	assert.Equal(t, x[2], 3)
	assert.Equal(t, (*ints)[2], 3)
	assert.Equal(t, v["e"], "e")
}
func TestDecode10(t *testing.T) {
	str := `{"name":"cj", "age":["1","2"],"e":"e"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	strings := v["age"].(*[]interface{})
	assert.Equal(t, (*strings)[0], "1")
	assert.Equal(t, (*strings)[1], "2")
	assert.Equal(t, v["e"], "e")
}
func TestDecode11(t *testing.T) {
	str := `{"name":"cj", "age":[null,null],"e":"e"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	strings := v["age"].(*[]interface{})
	assert.Equal(t, (*strings)[0], "")
	assert.Equal(t, (*strings)[1], "")
	assert.Equal(t, v["e"], "e")
}
func TestDecode12(t *testing.T) {
	str := `{"name":"cj", "age":[true,false],"e":"e"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	strings := v["age"].(*[]interface{})
	assert.Equal(t, (*strings)[0], true)
	assert.Equal(t, (*strings)[1], false)
	assert.Equal(t, v["e"], "e")
}
func TestDecode13(t *testing.T) {
	str := `{"e":"e","a":[{"b":{"c":"c"}}],"d":"d"}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	strings := v["a"].(*[]interface{})
	ob := (*strings)[0]
	b := ob.(map[string]interface{})["b"]
	bb := b.(map[string]interface{})
	assert.Equal(t, bb["c"], "c")
	assert.Equal(t, v["d"], "d")
	assert.Equal(t, v["e"], "e")
}
func TestDecode14(t *testing.T) {
	str := `{"name":"cj", "age":[{"a":"a"}, {"d":1}],"e":[1,2,3]}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	assert.Equal(t, v["name"], "cj")
	ages := v["age"].(*[]interface{})
	a := (*ages)[0].(map[string]interface{})
	assert.Equal(t, a["a"], "a")
	d := (*ages)[1].(map[string]interface{})
	assert.Equal(t, d["d"], 1)

	e := v["e"].(*[]interface{})
	assert.Equal(t, (*e)[0], 1)
	assert.Equal(t, (*e)[1], 2)
	assert.Equal(t, (*e)[2], 3)
}
func TestDecode15(t *testing.T) {
	str := `{"e":[1,[2,3]]}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	e := v["e"].(*[]interface{})
	assert.Equal(t, (*e)[0], 1)
	arr2 := (*e)[1].(*[]interface{})
	assert.Equal(t, (*arr2)[0], 2)
	assert.Equal(t, (*arr2)[1], 3)
}
func TestDecode16(t *testing.T) {
	str := `{"e":[1,[2,3],{"d":"d"}]}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	e := v["e"].(*[]interface{})
	assert.Equal(t, (*e)[0], 1)
	arr2 := (*e)[1].(*[]interface{})
	assert.Equal(t, (*arr2)[0], 2)
	assert.Equal(t, (*arr2)[1], 3)
	arr3 := (*e)[2].(map[string]interface{})
	assert.Equal(t, arr3["d"], "d")

}
func TestDecode17(t *testing.T) {
	str := `{"e":[1,[2,3],{"d":{"f":"f"}}]}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	e := v["e"].(*[]interface{})
	assert.Equal(t, (*e)[0], 1)
	arr2 := (*e)[1].(*[]interface{})
	assert.Equal(t, (*arr2)[0], 2)
	assert.Equal(t, (*arr2)[1], 3)
	arr3 := (*e)[2].(map[string]interface{})
	f := arr3["d"].(map[string]interface{})
	assert.Equal(t, f["f"], "f")

}

func TestDecode18(t *testing.T) {
	str := `{
   "glossary": {
       "title": "example glossary",
		"age":1,
		"long":99.99,
		"GlossDiv": {
           "title": "S",
			"GlossList": {
               "GlossEntry": {
                   "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                       "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML", true, null]
                   },
					"GlossSee": "markup"
               }
           }
       }
   }
}`
	decode, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(decode)
	v := decode.(map[string]interface{})
	glossary := v["glossary"].(map[string]interface{})
	assert.Equal(t, glossary["title"], "example glossary")
	assert.Equal(t, glossary["age"], 1)
	assert.Equal(t, glossary["long"], 99.99)
	glossDiv := glossary["GlossDiv"].(map[string]interface{})
	assert.Equal(t, glossDiv["title"], "S")
	glossList := glossDiv["GlossList"].(map[string]interface{})
	glossEntry := glossList["GlossEntry"].(map[string]interface{})
	assert.Equal(t, glossEntry["ID"], "SGML")
	assert.Equal(t, glossEntry["SortAs"], "SGML")
	assert.Equal(t, glossEntry["GlossTerm"], "Standard Generalized Markup Language")
	assert.Equal(t, glossEntry["Acronym"], "SGML")
	assert.Equal(t, glossEntry["Abbrev"], "ISO 8879:1986")
	glossDef := glossEntry["GlossDef"].(map[string]interface{})
	assert.Equal(t, glossDef["para"], "A meta-markup language, used to create markup languages such as DocBook.")
	glossSeeAlso := glossDef["GlossSeeAlso"].(*[]interface{})
	assert.Equal(t, (*glossSeeAlso)[0], "GML")
	assert.Equal(t, (*glossSeeAlso)[1], "XML")
	assert.Equal(t, (*glossSeeAlso)[2], true)
	assert.Equal(t, (*glossSeeAlso)[3], "")
	assert.Equal(t, glossEntry["GlossSee"], "markup")
}

func TestDecode19(t *testing.T) {
	str := `[1,2]`
	parse, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(parse)
	v := parse.(*[]interface{})
	assert.Equal(t, (*v)[0], 1)
	assert.Equal(t, (*v)[1], 2)
}
func TestDecode20(t *testing.T) {
	str := `{"a":10,"b":10.9}`
	parse, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(parse)
	v := parse.(map[string]interface{})
	assert.Equal(t, v["a"], 10)
	assert.Equal(t, v["b"], 10.9)
}
func TestDecode21(t *testing.T) {
	str := `{"a":10,"b":[20,21.1,22.2]}`
	parse, err := Decode(str)
	assert.Nil(t, err)
	fmt.Println(parse)
	v := parse.(map[string]interface{})
	b := v["b"].(*[]interface{})
	assert.Equal(t, v["a"], 10)
	assert.Equal(t, (*b)[0], 20)
	assert.Equal(t, (*b)[1], 21.1)
	assert.Equal(t, (*b)[2], 22.2)
}
func TestDecodeErr1(t *testing.T) {
	str := `{"e":tr"}`
	_, err := Decode(str)
	assert.NotNil(t, err)
	fmt.Println(err)
}
func TestDecodeErr2(t *testing.T) {
	str := `{"e":true,"a":}`
	_, err := Decode(str)
	assert.NotNil(t, err)
	fmt.Println(err)
}
func TestDecodeErr3(t *testing.T) {
	str := `{"e":true:"a":1}`
	_, err := Decode(str)
	assert.NotNil(t, err)
	fmt.Println(err)
}
func TestDecodeErr4(t *testing.T) {
	str := `{"a":10,"b":10.9.1}`
	_, err := Decode(str)
	assert.NotNil(t, err)
	fmt.Println(err)
}
func TestSlice(t *testing.T) {
	arr := make([][]int, 0)
	arr = append(arr, []int{1})
	v := &arr[len(arr)-1]
	*v = append(*v, 2)
	fmt.Println(arr)
}

func TestSlice2(t *testing.T) {
	arr2 := make([][]interface{}, 0)
	arr2 = append(arr2, []interface{}{1})
	v2 := &arr2[len(arr2)-1]
	*v2 = append(*v2, 2)
	println(arr2)
}

func TestJsonIndent(t *testing.T) {
	str := `[1,2,3]`
	indent, err := json.MarshalIndent(str, "", "\t")
	assert.Nil(t, err)
	fmt.Println(string(indent))
	str = `{"name":1,"obj":{"age":2,"b":2,"c":[1,2,3]}}`
	var t1 T
	err = json.Unmarshal([]byte(str), &t1)
	assert.Nil(t, err)
	fmt.Println(t1)
	marshalIndent, err := json.MarshalIndent(&t1, "", "\t")
	assert.Nil(t, err)
	fmt.Println(string(marshalIndent))
}

type T struct {
	Name int `json:"name"`
	Obj  struct {
		Age int   `json:"age"`
		B   int   `json:"b"`
		C   []int `json:"c"`
	} `json:"obj"`
}
