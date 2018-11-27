package types

import (
	"text/scanner"
)

type ListValue struct {
	At     scanner.Position
	Values []Value
}

func (self *ListValue) Position() scanner.Position {
	return self.At
}

func (self *ListValue) Type() Type {
	return List
}

func (self *ListValue) String() string {
	result := "["
	for i, value := range self.Values {
		if i != 0 {
			result += ", "
		}
		result += value.String()
	}
	result += "]"
	return result
}

func (self *ListValue) Clone(at scanner.Position) Value {
	clone := &ListValue{
		At:     at,
		Values: make([]Value, len(self.Values)),
	}
	for i, value := range self.Values {
		clone.Values[i] = value.Clone(at)
	}
	return clone
}

func (self *ListValue) ToFloat() float64 {
    return 0;
}
