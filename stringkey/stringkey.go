// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのKeyの実装例
// string型をそのままキーにしてある
// 標準パッケージのstrings.Compareの結果をそのままキーの比較の値として使っている
package stringkey

import (
	"strings"

	"github.com/neetsdkasu/avltree"
)

type StringKey string

func (key StringKey) CompareTo(other avltree.Key) avltree.KeyOrdering {
	s1 := string(key)
	s2 := string(other.(StringKey))
	return avltree.KeyOrdering(strings.Compare(s1, s2))
}

func (key StringKey) Copy() avltree.Key {
	return key
}
