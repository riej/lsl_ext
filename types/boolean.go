package types

import (
	"text/scanner"
)

type BooleanValue struct {
	At    scanner.Position
	Value bool
}

func (self *BooleanValue) Position() scanner.Position {
	return self.At
}

func (self *BooleanValue) Type() Type {
	return Integer
}

func (self *BooleanValue) String() string {
    if self.Value {
        return "TRUE"
    } else {
        return "FALSE"
    }
}

func (self *BooleanValue) Clone(at scanner.Position) Value {
	return &BooleanValue{
		At:    at,
		Value: self.Value,
	}
}

func (self *BooleanValue) ToFloat() float64 {
    if self.Value {
        return 1
    } else {
        return 0
    }
}
