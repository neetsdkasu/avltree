// 概要
// AVL木の構築と操作をするパッケージ？
//
//
// 注意点
// 本パッケージが正しくAVL木として機能しているかは未確認
// 性能(実行速度やメモリの扱い)が壊滅的なので実用性は皆無
// キーや値の型の検査などは一切していないので利用者側で統一させる必要がある
// キーの不変性も利用者側で確保する必要がある
// サブツリーやノードを操作するための処理は無い
//
//
// 利用方法
// 利用者側でインターフェースのRealTree,RealNode,Keyの実装を用意しそれを本パッケージの関数で操作などをする
//
// RealTree,RealNodeの実装例を以下のサブパッケージに置いてある
//  github.com/neetsdkasu/avltree/simpletree        最低限の実装のみ
//  github.com/neetsdkasu/avltree/standardtree      ノード数の保持や親ノード参照などの機能がある
//  github.com/neetsdkasu/avltree/immutabletree     木の構造の部分だけは不変性になるように実装されている(キーと値の不変性は取り扱わない)
//  github.com/neetsdkasu/avltree/intarraytree      int型の配列上に木が構築されるように実装(キーはintkeyの実装のみ、値もint型のみ)
//
// Keyの実装例を以下のサブパッケージに置いてある
//  github.com/neetsdkasu/avltree/intkey            int型をKeyとして使えるよう実装
//  github.com/neetsdkasu/avltree/stringkey         string型をKeyとして使えるよう実装
//
// 木の実装を内包し本パッケージの関数をメソッド経由で呼び出す、所謂"ラッパー"の実装例を以下のサブパッケージにおいてある
//  github.com/neetsdkasu/avltree/simplewrapper     簡易に実装したラッパー
//  github.com/neetsdkasu/avltree/intwrapper        キーも値もint型に強制するラッパー
//
//
// コード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			. "github.com/neetsdkasu/avltree/intkey"
//			"github.com/neetsdkasu/avltree/simpletree"
//		)
//		func Example() {
//			tree := simpletree.New(false)
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
//
// ラッパーを使ったコード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			. "github.com/neetsdkasu/avltree/intkey"
//			"github.com/neetsdkasu/avltree/simpletree"
//			"github.com/neetsdkasu/avltree/simplewrapper"
//		)
//		func Example_wrapper() {
//			tree := simpletree.New(false)
//			w := simplewrapper.New(tree)
//			w.Insert(IntKey(12), 345)
//			w.Insert(IntKey(67), 890)
//			w.Insert(IntKey(333), 666)
//			w.Insert(IntKey(-5), 12345)
//			w.Delete(IntKey(67))
//			w.Update(IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
//				newValue = oldValue.(int) * 3
//				return
//			})
//			if node := w.Find(IntKey(12)); node != nil {
//				fmt.Println("Find!", node.Key(), node.Value())
//			}
//			w.Iterate(func(node avltree.Node) (breakIteration bool) {
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
package avltree

// 木のノードを順に参照するIterate,RangeIterateの引数で渡す
// breakIterationをtrueにしたときにイテレーションを中断する
type IterateCallBack = func(node Node) (breakIteration bool)

// 指定のキーの値の更新をするUpdate,UpdateAll,UpdateRangeの引数で渡す
// newValueに値を代入することでキーに対する新しい値を設定する
// 値を変更しない場合はkeepOldValueをtrueに設定する
type UpdateValueCallBack = func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool)

// 木のノードを順に巡りノードの値を更新するUpdateIterate,UpdateRangeIterateの引数で渡す
// newValueに値を代入することでキーに対する新しい値を設定する
// 値を変更しない場合はkeepOldValueをtrueに設定する
// breakIterationをtrueにしたときにイテレーションを中断する
type UpdateIterateCallBack = func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool)

// 木のノードを順に巡りノードの削除判定をするDeleteIterate,DeleteRangeIterateの引数で渡す
// deleteNodeをtrueに設定すると当該ノードを削除する
// breakIterationをtrueにしたときにイテレーションを中断する
type DeleteIterateCallBack = func(key Key, value interface{}) (deleteNode, breakIteration bool)

// 指定のキーの値の更新またはノードを削除をするAlter,AlterAll,AlterRangeの引数で渡す
// return node.Keep()またはrequest.Keep()でノードに変更を加えない指示となる
// return node.Replace(newValue)またはrequest.Replace(newValue)で値をnewValueに更新する指示となる
// return node.Delete()またはrequest.Delete()でノードを削除をする指示となる
type AlterNodeCallBack = func(node AlterNode) (request AlterRequest)

// 木のノードを順に巡りノードの値の更新またはノードを削除をするAlterIterate,AlterRangeIterateの引数で渡す
// return node.Keep()またはrequest.Keep()でノードに変更を加えない指示となる
// return node.Replace(newValue)またはrequest.Replace(newValue)で値をnewValueに更新する指示となる
// return node.Delete()またはrequest.Delete()でノードを削除をする指示となる
// breakIterationをtrueにしたときにイテレーションを中断する
type AlterIterateCallBack = func(node AlterNode) (request AlterRequest, breakIteration bool)

// 木またはサブツリーのルートノードがノード数を保持する実装でその値を公開するためのメソッド
// このインターフェースが実装されている場合にCount,CountRangeの内部でNodeCountメソッドが呼び出される
type NodeCounter interface{ NodeCount() int }

// 子ノードが親のノードを参照できる実装でその参照を公開するためのメソッド
// 本パッケージ内で使用される予定は今のところ無い…
type ParentGetter interface {
	Node
	Parent() Node
}

// 木がこのインターフェースを実装している場合にDeleteなどの内部でRelaseNodeメソッドが呼び出される
// ノードが木から切り離された際にそのノードのインスタンスに対して何らかの処理をしたい場合に木側で実装する必要がある
// 例えばデータの切り離しやインスタンスの再利用やメモリの再利用・解放など(ノード本体のインスタンスやメモリのみの処理が実行されることが期待される)
// インターフェースRealNodeのNewNodeメソッドでのノード生成と対になる処理になることが期待される
// サブパッケージのintarraytreeではReleaseNodeでメモリの再利用のための処理を行っている
type NodeReleaser interface {
	RealTree
	ReleaseNode(node RealNode)
}

// 木のメモリの解放処理を行うためのメソッド(？)
// ReleaseTreeメソッドが呼び出された場合に木の再利用が不可能な状態であることが期待される(panicを起こすなど)
// このインターフェースが実装されている場合にReleaseの内部でReleaseTreeメソッドを呼び出す
type TreeReleaser interface {
	Tree
	ReleaseTree()
}

// このインターフェースが実装されている場合にClearの内部で木のノードを全て削除したあとでCleanUpTreeメソッドを呼び出す
// Delete系やAlter系で木のノードが全部削除されてもこのメソッドは呼び出されずClearのみで呼び出される
// 利用シーンは全く思い至らないのに何故か用意してしまった謎インターフェース
// サブパッケージのintarraytreeではCleanUpTreeでメモリ再利用のためのノードのチェーンを破棄するための処理を行っている
type TreeCleaner interface {
	Tree
	CleanUpTree()
}

// 木の公開用の基本的なインターフェース
// デフォルトでアクセスできる範囲を制限するためだけの用途
type Tree interface {
	// 木のルートノードを返却されることが期待される
	// 返却されるNodeはRealNodeを実装している必要がある
	// 木のルートノードが無い場合はreturn nilでnilを返却する必要がある
	Root() Node
}

// 本パッケージで操作するために木が実装するべきメソッド
type RealTree interface {
	Tree

	// Insertなどで木に新しいノードを作る場合に呼び出される
	// 引数で渡された情報を持つノードが返却されることが期待される
	// leftChildで対応するノードが無い場合は引数にnilが渡される
	// rightChildで対応するノードが無い場合は引数にnilが渡される
	NewNode(leftChild, rightChild Node, height int, key Key, value interface{}) RealNode

	// 木のルートノードを設定する場合に呼び出される
	// 木にノードが全く無くなった場合はnewRootにnilが渡される
	// 可変(mutable)の木の場合にはレシーバである木自身のインスタンスが返却されることが期待される
	// 不変(immutable)の木の場合にはレシーバと同じ設定を持ちnewRootのルートノードを持つ新しいインスタンスの木が返却されることが期待される
	SetRoot(newRoot RealNode) RealTree

	// 木が同一のキーのノードを複数保持できるかどうかを返す
	// trueの場合、木が同一キーのノードを複数保持することを許可することを表す
	// falseの場合、木が同一キーのノードを複数保持することを許可しないことを表す
	// 同一キーを許可する場合、キー指定操作(Findなど)において同一キーのノードのどのノードが取得されるかは不定である
	AllowDuplicateKeys() bool
}

type KeyAndValue interface {
	// ノードに設定されたキーを返却する
	Key() Key

	// ノードに設定された値を返却する
	Value() interface{}
}

// ノードの公開用の基本的なインターフェース
// デフォルトでアクセスできる範囲を制限するためだけの用途
type Node interface {
	KeyAndValue

	// このノードの左の子ノードが返却されることが期待される
	// 返却されるNodeはRealNodeを実装している必要がある
	// 左の子ノードが無い場合はreturn nilでnilを返却する必要がある
	LeftChild() Node

	// このノードの右の子のノードが返却されることが期待される
	// 返却されるNodeはRealNodeを実装している必要がある
	// 右の子ノードが無い場合はreturn nilでnilを返却する必要がある
	RightChild() Node

	// このノードに新しい値を設定する
	// 可変(mutable)の木の場合にはレシーバであるノード自身のインスタンスが返却されることが期待される
	// 不変(immutable)の木の場合にはレシーバと同じ設定を持ちnewValueの値を持つ新しいインスタンスのノードが返却されることが期待される
	// 返却されるNodeはRealNodeを実装している必要がある
	SetValue(newValue interface{}) Node
}

