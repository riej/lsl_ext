package nodes

import (
    "fmt"
    "strings"
)

type SwitchStatement struct {
    NodeBase

    Expression Expression
    Block *BlockStatement

    Cases []*CaseStatement

    HasBreak bool
}

func (self *SwitchStatement) NodeType() NodeType {
    return NodeStatement
}

func (self *SwitchStatement) StatementType() StatementType {
    return StatementSwitch
}

func (self *SwitchStatement) String() string {
    if self.Script.LegacySwitch {
        result := fmt.Sprintf("// switch(%s)\n", self.Expression)

        for i, child := range self.Cases {
            result += strings.Repeat(Indentation, self.IndentationLevel)
            if i > 0 {
                result += "else "
            }
            result += child.IfString() + "\n"
        }

        result += self.Block.String()
        if self.HasBreak {
            result += fmt.Sprintf("\n%s@%s;", strings.Repeat(Indentation, self.IndentationLevel), self.BreakLabel())
        }
        return result
    } else {
        result := fmt.Sprintf("// switch(%s)\n%s", self.Expression, strings.Repeat(Indentation, self.IndentationLevel))

        for i, child := range self.Cases {
            if i > 0 {
                result += " else "
            }
            result += child.String()
        }

        return result
    }
}

func (self *SwitchStatement) BreakLabel() string {
    return fmt.Sprintf("switch_end_%d", self.At.Offset)
}

func (self *SwitchStatement) NextCaseLabel(curr *CaseStatement) string {
    var prev *CaseStatement
    for _, child := range self.Cases {
        if prev == curr {
            return child.Label()
        }

        prev = child
    }

    return self.BreakLabel()
}

func (self *SwitchStatement) ConnectTree() {
    self.isValid = true

    self.Expression.SetParent(self)
    self.Expression.SetScope(self.Scope)
    self.Expression.SetScript(self.Script)
    self.Expression.ConnectTree()

    var defaultCase *CaseStatement
    var currCaseBody *BlockStatement

    for i, child := range self.Block.Children {
        switch child.StatementType() {
        case StatementCase:
            child.(*CaseStatement).Switch = self
            child.(*CaseStatement).Body = &BlockStatement{
                NoBraces: false,
            }
            currCaseBody = child.(*CaseStatement).Body
            currCaseBody.At = child.Position()

            if len(child.(*CaseStatement).Expressions) == 0 {
                defaultCase = child.(*CaseStatement)
            } else {
                self.Cases = append(self.Cases, child.(*CaseStatement))
            }
        default:
            if i == 0 {
                self.Script.AddError(self, child.Position(), "switch statement must start from case statement")
                self.isValid = false
            }

            if currCaseBody != nil {
                currCaseBody.Children = append(currCaseBody.Children, child)
            }
        }
    }

    if defaultCase != nil {
        self.Cases = append(self.Cases, defaultCase)
    }

    if self.Script.LegacySwitch {
        self.Block.SetParent(self)
        self.Block.SetScope(self.Scope)
        self.Block.SetScript(self.Script)
        self.Block.SetIndentationLevel(self.IndentationLevel)
        self.Block.NoBraces = true
        self.Block.ConnectTree()
    } else {
        for _, child := range self.Cases {
            child.SetParent(self)
            child.SetScope(self.Scope)
            child.SetScript(self.Script)
            child.SetIndentationLevel(self.IndentationLevel)
            child.ConnectTree()
        }
    }
}

func (self *SwitchStatement) GetChildren() []Node {
    if self.Script.LegacySwitch {
        return []Node{ self.Expression, self.Block }
    } else {
        result := []Node{}
        for _, child := range self.Cases {
            result = append(result, child)
        }
        return result
    }
}

func (self *SwitchStatement) RealNode() Node {
    return self
}
