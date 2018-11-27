package nodes

type Builtins struct {
	Functions []*Function
	Constants []*Variable
	Events    []*Function
}

func (self *Builtins) FindFunction(name string) *Function {
	for _, child := range self.Functions {
		if child.Name.String() == name {
			return child
		}
	}

	return nil
}

func (self *Builtins) FindConstant(name string) *Variable {
	for _, child := range self.Constants {
		if child.Name.String() == name {
			return child
		}
	}

	return nil
}

func (self *Builtins) FindEvent(name string) *Function {
	for _, child := range self.Events {
		if child.Name.String() == name {
			return child
		}
	}

	return nil
}

func (self *Builtins) Find(name string) Node {
	var child Node
	child = self.FindFunction(name)
	if child != nil {
		return child
	}
	child = self.FindConstant(name)
	if child != nil {
		return child
	}
	child = self.FindEvent(name)
	if child != nil {
		return child
	}
	return nil
}