// 本パッケージで操作するためにノードが実装するべきメソッド
type RealNode interface {
	Node

	// このノードの高さを返却する必要がある
	// 高さはRealTreeのNewNodeメソッドまたはRealNodeのSetChildre,Setメソッドで設定された値
	Height() int

	// このノードの新しい高さとこのノードの左右の子を設定する
	// 引数のnewHeightは計算後の値なのでノードはそのままの値を保持する必要がある
	// newLeftChildで対応するノードが無い場合は引数にnilが渡される
	// newRightChildで対応するノードが無い場合は引数にnilが渡される
	// 可変(mutable)の木の場合にはレシーバであるノード自身のインスタンスが返却されることが期待される
	// 不変(immutable)の木の場合にはレシーバと同じ設定を持ち渡された引数の情報を持つ新しいインスタンスのノードが返却されることが期待される
	SetChildren(newLeftChild, newRightChild Node, newHeight int) RealNode

	// このノードの新しい値とこのノードの新しい高さとこのノードの左右の子を設定する
	// 引数のnewHeightは計算後の値なのでノードはそのままの値を保持する必要がある
	// newLeftChildで対応するノードが無い場合は引数にnilが渡される
	// newRightChildで対応するノードが無い場合は引数にnilが渡される
	// 可変(mutable)の木の場合にはレシーバであるノード自身のインスタンスが返却されることが期待される
	// 不変(immutable)の木の場合にはレシーバと同じ設定を持ち渡された引数の情報を持つ新しいインスタンスのノードが返却されることが期待される
	Set(newLeftChild, newRightChild Node, newHeight int, newValue interface{}) RealNode
}

// キーの比較結果を表す型
type KeyOrdering int

const (
	// CompareToメソッドの引数のotherよりレシーバのほうが小さい
	LessThanOtherKey KeyOrdering = -1

	// CompareToメソッドの引数のotherとレシーバは同じキー
	EqualToOtherKey KeyOrdering = 0

	// CompareToメソッドの引数のotherよりレシーバのほうが大きい
	GreaterThanOtherKey KeyOrdering = 1
)

// AVL木でノードの比較に使うキーのインターフェース
type Key interface {
	// キーの順序(比較結果)を返却する必要がある
	CompareTo(other Key) KeyOrdering

	// キーのコピーが返却されることが期待される
	// RealTreeのNewNodeメソッドの引数に渡すときにこのCopyメソッドが呼ばれる
	Copy() Key
}

// 木に指定のキーと値を持つ新しいノードを追加する
// 既に指定キーが木に存在する場合は
// 木が同一キーのノードを許可しておらずreplaceIfExistsがfalseのときは木にノードを追加しない
// 木が同一キーのノードを許可しておらずreplaceIfExistsがtrueのときは既に存在するノードの値をvalueに置き換える
// 木が同一キーのノードを許可しておりreplaceIfExistsがfalseのときは木に新たにノードを追加する
// 木が同一キーのノードを許可しておりreplaceIfExistsがtrueのときは既に存在するノードのうち最初に見つかったノード(同一キーのノードのうち高さが最も高いノード)の値をvalueに置き換える
// 戻り値のmodifiedは木に変更があった場合はRealTreeのSetRootメソッドの戻り値となり、変更がない場合は引数のtreeがそのまま返却される
// 戻り値のokは木へのノード追加あるいはノードの値の更新があった場合にtrue、木に変化がなかった場合はfalseを返却する
func Insert(tree Tree, replaceIfExists bool, key Key, value interface{}) (modified Tree, ok bool) {
	realTree := tree.(RealTree)
	helper := insertHelper{
		&realTree,
		replaceIfExists,
		&key,
		&value,
	}
	if newRoot, ok := helper.insertTo(tree.Root()); ok {
		return realTree.SetRoot(newRoot), true
	} else {
		return tree, false
	}
}

// 指定のキーを持つノードを木から削除する(木から取り除く)
// 指定のキーが存在しない場合は木に変更は加えない
// 木が同一キーのノードを許可している場合は指定キーを持つノードのうち最初に見つかったノード(同一キーのノードのうち高さが最も高いノード)が削除される
// 戻り値のmodifiedは木に変更があった場合はRealTreeのSetRootメソッドの戻り値となり、変更がない場合は引数のtreeがそのまま返却される
// 戻り値のdeletedValueは削除したノードのキーと値を持っている
func Delete(tree Tree, key Key) (modified Tree, deleteValue KeyAndValue) {
	if newRoot, node, ok := removeNode(tree.Root(), key); ok {
		deleteValue = &keyAndValue{
			node.Key(),
			node.Value(),
		}
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
		realTree := tree.(RealTree)
		if root, ok := newRoot.(RealNode); ok {
			return realTree.SetRoot(root), deleteValue
		} else {
			return realTree.SetRoot(nil), deleteValue
		}
	} else {
		return tree, nil
	}
}

// 指定のキーを持つノードの値を変更する
// 指定のキーが存在しない場合は木に変更は加えない
// 木が同一キーのノードを許可している場合は指定キーを持つノードのうち最初に見つかったノード(同一キーのノードのうち高さが最も高いノード)が変更の対象となる
// 戻り値のmodifiedは変更があった場合はRealTreeのSetRootメソッドの戻り値となり、変更がない場合は引数のtreeがそのまま返却される
// 戻り値のokは変更があった場合にtrueとなる
// 変更があった場合とは、指定のキーを持つノードが存在し、かつコールバックの戻り値keepOldValueがfalseであったときのことを指す
// 変更が無かった場合とは、指定のキーを持つノードが存在しなかった場合もしくはコールバックの戻り値keepOldValueがtrueであった場合
func Update(tree Tree, key Key, callBack UpdateValueCallBack) (modified Tree, ok bool) {
	if newRoot, ok := updateValue(tree.Root(), key, callBack); ok {
		return tree.(RealTree).SetRoot(newRoot), true
	} else {
		return tree, false
	}
}

// 指定のキーを持つノードの値を別の値に置き換える
// 指定のキーが存在しない場合は木に変更は加えない
// 木が同一キーのノードを許可している場合は指定キーを持つノードのうち最初に見つかったノード(同一キーのノードのうち高さが最も高いノード)が変更の対象となる
// 戻り値のmodifiedは指定したキーが存在する場合はRealTreeのSetRootメソッドの戻り値となり、キーが存在しない場合は引数のtreeがそのまま返却される
// 戻り値のokは指定のキーを持つノードが存在する場合にtrueとなる
func Replace(tree Tree, key Key, newValue interface{}) (modified Tree, ok bool) {
	return Update(tree, key, func(key Key, oldValue interface{}) (interface{}, bool) {
		return newValue, false
	})
}

// Alterのコールバックの引数で受け取るノード情報
// コールバックの戻り値のAlterRequestを生成できるメソッドも実装している
type AlterNode interface {
	KeyAndValue

	// ノードに対して変更や削除の操作をしないことを指定するAlterRequestを生成する
	Keep() AlterRequest

	// ノードの値を指定した値へ置き換えることを指定するAlterRequestを生成する
	Replace(newValue interface{}) AlterRequest

	// ノードを削除することを指定するAlterRequestを生成する
	Delete() AlterRequest
}

// AlterNodeの実体
type alterNode struct {
	inner Node
}

// Alterのコールバックの戻り値で使用する
// ノードに対する処理の情報を保持する
type AlterRequest struct {
	replaceValue bool
	newValue     interface{}
	deleteNode   bool
}

// 指定したキーを持つノードの値を変更またはノードを削除する
// 指定のキーが存在しない場合は木に変更は加えない
// 木が同一キーのノードを許可している場合は指定キーを持つノードのうち最初に見つかったノード(同一キーのノードのうち高さが最も高いノード)が変更の対象となる
// 戻り値のmodifiedは対象のノードが存在しコールバックの戻り値で変更か削除を指定された場合はRealTreeのSetRootメソッドの戻り値となり、それ以外の場合は引数のtreeがそのまま返却される
// 戻り値のdeletedValueは削除したノードのキーと値を持っている
// 戻り値のokは対象のノードが存在しコールバックの戻り値で変更か削除を指定された場合にtrueとなり、それ以外の場合はfalseとなる
func Alter(tree Tree, key Key, callBack AlterNodeCallBack) (modified Tree, deletedValue KeyAndValue, ok bool) {
	if newRoot, deleted, ok := alter(tree.Root(), key, callBack); ok {
		if deleted != nil {
			deletedValue = &keyAndValue{
				deleted.Key(),
				deleted.Value(),
			}
			if releaser, ok := tree.(NodeReleaser); ok {
				releaser.ReleaseNode(deleted.(RealNode))
			}
		}
		if root, ok := newRoot.(RealNode); ok {
			return tree.(RealTree).SetRoot(root), deletedValue, true
		} else {
			return tree.(RealTree).SetRoot(nil), deletedValue, true
		}
	} else {
		return tree, nil, false
	}
}

// 木のインスタンスを再利用するために木が保持するノードを全て削除し空にする
// 戻り値のmodifiedはRealTreeのSetRootメソッドの戻り値となる
func Clear(tree Tree) (modified Tree) {
	if releaser, ok := tree.(NodeReleaser); ok {
		stack := []Node{tree.Root()}
		for 0 < len(stack) {
			newsize := len(stack) - 1
			node := stack[newsize]
			stack = stack[:newsize]
			if node == nil {
				continue
			}
			stack = append(stack, node.RightChild(), node.LeftChild())
			releaser.ReleaseNode(node.(RealNode))
		}
	}
	if cleaner, ok := tree.(TreeCleaner); ok {
		cleaner.CleanUpTree()
	}
	return tree.(RealTree).SetRoot(nil)
}

// 木のインスタンスを解放(または破棄)する
// 木がインターフェースTreeReleaserを実装している場合にのみ機能する
// 解放の全ての処理が終わったあと引数のtreeの参照先にはnilが代入される
func Release(tree *Tree) {
	Clear(*tree)
	if releaser, ok := (*tree).(TreeReleaser); ok {
		releaser.ReleaseTree()
	}
	*tree = nil
}

// 指定のキーを持つノードを取得する
// 指定のキーを持つノードが無い場合は戻り値nodeはnilになる
// 戻り値nodeは木の一部のままなのでこのnodeを編集すると木にも影響する
// 可変(mutable)の木の場合にはnodeの内容を変更する操作は木にも影響する
// 不変(immutable)の木の場合にはnodeのインターフェスNodeやRealNodeのメソッド呼び出しでは木に影響がないことが期待される
func Find(tree Tree, key Key) (node Node) {
	node = tree.Root()
	for node != nil {
		cmp := key.CompareTo(node.Key())
		switch {
		case cmp.LessThan():
			node = node.LeftChild()
		case cmp.GreaterThan():
			node = node.RightChild()
		default:
			return node
		}
	}
	return nil
}

// 木のノードを順に巡ってコールバックを呼び出す
// descOrderがfalseのときはキーの昇順
// descOrderがtrueのときはキーの降順
// 戻り値のokはコールバックから中断を要求されなかった場合はtrue、中断を要求された場合はfalse
func Iterate(tree Tree, descOrder bool, callBack IterateCallBack) (ok bool) {
	if descOrder {
		return descIterateNode(tree.Root(), callBack)
	} else {
		return ascIterateNode(tree.Root(), callBack)
	}
}

