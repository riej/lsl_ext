package nodes

import (
    "fmt"

	"../types"
)

type Function struct {
	NodeBase

	Name      *Identifier
	Type      types.Type
	Arguments []*Variable
	Body      *BlockStatement

	IsStateEvent bool


    IsUsed bool
}

func (self *Function) NodeType() NodeType {
	return NodeFunction
}

func (self *Function) String() string {
	result := ""
	if self.Type != types.Void {
		result += self.Type.String() + " "
	}
	result += self.Name.String() + "("
	for i, arg := range self.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	if self.Body != nil {
		result += " " + self.Body.String() + "\n"
	}

	return result
}

func (self *Function) ConnectTree() {
	self.Name.SetParent(self)
	self.Name.SetScope(self.Scope)
	self.Name.SetScript(self.Script)
	self.Name.ConnectTree()

	for _, child := range self.Arguments {
		child.SetParent(self)
		child.SetScope(self.Scope)
		child.SetScript(self.Script)
		child.ConnectTree()
	}

	currScope := self.Scope.Clone(self.At)
	currScope.Variables = append(self.Arguments, currScope.Variables...)

	self.Body.SetParent(self)
	self.Body.SetIndentationLevel(self.IndentationLevel)
	self.Body.SetScope(currScope)
	self.Body.SetScript(self.Script)
	self.Body.ConnectTree()

	if self.Parent == self.Script {
        if existing, _ := self.Script.Functions[self.Name.String()]; existing != self {
            self.Script.AddError(self, self.At, fmt.Sprintf("redeclared function \"%s\" (previously declared at %s)", self.Name, existing.Position()))
        }
	}
}

func (self *Function) ValueType() types.Type {
	return self.Type
}

func (self *Function) GetChildren() []Node {
    result := []Node{}
    result = append(result, self.Name)
    for _, child := range self.Arguments {
        result = append(result, child)
    }
    if self.Body != nil {
        result = append(result, self.Body)
    }
    return result
}

func (self *Function) RealNode() Node {
    return self
}
