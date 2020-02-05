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

package jnode

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestBasic(t *testing.T) {
	n := NewNode("hello")
	t.Log(n)
	asertJSON(t, n, `"hello"`)
	if n.AsText() != "hello" {
		t.Fail()
	}
	if n.AsBool() {
		t.Fail()
	}
	if n.AsInt() != 0 {
		t.Fail()
	}
	if n.AsFloat() != 0 {
		t.Fail()
	}
	if n.IsContainer() {
		t.Fail()
	}
	n = NewNode(nil)
	if !n.IsNull() || n.GetType() != Null {
		t.Errorf("null node bad")
	}
	if !MissingNode.IsMissing() || MissingNode.GetType() != Missing {
		t.Errorf("missing node bad")
	}
}

func TestMissing(t *testing.T) {
	if MissingNode.AsText() != "" {
		t.Error()
	}
	if !NewObjectNode().Path("Foo").IsMissing() {
		t.Error()
	}
}

func TestBadNode(t *testing.T) {
	assertPanic(t, func() { _ = NewNode([]int{1, 2, 3}) })
}

func TestBinary(t *testing.T) {
	pi := []byte("3.141")
	n := NewNode(pi)
	if b, _ := n.AsBinary(); !bytes.Equal(pi, b) {
		t.Fail()
	}
	if n.String() != `"My4xNDE="` {
		t.Errorf(n.String())
	}
	n, _ = FromJSON([]byte(`"My4xNDE"`))
	if b, _ := n.AsBinary(); !bytes.Equal(pi, b) {
		t.Fail()
	}
	n = NewNode(1)
	if _, err := n.AsBinary(); err == nil {
		t.Fail()
	}

}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestNumbers(t *testing.T) {
	n := NewNode(100)
	asertJSON(t, n, "100")
	if n.AsInt() != 100 || !n.AsBool() || n.AsText() != "100" ||
		n.String() != "100" || n.AsFloat() != 100.0 {
		t.Fail()
	}
	assert100(t, NewNode(int8(100)))
	assert100(t, NewNode(int16(100)))
	assert100(t, NewNode(int32(100)))
	assert100(t, NewNode(int64(100)))
	assert100(t, NewNode(uint(100)))
	assert100(t, NewNode(uint8(100)))
	assert100(t, NewNode(uint16(100)))
	assert100(t, NewNode(uint32(100)))
	assert100(t, NewNode(uint64(100)))
	assert100(t, NewNode(float32(100)))
	assert100(t, NewNode(float64(100)))
	assert100(t, NewNode("100"))
	assert0(t, NewObjectNode())
	assert0(t, NewArrayNode())
}

func assert0(t *testing.T, n *Node) {
	if n.AsInt() != 0 || n.AsFloat() != 0 || n.AsBool() {
		t.Errorf("%v is not 0", n)
	}
}

func assert100(t *testing.T, n *Node) {
	if n.AsInt() != 100 || n.AsFloat() != 100.0 || n.IsArray() ||
		n.IsObject() || len(n.Entries()) != 0 || len(n.Elements()) != 0 ||
		n.Size() != 0 {
		t.Errorf("%v is not 100", n)
	}
	if ty := n.GetType(); ty != Text {
		if n.AsText() != "100" || !n.AsBool() {
			t.Errorf("%v is not 100", n)
		}
	}
}

func TestBool(t *testing.T) {
	tr := NewNode(true)
	asertJSON(t, tr, "true")
	if !tr.AsBool() || tr.AsInt() != 1 || tr.AsFloat() != 1.0 {
		t.Fail()
	}
	fa := NewNode(false)
	if fa.AsBool() || fa.AsInt() != 0 || fa.AsFloat() != 0 {
		t.Fail()
	}
}

func TestFromJSON(t *testing.T) {
	n, _ := FromJSON([]byte(`
	{ "words": [ "one", "two", "three" ],
	  "numbers": [ 1, 2, 3 ],
	  "structs": [ { "four": 4 }, { "five": 5 }]
    }`))
	if n.Size() != 3 {
		t.Fail()
	}
	if !n.Path("words").IsArray() {
		t.Fail()
	}
	if n.Path("words").Get(2).AsText() != "three" {
		t.Fail()
	}
	if n.Path("structs").Get(1).Path("five").AsInt() != 5 {
		t.Fail()
	}
}

func asertJSON(t *testing.T, n *Node, value string) {
	b, err := json.Marshal(n)
	if err != nil {
		t.Error(err)
	}
	if string(b) != value {
		t.Errorf("%s != %s", string(b), value)
	}
	m := &Node{}
	if err := json.Unmarshal(b, m); err != nil {
		t.Error(err)
	}
	nt := n.GetType()
	mt := m.GetType()
	if nt != mt {
		t.Errorf("unmarshalled type %d is not %d", mt, nt)
	}
	b2, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	if string(b2) != string(b) {
		t.Errorf("json after round trip is not the same")
	}
}

