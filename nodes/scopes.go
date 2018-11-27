package nodes

import (
	"fmt"
	"text/scanner"
)

type Scope struct {
	At     scanner.Position
	Parent *Scope

	Variables []*Variable

	Children []*Scope
}

func (self *Scope) Position() scanner.Position {
	return self.At
}

func (self *Scope) Clone(at scanner.Position) *Scope {
	scope := &Scope{
		At:     at,
		Parent: self,

		Variables: append([]*Variable(nil), self.Variables...),
	}
	self.Children = append(self.Children, scope)
	return scope
}

func (self *Scope) AddVariable(variable *Variable) *Scope {
	cloned := self.Clone(variable.At)
	cloned.Variables = append([]*Variable{variable}, cloned.Variables...)
	return cloned
}

func (self *Scope) FindVariable(name string) *Variable {
	for _, node := range self.Variables {
		if node.Name.String() == name {
			return node
		}
	}

	return nil
}

func (self *Scope) Find(name string) Node {
	var node Node

	node = self.FindVariable(name)
	if node != nil {
		return node
	}

	return nil
}

func (self *Scope) DumpTree() {
	scope := self
	for scope != nil {
		fmt.Printf("-- %p %s\n", scope, scope.At)
		scope = scope.Parent
	}
}
