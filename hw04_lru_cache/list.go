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

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{Value: v, Next: l.front}

	if l.front != nil {
		l.front.wire(i)
	} else {
		l.back = i
	}
	l.front = i
	l.len++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{Value: v, Next: nil, Prev: l.back}
	if l.back != nil {
		i.wire(l.back)
	} else {
		l.front = i
	}
	l.back = i
	l.len++
	return i
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	prevItem := i.Prev
	nextItem := i.Next
	nextItem.wire(prevItem)

	if nextItem == nil {
		l.back = prevItem
	}
	if prevItem == nil {
		l.front = nextItem
	}
	i.Prev = nil
	i.Next = nil

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	prevItem := i.Prev
	if prevItem == nil {
		return
	}

	i.Next.wire(prevItem)
	if i.Next == nil {
		l.back = prevItem
	}

	l.front.wire(i)
	l.front = i
	i.Prev = nil
}

func (next *ListItem) wire(prev *ListItem) {
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
}
