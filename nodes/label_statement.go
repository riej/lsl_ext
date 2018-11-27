package nodes

import (
	"fmt"
)

type LabelStatement struct {
	NodeBase

	Name *Identifier
}

func (self *LabelStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *LabelStatement) StatementType() StatementType {
	return StatementLabel
}

func (self *LabelStatement) String() string {
	return fmt.Sprintf("@%s;", self.Name)
}

func (self *LabelStatement) ConnectTree() {
	self.Name.SetParent(self)
	self.Name.SetScope(self.Scope)
	self.Name.SetScript(self.Script)
	self.Name.ConnectTree()
}

func (self *LabelStatement) GetChildren() []Node {
    return []Node{ self.Name }
}
