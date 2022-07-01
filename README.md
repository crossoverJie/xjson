
`gjson` is a `JSON` parsing library, you can query `JSON` like `OOP`. 

```go
first := Get(`{"people":{"name":{"first":"bob"}}}`, "people.name.first")
assert.Equal(t, first.String(), "bob")
```

You can even perform arithmetic operations with `JSON`.

```go
str := `{"people":[{"bob":{"age":10}},{"alice":{"age":10}}]}`
age := gjson.GetWithArithmetic(str, "people[0].bob.age + people[1].alice.age")
assert.Equal(t, age.Int(), 20)
```


[![codecov](https://codecov.io/gh/crossoverJie/gjson/branch/main/graph/badge.svg?token=51WIOVFN95)](https://codecov.io/gh/crossoverJie/gjson)


# Usage

```shell
go get github.com/crossoverJie/gjson
```

```go
import "github.com/crossoverJie/gjson"
func TestJson(t *testing.T) {
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
	decode, err := gjson.Decode(str)
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
```

# Features
- [x] Support syntax: `gjson.Get("glossary.title")`
- [ ] Support arithmetic operators: `gjson.Get("glossary.age+long")`
- [ ] Resolve to struct
