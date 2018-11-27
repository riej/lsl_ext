package nodes

import (
	"fmt"
    "strings"
)

type DoStatement struct {
	NodeBase

	Body      Statement
	Condition Expression

    HasBreak bool
    HasContinue bool
}

func (self *DoStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *DoStatement) StatementType() StatementType {
	return StatementDo
}

func (self *DoStatement) String() string {
    switch self.Body.StatementType() {
    case StatementBlock:
        self.Body.(*BlockStatement).NoBraces = true

        result := "do {" + self.Body.String()

        if self.HasContinue {
            result += fmt.Sprintf("\n%s@%s;\n%s", strings.Repeat(Indentation, self.IndentationLevel + 1), self.ContinueLabel(), strings.Repeat(Indentation, self.IndentationLevel));
        }

        result += "}"

        if self.HasBreak {
            result += fmt.Sprintf("\n%s@%s;\n", strings.Repeat(Indentation, self.IndentationLevel), self.BreakLabel());
        }
        return result
    case StatementBreak:
        return fmt.Sprintf("// do break; while (%s);", self.Body, self.Condition)
    case StatementContinue:
        return fmt.Sprintf("do {} while (%s);", self.Body, self.Condition)
    default:
        return fmt.Sprintf("do %s while (%s);", self.Body, self.Condition)
    }
}

func (self *DoStatement) ContinueLabel() string {
    return fmt.Sprintf("do_body_end_%d", self.At.Offset)
}

func (self *DoStatement) BreakLabel() string {
    return fmt.Sprintf("do_end_%d", self.At.Offset)
}

func (self *DoStatement) ConnectTree() {
	self.Body.SetParent(self)
	self.Body.SetIndentationLevel(self.IndentationLevel)
	self.Body.SetScope(self.Scope)
	self.Body.SetScript(self.Script)
	self.Body.ConnectTree()

	self.Condition.SetParent(self)
	self.Condition.SetScope(self.Scope)
	self.Condition.SetScript(self.Script)
	self.Condition.ConnectTree()
}

func (self *DoStatement) GetChildren() []Node {
    return []Node{ self.Body, self.Condition }
}

func (self *DoStatement) RealNode() Node {
    return self
}