// キーの順序でlower以上upper以下の範囲にあるノード全てを指定の順序で取得する
// descOrderがfalseのときはキーの昇順
// descOrderがtrueのときはキーの降順
func Range(tree Tree, descOrder bool, lower, upper Key) (nodes []Node) {
	RangeIterate(tree, descOrder, lower, upper, func(node Node) (breakIteration bool) {
		nodes = append(nodes, node)
		return
	})
	return
}

// 木のノードをキーの順序でlower以上upper以下の範囲を順に巡ってコールバックを呼び出す
// descOrderがfalseのときはキーの昇順
// descOrderがtrueのときはキーの降順
// 戻り値のokはコールバックから中断を要求されなかった場合はtrue、中断を要求された場合はfalse
func RangeIterate(tree Tree, descOrder bool, lower, upper Key, callBack IterateCallBack) (ok bool) {
	if lower == nil && upper == nil {
		return Iterate(tree, descOrder, callBack)
	}
	bounds := newKeyBounds(lower, upper, tree.(RealTree).AllowDuplicateKeys())
	if descOrder {
		return descRangeNode(tree.Root(), bounds, callBack)
	} else {
		return ascRangeNode(tree.Root(), bounds, callBack)
	}
}

// 木を構成するノードの総数を求める
// 木またはノードがインターフェースNodeCounterを実装している場合はNodeCountメソッドの戻り値をそのまま返す
func Count(tree Tree) int {
	if counter, ok := tree.(NodeCounter); ok {
		return counter.NodeCount()
	} else {
		return countNode(tree.Root())
	}
}

// キーの順序でlower以上upper以下の範囲にあるノード総数を求める
func CountRange(tree Tree, lower, upper Key) int {
	if lower == nil && upper == nil {
		return Count(tree)
	}
	if tree.(RealTree).AllowDuplicateKeys() {
		return countExtendedRange(tree.Root(), lower, upper)
	} else {
		return countRange(tree.Root(), lower, upper)
	}
}

// キーの順序で最小となるキーを持つノードを返す
// 木にひとつもノードが無い場合はnilを返す
// 戻り値は木の一部のままなので編集すると木にも影響する
func Min(tree Tree) (minimum Node) {
	node := tree.Root()
	if node == nil {
		return nil
	}
	for {
		leftChild := node.LeftChild()
		if leftChild == nil {
			break
		}
		node = leftChild
	}
	return node
}

// キーの順序で最大となるキーを持つノードを返す
// 木にひとつもノードが無い場合はnilを返す
// 戻り値は木の一部のままなので編集すると木にも影響する
func Max(tree Tree) (maximum Node) {
	node := tree.Root()
	if node == nil {
		return nil
	}
	for {
		rightChild := node.RightChild()
		if rightChild == nil {
			break
		}
		node = rightChild
	}
	return node
}

// 同一キーを持つノードが複数ある場合のDeleteの強化版
// 指定したキーを持つ全てのノードを木から削除する
func DeleteAll(tree Tree, key Key) (modified Tree, deletedvalues []KeyAndValue) {
	return DeleteRange(tree, false, key, key)
}

// 同一キーを持つノードが複数ある場合のUpdateの強化版
// 指定したキーを持つ全てのノードに対してコールバックが呼び出される
func UpdateAll(tree Tree, key Key, callBack UpdateValueCallBack) (modified Tree, ok bool) {
	return UpdateRangeIterate(tree, false, key, key, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
		newValue, keepOldValue = callBack(key, oldValue)
		return
	})
}

// 同一キーを持つノードが複数ある場合のReplaceの強化版
// 指定したキーを持つ全てのノードの値をnewValueに変更する
func ReplaceAll(tree Tree, key Key, newValue interface{}) (modified Tree, ok bool) {
	return UpdateRangeIterate(tree, false, key, key, func(key Key, oldValue interface{}) (interface{}, bool, bool) {
		return newValue, false, false
	})
}

// 同一キーを持つノードが複数ある場合のAlterの強化版
// 指定したキーを持つ全てのノードに対してコールバックが呼び出される
func AlterAll(tree Tree, key Key, callBack AlterNodeCallBack) (modified Tree, deletedValues []KeyAndValue, ok bool) {
	return AlterRangeIterate(tree, false, key, key, func(node AlterNode) (request AlterRequest, breakIteration bool) {
		return callBack(node), false
	})
}

// 同一キーを持つノードが複数ある場合のFindの強化版
// 指定したキーを持つ全てのノードを取得する
// 指定のキーが存在しない場合はnilを返す
func FindAll(tree Tree, key Key) (nodes []Node) {
	// FindAllを頻繁に呼び出すでもない限りは
	// Range呼び出しのオーバーヘッドなんて気にするほどのものではないはず
	RangeIterate(tree, false, key, key, func(node Node) (breakIteration bool) {
		nodes = append(nodes, node)
		return
	})
	return nodes
}

// 同一キーを持つノードが複数ある場合のMinの強化版
// 最小のキーを持つ全てのノードを取得する
// 木にノードが存在しない場合はnilを返す
func MinAll(tree Tree) (minimums []Node) {
	minimum := Min(tree)
	if minimum == nil {
		return nil
	}
	key := minimum.Key()
	RangeIterate(tree, false, key, key, func(node Node) (breakIteration bool) {
		minimums = append(minimums, node)
		return
	})
	return minimums
}

// 同一キーを持つノードが複数ある場合のMaxの強化版
// 最大のキーを持つ全てのノードを取得する
// 木にノードが存在しない場合はnilを返す
func MaxAll(tree Tree) (maximums []Node) {
	maximum := Max(tree)
	if maximum == nil {
		return nil
	}
	key := maximum.Key()
	// 最大側は descOrder にすべきか…？
	// FindAllに揃えるなら ascOrder だが
	RangeIterate(tree, false, key, key, func(node Node) (breakIteration bool) {
		maximums = append(maximums, node)
		return
	})
	return maximums
}

// 戻り値で使うKeyAndValueの実体
type keyAndValue struct {
	key   Key
	value interface{}
}

func (kv *keyAndValue) Key() Key {
	return kv.key
}

func (kv *keyAndValue) Value() interface{} {
	return kv.value
}

// DeleteとIterateを組み合わせた感じ
// コールバックの戻り値で削除対象となるノードを決定していく
func DeleteIterate(tree Tree, descOrder bool, callBack DeleteIterateCallBack) (modified Tree, values []KeyAndValue) {
	var newRoot Node
	var deleted []Node
	if descOrder {
		newRoot, deleted, _ = descDeleteIterate(tree.Root(), callBack)
	} else {
		newRoot, deleted, _ = ascDeleteIterate(tree.Root(), callBack)
	}
	if len(deleted) == 0 {
		return tree, nil
	}
	for _, node := range deleted {
		values = append(values, &keyAndValue{
			node.Key(),
			node.Value(),
		})
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
	}
	if root, ok := newRoot.(RealNode); ok {
		return tree.(RealTree).SetRoot(root), values
	} else {
		return tree.(RealTree).SetRoot(nil), values
	}
}

// DeleteとRangeを組み合わせた感じ
// lower以上upper以下のキーを持つノード全てを削除する
func DeleteRange(tree Tree, descOrder bool, lower, upper Key) (modified Tree, values []KeyAndValue) {
	return DeleteRangeIterate(tree, descOrder, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
		deleteNode = true
		return
	})
}

// DeleteとRangeIterateを組み合わせた感じ
// lower以上upper以下のキーを持つノードに対して順にコールバックが呼ばれる
// コールバックの戻り値で削除対象となるノードを決定していく
func DeleteRangeIterate(tree Tree, descOrder bool, lower, upper Key, callBack DeleteIterateCallBack) (modified Tree, values []KeyAndValue) {
	if lower == nil && upper == nil {
		return DeleteIterate(tree, descOrder, callBack)
	}
	var newRoot Node
	var deleted []Node
	bounds := newKeyBounds(lower, upper, tree.(RealTree).AllowDuplicateKeys())
	if descOrder {
		newRoot, deleted, _ = descDeleteRange(tree.Root(), bounds, callBack)
	} else {
		newRoot, deleted, _ = ascDeleteRange(tree.Root(), bounds, callBack)
	}
	if len(deleted) == 0 {
		return tree, nil
	}
	for _, node := range deleted {
		values = append(values, &keyAndValue{
			node.Key(),
			node.Value(),
		})
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
	}
	if root, ok := newRoot.(RealNode); ok {
		return tree.(RealTree).SetRoot(root), values
	} else {
		return tree.(RealTree).SetRoot(nil), values
	}
}

