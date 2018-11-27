package nodes

import (
	"../types"
)

type ExpressionStatement struct {
	NodeBase

	Expression Expression
}

func (self *ExpressionStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *ExpressionStatement) StatementType() StatementType {
	return StatementExpression
}

func (self *ExpressionStatement) String() string {
	return self.Expression.String() + ";"
}

func (self *ExpressionStatement) ConnectTree() {
	self.Expression.SetParent(self)
	self.Expression.SetIndentationLevel(self.IndentationLevel)
    self.Expression.SetScope(self.Scope)
	self.Expression.SetScript(self.Script)
	self.Expression.ConnectTree()
}

func (self *ExpressionStatement) ValueType() types.Type {
	return self.Expression.ValueType()
}

func (self *ExpressionStatement) GetChildren() []Node {
    return []Node{ self.Expression }
}

func (self *ExpressionStatement) RealNode() Node {
    return self
}
