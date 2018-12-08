package nodes

import (
    "fmt"
    "strings"
)

type ForStatement struct {
	NodeBase

	Init      []Expression
	Condition Expression
	Loop      []Expression

    HasBreak bool
    HasContinue bool

	Body Statement
}

func (self *ForStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *ForStatement) StatementType() StatementType {
	return StatementFor
}

func (self *ForStatement) String() string {
	result := "for ("
	for i, expr := range self.Init {
		if i > 0 {
			result += ", "
		}
		result += expr.String()
	}

	if self.Condition == nil {
		result += "; TRUE;"
	} else {
		result += "; " + self.Condition.String() + ";"
	}

	for i, expr := range self.Loop {
		if i == 0 {
			result += " "
		} else {
			result += ", "
		}
		result += expr.String()
	}

	result += ") "

    switch self.Body.StatementType() {
    case StatementBlock:
        self.Body.(*BlockStatement).NoBraces = true

        result += "{"
        result += self.Body.String()

        if self.HasContinue {
            result += fmt.Sprintf("\n%s@%s;\n%s", strings.Repeat(Indentation, self.IndentationLevel + 1), self.ContinueLabel(), strings.Repeat(Indentation, self.IndentationLevel));
        }

        result += "}"

        if self.HasBreak {
            result += fmt.Sprintf("\n%s@%s;\n", strings.Repeat(Indentation, self.IndentationLevel), self.BreakLabel());
        }
    case StatementBreak:
        result = "// " + result + "break;"
    case StatementContinue:
        result += ";"
    default:
        result += fmt.Sprintf("\n%s%s\n", strings.Repeat(Indentation, self.IndentationLevel + 1), self.Body)
    }

	return result
}

func (self *ForStatement) ContinueLabel() string {
    return fmt.Sprintf("for_body_end_%d", self.At.Offset)
}

func (self *ForStatement) BreakLabel() string {
    return fmt.Sprintf("for_end_%d", self.At.Offset)
}

func (self *ForStatement) ConnectTree() {
	for _, child := range self.Init {
		child.SetParent(self)
		child.SetScope(self.Scope)
		child.SetScript(self.Script)
		child.ConnectTree()
	}

    if self.Condition != nil {
        self.Condition.SetParent(self)
        self.Condition.SetScope(self.Scope)
        self.Condition.SetScript(self.Script)
        self.Condition.ConnectTree()
    }

	for _, child := range self.Loop {
		child.SetParent(self)
		child.SetScope(self.Scope)
		child.SetScript(self.Script)
		child.ConnectTree()
	}

	self.Body.SetParent(self)
	self.Body.SetIndentationLevel(self.IndentationLevel)
	self.Body.SetScope(self.Scope)
	self.Body.SetScript(self.Script)
	self.Body.ConnectTree()
}

func (self *ForStatement) GetChildren() []Node {
    result := []Node{}
    for _, child := range self.Init {
        result = append(result, child)
    }
    if self.Condition != nil {
        result = append(result, self.Condition)
    }
    for _, child := range self.Loop {
        result = append(result, child)
    }
    result = append(result, self.Body)
    return result
}

func (self *ForStatement) RealNode() Node {
    return self
}
