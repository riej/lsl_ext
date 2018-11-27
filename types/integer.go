package types

import (
	"strconv"
	"text/scanner"
)

type IntegerValue struct {
	At    scanner.Position
	Value int32
	IsHex bool
}

func (self *IntegerValue) Position() scanner.Position {
	return self.At
}

func (self *IntegerValue) Type() Type {
	return Integer
}

func (self *IntegerValue) String() string {
	if self.IsHex {
		return "0x" + strconv.FormatInt(int64(self.Value), 16)
	} else {
		return strconv.FormatInt(int64(self.Value), 10)
	}
}

func (self *IntegerValue) Clone(at scanner.Position) Value {
	return &IntegerValue{
		At:    at,
		Value: self.Value,
        IsHex: self.IsHex,
	}
}

func (self *IntegerValue) ToFloat() float64 {
    return float64(self.Value)
}
