package nodes

import (
    "fmt"
    "strings"

    "../types"
)

type ListItemExpression struct {
    NodeBase

    LValue Node
    Type types.Type

    IsRange bool

    StartIndex Expression
    EndIndex Expression

    Item *Identifier // q[5].name
}

func (self *ListItemExpression) NodeType() NodeType {
    return NodeExpression
}

func (self *ListItemExpression) ExpressionType() ExpressionType {
    return ExpressionListItem
}

func (self *ListItemExpression) String() string {
    if self.LValue.ValueType() == types.String {
        return fmt.Sprintf("llGetSubString(%s, %s, %s)", self.LValue, self.StartIndex, self.EndIndex)
    } else if self.IsRange {
        return fmt.Sprintf("llList2List(%s, %s, %s)", self.LValue, self.StartIndex, self.EndIndex)
    }

    switch self.Type {
    case types.Integer:
        return fmt.Sprintf("llList2Integer(%s, %s)", self.LValue, self.StartIndex)
    case types.Float:
        return fmt.Sprintf("llList2Float(%s, %s)", self.LValue, self.StartIndex)
    case types.String:
        return fmt.Sprintf("llList2String(%s, %s)", self.LValue, self.StartIndex)
    case types.Key:
        return fmt.Sprintf("llList2Key(%s, %s)", self.LValue, self.StartIndex)
    case types.Vector:
        return fmt.Sprintf("(vector)llList2String(%s, %s)", self.LValue, self.StartIndex)
    case types.Rotation:
        return fmt.Sprintf("(rotation)llList2String(%s, %s)", self.LValue, self.StartIndex)
    default:
/*        if strings.HasSuffix(string(self.LValue.ValueType()), "[]") {
            sname := strings.TrimSuffix(string(self.LValue.ValueType()), "[]")
            if s, ok := self.Script.Structs[sname]; ok {
                stride := len(s.Fields)

                result := fmt.Sprintf("llList2List(%s, %d * ", self.LValue, stride)
                if self.StartIndex.ExpressionType() == ExpressionBinary {
                    result += fmt.Sprintf("(%s)", self.StartIndex)
                } else {
                    result += self.StartIndex.String()
                }
                result += fmt.Sprintf(", %d * ", stride)
                if self.EndIndex.ExpressionType() == ExpressionBinary {
                    result += fmt.Sprintf("(%s)", self.EndIndex)
                } else {
                    result += self.EndIndex.String()
                }
                result += fmt.Sprintf(" + %d)", stride - 1)
                return result
            }
        } else */if self.Type.IsList() {
            return fmt.Sprintf("llList2List(%s, %s, %s)", self.LValue, self.StartIndex, self.EndIndex)
        }

        return fmt.Sprintf("llList2String(%s, %s)", self.LValue, self.StartIndex)
    }
}

