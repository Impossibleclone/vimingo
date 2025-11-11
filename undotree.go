package main

type snapshot struct {
	lines  []string
	cursor Cursor
}

type Node struct {
	snapshot snapshot
	parent   *Node
	children []*Node
}
