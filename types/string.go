package types

import (
	"text/scanner"
)

type StringValue struct {
	At    scanner.Position
	Value string
}

func (self *StringValue) Position() scanner.Position {
	return self.At
}

func (self *StringValue) Type() Type {
	return String
}

func (self *StringValue) String() string {
	return self.Value
}

func (self *StringValue) Clone(at scanner.Position) Value {
	return &StringValue{
		At:    at,
		Value: self.Value,
	}
}

func (self *StringValue) ToFloat() float64 {
    return 0;
}
