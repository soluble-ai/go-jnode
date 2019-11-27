package jnode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	Unknown = iota
	Text
	Number
	Bool
	Binary
	Array
	Object
	Missing
	Null
)

// MissingNode represents a missing node.  Path() will return
// a MissingNode if the field is not found.
var MissingNode *Node = &Node{}

// NullNode represents the nil (json null) value.
var NullNode *Node = &Node{}

// Node represents a JSON value (text, bool, numeric, object, or array.)
type Node struct {
	value interface{}
}

// NewNode creates a Node from a simple value (nil, string, bool
// int, int8, int16, int32, int64, float32, or float64.)
func NewNode(value interface{}) *Node {
	switch v := value.(type) {
	case nil:
		return NullNode
	case string, bool, int, int8, int16, int32, int64, float32, float64, []byte:
		return &Node{v}
	default:
		panic("NewNode accepts only simple values")
	}
}

// NewObjectNode creates a Node that wraps a map[string]interface{}.
// Object nodes correspond to JSON objects.
func NewObjectNode() *Node {
	return &Node{make(map[string]interface{})}
}

// NewArrayNode creates a Node that wraps a []interface{}.
// Array nodes corresponds to JSON arrays.
func NewArrayNode() *Node {
	a := make([]interface{}, 0, 5)
	return &Node{&a}
}

// FromJSON creates a Node from JSON
func FromJSON(data []byte) (*Node, error) {
	n := &Node{nil}
	if err := json.Unmarshal(data, n); err == nil {
		return n, nil
	} else {
		return MissingNode, err
	}
}

// String returns the Node formatted as JSON
func (n *Node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func unpointSlices(value interface{}) interface{} {
	switch v := value.(type) {
	case *[]interface{}:
		for i, e := range *v {
			(*v)[i] = unpointSlices(e)
		}
		return *v
	case map[string]interface{}:
		for k, val := range v {
			v[k] = unpointSlices(val)
		}
		return v
	default:
		return v
	}
}

func pointSlices(value interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		for i, e := range v {
			v[i] = pointSlices(e)
		}
		return &v
	case map[string]interface{}:
		for k, val := range v {
			v[k] = pointSlices(val)
		}
		return v
	default:
		return v
	}
}

// MarshalJSON is the custom JSON marshaller for a Node
func (n *Node) MarshalJSON() ([]byte, error) {
	u := unpointSlices(n.value)
	defer pointSlices(n.value)
	return json.Marshal(u)
}

// UnmarshalJSON is the custom JSON unmarshaller for a Node
func (n *Node) UnmarshalJSON(b []byte) error {
	var value interface{}
	err := json.Unmarshal(b, &value)
	if err == nil {
		n.value = pointSlices(value)
	}
	return err
}

// GetType returns the type of a Node
func (n *Node) GetType() int {
	if n == MissingNode {
		return Missing
	}
	if n == NullNode {
		return Null
	}
	switch n.value.(type) {
	case *[]interface{}:
		return Array
	case map[string]interface{}:
		return Object
	case string:
		return Text
	case int, int8, int16, int32, int64, float32, float64:
		return Number
	case bool:
		return Bool
	case []byte:
		return Binary
	default:
		return Unknown
	}
}

// IsArray returns true if the Node is an Array
func (n *Node) IsArray() bool {
	return n.GetType() == Array
}

// IsObject returns true if the Node is an Object
func (n *Node) IsObject() bool {
	return n.GetType() == Object
}

// IsMissing returns true if the Node is missing
func (n *Node) IsMissing() bool {
	return n == MissingNode
}

// IsNull returns true if the Node is the NullNode
func (n *Node) IsNull() bool {
	return n == NullNode
}

// Size returns the length of an Array, the number of fields
// in an Object, or 0 otherwise.
func (n *Node) Size() int {
	switch n.GetType() {
	case Object:
		return len(n.toMap())
	case Array:
		return len(*n.toSlicePtr())
	default:
		return 0
	}
}

// AsText returns the text of a Node
func (n *Node) AsText() string {
	return fmt.Sprintf("%v", n.value)
}

