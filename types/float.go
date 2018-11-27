package types

import (
	"fmt"
	"text/scanner"
    "strings"
)

type FloatValue struct {
	At    scanner.Position
	Value float64
}

func (self *FloatValue) Position() scanner.Position {
	return self.At
}

func (self *FloatValue) Type() Type {
	return Float
}

func (self *FloatValue) String() string {
    str := fmt.Sprintf("%f", self.Value)
    if strings.Contains(str, ".") {
        str = strings.TrimRight(str, "0")
        if str[len(str) - 1] == '.' {
            str += "0"
        }
    }
    return str
}

func (self *FloatValue) Clone(at scanner.Position) Value {
	return &FloatValue{
		At:    at,
		Value: self.Value,
	}
}

func (self *FloatValue) Cast(_type Type, at scanner.Position) Value {
	switch _type {
	case Integer:
		return &IntegerValue{
			At:    at,
			Value: int32(self.Value),
		}
	case Float:
		return self.Clone(at)
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
	default:
		return nil
	}
}

func (self *FloatValue) ToFloat() float64 {
    return self.Value
}
