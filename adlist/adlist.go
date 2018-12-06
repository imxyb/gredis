package adlist

const (
	AL_START_HEAD = iota
	AL_START_TAIL
)

type ListMatch func(value interface{}, key interface{}) bool

// 链表节点
type ListNode struct {
	Prev  *ListNode
	Next  *ListNode
	Value interface{}
}

// 链表迭代器
type ListIter struct {
	next *ListNode
	direction int
}

// 双链表结构
type List struct {
	Head *ListNode
	Tail *ListNode
	Len  int64
	Match ListMatch
}

// 新建一个链表
func NewList() *List {
	return new(List)
}

func (list *List) AddNodeHead(value interface{}) {
	node := NewListNode(value)

	if list.Len == 0 {
		list.Head = node
		list.Tail = node
	} else {
		node.Prev = nil
		node.Next = list.Head
		list.Head.Prev = node
		list.Head = node
	}

	list.Len++
}

func (list *List) AddNodeTail(value interface{}) {
	node := NewListNode(value)

	if list.Len == 0 {
		list.Head = node
		list.Tail = node
	} else {
		node.Prev = list.Tail
		list.Tail.Next = node
		list.Tail = node
	}

	list.Len++
}

// 插入指定节点的前面或者后面
func (list *List) InsertNode(oldNode *ListNode, value interface{}, after int) {
	node := NewListNode(value)

	// 插入给定节点前面
	if after == 0 {
		node.Next = oldNode
		node.Prev = oldNode.Prev

		// 如果被插入的节点是头节点
		if list.Head == oldNode {
			list.Head = node
		}
	} else if after == 1 {         // 插入到给定节点的后面
		node.Prev = oldNode
		node.Next = oldNode.Next

		// 如果被插入的节点是尾节点
		if list.Tail == oldNode {
			list.Tail = node
		}
	}

	if node.Prev != nil {
		node.Prev.Next = node
	}

	if node.Next != nil {
		node.Next.Prev = node
	}

	list.Len++
}

func (list *List) SearchKey(key interface{}) *ListNode {
	iter := NewListIter(list, AL_START_HEAD)

	for {
		node := iter.Next()

		if node == nil {
			break
		}

		if list.Match != nil {
			if list.Match(node.Value, key) {
				return node
			}
		} else {
			if node.Value == key {
				return node
			}
		}
	}

	return nil
}

// 寻找对应索引的节点
func (list *List) Index(index int) *ListNode {
	var node *ListNode

	if index >= 0 {
		node = list.Head
		for i := 0; i < index; i++ {
			node = node.Next
		}
	} else {
		node = list.Tail
		index = (-index)-1
		for i := 0; i < index; i++ {
			node = node.Prev
		}
	}

	return node
}

// 删除一个节点
func (list *List) DelNode(node *ListNode) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		list.Head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		list.Tail = node.Prev
	}

	list.Len--
}

// 新建一个链表迭代器
func NewListIter(list *List, direction int) *ListIter {
	iter := new(ListIter)

	if direction == AL_START_HEAD {
		iter.next = list.Head
	} else if direction == AL_START_TAIL {
		iter.next = list.Tail
	}

	iter.direction = direction

	return iter
}

func (iter *ListIter) Next() *ListNode {
	current := iter.next

	if current != nil {
		if iter.direction == AL_START_HEAD {
			iter.next = current.Next
		} else if iter.direction == AL_START_TAIL {
			iter.next = current.Prev
		}
	}

	return current
}

func NewListNode(value interface{}) *ListNode {
	return &ListNode{
		Value:value,
	}
}

func ListDup(origin *List) *List {
	cpy := NewList()

	iter := NewListIter(origin, AL_START_HEAD)

	for {
		node := iter.Next()
		if node == nil {
			break
		}
		cpy.AddNodeTail(node.Value)
	}

	return cpy
}

