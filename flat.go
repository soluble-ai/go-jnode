package jnode

import "container/list"

type frame struct {
	obj    *Node
	name   string
	append bool
}

func (f *frame) setValue(value interface{}) {
	if f.append {
		a := f.obj.Path(f.name)
		if !a.IsArray() {
			a = f.obj.PutArray(f.name)
		}
		a.Append(value)
	} else {
		f.obj.Put(f.name, value)
	}
}

func pop(stack *list.List) *frame {
	e := stack.Front()
	defer stack.Remove(e)
	return e.Value.(*frame)
}

func peek(stack *list.List) *frame {
	return stack.Front().Value.(*frame)
}

func newFrame() *frame {
	return &frame{NewObjectNode(), "", false}
}

// FromFlatString converts a string into an object Node.  It accepts a string in the form:
//
//     field1=value1,field2=value2,....
//     field1={subobj_field1=value1,...}
//     field1[list_value1,list_value2]
//
// The first form sets the fields and values of the object.  The second form sets the
// fields and values of a nested object.  The final form sets the values of an array.
func FromFlatString(s string) *Node {
	var state = ' '
	var nameStart, valueStart int
	stack := list.New()
	stack.PushFront(newFrame())
	for i, ch := range s {
		if state == ' ' && ch == ' ' {
			continue
		}
		top := peek(stack)
		switch state {
		case ' ':
			nameStart = i
			state = '='
		case '=':
			if ch == '=' || ch == '[' {
				name := s[nameStart:i]
				top.name = name
				top.append = ch == '['
				valueStart = i + 1
				state = '{'
			}
		case '{':
			if ch == '{' {
				stack.PushFront(newFrame())
				state = ' '
			} else if ch != ' ' {
				state = ','
			}
		case ',':
			if ch == ',' || ch == '}' || (top.append && ch == ']') {
				top.setValue(s[valueStart:i])
				if top.append {
					if ch == ']' {
						top.append = false
						top.name = ""
						state = ' '
					}
				} else {
					state = ' '
					if ch == '}' {
						val := pop(stack)
						peek(stack).setValue(val.obj)
						state = ' '
					}
				}
			}
		}
	}
	if state == ',' {
		top := stack.Front().Value.(*frame)
		top.setValue(s[valueStart:])
	}
	return pop(stack).obj
}
