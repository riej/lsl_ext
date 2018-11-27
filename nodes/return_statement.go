package nodes

import (
	"fmt"
)

type ReturnStatement struct {
	NodeBase

	Value Expression
}

func (self *ReturnStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *ReturnStatement) StatementType() StatementType {
	return StatementReturn
}

func (self *ReturnStatement) String() string {
	if self.Value == nil {
		return "return;"
	} else {
		return fmt.Sprintf("return %s;", self.Value)
	}
}

func (self *ReturnStatement) ConnectTree() {
	if self.Value != nil {
		self.Value.SetParent(self)
		self.Value.SetScope(self.Scope)
		self.Value.SetScript(self.Script)
		self.Value.ConnectTree()
	}
}

func (self *ReturnStatement) GetChildren() []Node {
    if self.Value == nil {
        return []Node{}
    } else {
        return []Node{ self.Value }
    }
}

func (self *ReturnStatement) RealNode() Node {
    return self
}
