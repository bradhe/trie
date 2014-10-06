package trie

// An implementation of a trie that supports a few common operations, namely
// Lookup, Range, and Prefix. Pretty straight forward stuff. Figuring out the
// type of the resultant object is an exercise for the reader.
type Trie interface {
	Insert(key []byte, val interface{})
	Lookup(key []byte) interface{}
	Range(start, end []byte) map[string] interface{}
	Prefix(prefix []byte) map[string] interface{}
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

func maxString(str []byte) []byte {
	s := make([]byte, len(str))

	for i := range s {
		s[i] = 0xFF
	}

	return s
}

func minString(str []byte) []byte {
	s := make([]byte, len(str))

	// NOTE: This isn't strictly nescessary, it gets instantiated to 0s anyway.
	for i := range s {
		s[i] = 0x00
	}

	return s
}

func (self *trieImpl) getChildValues(res map[string] interface{}, prefix []byte) {
	if self.value != nil {
		res[string(prefix)] = self.value
	}

	for _, child := range self.children {
		k := append(prefix, child.key)
		child.getChildValues(res, k)
	}
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

func (self *trieImpl) doRange(start, end, prefix []byte, res map[string] interface{}) {
	if self.key == 0 {
		for _, child := range self.children {
			child.doRange(start, end, prefix, res)
		}

		return
	}

	// If both are empty then we completely matched.
	if len(start) < 1 && len(end) < 1 {
		self.getChildValues(res, append(prefix, self.key))
		return
	}

	var startb, endb byte

	if len(start) > 0 {
		startb = start[0]
		start = start[1:]
	}

	if len(end) > 0 {
		endb = end[0]
		end = end[1:]
	}

	if (startb == self.key) && (endb == self.key) {
		for _, child := range self.children {
			child.doRange(start, end, append(prefix, self.key), res)
		}

		return
	}

	prefix = append(prefix, self.key)

	if (startb == self.key) {
		if self.value != nil {
			res[string(prefix)] = self.value
		}

		for _, child := range self.children {
			child.doRange(start, maxString(start), prefix, res)
		}
	} else if (endb == self.key) {
		if self.value != nil {
			res[string(prefix)] = self.value
		}

		for _, child := range self.children {
			child.doRange(minString(end), end, prefix, res)
		}
	} else if between(self.key, startb, endb) {
		self.getChildValues(res, prefix)
	}
}

func (self *trieImpl) Range(start, end []byte) map[string] interface{} {
	results := make(map[string]interface{}, 0)
	self.doRange(start, end, []byte(""), results)
	return results
}

func (self *trieImpl) doPrefix(prefix, orig []byte, res map[string] interface{}) {
	if len(prefix) == 0 {
		self.getChildValues(res, orig)
		return
	}

	front := prefix[0]

	for _, trie := range self.children {
		if trie.key == front {
			trie.doPrefix(prefix[1:len(prefix)], append(orig, front), res)
			return
		}
	}
}

func (self *trieImpl) Prefix(prefix []byte) map[string] interface{} {
	res := make(map[string] interface{})
	self.doPrefix(prefix, []byte(""), res)
	return res
}

func New() Trie {
	trie := new(trieImpl)
	trie.children = make([]*trieImpl, 0)
	return trie
}
