package nodes

import (
	"fmt"

	"../types"
)

type TypecastExpression struct {
	NodeBase

	Type  types.Type
	Child Expression
}

func (self *TypecastExpression) NodeType() NodeType {
	return NodeExpression
}

// TypecastExpression can be treated at Typecast
func (self *TypecastExpression) ExpressionType() ExpressionType {
	return ExpressionTypecast
}

func (self *TypecastExpression) String() string {
    if types.Type(self.Type.String()) == types.Type(self.Child.ValueType().String()) {
        return self.Child.String()
    }

	return fmt.Sprintf("(%s)%s", self.Type, self.Child)
}

func (self *TypecastExpression) ConnectTree() {
	self.Child.SetParent(self)
    self.Child.SetIndentationLevel(self.IndentationLevel)
	self.Child.SetScope(self.Scope)
	self.Child.SetScript(self.Script)
	self.Child.ConnectTree()

    self.Child = Expression(self.Child.RealNode())

    if self.Child.ExpressionType() == ExpressionListItem {
        self.Child.(*ListItemExpression).Type = self.Type
    } else if self.Child.IsValid() {
		self.isValid = true
		childType := self.Child.ValueType()

		// http://wiki.secondlife.com/wiki/Typecast
		switch self.Type {
		case types.Integer, types.Boolean:
			switch childType {
			case types.Integer, types.Float, types.String, types.Boolean:
			default:
				self.isValid = false
			}
		case types.Float:
			switch childType {
			case types.Integer, types.Float, types.String, types.Boolean:
			default:
				self.isValid = false
			}
		case types.String:
			switch childType {
			case types.Integer, types.Float, types.String, types.Key, types.List, types.Vector, types.Rotation:
			default:
				self.isValid = false
			}
		case types.Key:
			switch childType {
			case types.String, types.Key:
			default:
				self.isValid = false
			}
		case types.List:
			/*switch childType {
			case types.Integer, types.Float, types.String, types.Key, types.List, types.Vector, types.Rotation:
			default:
				self.isValid = false
			}*/
		case types.Vector:
			switch childType {
			case types.String, types.Vector:
			default:
				self.isValid = false
			}
		case types.Rotation:
			switch childType {
			case types.String, types.Rotation:
			default:
				self.isValid = false
			}
		default:
			self.isValid = false
		}

		if !self.isValid {
			self.Script.AddError(self, self.At, "invalid typecast ("+childType.String()+" to "+self.Type.String()+")")
		}
	} else {
		self.isValid = false
	}
}

func (self *TypecastExpression) ValueType() types.Type {
	return self.Type

}

func (self *TypecastExpression) GetChildren() []Node {
    return []Node{ self.Child }
}

func (self *TypecastExpression) RealNode() Node {
    return self
}
