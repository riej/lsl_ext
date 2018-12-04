package nodes

import (
	"fmt"
)

type StateStatement struct {
	NodeBase

	Name *Identifier
}

func (self *StateStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *StateStatement) StatementType() StatementType {
	return StatementState
}

func (self *StateStatement) String() string {
	if self.Name == nil {
		return "state default;"
	} else {
		return fmt.Sprintf("state %s;", self.Name)
	}
}

func (self *StateStatement) ConnectTree() {
    name := "default"

	if self.Name != nil {
		self.Name.SetParent(self)
		self.Name.SetScope(self.Scope)
		self.Name.SetScript(self.Script)
		self.Name.ConnectTree()

        name = self.Name.String()
	}

    self.isValid = true

    state, _ := self.Script.States[name]
    if state == nil {
        self.Script.AddError(self, self.At, fmt.Sprintf("undefined state \"%s\"", name))
        self.isValid = false
    } else {
        state.IsUsed = true
    }
}

func (self *StateStatement) GetChildren() []Node {
    if self.Name == nil {
        return []Node{}
    } else {
        return []Node{ self.Name }
    }
}

func (self *StateStatement) RealNode() Node {
    return self
}
