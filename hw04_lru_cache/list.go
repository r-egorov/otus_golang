package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func newNode(v interface{}) *ListItem {
	res := new(ListItem)
	res.Value = v
	return res
}

func (l *list) isEmpty() bool {
	return l.front == nil
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	toInsert := newNode(v)

	insertFront := func(node *ListItem) {
		node.Next = l.front
		l.front.Prev = node
		l.front = node
	}
	l.insert(toInsert, insertFront)

	return toInsert
}

func (l *list) PushBack(v interface{}) *ListItem {
	toInsert := newNode(v)

	insertBack := func(node *ListItem) {
		node.Prev = l.back
		l.back.Next = node
		l.back = node
	}
	l.insert(toInsert, insertBack)

	return toInsert
}

func (l *list) insert(toInsert *ListItem, insertFrontOrBack func(*ListItem)) {
	if l.isEmpty() {
		l.front = toInsert
		l.back = toInsert
	} else {
		insertFrontOrBack(toInsert)
	}
	l.len++
}

func (l *list) Remove(i *ListItem) {
	if i == l.front {
		l.front = i.Next
	}
	if i == l.back {
		l.back = i.Prev
	}
	l.tiePrevAndNext(i)
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}

	if i == l.back {
		l.back = i.Prev
	}
	l.tiePrevAndNext(i)
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

func (l *list) tiePrevAndNext(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}
