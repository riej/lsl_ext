package nodes

import (
	"../types"
)

type Constant struct {
	NodeBase

	Value types.Value
}

func (self *Constant) NodeType() NodeType {
	return NodeConstant
}

// Constant can be treated as expression
func (self *Constant) ExpressionType() ExpressionType {
	return ExpressionConstant
}

func (self *Constant) String() string {
	return self.Value.String()
}

func (self *Constant) ConnectTree() {
}

func (self *Constant) ValueType() types.Type {
	return self.Value.Type()
}

func (self *Constant) RealNode() Node {
    return self
}
