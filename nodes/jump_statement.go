package nodes

import (
	"fmt"
)

type JumpStatement struct {
	NodeBase

	Name *Identifier
}

func (self *JumpStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *JumpStatement) StatementType() StatementType {
	return StatementJump
}

func (self *JumpStatement) String() string {
	return fmt.Sprintf("jump %s;", self.Name)
}

func (self *JumpStatement) ConnectTree() {
	self.Name.SetParent(self)
	self.Name.SetScope(self.Scope)
	self.Name.SetScript(self.Script)
	self.Name.ConnectTree()
}

func (self *JumpStatement) GetChildren() []Node {
    return []Node{ self.Name }
}

func (self *JumpStatement) RealNode() Node {
    return self
}
