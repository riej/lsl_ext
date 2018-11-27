package nodes

import (
    "strings"

	"../types"
)

type ListExpression struct {
	NodeBase

	Values []Expression
}

func (self *ListExpression) NodeType() NodeType {
	return NodeExpression
}

func (self *ListExpression) ExpressionType() ExpressionType {
	return ExpressionList
}

func (self *ListExpression) String() string {
    hasLineBreaks := false

	result := "["
	for i, value := range self.Values {
		if i > 0 {
			result += ", "
		}

        val := value.String()
        breakTheLine := false

        if self.Script.Builtins.FindConstant(val) != nil {
            if strings.HasPrefix(val, "HTTP_") {
                breakTheLine = true
            }

            switch val {
            case "PRIM_NAME", "PRIM_DESC", "PRIM_TYPE", "PRIM_SLICE", "PRIM_PHYSICS_SHAPE_TYPE", "PRIM_MATERIAL", "PRIM_PHYSICS", "PRIM_TEMP_ON_REZ", "PRIM_PHANTOM", "PRIM_POSITION", "PRIM_POS_LOCAL", "PRIM_ROTATION", "PRIM_ROT_LOCAL", "PRIM_SIZE", "PRIM_TEXTURE", "PRIM_TEXT", "PRIM_COLOR", "PRIM_BUMP_SHINY", "PRIM_POINT_LIGHT", "PRIM_FULLBRIGHT", "PRIM_FLEXIBLE", "PRIM_TEXGEN", "PRIM_GLOW", "PRIM_OMEGA", "PRIM_NORMAL", "PRIM_SPECULAR", "PRIM_ALPHA_MODE", "PRIM_LINK_TARGET", "PRIM_CAST_SHADOWS", "PRIM_TYPE_LEGACY", "PRIM_ALLOW_UNSIT", "PRIM_SCRIPTED_SIT_ONLY", "PRIM_SIT_TARGET":
                if self.Parent.ExpressionType() == ExpressionFunctionCall {
                    funcName := self.Parent.(*FunctionCallExpression).Name.String()
                    switch funcName {
                    case "llGetPrimitiveParams", "llGetLinkPrimitiveParams":
                    default:
                        breakTheLine = true
                    }
                } else {
                    breakTheLine = true
                }
            }
        }

        if breakTheLine {
            result += "\n" + strings.Repeat(Indentation, self.IndentationLevel + 1)
            hasLineBreaks = true
        }
		result += val
	}
    if hasLineBreaks {
        result += "\n" + strings.Repeat(Indentation, self.IndentationLevel)
    }

    result += "]"

	return result
}

func (self *ListExpression) ConnectTree() {
	for _, child := range self.Values {
		child.SetParent(self)
		child.SetScope(self.Scope)
		child.SetScript(self.Script)
		child.ConnectTree()
	}


    self.isValid = true
	for _, child := range self.Values {
        if !child.IsValid() {
            self.isValid = false
        } else if child.ValueType().IsList() {
			self.Script.AddError(child, child.Position(), "lists cannot contain other lists")
            self.isValid = false
		}
	}
}

func (self *ListExpression) ValueType() types.Type {
    return types.List
}

func (self *ListExpression) GetChildren() []Node {
    result := make([]Node, len(self.Values))
    for i, child := range self.Values {
        result[i] = child
    }
    return result
}

func (self *ListExpression) RealNode() Node {
    return self
}
