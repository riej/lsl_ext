package types

import (
	"fmt"
	"text/scanner"
)

type RotationValue struct {
	At scanner.Position
	X  float64
	Y  float64
	Z  float64
	S  float64
}

func (self *RotationValue) Position() scanner.Position {
	return self.At
}

func (self *RotationValue) Type() Type {
	return Rotation
}

func (self *RotationValue) String() string {
	return fmt.Sprintf("<%g, %g, %g, %g>", self.X, self.Y, self.Z, self.S)
}

func (self *RotationValue) Clone(at scanner.Position) Value {
	return &RotationValue{
		At: at,
		X:  self.X,
		Y:  self.Y,
		Z:  self.Z,
		S:  self.S,
	}
}

func (self *RotationValue) ToFloat() float64 {
    return 0;
}
