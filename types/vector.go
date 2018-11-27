package types

import (
	"fmt"
	"text/scanner"
)

type VectorValue struct {
	At scanner.Position
	X  float64
	Y  float64
	Z  float64
}

func (self *VectorValue) Position() scanner.Position {
	return self.At
}

func (self *VectorValue) Type() Type {
	return Vector
}

func (self *VectorValue) String() string {
	return fmt.Sprintf("<%g, %g, %g>", self.X, self.Y, self.Z)
}

func (self *VectorValue) Clone(at scanner.Position) Value {
	return &VectorValue{
		At: at,
		X:  self.X,
		Y:  self.Y,
		Z:  self.Z,
	}
}

func (self *VectorValue) Cast(_type Type, at scanner.Position) Value {
	switch _type {
	case String:
		return &StringValue{
			At:    at,
			Value: self.String(),
		}
	case List:
		return &ListValue{
			At: at,
			Values: []Value{
				self.Clone(at),
			},
		}
	case Vector:
		return self.Clone(at)
	default:
		return nil
	}
}

func (self *VectorValue) ToFloat() float64 {
    return 0
}
