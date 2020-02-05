// Copyright 2019 Soluble Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

The Put methods accept simple types, slices, maps and other
Node's.  For complex types the argument will be copied,
and it may be modified (see implementation note below.)

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

For this reason, when jnode accepts an []interface{} or
map[string]interface{} (via Put, Append, FromSlice, or FromMap)
it will rewrite any slices to pointers to slices.

During JSON marshalling the code walks through the value replacing
pointers to slices with slices, then undoing the replacement afterwards.
Likewise, during unmarshalling the code walks through the value replacing
slices with pointers to the slices.

*/
package jnode
