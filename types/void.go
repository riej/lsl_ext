package types

import (
	"text/scanner"
)

type VoidValue struct {
	At scanner.Position
}

func (self *VoidValue) Position() scanner.Position {
	return self.At
}

func (self *VoidValue) Type() Type {
	return Void
}

func (self *VoidValue) String() string {
	return ""
}

func (self *VoidValue) Clone(at scanner.Position) Value {
	return &VoidValue{
		At: at,
	}
}

func (self *VoidValue) Cast(_type Type, at scanner.Position) Value {
	return nil
}

func (self *VoidValue) Compatible(_type Type) bool {
	return false
}

func (self *VoidValue) ToFloat() float64 {
    return 0;
}
