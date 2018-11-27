package nodes

import (
    "fmt"
    "strings"
)

type CaseStatement struct {
    NodeBase

    Switch *SwitchStatement // set from top level

    Expressions []Expression
}

func (self *CaseStatement) NodeType() NodeType {
    return NodeStatement
}

func (self *CaseStatement) StatementType() StatementType {
    return StatementCase
}

func (self *CaseStatement) String() string {
    result := ""

    if len(self.Expressions) == 0 {
        result += "// default:"
    } else {
        expr := ""
        endLabel := ""
        if self.Switch != nil {
            expr = self.Switch.Expression.String()
            endLabel = self.Switch.BreakLabel()

            if self.Switch.Expression.ExpressionType() == ExpressionBinary {
                expr = "(" + expr + ")"
            }
        }

        result += "// case "
        for i, child := range self.Expressions {
            if i != 0 {
                result += ", "
            }

            result += child.String()
        }
        result += ":\n" + strings.Repeat(Indentation, self.IndentationLevel)
        result += "if ("

        if len(self.Expressions) == 1 {
            if self.Expressions[0].ExpressionType() == ExpressionBinary {
                result += fmt.Sprintf("%s != (%s)", expr, self.Expressions[0])
            } else {
                result += fmt.Sprintf("%s != %s", expr, self.Expressions[0])
            }
        } else {
            for i, child := range self.Expressions {
                if i != 0 {
                    result += " && "
                }

                if child.ExpressionType() == ExpressionBinary {
                    result += fmt.Sprintf("(%s != (%s))", expr, child)
                } else {
                    result += fmt.Sprintf("(%s != %s)", expr, child)
                }
            }
        }

        result += fmt.Sprintf(") jump %s;\n", endLabel)
    }

    return result
}

func (self *CaseStatement) ConnectTree() {
    self.isValid = true

    for _, child := range self.Expressions {
        child.SetParent(self)
        child.SetScope(self.Scope)
        child.SetScript(self.Script)
        child.ConnectTree()

        self.isValid = self.isValid || child.IsValid()
    }

    if self.Switch == nil {
        self.Script.AddError(self, self.At, "case statement outside of switch")
        self.isValid = false
    }
}

func (self *CaseStatement) GetChildren() []Node {
    result := []Node{}
    for _, child := range self.Expressions {
        result = append(result, child)
    }
    return result
}

func (self *CaseStatement) RealNode() Node {
    return self
}