// UpdateとIterateを組み合わせた感じ
// 巡っていく各ノードに対してコールバックが呼ばれる
func UpdateIterate(tree Tree, descOrder bool, callBack UpdateIterateCallBack) (modified Tree, ok bool) {
	if descOrder {
		if newRoot, updated, _ := descUpdateIterate(tree.Root(), callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	} else {
		if newRoot, updated, _ := ascUpdateIterate(tree.Root(), callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	}
}

// UpdateとRangeを組み合わせた感じ
// lower以上upper以下のキーを持つノードそれぞれに順にコールバックが呼ばれる
func UpdateRange(tree Tree, descOrder bool, lower, upper Key, callBack UpdateValueCallBack) (modified Tree, ok bool) {
	return UpdateRangeIterate(tree, descOrder, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
		newValue, keepOldValue = callBack(key, oldValue)
		return
	})
}

// UpdateとRangeIterateを組み合わせた感じ
// lower以上upper以下のキーを持つノードに順にコールバックが呼ばれる
func UpdateRangeIterate(tree Tree, descOrder bool, lower, upper Key, callBack UpdateIterateCallBack) (modified Tree, ok bool) {
	if lower == nil && upper == nil {
		return UpdateIterate(tree, descOrder, callBack)
	}
	bounds := newKeyBounds(lower, upper, tree.(RealTree).AllowDuplicateKeys())
	if descOrder {
		if newRoot, updated, _ := descUpdateRange(tree.Root(), bounds, callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	} else {
		if newRoot, updated, _ := ascUpdateRange(tree.Root(), bounds, callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	}
}

// ReplaceとRangeを組み合わせた感じ
// lower以上upper以下のキーを持つ全てのノードの値をnewValueに変更する
func ReplaceRange(tree Tree, lower, upper Key, newValue interface{}) (modified Tree, ok bool) {
	return UpdateRangeIterate(tree, false, lower, upper, func(key Key, oldValue interface{}) (interface{}, bool, bool) {
		return newValue, false, false
	})
}

// AlterとIterateを組み合わせた感じ
// 巡っていく各ノードに対してコールバックが呼ばれる
func AlterIterate(tree Tree, descOrder bool, callBack AlterIterateCallBack) (modified Tree, deletedValues []KeyAndValue, ok bool) {
	var newRoot Node
	var deleted []Node
	var anyChanged bool
	if descOrder {
		newRoot, deleted, anyChanged, _ = descAlterIterate(tree.Root(), callBack)
	} else {
		newRoot, deleted, anyChanged, _ = ascAlterIterate(tree.Root(), callBack)
	}
	if !anyChanged {
		return tree, nil, false
	}
	for _, node := range deleted {
		kv := &keyAndValue{node.Key(), node.Value()}
		deletedValues = append(deletedValues, kv)
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
	}
	if root, ok := newRoot.(RealNode); ok {
		return tree.(RealTree).SetRoot(root), deletedValues, true
	} else {
		return tree.(RealTree).SetRoot(nil), deletedValues, true
	}
}

// AlterとRangeを組み合わせた感じ
// lower以上upper以下のキーを持つノードそれぞれに順にコールバックが呼ばれる
func AlterRange(tree Tree, descOrder bool, lower, upper Key, callBack AlterNodeCallBack) (modified Tree, deletedValues []KeyAndValue, ok bool) {
	return AlterRangeIterate(tree, descOrder, lower, upper, func(node AlterNode) (request AlterRequest, breakIteration bool) {
		return callBack(node), false
	})
}

// UpdateとRangeIterateを組み合わせた感じ
// lower以上upper以下のキーを持つノードに順にコールバックが呼ばれる
func AlterRangeIterate(tree Tree, descOrder bool, lower, upper Key, callBack AlterIterateCallBack) (modified Tree, deletedValues []KeyAndValue, ok bool) {
	if lower == nil && upper == nil {
		return AlterIterate(tree, descOrder, callBack)
	}
	var newRoot Node
	var deleted []Node
	var anyChanged bool
	bounds := newKeyBounds(lower, upper, tree.(RealTree).AllowDuplicateKeys())
	if descOrder {
		newRoot, deleted, anyChanged, _ = descAlterRange(tree.Root(), bounds, callBack)
	} else {
		newRoot, deleted, anyChanged, _ = ascAlterRange(tree.Root(), bounds, callBack)
	}
	if !anyChanged {
		return tree, nil, false
	}
	for _, node := range deleted {
		deletedValues = append(deletedValues, &keyAndValue{
			node.Key(),
			node.Value(),
		})
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
	}
	if root, ok := newRoot.(RealNode); ok {
		return tree.(RealTree).SetRoot(root), deletedValues, true
	} else {
		return tree.(RealTree).SetRoot(nil), deletedValues, true
	}
}

// キーが比較対象(CompareToの引数のキー)より小さい場合にtrue、それ以外はfalse
func (ordering KeyOrdering) LessThan() bool {
	return int(ordering) < 0
}

// キーが比較対象(CompareToの引数のキー)以下の場合にtrue、それ以外はfalse
func (ordering KeyOrdering) LessThanOrEqualTo() bool {
	return int(ordering) <= 0
}

// キーが比較対象(CompareToの引数のキー)と等しい場合にtrue、それ以外はfalse
func (ordering KeyOrdering) EqualTo() bool {
	return ordering == EqualToOtherKey
}

// キーが比較対象(CompareToの引数のキー)と等しくない場合にtrue、それ以外はfalse
func (ordering KeyOrdering) NotEqualTo() bool {
	return ordering != EqualToOtherKey
}

// キーが比較対象(CompareToの引数のキー)より大きい場合にtrue、それ以外はfalse
func (ordering KeyOrdering) GreaterThan() bool {
	return 0 < int(ordering)
}

// キーが比較対象(CompareToの引数のキー)以上の場合にtrue、それ以外はfalse
func (ordering KeyOrdering) GreaterThanOrEqualTo() bool {
	return 0 <= int(ordering)
}

// LessThanメソッドとLessThanOrEqualToメソッドを引数のorEqualで使い分ける
// orEqualがtrueの場合、LessThanOrEqualToが呼ばれる
// orEqualがfalseの場合、LessThanが呼ばれる
func (ordering KeyOrdering) Less(orEqual bool) bool {
	if orEqual {
		return ordering.LessThanOrEqualTo()
	} else {
		return ordering.LessThan()
	}
}

// GreaterThanメソッドとGreaterThanOrEqualToメソッドを引数のorEqualで使い分ける
// orEqualがtrueの場合、GreaterThanOrEqualToが呼ばれる
// orEqualがfalseの場合、GreaterThanが呼ばれる
func (ordering KeyOrdering) Greater(orEqual bool) bool {
	if orEqual {
		return ordering.GreaterThanOrEqualTo()
	} else {
		return ordering.GreaterThan()
	}
}

// 引数のノードをルートとする木(もしくはサブツリー)のノード総数を数える
func countNode(node Node) int {
	if node == nil {
		return 0
	}
	// countNodeの再帰での呼び出し頻度を考えると
	// この判定は無駄が多く、コストが高くつく、かも
	// Treeを構成する一部のNodeにだけNodeCounterが実装されている可能性は低いと思う
	if counter, ok := node.(NodeCounter); ok {
		return counter.NodeCount()
	}
	return 1 + countNode(node.LeftChild()) + countNode(node.RightChild())
}

// 引数のノードをルートとする木(もしくはサブツリー)のキーがlower以上upper以下のノード総数を数える
// countRangeは同一キーを許可しない木に対して呼ばれる
func countRange(node Node, lower, upper Key) int {
	if node == nil {
		return 0
	}
	// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
	// lower == nil, upper < key    ... leftChild
	// lower == nil, key == upper   ... all(leftChild) key
	// lower == nil, key < upper    ... all(leftChild) key rightChild
	// upper == nil, lower < key    ... leftChild key all(rightChild)
	// upper == nil, lower == key   ... key all(rightChild)
	// upper == nil, key < lower    ... rightChild
	// lower < upper < key      ... leftChild
	// lower < upper == key     ... leftChild{upper->nil} key
	// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
	// key == lower < upper     ... key rightChild{lower->nil}
	// key < lower < upper      ... rightChild
	// lower == key == upper    ... key
	if lower == nil {
		if upper == nil {
			// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
			return countNode(node)
		}
		cmp := node.Key().CompareTo(upper)
		switch {
		case cmp.GreaterThan():
			// lower == nil, upper < key    ... leftChild
			return countRange(node.LeftChild(), lower, upper)
		case cmp.EqualTo():
			// lower == nil, key == upper   ... all(leftChild) key
			return countNode(node.LeftChild()) + 1
		case cmp.LessThan():
			// lower == nil, key < upper    ... all(leftChild) key rightChild
			return countNode(node.LeftChild()) + 1 + countRange(node.RightChild(), lower, upper)
		default:
			// ここには到達しないはず
			panic("unreachable?")
		}
	}
	if upper == nil {
		cmp := node.Key().CompareTo(lower)
		switch {
		case cmp.GreaterThan():
			// upper == nil, lower < key    ... leftChild key all(rightChild)
			return countRange(node.LeftChild(), lower, upper) + 1 + countNode(node.RightChild())
		case cmp.EqualTo():
			// upper == nil, lower == key   ... key all(rightChild)
			return 1 + countNode(node.RightChild())
		case cmp.LessThan():
			// upper == nil, key < lower    ... rightChild
			return countRange(node.RightChild(), lower, upper)
		default:
			// ここには到達しないはず
			panic("unreachable?")
		}
	}
	key := node.Key()
	cmpLower := key.CompareTo(lower)
	cmpUpper := key.CompareTo(upper)
	switch {
	case cmpUpper.GreaterThan():
		// lower < upper < key      ... leftChild
		return countRange(node.LeftChild(), lower, upper)
	case cmpLower.GreaterThan() && cmpUpper.EqualTo():
		// lower < upper == key     ... leftChild{upper->nil} key
		return countRange(node.LeftChild(), lower, nil) + 1
	case cmpLower.GreaterThan() && cmpUpper.LessThan():
		// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
		return countRange(node.LeftChild(), lower, nil) + 1 + countRange(node.RightChild(), nil, upper)
	case cmpLower.EqualTo() && cmpUpper.LessThan():
		// key == lower < upper     ... key rightChild{lower->nil}
		return 1 + countRange(node.RightChild(), nil, upper)
	case cmpLower.LessThan():
		// key < lower < upper      ... rightChild
		return countRange(node.RightChild(), lower, upper)
	case cmpLower.EqualTo() && cmpUpper.EqualTo():
		// lower == key == upper    ... key
		return 1
	}
	// 条件漏れが無ければ、ここには到達しないと思う…
	panic("unreachable?")
}

// 引数のノードをルートとする木(もしくはサブツリー)のキーがlower以上upper以下のノード総数を数える
// countExtendedRangeは同一キーを許可する木に対して呼ばれる
func countExtendedRange(node Node, lower, upper Key) int {
	if node == nil {
		return 0
	}
	// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
	// lower == nil, upper < key    ... leftChild
	// lower == nil, key == upper   ... all(leftChild) key rightChild
	// lower == nil, key < upper    ... all(leftChild) key rightChild
	// upper == nil, lower < key    ... leftChild key all(rightChild)
	// upper == nil, lower == key   ... leftChild key all(rightChild)
	// upper == nil, key < lower    ... rightChild
	// lower < upper < key      ... leftChild
	// lower < upper == key     ... leftChild{upper->nil} key rightChild{lower->nil}
	// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
	// key == lower < upper     ... leftChild{upper->nil} key rightChild{lower->nil}
	// key < lower < upper      ... rightChild
	// lower == key == upper    ... leftChild{upper->nil} key rightChild{lower->nil}
	if lower == nil {
		if upper == nil {
			// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
			return countNode(node)
		}
		if node.Key().CompareTo(upper).GreaterThan() {
			// lower == nil, upper < key    ... leftChild
			return countExtendedRange(node.LeftChild(), lower, upper)
		} else {
			// lower == nil, key == upper   ... all(leftChild) key rightChild
			// lower == nil, key < upper    ... all(leftChild) key rightChild
			return countNode(node.LeftChild()) + 1 + countExtendedRange(node.RightChild(), lower, upper)
		}
	}
	if upper == nil {
		if node.Key().CompareTo(lower).GreaterThanOrEqualTo() {
			// upper == nil, lower < key    ... leftChild key all(rightChild)
			// upper == nil, lower == key   ... leftChild key all(rightChild)
			return countExtendedRange(node.LeftChild(), lower, upper) + 1 + countNode(node.RightChild())
		} else {
			// upper == nil, key < lower    ... rightChild
			return countExtendedRange(node.RightChild(), lower, upper)
		}
	}
	key := node.Key()
	cmpLower := key.CompareTo(lower)
	cmpUpper := key.CompareTo(upper)
	switch {
	case cmpUpper.GreaterThan():
		// lower < upper < key      ... leftChild
		return countExtendedRange(node.LeftChild(), lower, upper)
	case cmpLower.LessThan():
		// key < lower < upper      ... rightChild
		return countExtendedRange(node.RightChild(), lower, upper)
	default:
		// lower < upper == key     ... leftChild{upper->nil} key rightChild{lower->nil}
		// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
		// key == lower < upper     ... leftChild{upper->nil} key rightChild{lower->nil}
		// lower == key == upper    ... leftChild{upper->nil} key rightChild{lower->nil}
		return countExtendedRange(node.LeftChild(), lower, nil) + 1 + countExtendedRange(node.RightChild(), nil, upper)
	}
}

// ノードの高さ情報を取得する
func getHeight(node Node) int {
	if node == nil {
		return 0
	} else {
		return node.(RealNode).Height()
	}
}

// 左右の子の高さの違いを表現する
type balance int

const (
	// 左右の子の高さに違いはないと見なすことを表現する
	balanced balance = 0

	// 左の子の高さが右の子の高さより高いことを表現する
	leftIsHigher balance = -1

	// 左の子の高さより右の子の高さが高いことを表現する
	rightIsHigher balance = 1
)

// AVL木における回転の必要性を取得する
func checkBalance(node RealNode) balance {
	if node == nil {
		return balanced
	}
	heightL := getHeight(node.LeftChild())
	heightR := getHeight(node.RightChild())
	switch {
	case heightL < heightR && heightL+1 < heightR:
		return rightIsHigher
	case heightL > heightR && heightL > heightR+1:
		return leftIsHigher
	default:
		return balanced
	}
	// 算術オーバーフローが怖いのか？
	// heightL - heightR
	// heightL + 1
	// heightR + 1
	// だが、現実的なAVL木の高さを考えるとせいぜい高くても33～34くらいでは？
	// 高さ10で完全二分木のノード総数は(2の10乗-1)個=1023個
	// 高さ32で4億個？最小構成のノードでもノード1つ14bytesくらいかもだろうし
	// (最小構成、leftChild,rightChild　32bitsサイズのアドレス値, height,data int8, key int32)
	// 4億個もあったら最小構成でも56GB以上のメモリ・ストレージを必要とするわけで…
	// それ以上だとメモリアドレスが64bitsサイズ、必要な容量がグっと増えるし・・・
	// そもオーバーフローするとしたら内部データが意図せず破壊された場合のみで
	// その場合は正常動作を保証する必要もないわけで…
}

// ２つのノードの高さに違いがあるかを判定する
func compareNodeHeight(leftNode, rightNode Node) balance {
	heightL := getHeight(leftNode)
	heightR := getHeight(rightNode)
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

// 左右の子ノードの高さに違いがあるかを判定する
func compareChildHeight(node Node) balance {
	if node == nil {
		return balanced
	} else {
		return compareNodeHeight(node.LeftChild(), node.RightChild())
	}
}

// 2つの引数のうち大きいほうの値を返す
func intMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// 左右の子のノードから親ノードの新しい高さを計算して返す
func calcNewHeight(leftChild, rightChild Node) int {
	leftHeight := getHeight(leftChild)
	rightHeight := getHeight(rightChild)
	return 1 + intMax(leftHeight, rightHeight)
}

// 対象のノードに新しい左右の子を設定する
func setChildren(root RealNode, leftChild, rightChild Node) RealNode {
	newHeight := calcNewHeight(leftChild, rightChild)
	return root.SetChildren(leftChild, rightChild, newHeight)
}

// 対象のノードに新しい左の子を設定する
func setLeftChild(root RealNode, newLeftChild Node) RealNode {
	return setChildren(root, newLeftChild, root.RightChild())
}

// 対象のノードに新しい右の子を設定する
func setRightChild(root RealNode, newRightChild Node) RealNode {
	return setChildren(root, root.LeftChild(), newRightChild)
}

// 対象のノードの新しい左右の子と新しい値を設定する
func resetNode(root RealNode, newLeftChild, newRightChild Node, newValue interface{}) RealNode {
	newHeight := calcNewHeight(newLeftChild, newRightChild)
	return root.Set(newLeftChild, newRightChild, newHeight, newValue)
}

// Insertの処理を補助する
type insertHelper struct {
	tree            *RealTree
	replaceIfExists bool
	key             *Key
	value           *interface{}
}

// RealTreeのNewNodeを呼び出す
func (helper *insertHelper) newNode() RealNode {
	return (*helper.tree).NewNode(nil, nil, 1, (*helper.key).Copy(), *helper.value)
}

// 挿入するキーと対象のノードのキーと比較する
func (helper *insertHelper) compareKey(node Node) KeyOrdering {
	return (*helper.key).CompareTo(node.Key())
}

// 木が同一キーを許可するかを取得する
func (helper *insertHelper) allowDuplicateKeys() bool {
	return (*helper.tree).AllowDuplicateKeys()
}

// 木(またはサブツリー)に新しいキーと値を挿入する処理
func (helper *insertHelper) insertTo(root Node) (newRoot RealNode, ok bool) {
	if root == nil {
		return helper.newNode(), true
	}
	cmp := helper.compareKey(root)
	switch {
	case cmp.LessThan(): // newKey < root.key
		if newLeftChild, ok := helper.insertTo(root.LeftChild()); ok {
			newRoot = setLeftChild(root.(RealNode), newLeftChild)
		} else {
			return root.(RealNode), false
		}
	case cmp.GreaterThan(): // root.key < newKey
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

// 木の回転の処理
func rotate(root RealNode) RealNode {
	// 無限ループは複数のノードをまとめて削除した場合のバランス崩れを解消するためのもの
	for {
		switch checkBalance(root) {
		case leftIsHigher:
			root = rotateRight(root)
		case rightIsHigher:
			root = rotateLeft(root)
		case balanced:
			return root
		default:
			panic("unreachable")
		}
	}
}

// 時計回りに木を回転させる
// 左側の子孫ノードが親に、親ノードが右側サブツリーのルートノードに、なる
func rotateRight(root RealNode) RealNode {
	oldRootLeftChild := root.LeftChild().(RealNode)
	if compareChildHeight(oldRootLeftChild) == rightIsHigher {
		newRoot := oldRootLeftChild.RightChild().(RealNode)
		tempLeftChild := newRoot.LeftChild()
		tempRightChild := newRoot.RightChild()
		newRootRightChild := rotate(setLeftChild(root, tempRightChild))
		newRootLeftChild := rotate(setRightChild(oldRootLeftChild, tempLeftChild))
		return setChildren(newRoot, newRootLeftChild, newRootRightChild)
	} else {
		tempRightChild := oldRootLeftChild.RightChild()
		newRootRightChild := rotate(setLeftChild(root, tempRightChild))
		return setRightChild(oldRootLeftChild, newRootRightChild)
	}
}

// 反時計回りに木を回転
// 右側の子孫ノードが親に、親ノードが左側サブツリーのルートノードに、なる
func rotateLeft(root RealNode) RealNode {
	oldRootRightChild := root.RightChild().(RealNode)
	if compareChildHeight(oldRootRightChild) == leftIsHigher {
		newRoot := oldRootRightChild.LeftChild().(RealNode)
		tempLeftChild := newRoot.LeftChild()
		tempRightChild := newRoot.RightChild()
		newRootLeftChild := rotate(setRightChild(root, tempLeftChild))
		newRootRightChild := rotate(setLeftChild(oldRootRightChild, tempRightChild))
		return setChildren(newRoot, newRootLeftChild, newRootRightChild)
	} else {
		tempLeftChild := oldRootRightChild.LeftChild()
		newLeftChild := rotate(setRightChild(root, tempLeftChild))
		return setLeftChild(oldRootRightChild, newLeftChild)
	}
}

// 指定キーを持つノードを木から取り除く
func removeNode(root Node, key Key) (newRoot, removed Node, ok bool) {
	if root == nil {
		return nil, nil, false
	}
	cmp := key.CompareTo(root.Key())
	switch {
	case cmp.LessThan(): // key < root.Key()
		if tempLeftChild, node, ok := removeNode(root.LeftChild(), key); ok {
			removed = node
			newRoot = setLeftChild(root.(RealNode), tempLeftChild)
		} else {
			return nil, nil, false
		}
	case cmp.GreaterThan(): // root.Key() < key
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
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
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

// 木(またはサブツリー)の最小のキーを持つノードを木(またはサブツリー)から取り除く
func removeMin(root Node) (newRoot, removed Node) {
	if root == nil {
		return nil, nil
	}
	leftChild := root.LeftChild()
	if leftChild == nil {
		newRoot = root.RightChild()
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, root
	}
	leftChild, removed = removeMin(leftChild)
	newRoot = setLeftChild(root.(RealNode), leftChild)
	newRoot = rotate(newRoot.(RealNode))
	return newRoot, removed
}

// 木(またはサブツリー)の最大のキーを持つノードを木(またはサブツリー)から取り除く
func removeMax(root Node) (newRoot, removed Node) {
	if root == nil {
		return nil, nil
	}
	rightChild := root.RightChild()
	if rightChild == nil {
		newRoot = root.LeftChild()
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, root
	}
	rightChild, removed = removeMax(rightChild)
	newRoot = setRightChild(root.(RealNode), rightChild)
	newRoot = rotate(newRoot.(RealNode))
	return newRoot, removed
}

// 昇順Iterateの中身
func ascIterateNode(node Node, callBack IterateCallBack) (ok bool) {
	if node == nil {
		return true
	}
	if !ascIterateNode(node.LeftChild(), callBack) {
		return false
	}
	if breakIteration := callBack(node); breakIteration {
		return false
	}
	return ascIterateNode(node.RightChild(), callBack)
}

// 降順Iterateの中身
func descIterateNode(node Node, callBack IterateCallBack) (ok bool) {
	if node == nil {
		return true
	}
	if !descIterateNode(node.RightChild(), callBack) {
		return false
	}
	if breakIteration := callBack(node); breakIteration {
		return false
	}
	return descIterateNode(node.LeftChild(), callBack)
}

// Range系,RangeIterate系において木を巡る判定に用いる
// keyBoundsのメソッドによる判定結果に対して実装されるメソッド
type boundsChecker interface {
	// 対象ノードの左の子(サブツリー)を含める必要がある場合はtrue
	includeLower() bool

	// 対象ノードを含める必要がある場合はtrue
	includeKey() bool

	// 対象ノードの右の子(サブツリー)を含める必要がある場合はtrue
	includeUpper() bool
}

// Range系,RangeIterate系において木を巡る判定に用いる
// Range系,RangeIterate系の引数のlower,upperによって示される範囲情報に実装されるメソッド
type keyBounds interface {
	// lower情報と対象のキーとを比較する
	checkLower(key Key) boundsChecker

	// upper情報と対象のキーとを比較する
	checkUpper(key Key) boundsChecker
}

// lower,upperによってしめさえる範囲情報を返す
// 引数のextendedは木が同一キーを許す場合はtrue、同一キーを許さない場合はfalseにする必要がある
// 同一キーの挿入箇所は右の子のサブツリーに挿入されるが、木の回転により同一キーが左の子のサブツリーに移動することがある
func newKeyBounds(lower, upper Key, extended bool) keyBounds {
	if lower == nil {
		return &upperBound{upper, extended}
	} else if upper == nil {
		return &lowerBound{lower, extended}
	} else {
		return &bothBounds{lower, upper, extended}
	}
}

// lower,upperともにKeyが指定された場合の範囲情報
type bothBounds struct {
	lower, upper Key
	ext          bool
}

func (bounds *bothBounds) checkLower(key Key) boundsChecker {
	return &lowerBoundsChecker{key.CompareTo(bounds.lower), bounds.ext}
}

func (bounds *bothBounds) checkUpper(key Key) boundsChecker {
	return &upperBoundsChecker{key.CompareTo(bounds.upper), bounds.ext}
}

// upperのみが指定された場合の範囲情報
type upperBound struct {
	upper Key
	ext   bool
}

func (bounds *upperBound) checkLower(key Key) boundsChecker {
	return noBoundsChecker{}
}

func (bounds *upperBound) checkUpper(key Key) boundsChecker {
	return &upperBoundsChecker{key.CompareTo(bounds.upper), bounds.ext}
}

// lowerのみが指定された場合の範囲情報
type lowerBound struct {
	lower Key
	ext   bool
}

func (bounds *lowerBound) checkLower(key Key) boundsChecker {
	return &lowerBoundsChecker{key.CompareTo(bounds.lower), bounds.ext}
}

func (bounds *lowerBound) checkUpper(key Key) boundsChecker {
	return noBoundsChecker{}
}

// 範囲の制限がない(すなわち全範囲を含む必要がある)という判定結果を表す
type noBoundsChecker struct{}

func (noBoundsChecker) includeLower() bool { return true }
func (noBoundsChecker) includeKey() bool   { return true }
func (noBoundsChecker) includeUpper() bool { return true }

// 上限側にだけ制限がある(下限には制限がない)という判定結果を表す
// 対象ノードとその右の子に対して含むかの判定が必要になる
type upperBoundsChecker struct {
	cmpUpper KeyOrdering
	ext      bool
}

func (checker *upperBoundsChecker) includeLower() bool {
	return true
}

func (checker *upperBoundsChecker) includeKey() bool {
	return checker.cmpUpper.LessThanOrEqualTo()
}

func (checker *upperBoundsChecker) includeUpper() bool {
	return checker.cmpUpper.Less(checker.ext)
}

// 下限側にだけ制限がある(上限には制限がない)という判定結果を表す
// 対象ノードとその左の子に対して含むかの判定が必要になる
type lowerBoundsChecker struct {
	cmpLower KeyOrdering
	ext      bool
}

func (checker *lowerBoundsChecker) includeLower() bool {
	return checker.cmpLower.Greater(checker.ext)
}

func (checker *lowerBoundsChecker) includeKey() bool {
	return checker.cmpLower.GreaterThanOrEqualTo()
}

func (checker *lowerBoundsChecker) includeUpper() bool {
	return true
}

// 昇順RangeIterateの中身
func ascRangeNode(node Node, bounds keyBounds, callBack IterateCallBack) (ok bool) {
	if node == nil {
		return true
	}
	key := node.Key()
	lower := bounds.checkLower(key)
	if lower.includeLower() {
		if !ascRangeNode(node.LeftChild(), bounds, callBack) {
			return false
		}
	}
	upper := bounds.checkUpper(key)
	if lower.includeKey() && upper.includeKey() {
		if breakIteration := callBack(node); breakIteration {
			return false
		}
	}
	if upper.includeUpper() {
		return ascRangeNode(node.RightChild(), bounds, callBack)
	} else {
		return true
	}
}

// 降順RangeIterateの中身
func descRangeNode(node Node, bounds keyBounds, callBack IterateCallBack) (ok bool) {
	if node == nil {
		return true
	}
	key := node.Key()
	upper := bounds.checkUpper(key)
	if upper.includeUpper() {
		if !descRangeNode(node.RightChild(), bounds, callBack) {
			return false
		}
	}
	lower := bounds.checkLower(key)
	if lower.includeKey() && upper.includeKey() {
		if breakIteration := callBack(node); breakIteration {
			return false
		}
	}
	if lower.includeLower() {
		return descRangeNode(node.LeftChild(), bounds, callBack)
	} else {
		return true
	}
}

// Updateの中身
func updateValue(node Node, key Key, callBack UpdateValueCallBack) (newNode RealNode, ok bool) {
	if node == nil {
		return nil, false
	}
	nodeKey := node.Key()
	cmp := key.CompareTo(nodeKey)
	switch {
	case cmp.LessThan():
		if leftChild, ok := updateValue(node.LeftChild(), key, callBack); ok {
			return setLeftChild(node.(RealNode), leftChild), true
		} else {
			return nil, false
		}
	case cmp.EqualTo():
		if newValue, keepOldValue := callBack(nodeKey, node.Value()); !keepOldValue {
			return node.SetValue(newValue).(RealNode), true
		} else {
			return nil, false
		}
	case cmp.GreaterThan():
		if rightChild, ok := updateValue(node.RightChild(), key, callBack); ok {
			return setRightChild(node.(RealNode), rightChild), true
		} else {
			return nil, false
		}
	default:
		panic("unreachable")
	}
}

// 昇順UpdateIterateの中身
func ascUpdateIterate(root Node, callBack UpdateIterateCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}

	leftChild, leftUpdated, breakIteration := ascUpdateIterate(root.LeftChild(), callBack)
	if breakIteration {
		if leftUpdated {
			newRoot = setLeftChild(root.(RealNode), leftChild)
		} else {
			newRoot = root.(RealNode)
		}
		return newRoot, leftUpdated, breakIteration
	}

	newValue, keepOldValue, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		switch {
		case !leftUpdated && keepOldValue:
			newRoot = root.(RealNode)
		case !leftUpdated && !keepOldValue:
			newRoot = root.SetValue(newValue).(RealNode)
		case leftUpdated && keepOldValue:
			newRoot = setLeftChild(root.(RealNode), leftChild)
		case leftUpdated && !keepOldValue:
			newRoot = resetNode(root.(RealNode), leftChild, root.RightChild(), newValue)
		}
		updated = leftUpdated || !keepOldValue
		return newRoot, updated, breakIteration
	}

	rightChild, rightUpdated, breakIteration := ascUpdateIterate(root.RightChild(), callBack)
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

// 降順UpdateIterateの中身
func descUpdateIterate(root Node, callBack UpdateIterateCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}

	rightChild, rightUpdated, breakIteration := descUpdateIterate(root.RightChild(), callBack)
	if breakIteration {
		if rightUpdated {
			newRoot = setRightChild(root.(RealNode), rightChild)
		} else {
			newRoot = root.(RealNode)
		}
		return newRoot, rightUpdated, breakIteration
	}

	newValue, keepOldValue, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		switch {
		case !rightUpdated && keepOldValue:
			newRoot = root.(RealNode)
		case !rightUpdated && !keepOldValue:
			newRoot = root.SetValue(newValue).(RealNode)
		case rightUpdated && keepOldValue:
			newRoot = setRightChild(root.(RealNode), rightChild)
		case rightUpdated && !keepOldValue:
			newRoot = resetNode(root.(RealNode), root.LeftChild(), rightChild, newValue)
		default:
			panic("unreachable")
		}
		updated = rightUpdated || !keepOldValue
		return newRoot, updated, breakIteration
	}

	leftChild, leftUpdated, breakIteration := descUpdateIterate(root.LeftChild(), callBack)
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

// 昇順UpdateRangeIterateの中身
func ascUpdateRange(root Node, bounds keyBounds, callBack UpdateIterateCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}
	var leftUpdated, keepOldValue, rightUpdated bool
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	lower := bounds.checkLower(key)
	if lower.includeLower() {
		leftChild, leftUpdated, breakIteration = ascUpdateRange(leftChild, bounds, callBack)
		if breakIteration {
			if leftUpdated {
				newRoot = setLeftChild(root.(RealNode), leftChild)
			} else {
				newRoot = root.(RealNode)
			}
			return newRoot, leftUpdated, breakIteration
		}
	}

	var newValue interface{}
	upper := bounds.checkUpper(key)
	if lower.includeKey() && upper.includeKey() {
		newValue, keepOldValue, breakIteration = callBack(root.Key(), root.Value())
		if breakIteration {
			switch {
			case !leftUpdated && keepOldValue:
				newRoot = root.(RealNode)
			case !leftUpdated && !keepOldValue:
				newRoot = root.SetValue(newValue).(RealNode)
			case leftUpdated && keepOldValue:
				newRoot = setLeftChild(root.(RealNode), leftChild)
			case leftUpdated && !keepOldValue:
				newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
			default:
				panic("unreachable")
			}
			updated = leftUpdated || !keepOldValue
			return newRoot, updated, breakIteration
		}
	} else {
		keepOldValue = true
	}

	if upper.includeUpper() {
		rightChild, rightUpdated, breakIteration = ascUpdateRange(rightChild, bounds, callBack)
	}
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

// 降順UpdateRangeIterateの中身
func descUpdateRange(root Node, bounds keyBounds, callBack UpdateIterateCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}
	var leftUpdated, keepOldValue, rightUpdated bool
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	upper := bounds.checkUpper(key)
	if upper.includeUpper() {
		rightChild, rightUpdated, breakIteration = descUpdateRange(rightChild, bounds, callBack)
		if breakIteration {
			if rightUpdated {
				newRoot = setRightChild(root.(RealNode), rightChild)
			} else {
				newRoot = root.(RealNode)
			}
			return newRoot, rightUpdated, breakIteration
		}
	}

	var newValue interface{}
	lower := bounds.checkLower(key)
	if lower.includeKey() && upper.includeKey() {
		newValue, keepOldValue, breakIteration = callBack(root.Key(), root.Value())
		if breakIteration {
			switch {
			case !rightUpdated && keepOldValue:
				newRoot = root.(RealNode)
			case !rightUpdated && !keepOldValue:
				newRoot = root.SetValue(newValue).(RealNode)
			case rightUpdated && keepOldValue:
				newRoot = setRightChild(root.(RealNode), rightChild)
			case rightUpdated && !keepOldValue:
				newRoot = resetNode(root.(RealNode), root.LeftChild(), rightChild, newValue)
			default:
				panic("unreachable")
			}
			updated = rightUpdated || !keepOldValue
			return newRoot, updated, breakIteration
		}
	} else {
		keepOldValue = true
	}

	if lower.includeLower() {
		leftChild, leftUpdated, breakIteration = descUpdateRange(leftChild, bounds, callBack)
	}
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

// 昇順DeleteIterateの中身
func ascDeleteIterate(root Node, callBack DeleteIterateCallBack) (newRoot Node, deleted []Node, breakIteration bool) {
	if root == nil {
		return nil, nil, false
	}

	leftChild, leftDeleted, breakIteration := ascDeleteIterate(root.LeftChild(), callBack)
	if breakIteration {
		if len(leftDeleted) > 0 {
			newRoot = setLeftChild(root.(RealNode), leftChild)
		} else {
			newRoot = root
		}
		newRoot = rotate(newRoot.(RealNode))
		return newRoot, leftDeleted, breakIteration
	}

	deleteRoot, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		deleted = leftDeleted
		if deleteRoot {
			deleted = append(deleted, root)
			rightChild := root.RightChild()
			if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
				leftChild, newRoot = removeMax(leftChild)
			} else {
				rightChild, newRoot = removeMin(rightChild)
			}
			if newRoot != nil {
				newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
			}
		} else if len(leftDeleted) > 0 {
			newRoot = setLeftChild(root.(RealNode), leftChild)
		} else {
			newRoot = root
		}
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, deleted, breakIteration
	}

	rightChild, rightDeleted, breakIteration := ascDeleteIterate(root.RightChild(), callBack)
	deleted = leftDeleted
	switch {
	case len(leftDeleted) == 0 && !deleteRoot && len(rightDeleted) == 0:
		newRoot = root
	case !deleteRoot:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case deleteRoot:
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, rightDeleted...)
	if newRoot != nil {
		newRoot = rotate(newRoot.(RealNode))
	}
	return newRoot, deleted, breakIteration
}

// 降順DeleteIterateの中身
func descDeleteIterate(root Node, callBack DeleteIterateCallBack) (newRoot Node, deleted []Node, breakIteration bool) {
	if root == nil {
		return nil, nil, false
	}

	rightChild, rightDeleted, breakIteration := descDeleteIterate(root.RightChild(), callBack)
	if breakIteration {
		if len(rightDeleted) > 0 {
			newRoot = setRightChild(root.(RealNode), rightChild)
		} else {
			newRoot = root
		}
		return rotate(newRoot.(RealNode)), rightDeleted, breakIteration
	}

	deleteRoot, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		deleted = rightDeleted
		if deleteRoot {
			deleted = append(deleted, root)
			leftChild := root.LeftChild()
			if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
				leftChild, newRoot = removeMax(leftChild)
			} else {
				rightChild, newRoot = removeMin(rightChild)
			}
			if newRoot != nil {
				newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
			}
		} else if len(rightDeleted) > 0 {
			newRoot = setRightChild(root.(RealNode), rightChild)
		} else {
			newRoot = root
		}
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, deleted, breakIteration
	}

	leftChild, leftDeleted, breakIteration := descDeleteIterate(root.LeftChild(), callBack)

	deleted = rightDeleted
	switch {
	case len(leftDeleted) == 0 && !deleteRoot && len(rightDeleted) == 0:
		newRoot = root
	case !deleteRoot:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case deleteRoot:
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, leftDeleted...)
	if newRoot != nil {
		newRoot = rotate(newRoot.(RealNode))
	}
	return newRoot, deleted, breakIteration
}

// 昇順DeleteRangeIterateの中身
func ascDeleteRange(root Node, bounds keyBounds, callBack DeleteIterateCallBack) (newRoot Node, deleted []Node, breakIteration bool) {
	if root == nil {
		return nil, nil, false
	}
	var deleteRoot bool
	var leftDeleted, rightDeleted []Node
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	lower := bounds.checkLower(key)
	if lower.includeLower() {
		leftChild, leftDeleted, breakIteration = ascDeleteRange(leftChild, bounds, callBack)
		if breakIteration {
			if len(leftDeleted) > 0 {
				newRoot = setLeftChild(root.(RealNode), leftChild)
			} else {
				newRoot = root
			}
			newRoot = rotate(newRoot.(RealNode))
			return newRoot, leftDeleted, breakIteration
		}
	}

	upper := bounds.checkUpper(key)
	if lower.includeKey() && upper.includeKey() {
		deleteRoot, breakIteration = callBack(root.Key(), root.Value())
		if breakIteration {
			deleted = leftDeleted
			if deleteRoot {
				deleted = append(deleted, root)
				rightChild := root.RightChild()
				if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
					leftChild, newRoot = removeMax(leftChild)
				} else {
					rightChild, newRoot = removeMin(rightChild)
				}
				if newRoot != nil {
					newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
				}
			} else if len(leftDeleted) > 0 {
				newRoot = setLeftChild(root.(RealNode), leftChild)
			} else {
				newRoot = root
			}
			if newRoot != nil {
				newRoot = rotate(newRoot.(RealNode))
			}
			return newRoot, deleted, breakIteration
		}
	}

	if upper.includeUpper() {
		rightChild, rightDeleted, breakIteration = ascDeleteRange(rightChild, bounds, callBack)
	}

	deleted = leftDeleted
	switch {
	case len(leftDeleted) == 0 && !deleteRoot && len(rightDeleted) == 0:
		newRoot = root
	case !deleteRoot:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case deleteRoot:
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, rightDeleted...)
	if newRoot != nil {
		newRoot = rotate(newRoot.(RealNode))
	}
	return newRoot, deleted, breakIteration
}

// 降順DeleteRangeIterateの中身
func descDeleteRange(root Node, bounds keyBounds, callBack DeleteIterateCallBack) (newRoot Node, deleted []Node, breakIteration bool) {
	if root == nil {
		return nil, nil, false
	}
	var deleteRoot bool
	var leftDeleted, rightDeleted []Node
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	upper := bounds.checkUpper(key)
	if upper.includeUpper() {
		rightChild, rightDeleted, breakIteration = descDeleteRange(rightChild, bounds, callBack)
		if breakIteration {
			if len(rightDeleted) > 0 {
				newRoot = setRightChild(root.(RealNode), rightChild)
			} else {
				newRoot = root
			}
			return rotate(newRoot.(RealNode)), rightDeleted, breakIteration
		}
	}

	lower := bounds.checkLower(key)
	if lower.includeKey() && upper.includeKey() {
		deleteRoot, breakIteration = callBack(root.Key(), root.Value())
		if breakIteration {
			deleted = rightDeleted
			if deleteRoot {
				deleted = append(deleted, root)
				leftChild := root.LeftChild()
				if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
					leftChild, newRoot = removeMax(leftChild)
				} else {
					rightChild, newRoot = removeMin(rightChild)
				}
				if newRoot != nil {
					newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
				}
			} else if len(rightDeleted) > 0 {
				newRoot = setRightChild(root.(RealNode), rightChild)
			} else {
				newRoot = root
			}
			if newRoot != nil {
				newRoot = rotate(newRoot.(RealNode))
			}
			return newRoot, deleted, breakIteration
		}
	}

	if lower.includeLower() {
		leftChild, leftDeleted, breakIteration = descDeleteRange(leftChild, bounds, callBack)
	}

	deleted = rightDeleted
	switch {
	case len(leftDeleted) == 0 && !deleteRoot && len(rightDeleted) == 0:
		newRoot = root
	case !deleteRoot:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case deleteRoot:
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, leftDeleted...)
	if newRoot != nil {
		newRoot = rotate(newRoot.(RealNode))
	}
	return newRoot, deleted, breakIteration
}

