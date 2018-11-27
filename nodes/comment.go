package nodes

import (
	"fmt"
	"strings"
)

type Comment struct {
	NodeBase

	Text     string
	IsCStyle bool
}

func (self *Comment) NodeType() NodeType {
	return NodeComment
}

// Comment also can be treated as statement
func (self *Comment) StatementType() StatementType {
	return StatementComment
}

func (self *Comment) String() string {
	text := strings.Trim(self.Text, " \t\n\r")
	if strings.Contains(text, "\n") {
		return fmt.Sprintf("/*\n%s\n*/", strings.Trim(self.Text, " \t\n\r"))
	} else {
		return fmt.Sprintf("// %s", text)
	}
}

func (self *Comment) ConnectTree() {
}

func (self *Comment) RealNode() Node {
    return self
}
