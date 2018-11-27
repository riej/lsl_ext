package nodes

import (
	"strings"
)

type IfStatement struct {
	NodeBase

	If   Expression
	Then Statement
	Else Statement
}

func (self *IfStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *IfStatement) StatementType() StatementType {
	return StatementIf
}

func (self *IfStatement) String() string {
	result := "if (" + self.If.String() + ") " + self.Then.String()
	if self.Else != nil {
		if self.Then.StatementType() == StatementBlock {
			result += " "
		} else {
			result += "\n" + strings.Repeat(Indentation, self.IndentationLevel)
		}
		result += "else " + self.Else.String()
	}

	return result
}

func (self *IfStatement) ConnectTree() {
	self.If.SetParent(self)
	self.If.SetScope(self.Scope)
	self.If.SetScript(self.Script)
	self.If.ConnectTree()

	if self.Then.StatementType() != StatementBlock {
        old := self.Then
		self.Then = &BlockStatement{
			Children: []Statement{
				self.Then,
			},
		}
        self.Then.SetPosition(old.Position())
	}

	self.Then.SetParent(self)
	self.Then.SetIndentationLevel(self.IndentationLevel)
	self.Then.SetScope(self.Scope)
	self.Then.SetScript(self.Script)
	self.Then.ConnectTree()

	if self.Else != nil {
		if self.Else.StatementType() != StatementBlock && self.Else.StatementType() != StatementIf {
            old := self.Else
			self.Else = &BlockStatement{
				Children: []Statement{
					self.Else,
				},
			}
            self.Else.SetPosition(old.Position())
		}

		self.Else.SetParent(self)
		self.Else.SetIndentationLevel(self.IndentationLevel)
		self.Else.SetScope(self.Scope)
		self.Else.SetScript(self.Script)
		self.Else.ConnectTree()
	}
}

func (self *IfStatement) GetChildren() []Node {
    result := []Node{ self.If, self.Then }
    if self.Else != nil {
        result = append(result, self.Else)
    }
    return result
}

func (self *IfStatement) RealNode() Node {
    return self
}