// AlterNodeにおいて実体のNodeを取得するためのバックドア
func (node *alterNode) Node() Node {
	return node.inner
}

func (node *alterNode) Key() Key {
	return node.inner.Key()
}

func (node *alterNode) Value() interface{} {
	return node.inner.Value()
}

func (*alterNode) Keep() (request AlterRequest) {
	return
}

func (*alterNode) Replace(newValue interface{}) (request AlterRequest) {
	request.replaceValue = true
	request.newValue = newValue
	return
}

func (*alterNode) Delete() (request AlterRequest) {
	request.deleteNode = true
	return
}

// ノードを変更も削除もせず維持することを指定する
func (request *AlterRequest) Keep() (ret AlterRequest) {
	*request = ret
	return
}

// ノードの値をnewValueに置き換えることを指定する
func (request *AlterRequest) Replace(newValue interface{}) (ret AlterRequest) {
	ret.replaceValue = true
	ret.newValue = newValue
	*request = ret
	return
}

// ノードを削除することを指定する
func (request *AlterRequest) Delete() (ret AlterRequest) {
	ret.deleteNode = true
	*request = ret
	return
}

// 維持指定だったかどうか
func (request *AlterRequest) isKeepRequest() bool {
	return !request.replaceValue && !request.deleteNode
}

