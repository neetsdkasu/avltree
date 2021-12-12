package intarraytree

import "github.com/neetsdkasu/avltree"

const (
	PositionRootPosition int = iota
	PositionAllowDuplicateKey
	PositionIdleNodePosition
	HeaderSize
)

const (
	OffsetLeftChildPosition int = iota
	OffsetRightChildPosition
	OffsetHeight
	OffsetParentPosition
	OffsetNodeCount
	OffsetKey
	OffsetValue
	NodeSize
)

const NodeIsNothing int = 0

const (
	DisallowDuplicateKeys int = 0
	AllowDuplicateKeys    int = 1
)

type IntArrayTree struct {
	Array []int
}

type IntArrayTreeNode struct {
	Tree     *IntArrayTree
	Position int
	Data     []int
}

func New(allowDuplicateKeys bool) avltree.Tree {
	tree := &IntArrayTree{make([]int, HeaderSize)}
	tree.Init(allowDuplicateKeys)
	return tree
}

func (tree *IntArrayTree) Init(allowDuplicateKeys bool) {
	array := tree.Array
	if len(array) < HeaderSize {
		var buf [HeaderSize]int
		array = append(array, buf[:]...)
	}
	array = array[:HeaderSize]
	array[PositionRootPosition] = NodeIsNothing
	if allowDuplicateKeys {
		array[PositionAllowDuplicateKey] = AllowDuplicateKeys
	} else {
		array[PositionAllowDuplicateKey] = DisallowDuplicateKeys
	}
	array[PositionIdleNodePosition] = NodeIsNothing
	tree.Array = array
}

func unwrap(node avltree.Node) int {
	if node == nil {
		return NodeIsNothing
	} else {
		return node.(*IntArrayTreeNode).Position
	}
}

func (node *IntArrayTreeNode) toNode() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node
	}
}

func (tree *IntArrayTree) getNode(position int) *IntArrayTreeNode {
	if position == NodeIsNothing {
		return nil
	} else {
		return &IntArrayTreeNode{
			Tree:     tree,
			Position: position,
			Data:     tree.Array[position : position+NodeSize],
		}
	}
}

func (tree *IntArrayTree) getRoot() *IntArrayTreeNode {
	return tree.getNode(tree.Array[PositionRootPosition])
}

// intarraytree.New()以外でIntArrayTreeが生成されたときの気休め保険
func (tree *IntArrayTree) init() bool {
	if len(tree.Array) < HeaderSize {
		tree.Init(true)
		return true
	} else {
		return false
	}
}

func (tree *IntArrayTree) Root() avltree.Node {
	tree.init()
	return tree.getRoot().toNode()
}

func (tree *IntArrayTree) ReleaseNode(node avltree.RealNode) {
	tree.init()
	position := unwrap(node)
	if position != NodeIsNothing {
		tree.Array[position] = tree.Array[PositionIdleNodePosition]
		tree.Array[PositionIdleNodePosition] = position
	}
}

func (tree *IntArrayTree) NewNode(leftChild, rightChild avltree.Node, height int, key avltree.Key, value interface{}) avltree.RealNode {
	tree.init()
	array := tree.Array
	newNodePosition := array[PositionIdleNodePosition]
	if newNodePosition == NodeIsNothing {
		var buf [NodeSize]int
		newNodePosition = len(array)
		array = append(array, buf[:]...)
		tree.Array = array
	} else {
		nextIdleNodePosition := array[newNodePosition]
		array[PositionIdleNodePosition] = nextIdleNodePosition
	}
	node := tree.getNode(newNodePosition)
	node.Data[OffsetLeftChildPosition] = unwrap(leftChild)
	node.Data[OffsetRightChildPosition] = unwrap(rightChild)
	node.Data[OffsetHeight] = height
	node.Data[OffsetParentPosition] = NodeIsNothing
	node.Data[OffsetNodeCount] = 1
	node.Data[OffsetKey] = int(key.(avltree.IntKey))
	node.Data[OffsetValue] = value.(int)
	node.resetNodeCount()
	return node
}

func (tree *IntArrayTree) SetRoot(newRoot avltree.RealNode) avltree.RealTree {
	tree.init()
	tree.Array[PositionRootPosition] = unwrap(newRoot)
	tree.getRoot().setParent(NodeIsNothing)
	return tree
}

func (tree *IntArrayTree) AllowDuplicateKeys() bool {
	tree.init()
	return tree.Array[PositionAllowDuplicateKey] == AllowDuplicateKeys
}

func (tree *IntArrayTree) NodeCount() int {
	tree.init()
	return tree.getRoot().NodeCount()
}

func (tree *IntArrayTree) CleanUpTree() {
	if tree.init() {
		return
	}
	tree.Init(tree.AllowDuplicateKeys())
}

func (node *IntArrayTreeNode) Key() avltree.Key {
	return avltree.IntKey(node.Data[OffsetKey])
}

func (node *IntArrayTreeNode) Value() interface{} {
	return node.Data[OffsetValue]
}

func (node *IntArrayTreeNode) getLeftChild() *IntArrayTreeNode {
	return node.Tree.getNode(node.Data[OffsetLeftChildPosition])
}

func (node *IntArrayTreeNode) LeftChild() avltree.Node {
	return node.getLeftChild().toNode()
}

func (node *IntArrayTreeNode) getRightChild() *IntArrayTreeNode {
	return node.Tree.getNode(node.Data[OffsetRightChildPosition])
}

func (node *IntArrayTreeNode) RightChild() avltree.Node {
	return node.getRightChild().toNode()
}

func (node *IntArrayTreeNode) SetValue(newValue interface{}) avltree.Node {
	node.Data[OffsetValue] = newValue.(int)
	return node
}

func (node *IntArrayTreeNode) setParent(position int) {
	if node != nil {
		node.Data[OffsetParentPosition] = position
	}
}

func (node *IntArrayTreeNode) Parent() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node.Tree.getNode(node.Data[OffsetParentPosition]).toNode()
	}
}

func (node *IntArrayTreeNode) NodeCount() int {
	if node == nil {
		return 0
	} else {
		return node.Data[OffsetNodeCount]
	}
}

func (node *IntArrayTreeNode) Height() int {
	return node.Data[OffsetHeight]
}

func (node *IntArrayTreeNode) resetNodeCount() {
	node.Data[OffsetNodeCount] = 1 +
		node.getLeftChild().NodeCount() +
		node.getRightChild().NodeCount()
}

func (node *IntArrayTreeNode) SetChildren(newLeftChild, newRightChild avltree.Node, newHeight int) avltree.RealNode {
	node.Data[OffsetLeftChildPosition] = unwrap(newLeftChild)
	node.Data[OffsetRightChildPosition] = unwrap(newRightChild)
	node.Data[OffsetHeight] = newHeight
	node.getLeftChild().setParent(node.Position)
	node.getRightChild().setParent(node.Position)
	node.resetNodeCount()
	return node
}

func (node *IntArrayTreeNode) Set(newLeftChild, newRightChild avltree.Node, newHeight int, newValue interface{}) avltree.RealNode {
	node.Data[OffsetValue] = newValue.(int)
	return node.SetChildren(newLeftChild, newRightChild, newHeight)
}
