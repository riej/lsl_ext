package nodes

import (
    "fmt"
    "strings"
)

type SwitchStatement struct {
    NodeBase

    Expression Expression
    Block *BlockStatement

    HasBreak bool
}

func (self *SwitchStatement) NodeType() NodeType {
    return NodeStatement
}

func (self *SwitchStatement) StatementType() StatementType {
    return StatementSwitch
}

func (self *SwitchStatement) String() string {
    result := fmt.Sprintf("// switch(%s) %s", self.Expression, strings.Repeat(Indentation, self.IndentationLevel))
    result += self.Block.String()
    if self.HasBreak {
        result += fmt.Sprintf("@%s;", self.BreakLabel())
    }
    return result
}

func (self *SwitchStatement) BreakLabel() string {
    return fmt.Sprintf("switch_end_%d", self.At.Offset)
}

func (self *SwitchStatement) ConnectTree() {
    self.Expression.SetParent(self)
    self.Expression.SetScope(self.Scope)
    self.Expression.SetScript(self.Script)
    self.Expression.ConnectTree()

    for _, child := range self.Block.Children {
        switch child.StatementType() {
        case StatementCase:
            child.(*CaseStatement).Switch = self
        }
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
