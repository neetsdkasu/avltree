package avltree

type NodeIteratorCallBack = func(node Node) bool

type NodeCounter interface{ NodeCount() int }
type NodeDeallocator interface{ ReleaseNode(node RealNode) }

type Tree interface {
	NewNode(leftChild, rightChild Node, height int, key Key, value interface{}) RealNode
	Root() Node
	SetRoot(newRoot RealNode) Tree
	AllowDuplicateKeys() bool
}

type Node interface {
	Key() Key
	Value() interface{}
	LeftChild() Node
	RightChild() Node
	SetValue(newValue interface{}) Node
}

type RealNode interface {
	Node
	Height() int
	SetLeftChild(newLeftChild Node, newHeight int) RealNode
	SetRightChild(newRightChild Node, newHeight int) RealNode
	SetChildren(newLeftChild, newRightChild Node, newHeight int) RealNode
}

type Key interface {
	CompareTo(other Key) int
}

func Insert(tree Tree, replaceIfExists bool, key Key, value interface{}) (modified Tree, ok bool) {
	helper := insertHelper{
		&tree,
		replaceIfExists,
		&key,
		&value,
	}
	if newRoot, ok := helper.insertTo(tree.Root()); ok {
		return tree.SetRoot(newRoot), true
	} else {
		return tree, false
	}
}

func Delete(tree Tree, key Key) (modified Tree, value interface{}, ok bool) {
	if newRoot, node, ok := removeNode(tree.Root(), key); ok {
		value = node.Value()
		if deallocator, ok := tree.(NodeDeallocator); ok {
			deallocator.ReleaseNode(node.(RealNode))
		}
		if root, ok := newRoot.(RealNode); ok {
			return tree.SetRoot(root), value, true
		} else {
			return tree.SetRoot(nil), value, true
		}
	} else {
		return tree, nil, false
	}
}

func Find(tree Tree, key Key) (node Node, ok bool) {
	node = tree.Root()
	for node != nil {
		cmp := key.CompareTo(node.Key())
		switch {
		case cmp < 0:
			node = node.LeftChild()
		case 0 < cmp:
			node = node.RightChild()
		default:
			return node, true
		}
	}
	return nil, false
}

func Iterate(tree Tree, descOrder bool, callBack NodeIteratorCallBack) bool {
	if descOrder {
		return descIterateNode(tree.Root(), callBack)
	} else {
		return ascIterateNode(tree.Root(), callBack)
	}
}

func Range(tree Tree, descOrder bool, lower, upper Key, callBack NodeIteratorCallBack) bool {
	if tree.AllowDuplicateKeys() {
		if descOrder {
			return descExtendedRangeNode(tree.Root(), lower, upper, callBack)
		} else {
			return ascExtendedRangeNode(tree.Root(), lower, upper, callBack)
		}
	} else {
		if descOrder {
			return descRangeNode(tree.Root(), lower, upper, callBack)
		} else {
			return ascRangeNode(tree.Root(), lower, upper, callBack)
		}
	}
}

func Count(tree Tree) int {
	if counter, ok := tree.(NodeCounter); ok {
		return counter.NodeCount()
	}
	root := tree.Root()
	if counter, ok := root.(NodeCounter); ok {
		return counter.NodeCount()
	} else {
		return countNode(root)
	}
}

func Min(tree Tree) (minimum Node, ok bool) {
	node := tree.Root()
	if node == nil {
		return nil, false
	}
	for {
		leftChild := node.LeftChild()
		if leftChild == nil {
			break
		}
		node = leftChild
	}
	return node, true
}

func Max(tree Tree) (maximum Node, ok bool) {
	node := tree.Root()
	if node == nil {
		return nil, false
	}
	for {
		rightChild := node.RightChild()
		if rightChild == nil {
			break
		}
		node = rightChild
	}
	return node, true
}

func countNode(node Node) int {
	if node == nil {
		return 0
	}
	return 1 + countNode(node.LeftChild()) + countNode(node.RightChild())
}

func getHeight(node Node) int {
	if node == nil {
		return 0
	} else {
		return node.(RealNode).Height()
	}
}

type balance int

const (
	balanced balance = iota
	leftIsHigher
	rightIsHigher
)

func checkBalance(node RealNode) balance {
	if node == nil {
		return balanced
	}
	heightL := getHeight(node.LeftChild())
	heightR := getHeight(node.RightChild())
	// 算術オーバーフローが怖いのかい？
	// heightL - heightR
	// heightL + 1
	// heightR + 1
	switch {
	case heightL < heightR && heightL+1 < heightR:
		return rightIsHigher
	case heightL > heightR && heightL > heightR+1:
		return leftIsHigher
	default:
		return balanced
	}
}

func compareChildHeight(node Node) balance {
	if node == nil {
		return balanced
	}
	heightL := getHeight(node.LeftChild())
	heightR := getHeight(node.RightChild())
	// 算術オーバーフローが怖いのかい？
	// heightL - heightR
	switch {
	case heightL < heightR:
		return rightIsHigher
	case heightL > heightR:
		return leftIsHigher
	default:
		return balanced
	}
}

func intMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func calcNewHeight(leftChild, rightChild Node) int {
	leftHeight := getHeight(leftChild)
	rightHeight := getHeight(rightChild)
	return 1 + intMax(leftHeight, rightHeight)
}

func setLeftChild(root RealNode, newLeftChild Node) RealNode {
	newHeight := calcNewHeight(newLeftChild, root.RightChild())
	return root.SetLeftChild(newLeftChild, newHeight)
}

func setRightChild(root RealNode, newRightChild Node) RealNode {
	newHeight := calcNewHeight(root.LeftChild(), newRightChild)
	return root.SetRightChild(newRightChild, newHeight)
}

type insertHelper struct {
	tree            *Tree
	replaceIfExists bool
	key             *Key
	value           *interface{}
}

func (helper *insertHelper) newNode() RealNode {
	return (*helper.tree).NewNode(nil, nil, 1, *helper.key, *helper.value)
}

func (helper *insertHelper) compareKey(node Node) int {
	return (*helper.key).CompareTo(node.Key())
}

func (helper *insertHelper) allowDuplicateKeys() bool {
	return (*helper.tree).AllowDuplicateKeys()
}

func (helper *insertHelper) insertTo(root Node) (newRoot RealNode, ok bool) {
	if root == nil {
		return helper.newNode(), true
	}
	cmp := helper.compareKey(root)
	switch {
	case cmp < 0: // newKey < root.key
		if newLeftChild, ok := helper.insertTo(root.LeftChild()); ok {
			newRoot = setLeftChild(root.(RealNode), newLeftChild)
		} else {
			return root.(RealNode), false
		}
	case 0 < cmp: // root.key < newKey
		if newRightChild, ok := helper.insertTo(root.RightChild()); ok {
			newRoot = setRightChild(root.(RealNode), newRightChild)
		} else {
			return root.(RealNode), false
		}
	default:
		if helper.replaceIfExists {
			newRoot = root.SetValue(*helper.value).(RealNode)
			return newRoot, true
		}
		if !helper.allowDuplicateKeys() {
			return root.(RealNode), false
		}
		if newRightChild, ok := helper.insertTo(root.RightChild()); ok {
			newRoot = setRightChild(root.(RealNode), newRightChild)
		} else {
			return root.(RealNode), false
		}
	}

	newRoot = rotate(newRoot)

	return newRoot, true
}

func rotate(root RealNode) RealNode {
	for {
		switch checkBalance(root) {
		case leftIsHigher:
			root = rotateRight(root)
		case rightIsHigher:
			root = rotateLeft(root)
		default:
			return root
		}
	}
}

func rotateRight(root RealNode) RealNode {
	oldRootLeftChild := root.LeftChild().(RealNode)
	if compareChildHeight(oldRootLeftChild) == rightIsHigher {
		newRoot := oldRootLeftChild.RightChild().(RealNode)
		tempLeftChild := newRoot.LeftChild()
		tempRightChild := newRoot.RightChild()
		newRootRightChild := setLeftChild(root, tempRightChild)
		newRootLeftChild := setRightChild(oldRootLeftChild, tempLeftChild)
		return setChildren(newRoot, newRootLeftChild, newRootRightChild)
	} else {
		tempRightChild := oldRootLeftChild.RightChild()
		newRootRightChild := setLeftChild(root, tempRightChild)
		return setRightChild(oldRootLeftChild, newRootRightChild)
	}
}

func rotateLeft(root RealNode) RealNode {
	oldRootRightChild := root.RightChild().(RealNode)
	if compareChildHeight(oldRootRightChild) == leftIsHigher {
		newRoot := oldRootRightChild.LeftChild().(RealNode)
		tempLeftChild := newRoot.LeftChild()
		tempRightChild := newRoot.RightChild()
		newRootLeftChild := setRightChild(root, tempLeftChild)
		newRootRightChild := setLeftChild(oldRootRightChild, tempRightChild)
		return setChildren(newRoot, newRootLeftChild, newRootRightChild)
	} else {
		tempLeftChild := oldRootRightChild.LeftChild()
		newLeftChild := setRightChild(root, tempLeftChild)
		return setLeftChild(oldRootRightChild, newLeftChild)
	}
}

func setChildren(root RealNode, leftChild, rightChild Node) RealNode {
	newHeight := 1 + getHeight(leftChild) + getHeight(rightChild)
	return root.SetChildren(leftChild, rightChild, newHeight)
}

