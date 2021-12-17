// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのKeyの実装例
// int型をそのままキーにしてある
// int型の自然順序をそのままキーの順序としてある
package intkey

import "github.com/neetsdkasu/avltree"

type IntKey int

func (key IntKey) CompareTo(other avltree.Key) avltree.KeyOrdering {
	v1 := int(key)
	v2 := int(other.(IntKey))
	switch {
	case v1 < v2:
		return avltree.LessThanOtherKey
	case v1 > v2:
		return avltree.GreaterThanOtherKey
	default:
		return avltree.EqualToOtherKey
	}
	// return v1 - v2 は 算術オーバーフローがこわい
}

func (key IntKey) Copy() avltree.Key {
	return key
}