func (self *ListItemExpression) ConnectTree() {
    self.LValue.SetParent(self)
    self.LValue.SetScope(self.Scope)
    self.LValue.SetScript(self.Script)
    self.LValue.ConnectTree()

    self.StartIndex.SetParent(self)
    self.StartIndex.SetScope(self.Scope)
    self.StartIndex.SetScript(self.Script)
    self.StartIndex.ConnectTree()

    if !self.IsRange || self.EndIndex == nil {
        self.EndIndex = self.StartIndex
    }

    if self.StartIndex != self.EndIndex {
        self.EndIndex.SetParent(self)
        self.EndIndex.SetScope(self.Scope)
        self.EndIndex.SetScript(self.Script)
        self.EndIndex.ConnectTree()
    }

    self.isValid = true

    if !self.StartIndex.ValueType().IsCompatible(types.Integer) {
        self.Script.AddError(self.StartIndex, self.StartIndex.Position(), "invalid array index")
        self.isValid = false
    } else if self.StartIndex != self.EndIndex && !self.EndIndex.ValueType().IsCompatible(types.Integer) {
        self.Script.AddError(self.EndIndex, self.EndIndex.Position(), "invalid array index")
        self.isValid = false
    }

    if self.LValue.IsValid() {
        ltype := self.LValue.ValueType()

        if strings.HasSuffix(string(ltype), "[]") {
            sname := strings.TrimSuffix(string(ltype), "[]")
            s, _ := self.Script.Structs[sname]
            if s == nil {
                self.Script.AddError(self, self.At, "undeclared struct \""+string(sname)+"\"")
                self.isValid = false
            } else {
                stride := int32(len(s.Fields))

                var field *Variable
                var index int32

                field = nil
                index = 0

                if self.Item != nil {
                    field = s.GetField(self.Item.String())
                    index = s.GetFieldIndex(self.Item.String())
                }

                if self.Item != nil && field == nil {
                    self.Script.AddError(self, self.At, "undeclared struct field \""+self.String()+"\"")
                    self.isValid = false
                } else {
                    oldSi := self.StartIndex

                    if self.StartIndex.ExpressionType() == ExpressionBinary {
                        self.StartIndex = &BracesExpression{
                            Child: self.StartIndex,
                        }
                        self.StartIndex.SetPosition(oldSi.Position())

                        self.StartIndex.SetParent(self)
                        self.StartIndex.SetScope(self.Scope)
                        self.StartIndex.SetScript(self.Script)
                        self.StartIndex.ConnectTree()
                    }

                    beMul := &BinaryExpression{
                        Operator: "*",
                        LValue: &Constant{
                            Value: &types.IntegerValue{ Value: stride },
                        },
                        RValue: self.StartIndex,
                    }
                    beMul.At = oldSi.Position()

                    if field == nil {
                        self.StartIndex = beMul
                    } else {
                        beAdd := &BinaryExpression{
                            Operator: "+",
                            LValue: beMul,
                            RValue: &Constant{
                                Value: &types.IntegerValue{ Value: index },
                            },
                        }
                        beAdd.At = oldSi.Position()

                        self.StartIndex = beAdd
                    }

                    self.StartIndex.SetParent(self)
                    self.StartIndex.SetScope(self.Scope)
                    self.StartIndex.SetScript(self.Script)
                    self.StartIndex.ConnectTree()




                    if field == nil {
                        oldEi := self.EndIndex

                        if self.EndIndex.ExpressionType() == ExpressionBinary {
                            self.EndIndex = &BracesExpression{
                                Child: self.EndIndex,
                            }
                            self.EndIndex.SetPosition(oldEi.Position())

                            self.EndIndex.SetParent(self)
                            self.EndIndex.SetScope(self.Scope)
                            self.EndIndex.SetScript(self.Script)
                            self.EndIndex.ConnectTree()
                        }

                        beMul := &BinaryExpression{
                            Operator: "*",
                            LValue: &Constant{
                                Value: &types.IntegerValue{ Value: stride },
                            },
                            RValue: self.EndIndex,
                        }
                        beMul.At = oldEi.Position()

                        beAdd := &BinaryExpression{
                            Operator: "+",
                            LValue: beMul,
                            RValue: &Constant{
                                Value: &types.IntegerValue{ Value: stride - 1 },
                            },
                        }
                        beAdd.At = oldEi.Position()

                        self.EndIndex = beAdd

                        self.EndIndex.SetParent(self)
                        self.EndIndex.SetScope(self.Scope)
                        self.EndIndex.SetScript(self.Script)
                        self.EndIndex.ConnectTree()
                    } else {
                        self.EndIndex = self.StartIndex
                    }
                }
            }
        } else if self.Item != nil {
            self.Script.AddError(self, self.At, "invalid list item access")
            self.isValid = false
        } else if !ltype.IsList() && ltype != types.String {
            self.Script.AddError(self, self.At, "invalid list item access (cannot access " + self.LValue.ValueType().String() + " items)")
            self.isValid = false
        }
    } else {
        self.isValid = false
    }
}

func (self *ListItemExpression) ValueType() types.Type {
    if self.LValue.ValueType() == types.String {
        return types.String
    } else if self.IsRange {
        return types.List
    }

    if self.Type == types.Unknown {
        return types.String
    }

    return self.Type
}

func (self *ListItemExpression) GetChildren() []Node {
    return []Node{ self.LValue, self.StartIndex, self.EndIndex }
}

func (self *ListItemExpression) RealNode() Node {
    return self
}
