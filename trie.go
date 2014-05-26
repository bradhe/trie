package trie

// An implementation of a trie that supports a few common operations, namely
// Lookup, Range, and Prefix. Pretty straight forward stuff. Figuring out the
// type of the resultant object is an exercise for the reader.
type Trie interface {
	Insert(key []byte, val interface{})
	Lookup(key []byte) interface{}
	Range(start, end []byte) []interface{}
	Prefix(prefix []byte) []interface{}
	Count() int
}

type trieImpl struct {
	key      byte
	value    interface{}
	parent   *trieImpl
	children []*trieImpl
}

func between(b, start, end byte) bool {
	return b >= start && b <= end
}

func (self *trieImpl) getChildValues() []interface{} {
	results := make([]interface{}, 0)

	if self.value != nil {
		results = append(results, self.value)
	}

	for _, child := range self.children {
		results = append(results, child.getChildValues()...)
	}

	return results
}

func (self *trieImpl) Count() int {
	totalChildren := 0

	for _, trie := range self.children {
		totalChildren += trie.Count()
	}

	return len(self.children) + totalChildren
}

func (self *trieImpl) Insert(key []byte, val interface{}) {
	// If we got to here that means that this element matches the key fully.
	if len(key) == 0 {
		self.value = val
		return
	}

	front := key[0]

	// No match yet, see if there is a child that matches.
	for _, trie := range self.children {
		if trie.key == front {
			trie.Insert(key[1:len(key)], val)
			return
		}
	}

	// Okay, let's insert a new one.
	trie := new(trieImpl)
	trie.parent = self
	trie.children = make([]*trieImpl, 0)
	trie.key = front

	self.children = append(self.children, trie)

	// Now let's drill in!
	trie.Insert(key[1:len(key)], val)
}

func (self *trieImpl) Lookup(key []byte) interface{} {
	if len(key) == 0 {
		return self.value
	}

	front := key[0]

	for _, trie := range self.children {
		if trie.key == front {
			return trie.Lookup(key[1:len(key)])
		}
	}

	// Didn't match anything.
	return nil
}

func (self *trieImpl) Range(start, end []byte) []interface{} {
	results := make([]interface{}, 0)

	// If both are empty then we completely matched.
	if len(start) < 1 && len(end) < 1 {
		return self.getChildValues()
	}

	var startb, endb byte

	if len(start) > 0 {
		startb = start[0]
	}

	if len(end) > 0 {
		endb = end[0]
	}

	if self.key == 0 || between(self.key, startb, endb) {
		if self.value != nil {
			results = append(results, self.value)
		}

		// Ternary would be nice here.
		var offset int

		if self.key > 0 {
			offset = 1
		}

		// Do a scan for all the children.
		for _, child := range self.children {
			subset := child.Range(start[offset:], end[offset:])
			results = append(results, subset...)
		}
	}

	return results
}

func (self *trieImpl) Prefix(prefix []byte) []interface{} {
	if len(prefix) == 0 {
		return self.getChildValues()
	}

	front := prefix[0]

	for _, trie := range self.children {
		if trie.key == front {
			return trie.Prefix(prefix[1:len(prefix)])
		}
	}

	// Didn't match anything.
	return nil
}

func New() Trie {
	trie := new(trieImpl)
	trie.children = make([]*trieImpl, 0)
	return trie
}
