package jnode

import (
	"fmt"
	"strings"
)

// FromFlatString converts a string into an object Node.  It accepts a string in the form:
//
//     field1=value1,field2=value2,....
//     field1={subobj_field1=value1,...}
//     field1=[list_value1,list_value2]
//
// The first form sets the fields and values of the object.  The second form sets the
// fields and values of a nested object.  The final form sets the values of an array.
func FromFlatString(s string) (*Node, error) {
	node := NewObjectNode()
	for len(s) > 0 {
		n, err := parseAssignment(node, s, ",\x00")
		if err != nil {
			return nil, err
		}
		if len(s) == n {
			break
		}
		s = s[n+1:]
	}
	return node, nil
}

func parseAssignment(node *Node, s, terms string) (int, error) {
	eq := strings.Index(s, "=")
	if eq < 0 {
		return 0, fmt.Errorf("missing '=' in assignment")
	}
	name := s[0:eq]
	val, n, err := parseValue(s[eq+1:], terms)
	if err != nil {
		return 0, err
	}
	node.Put(name, val)
	return eq + n + 1, nil
}

func parseValue(s, terms string) (interface{}, int, error) {
	if len(s) == 0 {
		return "", 0, nil
	}
	switch s[0] {
	case '[':
		val, n, err := parseArray(s[1:])
		return val, n + 1, err
	case '{':
		val, n, err := parseObject(s[1:])
		return val, n + 1, err
	default:
		for i, ch := range s {
			if strings.ContainsRune(terms, ch) {
				return s[0:i], i, nil
			}
		}
		if !strings.Contains(terms, "\x00") {
			return nil, 0, fmt.Errorf("unexpected end of input")
		}
		return s, len(s), nil
	}
}

func parseArray(s string) (interface{}, int, error) {
	a := NewArrayNode()
	k := 0
	for {
		val, n, err := parseValue(s[k:], ",]")
		if err != nil {
			return nil, 0, err
		}
		k += n
		a.Append(val)
		if s[k] == ']' {
			return a, k + 1, nil
		}
	}
}

func parseObject(s string) (interface{}, int, error) {
	obj := NewObjectNode()
	k := 0
	for {
		n, err := parseAssignment(obj, s[k:], ",}")
		if err != nil {
			return nil, 0, err
		}
		k += n
		if s[k] == '}' {
			return obj, k + 1, nil
		}
		k++
	}
}
