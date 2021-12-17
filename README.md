# avltree

### 概要

AVL木の構築と操作をするパッケージ？


### 注意点

 - 本パッケージが正しくAVL木として機能しているかは未確認
 - 性能(実行速度やメモリの扱い)が壊滅的なので実用性は皆無
 - キーや値の型の検査などは一切していないので利用者側で統一させる必要がある
 - キーの不変性も利用者側で確保する必要がある
 - サブツリーやノードを操作するための処理は無い


### 利用方法

利用者側でインターフェースの`RealTree`,`RealNode`,`Key`の実装を用意しそれを本パッケージの関数で操作などをする

`RealTree`,`RealNode`の実装例を以下のサブパッケージに置いてある

    github.com/neetsdkasu/avltree/simpletree        最低限の実装のみ
    github.com/neetsdkasu/avltree/standardtree      ノード数の保持や親ノード参照などの機能がある
    github.com/neetsdkasu/avltree/immutabletree     木の構造の部分だけは不変ぽくなるように実装されている(キーと値の不変性は取り扱わない)
    github.com/neetsdkasu/avltree/intarraytree      int型の配列上に木が構築されるように実装(キーはintkeyの実装のみ、値もint型のみ)


`Key`の実装例を以下のサブパッケージに置いてある

    github.com/neetsdkasu/avltree/intkey            int型をKeyとして使えるよう実装
    github.com/neetsdkasu/avltree/stringkey         string型をKeyとして使えるよう実装

木の実装を内包し本パッケージの関数をメソッド経由で呼び出す、所謂"ラッパー"の実装例を以下のサブパッケージに置いてある

    github.com/neetsdkasu/avltree/simplewrapper     簡易に実装したラッパー
    github.com/neetsdkasu/avltree/intwrapper        キーも値もint型に強制するラッパー

コード例
```go

import (
	"fmt"

	"github.com/neetsdkasu/avltree"
	. "github.com/neetsdkasu/avltree/intkey"
	"github.com/neetsdkasu/avltree/simpletree"
)

func Example() {
	tree := simpletree.New(false)
	avltree.Insert(tree, false, IntKey(12), 345)
	avltree.Insert(tree, false, IntKey(67), 890)
	avltree.Insert(tree, false, IntKey(333), 666)
	avltree.Insert(tree, false, IntKey(-5), 12345)
	avltree.Delete(tree, IntKey(67))
	avltree.Update(tree, IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
		newValue = oldValue.(int) * 3
		return
	})
	if node := avltree.Find(tree, IntKey(12)); node != nil {
		fmt.Println("Find!", node.Key(), node.Value())
	}
	avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
}
```

ラッパーを使ったコード例
```go

import (
	"fmt"

	"github.com/neetsdkasu/avltree"
	. "github.com/neetsdkasu/avltree/intkey"
	"github.com/neetsdkasu/avltree/simpletree"
	"github.com/neetsdkasu/avltree/simplewrapper"
)

func Example_wrapper() {
	tree := simpletree.New(false)
	w := simplewrapper.New(tree)
	w.Insert(IntKey(12), 345)
	w.Insert(IntKey(67), 890)
	w.Insert(IntKey(333), 666)
	w.Insert(IntKey(-5), 12345)
	w.Delete(IntKey(67))
	w.Update(IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
		newValue = oldValue.(int) * 3
		return
	})
	if node := w.Find(IntKey(12)); node != nil {
		fmt.Println("Find!", node.Key(), node.Value())
	}
	w.Iterate(func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
}
```
