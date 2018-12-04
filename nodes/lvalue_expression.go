package nodes

import (
	"../types"
)

type LValueExpression struct {
	NodeBase

	Name *Identifier
	Item *Identifier

    rnode Node
}

func (self *LValueExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *LValueExpression) ExpressionType() ExpressionType {
	return ExpressionLValue
}

func (self *LValueExpression) String() string {
    result := self.Name.String()

    if self.Item != nil {
        result += "." + self.Item.String()
    }

    return result
}

func (self *LValueExpression) GetStruct() *Struct {
    variable := self.Scope.FindVariable(self.Name.String())
    if variable == nil || !variable.Type.IsStruct() {
        return nil
    }

    s, _ := self.Script.Structs[string(variable.Type)]
    return s
}

func (self *LValueExpression) GetVariable() *Variable {
    return self.Scope.FindVariable(self.Name.String())
}

func (self *LValueExpression) ConnectTree() {
    self.rnode = self

	self.Name.SetParent(self)
	self.Name.SetScope(self.Scope)
	self.Name.SetScript(self.Script)
	self.Name.ConnectTree()

	if self.Item != nil {
		self.Item.SetParent(self)
		self.Item.SetScope(self.Scope)
		self.Item.SetScript(self.Script)
		self.Item.ConnectTree()
	}

	self.isValid = true

    variable := self.Scope.FindVariable(self.Name.String())
	if variable == nil {
		self.Script.AddError(self, self.At, "undeclared identifier \""+self.Name.String()+"\"")
		self.isValid = false
    } else {
        variable.IsUsed = true

        if variable.Type.IsStruct() && self.Item != nil {
            s, _ := self.Script.Structs[string(variable.Type)]
            if s == nil {
                self.Script.AddError(self, self.At, "undeclared struct \""+string(variable.Type)+"\"")
                self.isValid = false
            } else {
                field := s.GetField(self.Item.String())
                index := s.GetFieldIndex(self.Item.String())
                if field == nil {
                    self.Script.AddError(self, self.At, "undeclared struct field \""+self.String()+"\"")
                    self.isValid = false
                } else {
                    lvalue := &LValueExpression{
                        Name: self.Name,
                    }
                    lvalue.At = self.At

                    cindex := &Constant{
                        Value: &types.IntegerValue{
                            At: self.At,
                            Value: index,
                        },
                    }
                    cindex.At = self.At

                    self.rnode = &ListItemExpression{
                        LValue: lvalue,
                        Type: field.Type,
                        IsRange: false,
                        StartIndex: cindex,
                        EndIndex: cindex,
                    }
                    self.rnode.SetPosition(self.At)

                    self.rnode.SetParent(self)
                    self.rnode.SetScope(self.Scope)
                    self.rnode.SetScript(self.Script)
                    self.rnode.ConnectTree()
                }
            }
        } else if self.Item != nil {
            switch variable.Type {
            case types.Vector:
                switch self.Item.String() {
                case "x", "y", "z":
                default:
                    self.Script.AddError(self, self.At, "unknown field \""+self.String()+"\"")
                    self.isValid = false
                }
            case types.Rotation:
                switch self.Item.String() {
                case "x", "y", "z", "s":
                default:
                    self.Script.AddError(self, self.At, "unknown field \""+self.String()+"\"")
                    self.isValid = false
                }
            default:
                self.Script.AddError(self, self.At, "unknown field \""+self.String()+"\"")
                self.isValid = false
            }
        }
    }
}

func (self *LValueExpression) ValueType() types.Type {
	if !self.isValid {
		return types.Unknown
	}

	if self.Item != nil {
        s := self.GetStruct()
        if s != nil {
            f := s.GetField(self.Item.String())
            if f == nil {
                return types.Unknown
            } else {
                return f.ValueType()
            }
        }

		return types.Float
	}

	variable := self.Scope.FindVariable(self.Name.String())
	return variable.Type
}

func (self *LValueExpression) GetChildren() []Node {
    result := []Node{ self.Name }
    if self.Item != nil {
        result = append(result, self.Item)
    }
    if self.rnode != self {
        result = append(result, self.rnode)
    }
    return result
}

func (self *LValueExpression) RealNode() Node {
    return self.rnode
}
