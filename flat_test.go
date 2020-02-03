package jnode

import (
	"testing"
)

func TestFromFlat(t *testing.T) {
	n, err := FromFlatString("a=b,c=5")
	if err != nil {
		t.Error(err)
	}
	if !n.IsObject() || n.Size() != 2 || n.Path("a").AsText() != "b" {
		t.Fail()
	}
}

func TestFromFlatNested(t *testing.T) {
	n, err := FromFlatString("obj={b=1,c=2}")
	if err != nil {
		t.Error(err)
	}
	if !n.IsObject() || n.Size() != 1 || n.Path("obj").Size() != 2 || n.Path("obj").Path("c").AsInt() != 2 {
		t.Error(n)
	}
}

func TestFromFlatArray(t *testing.T) {
	n, err := FromFlatString("a=[hello]")
	if err != nil {
		t.Error(err)
	}
	if !n.IsObject() || n.Size() != 1 ||
		!n.Path("a").IsArray() || n.Path("a").Size() != 1 || n.Path("a").Get(0).AsText() != "hello" {
		t.Error(n)
	}

}

func TestComplex(t *testing.T) {
	n, err := FromFlatString("containers=[{image=perl}]")
	if err != nil {
		t.Error(err)
	}
	if !n.IsObject() || n.Path("containers").Get(0).Path("image").AsText() != "perl" {
		t.Error(n)
	}
}
