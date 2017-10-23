package main

// PositionNodeType type of Position Node
type PositionNodeType uint

const (
	// MoveNode move type of node in position tree
	MoveNode PositionNodeType = 1
	// CommentNode comment type of node in position tree
	CommentNode PositionNodeType = 2
	// VariationNode variation type of node in position tree
	VariationNode PositionNodeType = 3
)

// PositionNode struct
type PositionNode struct {
	nodetype       PositionNodeType
	variationdepth int
	comment        string
	move           MoveType
	potential      []MoveType
}
