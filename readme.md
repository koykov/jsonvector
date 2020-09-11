# JSON vector

JSON parser with minimum memory consumption. This project is a part of policy
to reducing memory consumption and amount of pointers.

As we know, JSON is a data interchange format that stores data as a tree.
And all known json parsers reproduces that tree in a memory somehow or other. 

This parser uses different way: it stores all parsed JSON nodes in a special array (vector).
That way protects from redundant memory allocations and reduces pointers.
In fact, vector contains only two pointers (array of nodes and array of indexes).
GC omits checking of that type of structs.

### Comparison

All known parsers has the following data structures (approximately) to represent JSON node value:
```go
type Node struct {
	typ Type    // [null, object, array, string, number, true, false]
	obj Object  // one pointer to slice and N*2 pointers inside KeyValue struct, see below
	arr []*Node // one pointer for the slice and N pointers for each array item
	str string  // one pointer
}

type Object []KeyValue

type KeyValue struct {
	key string // one pointer
	val *Node  // one pointer to node
}
```

As you see during parsing will be produced tons of pointers. Better consider this with an example.
Lets we have two JSON:
```json
{
  "a": true,
  "b": {
    "c": "foo",
    "d": [
      5,
      3.1415,
      812.48927
    ]
  }
}
```
and
```json
{
  "a": {
    "b": {
      "c": [
        "lorem ipsum",
        "dolor sit",
        "amet"
      ],
      "d": null
    }
  },
  "h": "bar"
}
```

See the picture below to see difference between classic JSON parser and vector:
