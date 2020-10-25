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

Majority of JSON parsers will build array of nodes like:

| 0               | 1          | 2               | 3           | 4                   | 5           | 6           | 7              |
|-----------------|------------|-----------------|-------------|---------------------|-------------|-------------|----------------|
| type: obj       | type: bool | type: obj       | type: str   | type: arr           | type: num   | type: num   | type: num      |
| key: ""         | key: "a"   | key: ""         | key: "c"    | key: ""             | key: ""     | key: ""     | key: ""        |
| str: ""         | str: ""    | str: ""         | str: "foo"  | str: ""             | str: ""     | str: ""     | str: ""        |
| num: 0          | num: 0     | num: 0          | num: 0      | num: 0              | num: 5      | num: 3.1415 | num: 812.48927 |
| bool: false     | bool: true | bool: false     | bool: false | bool: false         | bool: false | bool: false | bool: false    |
| child: [*1, *2] | child: []  | child: [*3, *4] | child: []   | child: [*5, *6, *7] | child: []   | child: []   | child: []      |

As you can see, independent of JSON node type, each parsed node contains at least 3 poiters:
* key (string)
* str (string)
* child (slice of node pointers)

JSON vector has different approach and build the following array of nodes and index:

Vector:

| 0          | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
|------------|---|---|---|---|---|---|---|
| type: obj  | type: bool | | | | | | |
| key pos: 0 | key pos:   | | | | | | |
| key len: 0 | key len:   | | | | | | |
| val pos: 0 | val pos:   | | | | | | |
| val len: 0 | val len:   | | | | | | |
| depth: 0   | depth: 1   | | | | | | |
| idx pos: 0 | idx pos: 0 | | | | | | |
| idx len: 2 | idx len: 0 | | | | | | |

Index (X-axis means depth):

| 0 | 1 | 2 | 3 |
|---|---|---|---|
| 0 ||||
| - ||||
| - ||||
| - ||||
| - ||||
| - ||||
| - ||||
