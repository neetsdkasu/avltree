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
