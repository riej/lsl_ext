package nodes

import (
	"fmt"
	"reflect"
	"text/scanner"

	"../types"
)

type NodeType int

const (
	NodeUnknown NodeType = iota

	NodeScript

	NodeComment

	NodeIdentifier
	NodeConstant

	NodeVariable
	NodeFunction
    NodeStruct

	NodeState

	NodeExpression
	NodeStatement
)

type Node interface {
	Position() scanner.Position
    SetPosition(at scanner.Position)

	NodeType() NodeType
	StatementType() StatementType
	ExpressionType() ExpressionType

	String() string

	ConnectTree() // connect children to parent and set their indentation levels

	GetIndentationLevel() int
	SetIndentationLevel(value int)

	GetParent() Node
	SetParent(parent Node)

    GetChildren() []Node

	GetScope() *Scope
	SetScope(value *Scope)

	GetScript() *Script
	SetScript(script *Script)

	ValueType() types.Type
	IsValid() bool

    RealNode() Node
}

type BreakableNode interface {
    Node

    BreakLabel() string
}

type ContinueableNode interface {
    Node

    ContinueLabel() string
}

type NodeBase struct {
	At               scanner.Position
	IndentationLevel int
	Parent           Node
	Scope            *Scope
	Script           *Script
	isValid          bool
}

func (self *NodeBase) Position() scanner.Position {
	return self.At
}

func (self *NodeBase) GetIndentationLevel() int {
	return self.IndentationLevel
}

func (self *NodeBase) SetIndentationLevel(value int) {
	self.IndentationLevel = value
}

func (self *NodeBase) GetParent() Node {
	return self.Parent
}

func (self *NodeBase) SetParent(parent Node) {
	self.Parent = parent
}

func (self *NodeBase) NodeType() NodeType {
	return NodeUnknown
}

func (self *NodeBase) ExpressionType() ExpressionType {
	return ExpressionUnknown
}

func (self *NodeBase) StatementType() StatementType {
	return StatementUnknown
}

func (self *NodeBase) ConnectTree() {
}

func (self *NodeBase) String() string {
	return ""
}

func (self *NodeBase) GetScope() *Scope {
	return self.Scope
}

func (self *NodeBase) SetScope(value *Scope) {
	self.Scope = value
}

func (self *NodeBase) GetScript() *Script {
	var node Node
	node = self
	for node.NodeType() != NodeScript {
		node = node.GetParent()
	}
	return node.(*Script)
}

func (self *NodeBase) SetScript(script *Script) {
	self.Script = script
}

func (self *NodeBase) ValueType() types.Type {
	return types.Unknown
}

func (self *NodeBase) IsValid() bool {
	return self.isValid
}

func (self *NodeBase) DumpTree() {
	var node Node
	node = self
	fmt.Printf("-------\n")
	for node != nil {
		fmt.Printf("%s %s (scope %p, parent %p)\n", node.Position(), reflect.TypeOf(node).Elem().Name(), node.GetScope(), node.GetScope().Parent)
        str := node.String()
        if len(str) < 500 {
            fmt.Println(str)
        }
		node = node.GetParent()
	}
}

func (self *NodeBase) SetPosition(at scanner.Position) {
    self.At = at
}

func (self *NodeBase) GetChildren() []Node {
    return []Node{}
}

func (self *NodeBase) RealNode() Node {
    return self
}
