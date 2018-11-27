package nodes

import (
	"fmt"
)

type VectorExpression struct {
	NodeBase

	X Expression
	Y Expression
	Z Expression
}

func (self *VectorExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *VectorExpression) ExpressionType() ExpressionType {
	return ExpressionVector
}

func (self *VectorExpression) String() string {
	return fmt.Sprintf("<%s, %s, %s>", self.X, self.Y, self.Z)
}

func (self *VectorExpression) ConnectTree() {
	self.X.SetParent(self)
	self.X.SetScope(self.Scope)
	self.X.SetScript(self.Script)
	self.X.ConnectTree()

	self.Y.SetParent(self)
	self.Y.SetScope(self.Scope)
	self.Y.SetScript(self.Script)
	self.Y.ConnectTree()

	self.Z.SetParent(self)
	self.Z.SetScope(self.Scope)
	self.Z.SetScript(self.Script)
	self.Z.ConnectTree()
}

func (self *VectorExpression) GetChildren() []Node {
    return []Node{ self.X, self.Y, self.Z }
}

func (self *VectorExpression) RealNode() Node {
    return self
}
