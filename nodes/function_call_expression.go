package nodes

import (
    "fmt"
	"../types"
)

type FunctionCallExpression struct {
	NodeBase

	Name      *Identifier
	Arguments []Expression
}

func (self *FunctionCallExpression) NodeType() NodeType {
	return NodeExpression
}

// FunctionCallExpression can be treated at Typecast
func (self *FunctionCallExpression) ExpressionType() ExpressionType {
	return ExpressionFunctionCall
}

func (self *FunctionCallExpression) String() string {
	result := self.Name.String() + "("
	for i, arg := range self.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"

	return result
}

func (self *FunctionCallExpression) ConnectTree() {
	self.Name.SetParent(self)
	self.Name.SetScope(self.Scope)
	self.Name.SetScript(self.Script)
	self.Name.ConnectTree()

    realArgs := make([]Expression, len(self.Arguments))

	for i, child := range self.Arguments {
		child.SetParent(self)
		child.SetIndentationLevel(self.IndentationLevel)
		child.SetScope(self.Scope)
		child.SetScript(self.Script)
		child.ConnectTree()

        realArgs[i] = child.RealNode()
	}

    self.Arguments = realArgs



	self.isValid = true
	f, ok := self.Script.Functions[self.Name.String()]
	if !ok {
		self.Script.AddError(self, self.At, "undeclared function \""+self.Name.String()+"\"")
		self.isValid = false
	} else {
        f.IsUsed = true

        if len(self.Arguments) != len(f.Arguments) {
    		self.Script.AddError(self, self.At, fmt.Sprintf("invalid arguments count (expected %d, found %d)", len(f.Arguments), len(self.Arguments)))
            self.isValid = false
        } else {
            for i, child := range self.Arguments {
                if child.IsValid() {
                    if child.ExpressionType() == ExpressionListItem {
                        child.(*ListItemExpression).Type = f.Arguments[i].ValueType()
                    } else {
                        ltype := child.ValueType()
                        rtype := f.Arguments[i].ValueType()
                        if !ltype.IsCompatible(rtype) {
                            self.Script.AddError(self, child.Position(), "type mismatch (" + ltype.String() + " and " + rtype.String() + ")")
                            self.isValid = false
                        }
                    }
                } else {
                    self.isValid = false
                }
            }
        }
    }
}

func (self *FunctionCallExpression) ValueType() types.Type {
    if f, ok := self.Script.Functions[self.Name.String()]; ok {
        return f.Type
    } else {
        return types.Unknown
    }
}

func (self *FunctionCallExpression) GetChildren() []Node {
    result := []Node{}
    result = append(result, self.Name)
    for _, child := range self.Arguments {
        result = append(result, child)
    }
    return result
}

func (self *FunctionCallExpression) RealNode() Node {
    return self
}
