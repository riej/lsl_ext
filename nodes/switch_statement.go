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
    self.Expression.SetParent(self)
    self.Expression.SetScope(self.Scope)
    self.Expression.SetScript(self.Script)
    self.Expression.ConnectTree()

    var defaultCase *CaseStatement

    for _, child := range self.Block.Children {
        switch child.StatementType() {
        case StatementCase:
            child.(*CaseStatement).Switch = self

            if len(child.(*CaseStatement).Expressions) == 0 {
                defaultCase = child.(*CaseStatement)
            } else {
                self.Cases = append(self.Cases, child.(*CaseStatement))
            }
        }
    }

    if defaultCase != nil {
        self.Cases = append(self.Cases, defaultCase)
    }

    self.Block.SetParent(self)
    self.Block.SetScope(self.Scope)
    self.Block.SetScript(self.Script)
    self.Block.SetIndentationLevel(self.IndentationLevel)
    self.Block.NoBraces = true
    self.Block.ConnectTree()
}

func (self *SwitchStatement) GetChildren() []Node {
    return []Node{ self.Expression, self.Block }
}

func (self *SwitchStatement) RealNode() Node {
    return self
}