// AsInt returns the value of a Node as an int.  For a String,
// returns the value of strconv.Atoi or 0.  For bool returns 0 or 1.
// Otherwise returns 0.
func (n *Node) AsInt() int {
	switch v := n.value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case float32:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

// AsFloat returns the Node as a float64.  For string values
// it parses the string with strconv.ParseFloat or returns 0.
// For bool returns 0 or 1.  Otherwise returns 0.
func (n *Node) AsFloat() float64 {
	switch v := n.value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if i, err := strconv.ParseFloat(v, 64); err == nil {
			return i
		}
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

// AsBool returns the boolean value of a Node.  For strings,
// returns true if the string is equal to "true" (ignoring case).
// For numeric types, returns true if AsInt() != 0.
func (n *Node) AsBool() bool {
	switch v := n.value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(n.AsText(), "true")
	default:
		return n.AsInt() != 0
	}
}

// AsBinary returns the binary value of a Node.  If
// the Node is string, it attempts to base64 decode it.
func (n *Node) AsBinary() ([]byte, error) {
	switch n.GetType() {
	case Binary:
		return n.value.([]byte), nil
	case Text:
		s := n.value.(string)
		return base64.RawStdEncoding.DecodeString(s)
	default:
		return nil, fmt.Errorf("Node is not binary")
	}
}

func denode(value interface{}) interface{} {
	switch v := value.(type) {
	case *Node:
		return v.value
	default:
		return value
	}
}

// Put sets the value of a field in an Object Node.  Returns
// the Node (for chaining).  Panics if the Node is not an Object.
func (n *Node) Put(name string, value interface{}) *Node {
	if m, err := n.PutE(name, value); err == nil {
		return m
	} else {
		panic(err)
	}
}

// PutE sets the value of a field in an Object Node, or returns
// an error if the Node is not an Object.
func (n *Node) PutE(name string, value interface{}) (*Node, error) {
	if !n.IsObject() {
		return nil, fmt.Errorf("not an object")
	}
	n.toMap()[name] = denode(value)
	return n, nil
}

// PutObject sets the value of a field to a new Object Node
// and returns the new Object Node.  Panics if the Node is not
// an Object.
func (n *Node) PutObject(name string) *Node {
	if o, err := n.PutObjectE(name); err == nil {
		return o
	} else {
		panic(err)
	}
}

// PutObjectE sets the value of a field to a new Object Node
// and returns the new Object Node.  Returns an error if the
// Node is not an Object.
func (n *Node) PutObjectE(name string) (*Node, error) {
	if !n.IsObject() {
		return nil, fmt.Errorf("not an object")
	}
	o := make(map[string]interface{})
	n.toMap()[name] = o
	return &Node{o}, nil
}

// Entries returns the entries of an Object Node (or
// an empty map if the Node is not an Object.)
func (n *Node) Entries() map[string]*Node {
	if !n.IsObject() {
		return make(map[string]*Node)
	}
	m := n.toMap()
	e := make(map[string]*Node, len(m))
	for k, v := range m {
		e[k] = &Node{v}
	}
	return e
}

func (n *Node) toMap() map[string]interface{} {
	return n.value.(map[string]interface{})
}

func (n *Node) toSlicePtr() *[]interface{} {
	return n.value.(*[]interface{})
}

// PutArray sets a field in an Object Node to a new Array Node,
// and returns the new Array Node.  Panics if the Node is not
// an Object.
func (n *Node) PutArray(name string) *Node {
	if a, err := n.PutArrayE(name); err == nil {
		return a
	} else {
		panic(err)
	}
}

// PutArrayE sets a field in an Object Node to a new Array Node,
// and returns the new Array Node.  Returns an error if the
// Node is not an Object.
func (n *Node) PutArrayE(name string) (*Node, error) {
	if !n.IsObject() {
		return nil, fmt.Errorf("not an object")
	}
	a := make([]interface{}, 0, 5)
	n.toMap()[name] = &a
	return &Node{&a}, nil
}

// Append adds a new element to an Array Node and returns
// the Array.  Panics if the Node is not an Array.
func (n *Node) Append(value interface{}) *Node {
	if a, err := n.AppendE(value); err == nil {
		return a
	} else {
		panic(err)
	}
}

// AppendE adds a new element to an Array Node and returns
// the Array.  Returns an error if the Node is not an Array.
func (n *Node) AppendE(value interface{}) (*Node, error) {
	if !n.IsArray() {
		return nil, fmt.Errorf("node is not an array")
	}
	a := n.toSlicePtr()
	*a = append(*a, denode(value))
	return n, nil
}

// Elements returns the elements of an Array Node (or an
// empty slice if the Node is not an Array.)
func (n *Node) Elements() []*Node {
	if !n.IsArray() {
		return make([]*Node, 0)
	}
	a := *n.toSlicePtr()
	e := make([]*Node, len(a))
	for i, v := range a {
		e[i] = &Node{v}
	}
	return e
}

// Get returns the i-th element of an Array Node.  If
// the Node is not an Array, or the index is beyond the bounds
// of the Array, MissingNode is returned
func (n *Node) Get(i int) *Node {
	if i < 0 || !n.IsArray() {
		return MissingNode
	}
	a := *n.toSlicePtr()
	if i >= len(a) {
		return MissingNode
	}
	return &Node{a[i]}
}

// Path returns the value of a field of an Object Node.
// Returns MissingNode if the Node is not an Object or
// the field is not present
func (n *Node) Path(name string) *Node {
	if !n.IsObject() {
		return MissingNode
	}
	if value, ok := n.toMap()[name]; ok {
		return &Node{value}
	} else {
		return MissingNode
	}
}