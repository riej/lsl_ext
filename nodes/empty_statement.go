package nodes

type EmptyStatement struct {
	NodeBase
}

func (self *EmptyStatement) NodeType() NodeType {
	return NodeStatement
}

func (self *EmptyStatement) StatementType() StatementType {
	return StatementEmpty
}

func (self *EmptyStatement) String() string {
	return ";"
}

func (self *EmptyStatement) ConnectTree() {
}

func (self *EmptyStatement) RealNode() Node {
    return self
}
