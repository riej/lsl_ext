package nodes

import (
	"fmt"

	"../types"
)

type UnaryExpression struct {
	NodeBase

	Operator string
	RValue   Expression

	IsPostfix bool // false = ++a, true = a++
}

func (self *UnaryExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *UnaryExpression) ExpressionType() ExpressionType {
	return ExpressionUnary
}

func (self *UnaryExpression) String() string {
	if self.IsPostfix {
		return fmt.Sprintf("%s%s", self.RValue, self.Operator)
	} else {
		return fmt.Sprintf("%s%s", self.Operator, self.RValue)
	}
}

func (self *UnaryExpression) ConnectTree() {
	self.RValue.SetParent(self)
    self.RValue.SetIndentationLevel(self.IndentationLevel)
	self.RValue.SetScope(self.Scope)
	self.RValue.SetScript(self.Script)
	self.RValue.ConnectTree()

    self.RValue = Expression(self.RValue.RealNode())

	self.isValid = true

    if self.RValue.ExpressionType() == ExpressionListItem {
        self.RValue.(*ListItemExpression).Type = types.Integer
    } else if self.RValue.IsValid() && !self.RValue.ValueType().IsCompatible(types.Float) {
		self.Script.AddError(self, self.At, "type mismatch")
		self.isValid = false
	}
}

func (self *UnaryExpression) ValueType() types.Type {
	return self.RValue.ValueType()
}

func (self *UnaryExpression) GetChildren() []Node {
    return []Node{ self.RValue }
}

func (self *UnaryExpression) RealNode() Node {
    return self
}
