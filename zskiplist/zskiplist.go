package zskiplist

import (
	"gredis/consts"
	"gredis/gredis"
	"math/rand"
	"time"
)

// 跳跃表节点
type Node struct {
	Obj      *gredis.Object
	Score    float64
	Backward *Node
	Level    []Level
}

// 层
type Level struct {
	Forward *Node
	Span    uint
}

// 跳跃表
type List struct {
	Header *Node
	Tail   *Node
	Length int64
	Level  int
}

// 范围
type ZRangeSpec struct {
	Min   float64
	Max   float64
	Minex bool
	Maxex bool
}

func NewNode(level int, score float64, obj *gredis.Object) *Node {
	node := new(Node)
	node.Score = score
	node.Obj = obj
	node.Level = make([]Level, level)
	return node
}

func NewList() *List {
	list := new(List)
	list.Level = 1

	list.Header = NewNode(consts.ZskipListMaxLevel, 0, nil)
	for i := 0; i < consts.ZskipListMaxLevel; i++ {
		list.Header.Level[i].Forward = nil
		list.Header.Level[i].Span = 0
	}

	return list
}

func (list *List) Insert(score float64, obj *gredis.Object) (*Node, error) {
	rank := make([]uint, consts.ZskipListMaxLevel)
	update := make([]*Node, consts.ZskipListMaxLevel)

	x := list.Header

	for i := list.Level - 1; i >= 0; i-- {
		if i == list.Level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// 说明要比头节点的下一个节点排序等级更高
		for x.Level[i].Forward != nil &&
			(x.Level[i].Forward.Score < score ||
				(x.Level[i].Forward.Score == score &&
					gredis.CompareStringObject(x.Level[i].Forward.Obj, obj) < 0)) {
			rank[i] += x.Level[i].Span

			x = x.Level[i].Forward
		}
		update[i] = x
	}

	level := randomLevel()

	if level > list.Level {
		for i := list.Level; i < level; i++ {
			rank[i] = 0
			update[i] = list.Header
			update[i].Level[i].Span = uint(list.Length)
		}

		list.Level = level
	}

	newNode := NewNode(level, score, obj)
	for i := 0; i < level; i++ {
		newNode.Level[i].Forward = update[i].Level[i].Forward
		update[i].Level[i].Forward = newNode
		newNode.Level[i].Span = update[i].Level[i].Span - (rank[0] - rank[i])
		update[i].Level[i].Span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < list.Level; i++ {
		update[i].Level[i].Span++
	}

	if update[0] != list.Header {
		newNode.Backward = update[0]
	}

	if newNode.Level[0].Forward != nil {
		newNode.Level[0].Forward.Backward = newNode
	} else {
		list.Tail = newNode
	}

	list.Length++

	return newNode, nil
}

func (list *List) GetRank(score float64, obj *gredis.Object) uint {
	var rank uint
	x := list.Header

	for i := list.Level - 1; i >= 0; i-- {
		for x.Level[i].Forward != nil &&
			(x.Level[i].Forward.Score < score ||
				(x.Level[i].Forward.Score == score &&
					gredis.CompareStringObject(x.Level[i].Forward.Obj, obj) <= 0)) {
			rank += x.Level[i].Span
			x = x.Level[i].Forward
		}
	}

	if x.Obj != nil && gredis.EqualStringObject(x.Obj, obj) {
		return rank
	}

	return 0
}

func (list *List) Delete(score float64, obj *gredis.Object) bool {
	var x *Node
	var update []*Node

	x = list.Header
	for i := list.Level - 1; i >= 0; i-- {
		for x.Level[i].Forward != nil &&
			(x.Level[i].Forward.Score < score ||
				(x.Level[i].Forward.Score == score &&
					gredis.CompareStringObject(x.Level[i].Forward.Obj, obj) < 0)) {
			x = x.Level[i].Forward
		}

		update = append(update, x)
	}

	x = x.Level[0].Forward
	if x != nil && x.Score == score && gredis.EqualStringObject(x.Obj, obj) {
		list.DeleteNode(x, update)
		return true
	}

	return false
}

func (list *List) DeleteNode(delNode *Node, update []*Node) {
	for i := 0; i < list.Level; i++ {
		if update[i].Level[i].Forward == delNode {
			update[i].Level[i].Span += delNode.Level[i].Span - 1
			update[i].Level[i].Forward = delNode.Level[i].Forward
		} else {
			update[i].Level[i].Span -= 1
		}
	}

	if delNode.Level[0].Forward != nil {
		delNode.Level[0].Forward.Backward = delNode.Backward
	} else {
		list.Tail = delNode.Backward
	}

	// 如果被删除的那个节点是最高level，那么要相应把level减少
	for list.Level > 1 && list.Header.Level[list.Level-1].Forward == nil {
		list.Level--
	}

	list.Length--
}

func (list *List) GetElementByRank(rank uint) *Node {
	var traversed uint
	x := list.Header

	for i := list.Level - 1; i >= 0; i-- {
		for x.Level[i].Forward != nil && (traversed+x.Level[i].Span <= rank) {
			traversed += x.Level[i].Span
			x = x.Level[i].Forward
		}
	}
	if traversed == rank {
		return x
	}

	return nil
}

func (list *List) IsInRange(rangeSpec *ZRangeSpec) bool {
	if rangeSpec.Max < rangeSpec.Min ||
		(rangeSpec.Min == rangeSpec.Max && (rangeSpec.Minex || rangeSpec.Maxex)) {
		return false
	}

	max := list.Tail
	if max == nil || !zslValueGteMin(max.Score, rangeSpec) {
		return false
	}

	min := list.Header.Level[0].Forward
	if min == nil || !zslValueLteMax(min.Score, rangeSpec) {
		return false
	}

	return true
}

func zslValueGteMin(value float64, rangeSpec *ZRangeSpec) bool {
	if rangeSpec.Minex {
		return value > rangeSpec.Min
	}
	return value >= rangeSpec.Min
}

func zslValueLteMax(value float64, rangeSpec *ZRangeSpec) bool {
	if rangeSpec.Maxex {
		return value < rangeSpec.Max
	}
	return value <= rangeSpec.Max
}
func randomLevel() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(consts.ZskipListMaxLevel + 1)
}
