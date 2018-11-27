package nodes

import (
	"fmt"
)

type RotationExpression struct {
	NodeBase

	X Expression
	Y Expression
	Z Expression
	S Expression
}

func (self *RotationExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *RotationExpression) ExpressionType() ExpressionType {
	return ExpressionRotation
}

func (self *RotationExpression) String() string {
	return fmt.Sprintf("<%s, %s, %s, %s>", self.X, self.Y, self.Z, self.S)
}

func (self *RotationExpression) ConnectTree() {
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

	self.S.SetParent(self)
	self.S.SetScope(self.Scope)
	self.S.SetScript(self.Script)
	self.S.ConnectTree()
}

func (self *RotationExpression) GetChildren() []Node {
    return []Node{ self.X, self.Y, self.Z, self.S }
}

func (self *RotationExpression) RealNode() Node {
    return self
}