func TestObject(t *testing.T) {
	n := NewObjectNode()
	n.Put("greeting", "hello")
	if n.Path("greeting").AsText() != "hello" {
		t.Fail()
	}
	n.PutObject("container").Put("order", 1).Put("value", true)
	if !n.Path("container").Path("value").AsBool() {
		t.Fail()
	}
	if !n.IsContainer() {
		t.Fail()
	}
	e := n.Entries()
	if len(e) != 2 || e["greeting"] == nil {
		t.Fail()
	}
	n.Put("empty", nil)
	if !n.Path("empty").IsNull() || n.Path("empty").GetType() != Null {
		t.Fail()
	}
	x := NewNode("hello")
	assertPanic(t, func() { x.Put("x", 1) })
	assertPanic(t, func() { x.PutObject("x") })
	assertPanic(t, func() { x.PutArray("x") })
	if !n.Path("foo").IsMissing() || !x.Path("foo").IsMissing() {
		t.Fail()
	}
}

func TestUnwrapNodes(t *testing.T) {
	n := NewObjectNode()
	n.Put("one", NewArrayNode().Append(1))
	if n.String() != `{"one":[1]}` {
		t.Fail()
	}
	if n.Path("one").Get(0).AsInt() != 1 {
		t.Fail()
	}
	if _, ok := n.Unwrap().(map[string]interface{}); !ok {
		t.Error("unwrap didn't")
	}
}

func TestArray(t *testing.T) {
	n := NewArrayNode()
	n.Append("hello")
	n.Append("world")
	if n.Size() != 2 {
		t.Fail()
	}
	e := n.Elements()
	if e[0].AsText() != "hello" {
		t.Fail()
	}
	if n.Get(1).AsText() != "world" {
		t.Fail()
	}
	if !n.Get(-1).IsMissing() || !n.Get(100).IsMissing() {
		t.Errorf("out of bounds wasn't missing")
	}
	n.Set(1, "everyone").Set(0, "howdy")
	if n.String() != `["howdy","everyone"]` {
		t.Fail()
	}
	n.Append([]interface{}{NewNode("hello"), 5})
	if n.Size() != 4 || n.Get(2).AsText() != "hello" || n.Get(3).AsInt() != 5 {
		t.Errorf("append didn't work")
	}
	assertPanic(t, func() { n.Set(-1, "fail") })
	x := NewNode("hello")
	assertPanic(t, func() { x.Append("foo") })
	assertPanic(t, func() { x.Set(0, "foo") })
}

func TestPutArray(t *testing.T) {
	n := NewObjectNode()
	a := n.PutArray("list")
	if n.ToMap()["list"] != a.value {
		t.Error("array value not correct")
	}
	if n.Size() != 1 {
		t.Error("object size not 1")
	}
	a.Append(1).Append(2)
	if a.Size() != 2 {
		t.Error("array size not 2")
	}
	asertJSON(t, n, `{"list":[1,2]}`)
}

func TestUnwrap(t *testing.T) {
	n := NewArrayNode().Append(1).Append(2)
	u := n.Unwrap().([]interface{})
	if len(u) != 2 || u[0].(int) != 1 {
		t.Fail()
	}
}

func TestMarshal(t *testing.T) {
	n := NewObjectNode()
	if err := json.Unmarshal([]byte(`{"username":"foo","token":"bar"}`), n); err != nil {
		t.Error(err)
	}
	if n.Path("username").AsText() != "foo" || n.Path("token").AsText() != "bar" {
		t.Error("node bad")
	}
}

func TestAppendSlice(t *testing.T) {
	n := NewArrayNode()
	n.Append([]string{"hello", "world"})
	if n.Size() != 2 {
		t.Errorf("size %d is not 2", n.Size())
	}
	if n.Get(0).AsText() != "hello" || n.Get(1).AsText() != "world" {
		t.Error("append didn't work")
	}
	n.Append([]int{1, 2, 3})
	if n.String() != `["hello","world",1,2,3]` {
		t.Errorf("%s isn't right", n.String())
	}
}

func TestAddArray(t *testing.T) {
	n := NewObjectNode()
	n.Put("a", NewArrayNode())
	n.Path("a").Append("hello").Append("world")
	if n.Path("a").Size() != 2 {
		t.Errorf("array doesn't contain 2 elements")
	}
	if n.String() != `{"a":["hello","world"]}` {
		t.Error("string form wrong")
	}
}

func TestFromMap(t *testing.T) {
	x := map[string]interface{}{"a": 1, "b": []interface{}{2, 3, 4}, "hello": "world"}
	n := FromMap(x)
	if n.Size() != 3 || n.Path("hello").AsText() != "world" || n.Path("b").Size() != 3 {
		t.Error("FromMap didn't work")
	}
}

func TestNull(t *testing.T) {
	var n *Node
	if !n.IsNull() || n.String() != "null" {
		t.Fail()
	}
	n = NewObjectNode().Put("x", nil)
	if !n.Path("x").IsNull() || n.String() != `{"x":null}` {
		t.Error(n)
	}
}

func TestFromSlice(t *testing.T) {
	a := FromSlice([]string{"hello", "world"})
	if !a.IsArray() || a.Size() != 2 || a.Get(1).AsText() != "world" {
		t.Error(a)
	}
	n := a.AppendObject()
	if !n.IsObject() {
		t.Fatal(a)
	}
	n.Put("one", 1)
	if a.Size() != 3 || a.Get(2).Path("one").AsInt() != 1 {
		t.Error(a)
	}
}

func TestRemove(t *testing.T) {
	var n *Node
	n.Remove("greeting")
	n = NewObjectNode().Put("greeting", "hello").Put("subject", "world")
	n.Remove("greeting")
	if !n.Path("greeting").IsMissing() || n.Path("subject").IsMissing() {
		t.Error(n)
	}
}
