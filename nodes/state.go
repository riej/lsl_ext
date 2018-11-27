package nodes

import (
	"strings"
)

type State struct {
	NodeBase

	Name   *Identifier
	Events []*Function
}

func (self *State) NodeType() NodeType {
	return NodeState
}

func (self *State) String() string {
	result := ""
	if self.Name == nil {
		result += "default"
	} else {
		result += self.Name.String()
	}
	result += " {\n"
	for i, event := range self.Events {
		if i > 0 {
			result += "\n"
		}
		result += strings.Repeat(Indentation, event.GetIndentationLevel())
		result += event.String()
	}
	result += "}\n"

	return result
}

func (self *State) StateName() string {
	if self.Name == nil {
		return "default"
	} else {
		return self.Name.String()
	}
}

func (self *State) ConnectTree() {
	if self.Name != nil {
		self.Name.SetParent(self)
		self.Name.SetScope(self.Scope)
		self.Name.SetScript(self.Script)
		self.Name.ConnectTree()
	}

	for _, child := range self.Events {
		child.SetParent(self)
		child.SetIndentationLevel(self.IndentationLevel + 1)
		child.SetScope(self.Scope) // events can call other events
		child.SetScript(self.Script)
		child.IsStateEvent = true
		child.ConnectTree()
	}

	var name string
	if self.Name == nil {
		name = "default"
	} else {
		name = self.Name.String()
		if name == "default" {
			self.Script.AddError(self, self.At, "invalid state name \""+name+"\"")
			return
		}
	}

    if existing, _ := self.Script.States[name]; existing != self {
		self.Script.AddError(self, self.At, "duplicate state \""+name+"\" (previously declared at " + existing.Position().String() + ")")
	}
}

func (self *State) GetChildren() []Node {
    result := []Node{}
    if self.Name != nil {
        result = append(result, self.Name)
    }
    for _, child := range self.Events {
        result = append(result, child)
    }
    return result
}

func (self *State) RealNode() Node {
    return self
}
