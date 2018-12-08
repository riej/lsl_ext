package nodes

import (
	"errors"
	"text/scanner"
    "fmt"
    "strings"
)

type ScriptError struct {
	At    scanner.Position
	Node  Node
	Error error
}

type Script struct {
	NodeBase

	Filename string

	Builtins *Builtins

	Globals []Node

    Identifiers map[string]Node

    Variables map[string]*Variable
	Functions map[string]*Function
    States map[string]*State
    Structs map[string]*Struct


    SkipUnused bool


	Errors []ScriptError
}

func (self *Script) NodeType() NodeType {
	return NodeScript
}

func (self *Script) String() string {
	var prev Node
	result := ""
	for _, child := range self.Globals {
		if child.StatementType() == StatementEmpty {
            prev = child
			continue
		}

        if child.NodeType() == NodeStruct {
            continue
        }

        if self.SkipUnused {
            switch child.NodeType() {
            case NodeFunction:
                if !child.(*Function).IsUsed {
                    continue
                }
            case NodeVariable:
                if !child.(*Variable).IsUsed {
                    continue
                }
            case NodeState:
                if !child.(*State).IsUsed {
                    continue
                }
            }
        }

		if prev != nil {
            if child.StatementType() == StatementComment {
                if child.Position().Line == prev.Position().Line {
                    result = strings.TrimRight(result, "\n\r\t ") + " "
                } else {
                    result = strings.TrimRight(result, "\n\r\t ") + "\n\n"
                }
            } else if prev.StatementType() != child.StatementType() && prev.StatementType() != StatementComment {
				result += "\n"
			} else {
			}
			/*
				if prev.NodeType() != NodeComment && child.NodeType() == NodeComment && prev.Position().Line == child.Position().Line {
					result += " "
				} else {
					result += "\n"
				}*/
		}

		result += child.String() + "\n"

		prev = child
	}

	return result
}

func (self *Script) ConnectTree() {
	self.Scope = &Scope{
		Variables: self.Builtins.Constants,
	}
	currScope := self.Scope

    self.Identifiers = make(map[string]Node)
    self.Variables = make(map[string]*Variable)
    self.Functions = make(map[string]*Function)
    self.States = make(map[string]*State)
    self.Structs = make(map[string]*Struct)

	self.Script = self



    if self.Builtins != nil {
        for _, child := range self.Builtins.Functions {
            name := child.Name.String()
            self.Identifiers[name] = child
            self.Functions[name] = child
        }
        for _, child := range self.Builtins.Constants {
            name := child.Name.String()
            self.Identifiers[name] = child
            self.Variables[name] = child
        }
    }



    var name string
	for _, child := range self.Globals {
		switch child.NodeType() {
        case NodeVariable:
            name = child.(*Variable).Name.String()
            if existing, ok := self.Identifiers[name]; ok {
                self.AddError(child, child.Position(), fmt.Sprintf("redeclared identifier \"%s\" (previously declared at %s)", name, existing.Position()))
            } else {
                self.Identifiers[name] = child
                self.Variables[name] = child.(*Variable)

                currScope.Variables = append(currScope.Variables, child.(*Variable))
            }
        case NodeFunction:
            name = child.(*Function).Name.String()
            if existing, ok := self.Identifiers[name]; ok {
                self.AddError(child, child.Position(), fmt.Sprintf("redeclared identifier \"%s\" (previously declared at %s)", name, existing.Position()))
            } else {
                self.Identifiers[name] = child
                self.Functions[name] = child.(*Function)
            }
        case NodeState:
            name = child.(*State).StateName()
            if existing, ok := self.Identifiers[name]; ok {
                self.AddError(child, child.Position(), fmt.Sprintf("redeclared identifier \"%s\" (previously declared at %s)", name, existing.Position()))
            } else {
                self.Identifiers[name] = child
                self.States[name] = child.(*State)
            }
        case NodeStruct:
            name = child.(*Struct).Name.String()
            if existing, ok := self.Identifiers[name]; ok {
                self.AddError(child, child.Position(), fmt.Sprintf("redeclared identifier \"%s\" (previously declared at %s)", name, existing.Position()))
            } else {
                self.Identifiers[name] = child
                self.Structs[name] = child.(*Struct)
            }
		}
	}

	for _, child := range self.Globals {
		child.SetParent(self)
		child.SetScope(currScope)
		child.SetScript(self.Script)
		child.ConnectTree()
	}
}

func (self *Script) AddError(node Node, at scanner.Position, message string) {
	self.Errors = append(self.Errors, ScriptError{
		Node:  node,
		At:    at,
		Error: errors.New(message),
	})
}

func (self *Script) GetChildren() []Node {
    return self.Globals
}

func (self *Script) RealNode() Node {
    return self
}
