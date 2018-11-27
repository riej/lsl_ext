package nodes

import (
	"fmt"
    "strings"

	"../types"
)

type LengthExpression struct {
	NodeBase

	RValue   Expression
}

func (self *LengthExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *LengthExpression) ExpressionType() ExpressionType {
	return ExpressionLength
}

func (self *LengthExpression) String() string {
    switch types.Type(self.RValue.ValueType().String()) {
    case types.String:
        return fmt.Sprintf("llStringLength(%s)", self.RValue);
    case types.List:
        rtype := string(self.RValue.ValueType())
        if strings.HasSuffix(rtype, "[]") {
            sname := strings.TrimSuffix(string(rtype), "[]")
            s, _ := self.Script.Structs[sname]
            if s == nil {
                self.Script.AddError(self, self.At, "undeclared struct \""+string(sname)+"\"")
                self.isValid = false
            } else {
                return fmt.Sprintf("(llGetListLength(%s) / %d)", self.RValue, len(s.Fields))
            }
        }

        return fmt.Sprintf("llGetListLength(%s)", self.RValue);
    default:
        return ""
    }
}

func (self *LengthExpression) ConnectTree() {
	self.RValue.SetParent(self)
    self.RValue.SetIndentationLevel(self.IndentationLevel)
	self.RValue.SetScope(self.Scope)
	self.RValue.SetScript(self.Script)
	self.RValue.ConnectTree()

    self.RValue = Expression(self.RValue.RealNode())

	self.isValid = true

    if self.RValue.IsValid() {
        switch types.Type(self.RValue.ValueType().String()) {
        case types.String, types.List:
        default:
            self.Script.AddError(self, self.At, "type mismatch (expected string or list, got " + string(self.RValue.ValueType()) + ")")
            self.isValid = false
        }
    } else {
        self.isValid = false
    }
}

func (self *LengthExpression) ValueType() types.Type {
	return types.Integer
}

func (self *LengthExpression) GetChildren() []Node {
    return []Node{ self.RValue }
}

func (self *LengthExpression) RealNode() Node {
    return self
}
