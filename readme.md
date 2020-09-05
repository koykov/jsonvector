# JSON vector

JSON parser with minimum memory consumption. This project is a part of policy
to reducing memory consumption and amount of pointers.

As we know, JSON is a data interchange format that stores data as a tree.
And all known json parsers reproduces that tree in a memory. 

This parser uses different way: it stores all parsed JSON nodes in a special array (vector).
That way protects from redundant memory allocations and reduces pointers.
In fact, vector contains only two pointers (array of nodes and array of indexes).
GC omits checking of that type of structs.