func removeNode(root Node, key Key) (newRoot, removed Node, ok bool) {
	if root == nil {
		return nil, nil, false
	}
	cmp := key.CompareTo(root.Key())
	switch {
	case cmp < 0: // key < root.Key()
		if tempLeftChild, node, ok := removeNode(root.LeftChild(), key); ok {
			removed = node
			newRoot = setLeftChild(root.(RealNode), tempLeftChild)
		} else {
			return nil, nil, false
		}
	case 0 < cmp: // root.Key() < key
		if tempRightChild, node, ok := removeNode(root.RightChild(), key); ok {
			removed = node
			newRoot = setRightChild(root.(RealNode), tempRightChild)
		} else {
			return nil, nil, false
		}
	default: // just target node
		removed = root
		leftChild := root.LeftChild()
		rightChild := root.RightChild()
		if compareChildHeight(root) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot == nil {
			// no children
			// height(root) == 1
			// root.LeftChild() == nil
			// root.RightChild() == nil
			return nil, removed, true
		}
		newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
	}

	newRoot = rotate(newRoot.(RealNode))

	return newRoot, removed, true
}

func removeMin(root Node) (newRoot, removed Node) {
	if root == nil {
		return nil, nil
	}
	leftChild := root.LeftChild()
	if leftChild == nil {
		newRoot = root.RightChild()
		return newRoot, root
	}
	leftChild, removed = removeMin(leftChild)
	newRoot = setLeftChild(root.(RealNode), leftChild)
	newRoot = rotate(newRoot.(RealNode))
	return newRoot, removed
}

func removeMax(root Node) (newRoot, removed Node) {
	if root == nil {
		return nil, nil
	}
	rightChild := root.RightChild()
	if rightChild == nil {
		newRoot = root.LeftChild()
		return newRoot, root
	}
	rightChild, removed = removeMax(rightChild)
	newRoot = setRightChild(root.(RealNode), rightChild)
	newRoot = rotate(newRoot.(RealNode))
	return newRoot, removed
}

func ascIterateNode(node Node, callBack NodeIteratorCallBack) bool {
	if node == nil {
		return false
	}
	if ascIterateNode(node.LeftChild(), callBack) {
		return true
	}
	if callBack(node) {
		return true
	}
	return ascIterateNode(node.RightChild(), callBack)
}

func descIterateNode(node Node, callBack NodeIteratorCallBack) bool {
	if node == nil {
		return false
	}
	if descIterateNode(node.RightChild(), callBack) {
		return true
	}
	if callBack(node) {
		return true
	}
	return descIterateNode(node.LeftChild(), callBack)
}

func ascRangeNode(node Node, lower, upper Key, callBack NodeIteratorCallBack) bool {
	if node == nil {
		return false
	}
	key := node.Key()
	cmpLower := key.CompareTo(lower)
	if 0 < cmpLower {
		if ascRangeNode(node.LeftChild(), lower, upper, callBack) {
			return true
		}
	}
	cmpUpper := key.CompareTo(upper)
	if 0 <= cmpLower && cmpUpper <= 0 {
		if callBack(node) {
			return true
		}
	}
	if cmpUpper < 0 {
		return ascRangeNode(node.RightChild(), lower, upper, callBack)
	} else {
		return false
	}
}

func descRangeNode(node Node, lower, upper Key, callBack NodeIteratorCallBack) bool {
	if node == nil {
		return false
	}
	key := node.Key()
	cmpUpper := key.CompareTo(upper)
	if cmpUpper < 0 {
		if descRangeNode(node.RightChild(), lower, upper, callBack) {
			return true
		}
	}
	cmpLower := key.CompareTo(lower)
	if 0 <= cmpLower && cmpUpper <= 0 {
		if callBack(node) {
			return true
		}
	}
	if 0 < cmpLower {
		return descRangeNode(node.LeftChild(), lower, upper, callBack)
	} else {
		return false
	}
}

func ascExtendedRangeNode(node Node, lower, upper Key, callBack NodeIteratorCallBack) bool {
	if node == nil {
		return false
	}
	key := node.Key()
	cmpLower := key.CompareTo(lower)
	if 0 <= cmpLower {
		if ascExtendedRangeNode(node.LeftChild(), lower, upper, callBack) {
			return true
		}
	}
	cmpUpper := key.CompareTo(upper)
	if 0 <= cmpLower && cmpUpper <= 0 {
		if callBack(node) {
			return true
		}
	}
	if cmpUpper <= 0 {
		return ascExtendedRangeNode(node.RightChild(), lower, upper, callBack)
	} else {
		return false
	}
}

func descExtendedRangeNode(node Node, lower, upper Key, callBack NodeIteratorCallBack) bool {
	if node == nil {
		return false
	}
	key := node.Key()
	cmpUpper := key.CompareTo(upper)
	if cmpUpper <= 0 {
		if descExtendedRangeNode(node.RightChild(), lower, upper, callBack) {
			return true
		}
	}
	cmpLower := key.CompareTo(lower)
	if 0 <= cmpLower && cmpUpper <= 0 {
		if callBack(node) {
			return true
		}
	}
	if 0 <= cmpLower {
		return descExtendedRangeNode(node.LeftChild(), lower, upper, callBack)
	} else {
		return false
	}
}
