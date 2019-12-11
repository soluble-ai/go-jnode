// Code generated by "stringer -type=NodeType"; DO NOT EDIT.

package jnode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Text-1]
	_ = x[Number-2]
	_ = x[Bool-3]
	_ = x[Binary-4]
	_ = x[Array-5]
	_ = x[Object-6]
	_ = x[Missing-7]
	_ = x[Null-8]
}

const _NodeType_name = "UnknownTextNumberBoolBinaryArrayObjectMissingNull"

var _NodeType_index = [...]uint8{0, 7, 11, 17, 21, 27, 32, 38, 45, 49}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}