package nodes

import (
	"fmt"
//    "strings"

	"../types"
)

type BinaryExpression struct {
	NodeBase

	Operator string
	LValue   Expression
	RValue   Expression
}

func (self *BinaryExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *BinaryExpression) ExpressionType() ExpressionType {
	return ExpressionBinary
}

func (self *BinaryExpression) String() string {
    if self.LValue.ExpressionType() == ExpressionListItem {
        li := self.LValue.(*ListItemExpression)
        liltype := li.LValue.ValueType()

        switch self.Operator {
        case "=":
            if liltype == types.String {
                if self.RValue.ValueType() == types.String {
                    return fmt.Sprintf("%s = llInsertString(llDeleteSubString(%s, %s, %s), %s, %s)", li.LValue, li.LValue, li.StartIndex, li.EndIndex, li.StartIndex, self.RValue)
                } else {
                    return fmt.Sprintf("%s = llInsertString(llDeleteSubString(%s, %s, %s), %s, (string)%s)", li.LValue, li.LValue, li.StartIndex, li.EndIndex, li.StartIndex, self.RValue)
                }
/*            } else if strings.HasSuffix(string(liltype), "[]") {
                sname := strings.TrimSuffix(string(liltype), "[]")
                if s, ok := self.Script.Structs[sname]; ok {
                    stride := len(s.Fields)

                    result := fmt.Sprintf("%s = llListReplaceList(%s, ", li.LValue, li.LValue)
                    if self.RValue.ValueType().IsList() {
                        result += self.RValue.String()
                    } else {
                        result += fmt.Sprintf("[%s]", self.RValue)
                    }
                    result += fmt.Sprintf(", %d * ", stride)
                    if li.StartIndex.ExpressionType() == ExpressionBinary {
                        result += fmt.Sprintf("(%s)", li.StartIndex)
                    } else {
                        result += li.StartIndex.String()
                    }
                    result += fmt.Sprintf(", %d * ", stride)
                    if li.EndIndex.ExpressionType() == ExpressionBinary {
                        result += fmt.Sprintf("(%s)", li.EndIndex)
                    } else {
                        result += li.EndIndex.String()
                    }
                    result += fmt.Sprintf(" + %d)", stride - 1)
                    return result
                }*/
            } else if self.RValue.ValueType().IsList() {
                return fmt.Sprintf("%s = llListReplaceList(%s, %s, %s, %s)", li.LValue, li.LValue, self.RValue, li.StartIndex, li.EndIndex)
            } else {
                return fmt.Sprintf("%s = llListReplaceList(%s, [%s], %s, %s)", li.LValue, li.LValue, self.RValue, li.StartIndex, li.EndIndex)
            }
        case "+=":
            return fmt.Sprintf("%s = llListReplaceList(%s, [%s + %s], %s, %s)", li.LValue, li.LValue, li, self.RValue, li.StartIndex, li.EndIndex)
        case "-=":
            return fmt.Sprintf("%s = llListReplaceList(%s, [%s - %s], %s, %s)", li.LValue, li.LValue, li, self.RValue, li.StartIndex, li.EndIndex)
        case "*=":
            return fmt.Sprintf("%s = llListReplaceList(%s, [%s * %s], %s, %s)", li.LValue, li.LValue, li, self.RValue, li.StartIndex, li.EndIndex)
        case "/=":
            return fmt.Sprintf("%s = llListReplaceList(%s, [%s / %s], %s, %s)", li.LValue, li.LValue, li, self.RValue, li.StartIndex, li.EndIndex)
        case "%=":
            return fmt.Sprintf("%s = llListReplaceList(%s, [%s %% %s], %s, %s)", li.LValue, li.LValue, li, self.RValue, li.StartIndex, li.EndIndex)
        }
    }


	return fmt.Sprintf("%s %s %s", self.LValue, self.Operator, self.RValue)
}

func (self *BinaryExpression) ConnectTree() {
	self.LValue.SetParent(self)
    self.LValue.SetIndentationLevel(self.IndentationLevel)
	self.LValue.SetScope(self.Scope)
	self.LValue.SetScript(self.Script)
	self.LValue.ConnectTree()

	self.RValue.SetParent(self)
    self.RValue.SetIndentationLevel(self.IndentationLevel)
	self.RValue.SetScope(self.Scope)
	self.RValue.SetScript(self.Script)
	self.RValue.ConnectTree()

	self.isValid = true


    self.LValue = Expression(self.LValue.RealNode())
    self.RValue = Expression(self.RValue.RealNode())

    if self.LValue.ExpressionType() == ExpressionListItem {
        if self.RValue.ExpressionType() == ExpressionListItem {
            // TODO
            self.LValue.(*ListItemExpression).Type = self.RValue.ValueType()
        } else {
            self.LValue.(*ListItemExpression).Type = self.RValue.ValueType()
        }
    } else if self.RValue.ExpressionType() == ExpressionListItem {
        self.RValue.(*ListItemExpression).Type = self.LValue.ValueType()
    }

    if self.LValue.ExpressionType() == ExpressionListItem {
        li := self.LValue.(*ListItemExpression)

        switch self.Operator {
        case "+=", "-=", "*=", "/=", "%=":
            if li.IsRange || !li.LValue.ValueType().IsList() || self.RValue.ValueType().IsList() {
                self.Script.AddError(self, self.At, "invalid list item assignment")
                self.isValid = false
            } else if li.Type == types.Unknown {
                if self.Operator == "%=" {
                    li.Type = types.Integer
                } else {
                    li.Type = types.Float
                }
            }
        }
    }

	if self.LValue.IsValid() && self.RValue.IsValid() {
        ltype := self.LValue.ValueType()
        rtype := self.RValue.ValueType()
		if !ltype.IsCompatible(rtype) {
			self.Script.AddError(self, self.At, "type mismatch (" + ltype.String() + " and " + rtype.String() + ")")
			self.isValid = false
		}
	}

    switch self.Operator {
    case "=", "+=", "-=", "*=", "/=", "%=":
        switch self.LValue.ExpressionType() {
        case ExpressionLValue:
            variable := self.LValue.(*LValueExpression).GetVariable()
            if variable != nil && variable.IsConstant {
                self.Script.AddError(self, self.At, "cannot assign constant value")
                self.isValid = false
            }
        case ExpressionListItem:
        default:
            self.Script.AddError(self, self.At, "invalid assignment")
            self.isValid = false
        }
    }
}

func (self *BinaryExpression) ValueType() types.Type {
    switch self.Operator {
    case "==", "!=", ">", "<", ">=", "<=":
        return types.Boolean
    }

	ltype := self.LValue.ValueType()
	rtype := self.RValue.ValueType()

	if ltype == types.Float || rtype == types.Float {
		return types.Float
	}

    if ltype == rtype && ltype.IsStruct() {
        return ltype
    }

	if ltype.IsList() || rtype.IsList() {
		return types.List
	}

	return ltype
}

func (self *BinaryExpression) GetChildren() []Node {
    return []Node{ self.LValue, self.RValue }
}

func (self *BinaryExpression) RealNode() Node {
    return self
}
