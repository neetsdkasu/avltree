// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのRealTree,RealNodeの実装例
// int型の可変長配列(スライス)上にAVL木を構築する
//
// 扱えるキーはintkeyのIntKeyのみ
// 扱える値はint型のみ
//
// int型は32bitマシンだと4bytes、64bitマシンだと8bytesになるので実行環境によって必要メモリ量が変わる
//
// 配列の最初の３要素は以下の木に関する情報を保持
// + ルートノードのインデックス
// + 同一キーを許可するかどうかの値
// + 再利用可能なノードのインデックス
//
// １つのノードは７個分の要素で構成され、以下の情報を保持
// + 左の子ノードのインデックス
// + 右の子ノードのインデックス
// + 木におけるノードの高さ
// + 親ノードのインデックス
// + 左右の子孫も合わせたノード総数
// + ノードのキー
// + ノードの値
//
// コード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			"github.com/neetsdkasu/avltree/intarraytree"
//			. "github.com/neetsdkasu/avltree/intkey"
//		)
//		func Example_intarraytree() {
//			tree := intarraytree.New(false)
//			avltree.Insert(tree, false, IntKey(12), 345)
//			avltree.Insert(tree, false, IntKey(67), 890)
//			avltree.Insert(tree, false, IntKey(333), 666)
//			avltree.Insert(tree, false, IntKey(-5), 12345)
//			avltree.Delete(tree, IntKey(67))
//			avltree.Update(tree, IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
//				newValue = oldValue.(int) * 3
//				return
//			})
//			if node := avltree.Find(tree, IntKey(12)); node != nil {
//				fmt.Println("Find!", node.Key(), node.Value())
//			}
//			avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
//				fmt.Println("Iterate!", node.Key(), node.Value())
//				return
//			})
//			// Output:
//			// Find! 12 345
//			// Iterate! -5 12345
//			// Iterate! 12 345
//			// Iterate! 333 1998
//		}
//
package intarraytree

import (
	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/intkey"
)

const (
	PositionRootPosition int = iota
	PositionDuplicateKeysBehavior
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
}

func New(allowDuplicateKeys bool) avltree.Tree {
	return NewWithInitialCapacity(HeaderSize, allowDuplicateKeys)
}

func NewWithInitialCapacity(initialCapacity int, allowDuplicateKeys bool) avltree.Tree {
	tree := &IntArrayTree{make([]int, initialCapacity)}
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
		array[PositionDuplicateKeysBehavior] = AllowDuplicateKeys
	} else {
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
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
	node.set(OffsetLeftChildPosition, unwrap(leftChild))
	node.set(OffsetRightChildPosition, unwrap(rightChild))
	node.set(OffsetHeight, height)
	node.set(OffsetParentPosition, NodeIsNothing)
	node.set(OffsetNodeCount, 1)
	node.set(OffsetKey, int(key.(intkey.IntKey)))
	node.set(OffsetValue, value.(int))
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
	return tree.Array[PositionDuplicateKeysBehavior] == AllowDuplicateKeys
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
	return intkey.IntKey(node.get(OffsetKey))
}

func (node *IntArrayTreeNode) Value() interface{} {
	return node.get(OffsetValue)
}

func (node *IntArrayTreeNode) get(offset int) int {
	return node.Tree.Array[node.Position+offset]
}
func (node *IntArrayTreeNode) set(offset, value int) {
	node.Tree.Array[node.Position+offset] = value
}

func (node *IntArrayTreeNode) getLeftChild() *IntArrayTreeNode {
	return node.Tree.getNode(node.get(OffsetLeftChildPosition))
}

func (node *IntArrayTreeNode) LeftChild() avltree.Node {
	return node.getLeftChild().toNode()
}

func (node *IntArrayTreeNode) getRightChild() *IntArrayTreeNode {
	return node.Tree.getNode(node.get(OffsetRightChildPosition))
}

func (node *IntArrayTreeNode) RightChild() avltree.Node {
	return node.getRightChild().toNode()
}

func (node *IntArrayTreeNode) SetValue(newValue interface{}) avltree.Node {
	node.set(OffsetValue, newValue.(int))
	return node
}

func (node *IntArrayTreeNode) setParent(position int) {
	if node != nil {
		node.set(OffsetParentPosition, position)
	}
}

func (node *IntArrayTreeNode) Parent() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node.Tree.getNode(node.get(OffsetParentPosition)).toNode()
	}
}

func (node *IntArrayTreeNode) NodeCount() int {
	if node == nil {
		return 0
	} else {
		return node.get(OffsetNodeCount)
	}
}

func (node *IntArrayTreeNode) Height() int {
	return node.get(OffsetHeight)
}

func (node *IntArrayTreeNode) resetNodeCount() {
	node.set(OffsetNodeCount, 1+node.getLeftChild().NodeCount()+node.getRightChild().NodeCount())
}

func (node *IntArrayTreeNode) SetChildren(newLeftChild, newRightChild avltree.Node, newHeight int) avltree.RealNode {
	node.set(OffsetLeftChildPosition, unwrap(newLeftChild))
	node.set(OffsetRightChildPosition, unwrap(newRightChild))
	node.set(OffsetHeight, newHeight)
	node.getLeftChild().setParent(node.Position)
	node.getRightChild().setParent(node.Position)
	node.resetNodeCount()
	return node
}

func (node *IntArrayTreeNode) Set(newLeftChild, newRightChild avltree.Node, newHeight int, newValue interface{}) avltree.RealNode {
	node.SetValue(newValue)
	return node.SetChildren(newLeftChild, newRightChild, newHeight)
}
