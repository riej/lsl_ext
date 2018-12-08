package nodes

import (
	"fmt"
    "strings"
)

type WhileStatement struct {
	NodeBase

	Condition Expression
	Body      Statement

    HasBreak bool
    HasContinue bool
}

func (self *WhileStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *WhileStatement) StatementType() StatementType {
	return StatementWhile
}

func (self *WhileStatement) String() string {
    result := "while (" + self.Condition.String() + ")"
    switch self.Body.StatementType() {
    case StatementBlock:
        self.Body.(*BlockStatement).NoBraces = true

        result += "{"
        result += self.Body.String()

        if self.HasContinue {
            result += fmt.Sprintf("\n%s@%s;\n%s", strings.Repeat(Indentation, self.IndentationLevel + 1), self.ContinueLabel(), strings.Repeat(Indentation, self.IndentationLevel))
        }

        result += "}"

        if self.HasBreak {
            result += fmt.Sprintf("\n%s@%s;\n", strings.Repeat(Indentation, self.IndentationLevel), self.BreakLabel())
        }
    case StatementBreak:
        result = "// " + result + " break;"
    case StatementContinue:
        result += ";"
    default:
        result += fmt.Sprintf("\n%s%s\n", strings.Repeat(Indentation, self.IndentationLevel + 1), self.Body)
    }

    return result
}

func (self *WhileStatement) ContinueLabel() string {
    return fmt.Sprintf("while_body_end_%d", self.At.Offset)
}

func (self *WhileStatement) BreakLabel() string {
    return fmt.Sprintf("while_end_%d", self.At.Offset)
}

func (self *WhileStatement) ConnectTree() {
	self.Condition.SetParent(self)
	self.Condition.SetScope(self.Scope)
	self.Condition.SetScript(self.Script)
	self.Condition.ConnectTree()

	self.Body.SetParent(self)
	self.Body.SetIndentationLevel(self.IndentationLevel)
	self.Body.SetScope(self.Scope)
	self.Body.SetScript(self.Script)
	self.Body.ConnectTree()
}

func (self *WhileStatement) GetChildren() []Node {
    return []Node{ self.Condition, self.Body }
}

func (self *WhileStatement) RealNode() Node {
    return self
}
