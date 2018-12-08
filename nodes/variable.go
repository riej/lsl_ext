package nodes

import (
    "strings"

	"../types"
)

type Variable struct {
	NodeBase

	Name   *Identifier
	Type   types.Type
	RValue Node

	IsArgument bool // is function argument?
    IsConstant bool


    IsUsed bool
}

func (self *Variable) NodeType() NodeType {
	return NodeVariable
}

// Variable declaration also can be statement
func (self *Variable) StatementType() StatementType {
	return StatementVariable
}

func (self *Variable) String() string {
	result := self.Type.String() + " " + self.Name.String()
	if self.RValue != nil {
		result += " = " + self.RValue.String()
	}

	if !self.IsArgument {
		result += ";"
	}

	return result
}

func (self *Variable) ValueString() string {
    if self.RValue != nil {
        return self.RValue.String()
    }

    switch self.Type {
    case types.String:
        return "\"\""
    case types.Key:
        return "NULL_KEY"
    case types.Integer, types.Float:
        return "0"
    case types.Vector:
        return "<0.0, 0.0, 0.0>"
    case types.Rotation:
        return "<0.0, 0.0, 0.0, 0.0>"
    default:
        return "[]"
    }
}

func (self *Variable) ConnectTree() {
	self.Name.SetParent(self)
	self.Name.SetScope(self.Scope)
	self.Name.SetScript(self.Script)
	self.Name.ConnectTree()

	if self.RValue != nil {
        if self.RValue.ExpressionType() == ExpressionListItem {
            self.RValue.(*ListItemExpression).Type = self.Type
        }

		self.RValue.SetParent(self)
		self.RValue.SetScope(self.Scope)
		self.RValue.SetScript(self.Script)
		self.RValue.ConnectTree()

        if self.IsConstant && self.Type == types.Unknown {
            self.Type = self.RValue.ValueType()
        }
	}

	self.isValid = true

    if self.IsConstant && self.RValue == nil {
		self.Script.AddError(self, self.At, "constant must have value")
		self.isValid = false
    }

    if strings.HasSuffix(string(self.Type), "[]") {
        sname := strings.TrimSuffix(string(self.Type), "[]")
        if _, ok := self.Script.Structs[sname]; !ok {
            self.Script.AddError(self, self.At, "undeclared identifier \"" + sname + "\"")
            self.isValid = false
        }
    }

	if self.RValue != nil && self.RValue.IsValid() && !self.Type.IsCompatible(self.RValue.ValueType()) {
		self.Script.AddError(self, self.At, "type mismatch")
		self.isValid = false
	}
}

func (self *Variable) ValueType() types.Type {
    return self.Type
}

func (self *Variable) GetChildren() []Node {
    result := []Node{ self.Name }
    if self.RValue != nil {
        result = append(result, self.RValue)
    }
    return result
}

func (self *Variable) RealNode() Node {
    return self
}
