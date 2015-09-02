package trie

// An implementation of a trie that supports a few common operations, namely
// Lookup, Range, and Prefix. Pretty straight forward stuff. Figuring out the
// type of the resultant object is an exercise for the reader.
type Trie interface {
	Insert(key []byte, val interface{})
	Lookup(key []byte) interface{}
	Range(start, end []byte) map[string]interface{}
	RangeN(start, end []byte, n int) map[string]interface{}
	OffsetRangeN(offset, start, end []byte, n int) map[string]interface{}
	Prefix(prefix []byte) map[string]interface{}
	PrefixN(prefix []byte, n int) map[string]interface{}
	OffsetPrefixN(offset, prefix []byte, n int) map[string]interface{}
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

func minOf(i1, i2 int) int {
	if i1 < i2 {
		return i1
	}

	return i2
}

func isOffsetLesser(offset, orig []byte) bool {
	// If there is no offset, then we don't want to do anything.
	if len(offset) < 1 {
		return true
	}

	// We want the min of
	l1 := len(orig)
	l2 := len(offset)

	if l1 == l2 {
		// We want to determine if they're perfectly equal at the same time. We'll do
		// an equals comparison along the way to determine if
		equal := true

		for i := 0; i < l1; i++ {
			if offset[i] == orig[i] {
				continue
			}

			// They aren't equal, so let's flip the switch.
			equal = false

			// If any part of the prefix is greater, than that means we're good to
			// go.
			if offset[i] < orig[i] {
				return true
			}

			// Otherwise, we get out of here.
			return false
		}

		// If they are equal, we don't want to do this work.
		if equal {
			return false
		}
	} else {
		// They aren't equal in size, so we'll compare as many as we can.
		n := minOf(l1, l2)

		for i := 0; i < n; i++ {
			// If any one is greater, then we're done here.
			if orig[i] > offset[i] {
				return true
			}

			if offset[i] > orig[i] {
				return false
			}
		}
	}

	// If we made it to here, it means that the orig >= offset, which means we
	// can keep rollin'
	return true
}

func (self *trieImpl) getChildValues(res map[string]interface{}, prefix, offset []byte, n *int) {
	// If n is exactly 0, we can finish visiting child nodes.
	if *n == 0 {
		return
	}

	if self.value != nil {
		if isOffsetLesser(offset, prefix) {
			res[string(prefix)] = self.value

			// One less to go.
			if *n > 0 {
				*n -= 1
			}
		}
	}

	for _, child := range self.children {
		k := append(prefix, child.key)
		child.getChildValues(res, k, offset, n)
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

func (self *trieImpl) doRange(offset, start, end, prefix []byte, res map[string]interface{}, n *int) {
	if self.key == 0 {
		for _, child := range self.children {
			child.doRange(offset, start, end, prefix, res, n)
		}

		return
	}

	// If both are empty then we completely matched.
	if len(start) < 1 && len(end) < 1 {
		self.getChildValues(res, append(prefix, self.key), offset, n)
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
			child.doRange(offset, start, end, append(prefix, self.key), res, n)
		}

		return
	}

	// Is there more work left to do?
	if *n == 0 {
		return
	}

	prefix = append(prefix, self.key)

	if startb == self.key {
		if self.value != nil && isOffsetLesser(offset, prefix) {
			res[string(prefix)] = self.value

			if *n > 0 {
				*n -= 1
			}
		}

		// Nothing left to do.
		if *n == 0 {
			return
		}

		for _, child := range self.children {
			child.doRange(offset, start, maxString(start), prefix, res, n)
		}
	} else if endb == self.key {
		if self.value != nil && isOffsetLesser(offset, prefix) {
			res[string(prefix)] = self.value

			if *n > 0 {
				*n -= 1
			}
		}

		// Nothing left to do.
		if *n == 0 {
			return
		}

		for _, child := range self.children {
			child.doRange(offset, minString(end), end, prefix, res, n)
		}
	} else if between(self.key, startb, endb) {
		self.getChildValues(res, prefix, offset, n)
	}
}

func (self *trieImpl) Range(start, end []byte) map[string]interface{} {
	return self.RangeN(start, end, -1)
}

func (self *trieImpl) RangeN(start, end []byte, n int) map[string]interface{} {
	return self.OffsetRangeN([]byte{}, start, end, n)
}

func (self *trieImpl) OffsetRangeN(offset, start, end []byte, n int) map[string]interface{} {
	results := make(map[string]interface{}, 0)
	self.doRange(offset, start, end, []byte(""), results, &n)
	return results
}

func (self *trieImpl) doPrefix(offset, prefix, orig []byte, res map[string]interface{}, n *int) {
	// If there is an offset, we'll navigate down to it before we do our own
	// thing.
	if len(prefix) == 0 {
		self.getChildValues(res, orig, offset, n)
		return
	}

	// Are we don't here?
	if *n == 0 {
		return
	}

	front := prefix[0]

	for _, trie := range self.children {
		if trie.key == front {
			trie.doPrefix(offset, prefix[1:len(prefix)], append(orig, front), res, n)
			return
		}
	}
}

func (self *trieImpl) Prefix(prefix []byte) map[string]interface{} {
	return self.PrefixN(prefix, -1)
}

func (self *trieImpl) PrefixN(prefix []byte, n int) map[string]interface{} {
	return self.OffsetPrefixN([]byte{}, prefix, n)
}

func (self *trieImpl) OffsetPrefixN(offset, prefix []byte, n int) map[string]interface{} {
	res := make(map[string]interface{})
	self.doPrefix(offset, prefix, []byte{}, res, &n)
	return res
}

func New() Trie {
	trie := new(trieImpl)
	trie.children = make([]*trieImpl, 0)
	return trie
}