// 変更指定だったかどうか
func (request *AlterRequest) isReplaceRequest() bool {
	return request.replaceValue && !request.deleteNode
}

// 削除指定だったかどうか
func (request *AlterRequest) isDeleteRequest() bool {
	return !request.replaceValue && request.deleteNode
}

// Alterの中身
func alter(node Node, key Key, callBack AlterNodeCallBack) (newNode, deleted Node, ok bool) {
	if node == nil {
		// nodeの返却は必要か？
		return node, nil, false
	}
	cmp := key.CompareTo(node.Key())
	switch {
	case cmp.LessThan():
		if newLeftChild, deleted, ok := alter(node.LeftChild(), key, callBack); ok {
			newNode = rotate(setLeftChild(node.(RealNode), newLeftChild))
			return newNode, deleted, ok
		} else {
			return node, nil, false
		}
	case cmp.GreaterThan():
		if newRightChild, deleted, ok := alter(node.RightChild(), key, callBack); ok {
			newNode = rotate(setRightChild(node.(RealNode), newRightChild))
			return newNode, deleted, ok
		} else {
			return node, nil, false
		}
	}
	request := callBack(&alterNode{node})
	switch {
	case request.isKeepRequest():
		return node, nil, false
	case request.isReplaceRequest():
		newNode = node.SetValue(request.newValue)
		return newNode, nil, true
	case request.isDeleteRequest():
		deleted = node
		leftChild := node.LeftChild()
		rightChild := node.RightChild()
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newNode = removeMax(leftChild)
		} else {
			rightChild, newNode = removeMin(rightChild)
		}
		if newNode != nil {
			newNode = rotate(setChildren(newNode.(RealNode), leftChild, rightChild))
		}
		return newNode, deleted, true
	default:
		panic("unreachable")
	}
}

