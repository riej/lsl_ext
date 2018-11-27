package nodes

import (
    "fmt"
)

type ContinueStatement struct {
    NodeBase

    Continueable ContinueableNode
}

func (self *ContinueStatement) NodeType() NodeType {
    return NodeStatement
}

func (self *ContinueStatement) StatementType() StatementType {
    return StatementContinue
}

func (self *ContinueStatement) String() string {
    if self.Continueable == nil {
        return ""
    }

    return fmt.Sprintf("jump %s; // continue", self.Continueable.ContinueLabel())
}

func (self *ContinueStatement) ConnectTree() {
    self.isValid = true

    node := self.Parent
    for node != nil {
        switch node.StatementType() {
        case StatementFor:
            node.(*ForStatement).HasContinue = true
            self.Continueable = node.(*ForStatement)
            return
        case StatementWhile:
            node.(*WhileStatement).HasContinue = true
            self.Continueable = node.(*WhileStatement)
            return
        case StatementDo:
            node.(*DoStatement).HasContinue = true
            self.Continueable = node.(*DoStatement)
            return
        }

        node = node.GetParent()
    }

    self.Script.AddError(self, self.At, "continue statement outside of loop/switch")
    self.isValid = false
}

func (self *ContinueStatement) GetChildren() []Node {
    return []Node{}
}

func (self *ContinueStatement) RealNode() Node {
    return self
}
