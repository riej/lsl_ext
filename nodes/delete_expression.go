package nodes

import (
    "fmt"

    "../types"
)


type DeleteExpression struct {
    NodeBase

    RValue *ListItemExpression
}

func (self *DeleteExpression) NodeType() NodeType {
    return NodeExpression
}

func (self *DeleteExpression) ExpressionType() ExpressionType {
    return ExpressionDelete
}

func (self *DeleteExpression) String() string {
    if self.RValue == nil {
        return ""
    }

    li := self.RValue
    result := ""

    if li.LValue.ValueType() == types.String {
        result += fmt.Sprintf("llDeleteSubString(%s, %s, %s)", li.LValue, li.StartIndex, li.EndIndex)
    } else {
        result += fmt.Sprintf("llDeleteSubList(%s, %s, %s)", li.LValue, li.StartIndex, li.EndIndex)
    }

    if self.Parent.StatementType() == StatementExpression {
        result = fmt.Sprintf("%s = %s", li.LValue, result)
    }

    return result
}

func (self *DeleteExpression) ConnectTree() {
	self.RValue.SetParent(self)
    self.RValue.SetIndentationLevel(self.IndentationLevel)
	self.RValue.SetScope(self.Scope)
	self.RValue.SetScript(self.Script)
	self.RValue.ConnectTree()

    self.isValid = self.RValue.IsValid()
}

func (self *DeleteExpression) ValueType() types.Type {
    if self.RValue.LValue.ValueType() == types.String {
        return types.String
    } else {
        return types.List
    }
}

func (self *DeleteExpression) GetChildren() []Node {
    return []Node{ self.RValue }
}

func (self *DeleteExpression) RealNode() Node {
    return self
}
