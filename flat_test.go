package jnode

import "testing"

import "container/list"

func TestPop(t *testing.T) {
	stack := list.New()
	stack.PushFront(newFrame())
	f := pop(stack)
	if f == nil || stack.Len() != 0 {
		t.Fail()
	}
}

func TestFromFlat(t *testing.T) {
	n := FromFlatString("a=b,c=5")
	if !n.IsObject() || n.Size() != 2 || n.Path("a").AsText() != "b" {
		t.Fail()
	}
}

func TestFromFlatNested(t *testing.T) {
	n := FromFlatString("obj={b=1,c=2}")
	if !n.IsObject() || n.Size() != 1 || n.Path("obj").Size() != 2 || n.Path("obj").Path("c").AsInt() != 2 {
		t.Error(n)
	}
}

func TestFromFlatArray(t *testing.T) {
	n := FromFlatString("a[hello]")
	if !n.IsObject() || n.Size() != 1 ||
		!n.Path("a").IsArray() || n.Path("a").Size() != 1 || n.Path("a").Get(0).AsText() != "hello" {
		t.Error(n)
	}

}
