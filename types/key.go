package types

import "text/scanner"

type KeyValue struct {
	At    scanner.Position
	Value string
}

func (self *KeyValue) Position() scanner.Position {
	return self.At
}

func (self *KeyValue) Type() Type {
	return Key
}

func (self *KeyValue) String() string {
	return self.Value
}

func (self *KeyValue) Clone(at scanner.Position) Value {
	return &KeyValue{
		At:    at,
		Value: self.Value,
	}
}

func (self *KeyValue) ToFloat() float64 {
    return 0;
}
