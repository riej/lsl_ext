package nodes

import (
    "fmt"
)

type Struct struct {
    NodeBase

    Name *Identifier

    Fields []*Variable
}

func (self *Struct) NodeType() NodeType {
    return NodeStruct
}

func (self *Struct) String() string {
    return ""
}

func (self *Struct) ConnectTree() {
    self.isValid = true
    names := make(map[string]*Variable)

    for _, child := range self.Fields {
        child.SetParent(self)
        child.SetScope(self.Scope)
        child.SetScript(self.Script)
        child.ConnectTree()

        if existing, ok := names[child.Name.String()]; ok {
            self.isValid = false
            self.Script.AddError(child, child.At, fmt.Sprintf("redeclared struct field (previous declared at %s)", existing.Position()))
        } else {
            names[child.Name.String()] = child
        }

        if child.ValueType().IsList() {
            self.isValid = false
            self.Script.AddError(child, child.At, "structs cannot contain fields with type of list or struct")
        }
    }
}

func (self *Struct) Children() []Node {
    result := make([]Node, 0)
    for _, child := range self.Fields {
        result = append(result, child)
    }
    return result
}

func (self *Struct) GetField(name string) *Variable {
    for _, child := range self.Fields {
        if child.Name.String() == name {
            return child
        }
    }

    return nil
}

func (self *Struct) GetFieldIndex(name string) int32 {
    for i, child := range self.Fields {
        if child.Name.String() == name {
            return int32(i)
        }
    }

    return -1
}

func (self *Struct) RealNode() Node {
    return self
}