// 昇順AlterIterateの中身
func ascAlterIterate(root Node, callBack AlterIterateCallBack) (newRoot Node, deleted []Node, anyChanged, breakIteration bool) {
	if root == nil {
		return nil, nil, false, false
	}

	leftChild, leftDeleted, leftAnyChanged, breakIteration := ascAlterIterate(root.LeftChild(), callBack)
	if breakIteration {
		if leftAnyChanged {
			newRoot = rotate(setLeftChild(root.(RealNode), leftChild))
		} else {
			newRoot = root
		}
		return newRoot, leftDeleted, leftAnyChanged, breakIteration
	}

	request, breakIteration := callBack(&alterNode{root})
	if breakIteration {
		deleted = leftDeleted
		switch {
		case !leftAnyChanged && request.isKeepRequest():
			newRoot = root
		case leftAnyChanged && request.isKeepRequest():
			newRoot = rotate(setLeftChild(root.(RealNode), leftChild))
		case !leftAnyChanged && request.isReplaceRequest():
			newRoot = root.SetValue(request.newValue)
		case leftAnyChanged && request.isReplaceRequest():
			newValue := request.newValue
			rightChild := root.RightChild()
			newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
		case request.isDeleteRequest():
			deleted = append(deleted, root)
			rightChild := root.RightChild()
			if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
				leftChild, newRoot = removeMax(leftChild)
			} else {
				rightChild, newRoot = removeMin(rightChild)
			}
			if newRoot != nil {
				newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
			}
		default:
			panic("unreachable")
		}
		anyChanged = leftAnyChanged || !request.isKeepRequest()
		return newRoot, deleted, anyChanged, breakIteration
	}

	rightChild, rightDeleted, rightAnyChanged, breakIteration := ascAlterIterate(root.RightChild(), callBack)

	deleted = leftDeleted
	switch {
	case !leftAnyChanged && request.isKeepRequest() && !rightAnyChanged:
		newRoot = root
	case request.isKeepRequest():
		newRoot = rotate(setChildren(root.(RealNode), leftChild, rightChild))
	case request.isReplaceRequest():
		newValue := request.newValue
		newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
	case request.isDeleteRequest():
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, rightDeleted...)
	anyChanged = leftAnyChanged || !request.isKeepRequest() || rightAnyChanged
	return newRoot, deleted, anyChanged, breakIteration
}

