package nodes

import (
    "fmt"
	"strings"

	"../types"
)

type BlockStatement struct {
	NodeBase

    NoBraces bool
	Children []Statement
}

func (self *BlockStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *BlockStatement) StatementType() StatementType {
	return StatementBlock
}

func (self *BlockStatement) String() string {
	var prev Statement

    result := ""

    if self.NoBraces {
        result += "\n"
    } else {
    	result += "{\n"
    }

	for _, child := range self.Children {
		if child.StatementType() == StatementEmpty {
            prev = child
			continue
		}

        if self.Script.SkipUnused {
            switch child.NodeType() {
            case NodeVariable:
                if !child.(*Variable).IsUsed {
                    continue
                }
            }
        }

		if prev != nil {
			switch prev.StatementType() {
			case StatementIf, StatementFor, StatementDo, StatementWhile:
				result += "\n"
			case StatementExpression:
				prevExpr := prev.(*ExpressionStatement).Expression
				if child.StatementType() == StatementExpression {
					childExpr := child.(*ExpressionStatement).Expression

					if prevExpr.ExpressionType() != childExpr.ExpressionType() {
						result += "\n"
					}
				} else {
					result += "\n"
				}
			default:
                if prev.StatementType() != child.StatementType() {
					result += "\n"
				}
			}
		}

        switch child.StatementType() {
        case StatementComment:
            if prev != nil && child.Position().Line == prev.Position().Line {
                result = strings.TrimRight(result, "\n\r\t ") + " "
                result += child.String() + "\n"
            } else {
                result = strings.TrimRight(result, "\n\r\t ") + "\n\n"
                result += strings.Repeat(Indentation, child.GetIndentationLevel())
                result += child.String() + "\n"
            }
        case StatementCase:
    		result += strings.Repeat(Indentation, child.GetIndentationLevel() - 1)
		    result += child.String()
        default:
    		result += strings.Repeat(Indentation, child.GetIndentationLevel())
		    result += child.String() + "\n"
        }

		prev = child
	}

	result += strings.Repeat(Indentation, self.IndentationLevel)

    if !self.NoBraces {
    	result += "}"
    }

	return result
}

func (self *BlockStatement) ConnectTree() {
    self.isValid = true

	currScope := self.Scope

	for _, child := range self.Children {
		child.SetParent(self)
		child.SetIndentationLevel(self.IndentationLevel + 1)
		child.SetScope(currScope)
		child.SetScript(self.Script)
		child.ConnectTree()

        self.isValid = self.isValid || child.IsValid()

		switch child.NodeType() {
		case NodeVariable:
            existing := currScope.FindVariable(child.(*Variable).Name.String())
            if existing != nil && existing.Parent == self {
                self.isValid = false
                self.Script.AddError(child, child.Position(), fmt.Sprintf("redeclared variable \"%s\" (previously declared at %s)", child.(*Variable).Name, existing.Position()))
            }
			currScope = currScope.AddVariable(child.(*Variable))
		}
	}
}

func (self *BlockStatement) ValueType() types.Type {
	return types.Unknown
}

func (self *BlockStatement) GetChildren() []Node {
    result := make([]Node, len(self.Children))
    for i, child := range self.Children {
        result[i] = child
    }
    return result
}

func (self *BlockStatement) RealNode() Node {
    return self
}
