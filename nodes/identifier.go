package nodes

type Identifier struct {
	NodeBase

	Name string
}

func (self *Identifier) NodeType() NodeType {
	return NodeIdentifier
}

// Identifier can be treated as LValue expression
func (self *Identifier) ExpressionType() ExpressionType {
	return ExpressionLValue
}

func (self *Identifier) String() string {
	return self.Name
}

func (self *Identifier) ConnectTree() {
}

func (self *Identifier) RealNode() Node {
    return self
}