// 降順AlterIterateの中身
func descAlterIterate(root Node, callBack AlterIterateCallBack) (newRoot Node, deleted []Node, anyChanged, breakIteration bool) {
	if root == nil {
		return nil, nil, false, false
	}

	rightChild, rightDeleted, rightAnyChanged, breakIteration := descAlterIterate(root.RightChild(), callBack)
	if breakIteration {
		if rightAnyChanged {
			newRoot = rotate(setRightChild(root.(RealNode), rightChild))
		} else {
			newRoot = root
		}
		return newRoot, rightDeleted, rightAnyChanged, breakIteration
	}

	request, breakIteration := callBack(&alterNode{root})
	if breakIteration {
		deleted = rightDeleted
		switch {
		case !rightAnyChanged && request.isKeepRequest():
			newRoot = root
		case rightAnyChanged && request.isKeepRequest():
			newRoot = rotate(setRightChild(root.(RealNode), rightChild))
		case !rightAnyChanged && request.isReplaceRequest():
			newRoot = root.SetValue(request.newValue)
		case rightAnyChanged && request.isReplaceRequest():
			newValue := request.newValue
			leftChild := root.LeftChild()
			newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
		case request.isDeleteRequest():
			deleted = append(deleted, root)
			leftChild := root.LeftChild()
			if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
				leftChild, newRoot = removeMax(leftChild)
			} else {
				rightChild, newRoot = removeMin(rightChild)
			}
			if newRoot != nil {
				newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
			}
		default:
			panic("unreachable")
		}
		anyChanged = rightAnyChanged || !request.isKeepRequest()
		return newRoot, deleted, anyChanged, breakIteration
	}

	leftChild, leftDeleted, leftAnyChanged, breakIteration := descAlterIterate(root.LeftChild(), callBack)

	deleted = rightDeleted
	switch {
	case !leftAnyChanged && request.isKeepRequest() && !rightAnyChanged:
		newRoot = root
	case request.isKeepRequest():
		newRoot = rotate(setChildren(root.(RealNode), leftChild, rightChild))
	case request.isReplaceRequest():
		newValue := request.newValue
		newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
	case request.isDeleteRequest():
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, leftDeleted...)
	anyChanged = leftAnyChanged || !request.isKeepRequest() || rightAnyChanged
	return newRoot, deleted, anyChanged, breakIteration
}

// 昇順AlterRangeIterateの中身
func ascAlterRange(root Node, bounds keyBounds, callBack AlterIterateCallBack) (newRoot Node, deleted []Node, anyChanged, breakIteration bool) {
	if root == nil {
		return nil, nil, false, false
	}
	var leftDeleted, rightDeleted []Node
	var leftAnyChanged, rightAnyChanged bool
	var request AlterRequest
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	lower := bounds.checkLower(key)
	if lower.includeLower() {
		leftChild, leftDeleted, leftAnyChanged, breakIteration = ascAlterRange(leftChild, bounds, callBack)
		if breakIteration {
			if leftAnyChanged {
				newRoot = rotate(setLeftChild(root.(RealNode), leftChild))
			} else {
				newRoot = root
			}
			return newRoot, leftDeleted, leftAnyChanged, breakIteration
		}
	}

	upper := bounds.checkUpper(key)
	if lower.includeKey() && upper.includeKey() {
		request, breakIteration = callBack(&alterNode{root})
		if breakIteration {
			deleted = leftDeleted
			switch {
			case !leftAnyChanged && request.isKeepRequest():
				newRoot = root
			case leftAnyChanged && request.isKeepRequest():
				newRoot = rotate(setLeftChild(root.(RealNode), leftChild))
			case !leftAnyChanged && request.isReplaceRequest():
				newRoot = root.SetValue(request.newValue)
			case leftAnyChanged && request.isReplaceRequest():
				newValue := request.newValue
				newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
			case request.isDeleteRequest():
				deleted = append(deleted, root)
				if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
					leftChild, newRoot = removeMax(leftChild)
				} else {
					rightChild, newRoot = removeMin(rightChild)
				}
				if newRoot != nil {
					newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
				}
			default:
				panic("unreachable")
			}
			anyChanged = leftAnyChanged || !request.isKeepRequest()
			return newRoot, deleted, anyChanged, breakIteration
		}
	}

	if upper.includeUpper() {
		rightChild, rightDeleted, rightAnyChanged, breakIteration = ascAlterRange(rightChild, bounds, callBack)
	}

	deleted = leftDeleted
	switch {
	case !leftAnyChanged && request.isKeepRequest() && !rightAnyChanged:
		newRoot = root
	case request.isKeepRequest():
		newRoot = rotate(setChildren(root.(RealNode), leftChild, rightChild))
	case request.isReplaceRequest():
		newValue := request.newValue
		newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
	case request.isDeleteRequest():
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, rightDeleted...)
	anyChanged = leftAnyChanged || !request.isKeepRequest() || rightAnyChanged
	return newRoot, deleted, anyChanged, breakIteration
}

// 降順AlterRangeIterateの中身
func descAlterRange(root Node, bounds keyBounds, callBack AlterIterateCallBack) (newRoot Node, deleted []Node, anyChanged, breakIteration bool) {
	if root == nil {
		return nil, nil, false, false
	}
	var leftDeleted, rightDeleted []Node
	var leftAnyChanged, rightAnyChanged bool
	var request AlterRequest
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	upper := bounds.checkUpper(key)
	if upper.includeUpper() {
		rightChild, rightDeleted, rightAnyChanged, breakIteration = descAlterRange(rightChild, bounds, callBack)
		if breakIteration {
			if rightAnyChanged {
				newRoot = rotate(setRightChild(root.(RealNode), rightChild))
			} else {
				newRoot = root
			}
			return newRoot, rightDeleted, rightAnyChanged, breakIteration
		}

	}

	lower := bounds.checkLower(key)
	if lower.includeKey() && upper.includeKey() {
		request, breakIteration = callBack(&alterNode{root})
		if breakIteration {
			deleted = rightDeleted
			switch {
			case !rightAnyChanged && request.isKeepRequest():
				newRoot = root
			case rightAnyChanged && request.isKeepRequest():
				newRoot = rotate(setRightChild(root.(RealNode), rightChild))
			case !rightAnyChanged && request.isReplaceRequest():
				newRoot = root.SetValue(request.newValue)
			case rightAnyChanged && request.isReplaceRequest():
				newValue := request.newValue
				newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
			case request.isDeleteRequest():
				deleted = append(deleted, root)
				if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
					leftChild, newRoot = removeMax(leftChild)
				} else {
					rightChild, newRoot = removeMin(rightChild)
				}
				if newRoot != nil {
					newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
				}
			default:
				panic("unreachable")
			}
			anyChanged = rightAnyChanged || !request.isKeepRequest()
			return newRoot, deleted, anyChanged, breakIteration
		}
	}

	if lower.includeLower() {
		leftChild, leftDeleted, leftAnyChanged, breakIteration = descAlterRange(leftChild, bounds, callBack)
	}

	deleted = rightDeleted
	switch {
	case !leftAnyChanged && request.isKeepRequest() && !rightAnyChanged:
		newRoot = root
	case request.isKeepRequest():
		newRoot = rotate(setChildren(root.(RealNode), leftChild, rightChild))
	case request.isReplaceRequest():
		newValue := request.newValue
		newRoot = rotate(resetNode(root.(RealNode), leftChild, rightChild, newValue))
	case request.isDeleteRequest():
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = rotate(setChildren(newRoot.(RealNode), leftChild, rightChild))
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, leftDeleted...)
	anyChanged = leftAnyChanged || !request.isKeepRequest() || rightAnyChanged
	return newRoot, deleted, anyChanged, breakIteration
}
