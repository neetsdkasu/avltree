package linkedtree

import "avltree"

type (
	Key      = avltree.Key
	Node     = avltree.Node
	RealNode = avltree.RealNode
	Tree     = avltree.Tree
)

type LinkedTree struct {
	root               *linkedTreeNode
	allowDuplicateKeys bool
}

type linkedTreeNode struct {
	leftChild  *linkedTreeNode
	rightChild *linkedTreeNode
	height     int
	parent     *linkedTreeNode
	nodeCount  int
	key        Key
	value      interface{}
}

func NewLinkedTree(allowDuplicateKeys bool) *LinkedTree {
	return &LinkedTree{
		nil, // root
		allowDuplicateKeys,
	}
}

func unwrap(node Node) *linkedTreeNode {
	if ltNode, ok := node.(*linkedTreeNode); ok {
		return ltNode
	} else {
		return nil
	}
}

func (node *linkedTreeNode) toNode() Node {
	if node == nil {
		return nil
	} else {
		return node
	}
}

func (tree *LinkedTree) NodeCount() int {
	return tree.root.NodeCount()
}

func (tree *LinkedTree) NewNode(leftChild, rightChild Node, height int, key Key, value interface{}) RealNode {
	node := &linkedTreeNode{
		unwrap(leftChild),
		unwrap(rightChild),
		height,
		nil, // parent
		1,   // nodeCount
		key,
		value,
	}
	node.resetNodeCount()
	return node
}

func (tree *LinkedTree) Root() Node {
	return tree.root.toNode()
}

func (tree *LinkedTree) SetRoot(newRoot RealNode) Tree {
	tree.root = unwrap(newRoot)
	tree.root.SetParent(nil)
	return tree
}

func (tree *LinkedTree) AllowDuplicateKeys() bool {
	return tree.allowDuplicateKeys
}

func (node *linkedTreeNode) NodeCount() int {
	if node == nil {
		return 0
	} else {
		return node.nodeCount
	}
}

func (node *linkedTreeNode) Key() Key {
	return node.key
}

func (node *linkedTreeNode) Value() interface{} {
	return node.value
}

func (node *linkedTreeNode) Height() int {
	return node.height
}

func (node *linkedTreeNode) LeftChild() Node {
	return node.leftChild.toNode()
}

func (node *linkedTreeNode) RightChild() Node {
	return node.rightChild.toNode()
}

func (node *linkedTreeNode) Parent() Node {
	if node == nil {
		return nil
	} else {
		return node.parent.toNode()
	}
}

func (node *linkedTreeNode) SetValue(newValue interface{}) Node {
	node.value = newValue
	return node
}

func (node *linkedTreeNode) SetParent(newParent Node) RealNode {
	if node != nil {
		node.parent = unwrap(newParent)
	}
	return node
}

func (node *linkedTreeNode) resetNodeCount() {
	node.nodeCount = 1 + node.leftChild.NodeCount() + node.rightChild.NodeCount()
}

func (node *linkedTreeNode) SetLeftChild(newLeftChild Node, newHeight int) RealNode {
	node.leftChild = unwrap(newLeftChild)
	node.leftChild.SetParent(node)
	node.height = newHeight
	node.resetNodeCount()
	return node
}

func (node *linkedTreeNode) SetRightChild(newRightChild Node, newHeight int) RealNode {
	node.rightChild = unwrap(newRightChild)
	node.rightChild.SetParent(node)
	node.height = newHeight
	node.resetNodeCount()
	return node
}

func (node *linkedTreeNode) SetChildren(newLeftChild, newRightChild Node, newHeight int) RealNode {
	node.leftChild = unwrap(newLeftChild)
	node.rightChild = unwrap(newRightChild)
	node.leftChild.SetParent(node)
	node.rightChild.SetParent(node)
	node.height = newHeight
	node.resetNodeCount()
	return node
}
