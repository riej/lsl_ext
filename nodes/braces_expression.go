package nodes

import (
	"fmt"

	"../types"
)

type BracesExpression struct {
	NodeBase

	Child Expression
}

func (self *BracesExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *BracesExpression) ExpressionType() ExpressionType {
	return ExpressionBraces
}

func (self *BracesExpression) String() string {
    if self.Child.ExpressionType() == ExpressionBraces {
        return self.Child.String()
    } else {
    	return fmt.Sprintf("(%s)", self.Child)
    }
}

func (self *BracesExpression) ConnectTree() {
	self.Child.SetParent(self)
	self.Child.SetScope(self.Scope)
	self.Child.SetScript(self.Script)
	self.Child.ConnectTree()
}

func (self *BracesExpression) ValueType() types.Type {
	return self.Child.ValueType()
}

func (self *BracesExpression) GetChildren() []Node {
    return []Node{ self.Child }
}

func (self *BracesExpression) RealNode() Node {
    return self
}
