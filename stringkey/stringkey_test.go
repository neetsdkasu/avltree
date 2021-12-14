package stringkey

import (
	"testing"
	"testing/quick"

	"github.com/neetsdkasu/avltree"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

func TestStringKey(t *testing.T) {
	f := func(k1, k2 string) bool {
		var key1 avltree.Key = StringKey(k1)
		var key2 avltree.Key = StringKey(k2)
		switch key1.CompareTo(key2) {
		case avltree.LessThanOtherKey:
			return k1 < k2
		case avltree.EqualToOtherKey:
			return k1 == k2
		case avltree.GreaterThanOtherKey:
			return k1 > k2
		default:
			return false
		}
	}

	if err := quick.Check(f, cfg1000); err != nil {
		t.Fatal(err)
	}
}
