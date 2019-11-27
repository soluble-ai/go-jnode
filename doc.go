/*
The jnode package handles generic JSON documents building on the
support provided by Go's "encoding/json".

The API is loosely based on the Java's
https://github.com/FasterXML/jackson-databind.

Creating Node's

The Node type represents a generic JSON node (e.g. text, bool, null,
int, float, array, or object.)  Here's some ways to create a Node:

	t := jnode.NewNode("hello")
	i := jnode.NewNode(10)

	// return values chain
	o := jnode.NewObject().Put("greeting", "hello").Put("subject", "world")
	a := jnode.NewArray().Append(1).Append(2).Append(3)

	o.PutArray("list").Append(true).Append(100.0)
	o.PutObject("struct").Put("one", 1)

	// Put/Append can take other nodes or simple types
	a.Append(jnode.NewObject().Put("two", 2))

	// Can also build from JSON
	n, _ := json.FromJSON([]byte(`{"three": 3}`)

Navigation

For Object Node's, use Path:

	o.Path("greeting").AsText()         // "hello"

	o.Path("x").Path("y").AsText()      // "", missing paths return jnode.MissingNode

	jnode.NewNode("hello").Path("foo")  // also returns jnode.MissingNode

All elements of an Object can be accessed via Entries().
All elements of an Array Node can be accessed via Elements().
Both methods return empty maps or slices if the Node is not an Object or Array
respectively.

JSON Marshal

A Node's String() method returns JSON:

	o := jnode.NewObject().put("one", 1)
	fmt.Println(o)  // {"one":1}

Note that because jnode is building on the generic JSON support in Go, the
order of fields in an Object is unpredictable.

Implementation Note

A Node holds an interface{}, which stores the unwrapped Go value.
(A Node does not contain another Node, only the basic generic JSON values.)
For Array values ([]interface{}) the Node holds a pointer to the slice
instead of the slice directly.  (This is necessary because Append creates
new slices.)

During JSON marshalling the code walks through the value replacing
pointers to slices with slices, then undoing the replacement afterwards.
Likewise, during unmarshalling the code walks through the value replacing
slices with pointers to the slices.

*/
package jnode
