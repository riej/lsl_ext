package nodes

import (
    "fmt"
)

type BreakStatement struct {
    NodeBase

    Breakable BreakableNode
}

func (self *BreakStatement) NodeType() NodeType {
    return NodeStatement
}

func (self *BreakStatement) StatementType() StatementType {
    return StatementBreak
}

func (self *BreakStatement) String() string {
    if self.Breakable == nil {
        return ""
    }

    if self.Breakable.StatementType() == StatementSwitch && !self.Script.LegacySwitch {
        return ""
    }

    return fmt.Sprintf("jump %s; // break", self.Breakable.BreakLabel())
}

func (self *BreakStatement) ConnectTree() {
    self.isValid = true

    node := self.Parent
    for node != nil {
        switch node.StatementType() {
        case StatementSwitch:
            node.(*SwitchStatement).HasBreak = true
            self.Breakable = node.(*SwitchStatement)
            return
        case StatementFor:
            node.(*ForStatement).HasBreak = true
            self.Breakable = node.(*ForStatement)
            return
        case StatementWhile:
            node.(*WhileStatement).HasBreak = true
            self.Breakable = node.(*WhileStatement)
            return
        case StatementDo:
            node.(*DoStatement).HasBreak = true
            self.Breakable = node.(*DoStatement)
            return
        }

        node = node.GetParent()
    }

    self.Script.AddError(self, self.At, "break statement outside of loop/switch")
    self.isValid = false
}

func (self *BreakStatement) GetChildren() []Node {
    return []Node{}
}

func (self *BreakStatement) RealNode() Node {
    return self
}
