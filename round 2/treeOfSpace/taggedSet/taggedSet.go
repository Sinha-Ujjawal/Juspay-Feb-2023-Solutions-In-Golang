package taggedSet

type TaggedSet[T comparable, S comparable] struct {
	store map[T]map[S]bool
	size  uint
}

func New[T comparable, S comparable]() TaggedSet[T, S] {
	return TaggedSet[T, S]{store: nil, size: 0}
}

func (m *TaggedSet[T, S]) Size() uint {
	return m.size
}

func (m *TaggedSet[T, S]) NumTags() uint {
	return uint(len(m.store))
}

func (m *TaggedSet[T, S]) Contains(tag T) bool {
	if m.store == nil {
		return false
	}
	_, ok := m.store[tag]
	return ok
}

func (m *TaggedSet[T, S]) Lookup(tag T) []S {
	if m.store == nil {
		return nil
	}
	items, ok := m.store[tag]
	if !ok {
		return nil
	}
	var ret []S
	for item := range items {
		ret = append(ret, item)
	}
	return ret
}

func (m *TaggedSet[T, S]) AddEntry(tag T, item S) {
	if m.store == nil {
		m.store = make(map[T]map[S]bool)
	}
	items, ok := m.store[tag]
	if !ok {
		m.store[tag] = map[S]bool{item: true}
		m.size++
		return
	}
	_, ok = items[item]
	if !ok {
		items[item] = true
		m.size++
	}
}

func (m *TaggedSet[T, S]) RemoveEntry(tag T, item S) {
	if m.store == nil {
		return
	}
	items, ok := m.store[tag]
	if !ok {
		return
	}
	delete(items, item)
	m.size--
	if len(items) == 0 {
		delete(m.store, tag)
	}
}
