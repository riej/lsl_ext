package nodes

import (
    "fmt"

    "../types"
)

type StructExpression struct {
    NodeBase

    Name *Identifier
    Fields []*Variable
}

func (self *StructExpression) NodeType() NodeType {
    return NodeExpression
}

func (self *StructExpression) ExpressionType() ExpressionType {
    return ExpressionStruct
}

func (self *StructExpression) String() string {
    s, ok := self.Script.Structs[self.Name.String()]
    if !ok {
        return "[]"
    }

    result := "["
    for i, field := range s.Fields {
        if i != 0 {
            result += ", "
        }

        found := false
        for _, child := range self.Fields {
            if child.Name.String() == field.Name.String() {
                result += child.ValueString()

                found = true
                break
            }
        }

        if !found {
            result += field.ValueString()
        }
    }

    result += "]"

    return result
}

func (self *StructExpression) ConnectTree() {
    self.isValid = true

    if self.Name == nil {
        if self.Parent.NodeType() == NodeVariable {
            if self.Parent.(*Variable).Type.IsStruct() {
                self.Name = &Identifier{
                    Name: string(self.Parent.(*Variable).Type),
                }
            } else {
                self.isValid = false
                self.Script.AddError(self, self.At, "invalid struct expression")
                return
            }
        } else {
            self.isValid = false
            self.Script.AddError(self, self.At, "invalid struct expression")
            return
        }
    }

    self.Name.SetParent(self)
    self.Name.SetScope(self.Scope)
    self.Name.SetScript(self.Script)
    self.Name.ConnectTree()

    s, ok := self.Script.Structs[self.Name.String()]
    if !ok {
        self.isValid = false
        self.Script.AddError(self, self.At, fmt.Sprintf("undeclared struct \"%s\"", self.Name))
    }

    names := make(map[string]*Variable)
    for _, child := range self.Fields {
        child.SetParent(self)
        child.SetScope(self.Scope)
        child.SetScript(self.Script)

        if s != nil {
            f := s.GetField(child.Name.String())
            if f == nil {
                self.isValid = false
                self.Script.AddError(child, child.Position(), fmt.Sprintf("undeclared struct field \"%s\"", child.Name))
            } else {
                child.Type = f.Type
            }
        }

        if existing, ok := names[child.Name.String()]; ok {
            self.isValid = false
            self.Script.AddError(child, child.Position(), fmt.Sprintf("redeclared value \"%s\" (previously declared at %s)", child.Name, existing.Position()))
        }

        child.ConnectTree()

        names[child.Name.String()] = child
    }
}

func (self *StructExpression) ValueType() types.Type {
    if self.Name == nil {
        return types.Unknown
    }

    return types.Type(self.Name.String())
}

func (self *StructExpression) Children() []Node {
    result := make([]Node, 0)

    if self.Name != nil {
        result = append(result, self.Name)
    }

    for _, child := range self.Fields {
        result = append(result, child)
    }
    return result
}

func (self *StructExpression) RealNode() Node {
    return self
}
