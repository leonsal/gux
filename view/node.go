package view

type INode interface {
	Children() []INode
}

type Node struct {
}
