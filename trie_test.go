package trie

import (
	"testing"
)

func TestTrieLookup(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	val := trie.Lookup([]byte("test1"))

	if val == nil {
		t.Fatalf(`Expected val not to be nil, was nil.`)
	}

	if val != "Hello" {
		t.Fatalf(`Expected val not to be "Hello", was "%s".`, val)
	}
}

func TestTrieRangeInclusive(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test1"), []byte("test2"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["test1"] != "Hello" {
		t.Fatalf(`Expected "test1" to be "Hello", got "%v"`, vals["test1"])
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["test2"])
	}
}

func TestTrieRangeWithCommonSuffix(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test:2015-04-30:test"), "Hello")
	trie.Insert([]byte("test:2015-05-01:test"), "Hello")
	trie.Insert([]byte("test:2015-05-02:test"), "Hello")
	trie.Insert([]byte("test:2015-05-03:test"), "Hello")
	trie.Insert([]byte("test:2015-05-04:test"), "Hello")

	vals := trie.Range([]byte("test:2015-05-01:test"), []byte("test:2015-05-04:test"))

	if len(vals) != 4 {
		t.Fatalf(`Expected length of val to be 4, got %d.`, len(vals))
	}
}

func TestTrieRangeExclusive(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test0"), []byte("test3"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}
}

func TestTrieRangePrefixed(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test"), []byte("tesu"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["test1"] != "Hello" {
		t.Fatalf(`Expected "test1" to be "Hello", got "%v"`, vals["test1"])
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["test2"])
	}
}

func TestTrieRangeOutside(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test3"), []byte("test4"))

	if len(vals) != 0 {
		t.Fatalf(`Expected length of val to be 0, got %d.`, len(vals))
	}
}

func TestTrieRangePartiallyOutside(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test1"), []byte("test4"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["test1"] != "Hello" {
		t.Fatalf(`Expected "test1" to be "Hello", got "%v"`, vals["test1"])
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["test2"])
	}
}

func TestTrieOffsetRangeNPartiallyOutside(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	offset := []byte("test1")
	vals := trie.OffsetRangeN(offset, []byte("test1"), []byte("test4"), 2)

	if len(vals) != 1 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["test2"])
	}
}

func TestTrieRangeUnbalancedEnding(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test"), []byte("test4"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["test1"] != "Hello" {
		t.Fatalf(`Expected "test1" to be "Hello", got "%v"`, vals["test1"])
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["test2"])
	}
}

func TestTrieOffsetRangeNUnbalancedEnding(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	offset := []byte("test1")
	vals := trie.OffsetRangeN(offset, []byte("test"), []byte("test4"), 2)

	if len(vals) != 1 {
		t.Fatalf(`Expected length of val to be 1, got %d.`, len(vals))
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["test2"])
	}
}

func TestTrieRangeNFullyBetween(t *testing.T) {
	trie := setupTrie()

	vals := trie.RangeN([]byte("20140901"), []byte("20140911"), 5)

	if len(vals) != 5 {
		t.Fatalf(`Expected length of val to be 5, got %d.`, len(vals))
	}
}

func TestTrieOffsetRangeNFullyBetween(t *testing.T) {
	trie := setupTrie()

	offset := []byte("20140905")
	vals := trie.OffsetRangeN(offset, []byte("20140901"), []byte("20140911"), 5)

	if len(vals) != 5 {
		t.Logf("%v", vals)
		t.Fatalf(`Expected length of val to be 5, got %d.`, len(vals))
	}
}

func TestTrieRangeFullyBetween(t *testing.T) {
	trie := setupTrie()

	// NOTE: This is invalid as endKey technically comes before startKey
	vals := trie.Range([]byte("20140901"), []byte("20140911"))

	if len(vals) != 11 {
		t.Fatalf(`Expected length of val to be 11, got %d.`, len(vals))
	}
}

func TestTriePrefix(t *testing.T) {
	trie := New()
	trie.Insert([]byte("table2#test1"), "Hello")
	trie.Insert([]byte("table2#test2"), "World")

	vals := trie.Prefix([]byte("table2"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	// We also want to test that it got the write keys back.
	if vals["table2#test1"] != "Hello" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test1" to be "Hello", got %s.`, vals["table2#test1"])
	}

	if vals["table2#test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got %s.`, vals["table2#test2"])
	}
}

func TestTriePrefixN(t *testing.T) {
	trie := New()
	trie.Insert([]byte("table2#test1"), "Hello")
	trie.Insert([]byte("table2#test2"), "World")

	vals := trie.PrefixN([]byte("table2"), 1)

	if len(vals) != 1 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	// We also want to test that it got the write keys back.
	if vals["table2#test1"] != "Hello" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test1" to be "Hello", got %s.`, vals["table2#test1"])
	}
}

func TestTrieOffsetPrefixN(t *testing.T) {
	trie := New()
	trie.Insert([]byte("table2#test1"), "Hello")
	trie.Insert([]byte("table2#test2"), "World")
	trie.Insert([]byte("table2#test3"), "Yes")
	trie.Insert([]byte("table2#test4"), "Again")
	trie.Insert([]byte("table2#test5"), "Blah")

	// Throw in a random out-of-bounds value to make sure we don't get outside
	// the prefix.
	trie.Insert([]byte("table3#test2"), "Hallooo")

	vals := trie.OffsetPrefixN([]byte("table2#test2"), []byte("table2"), 10)

	if len(vals) != 3 {
		t.Logf("%v", vals)
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["table2#test3"] != "Yes" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test3" to be "Yes", got %s.`, vals["table2#test3"])
	}

	if vals["table2#test4"] != "Again" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test4" to be "Again", got %s.`, vals["table2#test4"])
	}

	if vals["table2#test5"] != "Blah" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test5" to be "Blah", got %s.`, vals["table2#test5"])
	}
}

func TestTriePrefixNReturnsEarlyIfThereAreMissingValues(t *testing.T) {
	trie := New()
	trie.Insert([]byte("table2#test1"), "Hello")
	trie.Insert([]byte("table2#test2"), "World")
	trie.Insert([]byte("table2#test3"), "Yes")

	vals := trie.PrefixN([]byte("table2"), 5)

	if len(vals) != 3 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	// We also want to test that it got the write keys back.
	if vals["table2#test1"] != "Hello" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test1" to be "Hello", got %s.`, vals["table2#test1"])
	}

	if vals["table2#test2"] != "World" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test2" to be "World", got %s.`, vals["table2#test2"])
	}

	if vals["table2#test3"] != "Yes" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test3" to be "Yes", got %s.`, vals["table2#test3"])
	}
}

func TestTriePrefixWithExactMatch(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Prefix([]byte("test"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	// We also want to test that it got the write keys back.
	if vals["test"] != "Hello" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test" to be "Hello", got %s.`, vals["test"])
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got %s.`, vals["test2"])
	}
}

func TestTriePrefixWithLongTails(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test"), "Hello")
	trie.Insert([]byte("test2"), "World")
	trie.Insert([]byte("test again"), "Once")
	trie.Insert([]byte("test again wow"), "Again")

	vals := trie.Prefix([]byte("test"))

	if len(vals) != 4 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	// We also want to test that it got the write keys back.
	if vals["test again wow"] != "Again" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test again wow" to be "Again", got %s.`, vals["test again wow"])
	}
}

func TestTrieRangeWithCloseNeighbors(t *testing.T) {
	trie := New()
	trie.Insert([]byte("prefix1:prefix2:2015-05-01"), "Hello")
	trie.Insert([]byte("prefix1:prefix200:2015-05-01"), "What")
	trie.Insert([]byte("prefix1:prefix2:2015-05-30"), "Friend")

	vals := trie.Range([]byte("prefix1:prefix2:2015-05-01"), []byte("prefix1:prefix2:2015-05-30"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	if vals["prefix1:prefix2:2015-05-01"] != "Hello" {
		t.Fatalf(`Expected "test1" to be "Hello", got "%v"`, vals["prefix1:prefix2:2015-05-01"])
	}

	if vals["prefix1:prefix2:2015-05-30"] != "Friend" {
		t.Fatalf(`Expected "test2" to be "World", got "%v"`, vals["prefix1:prefix2:2015-05-30"])
	}
}

func BenchmarkTrieLookup(b *testing.B) {
	trie := New()
	keys := generateKeys(6, "")

	for i, k := range keys {
		trie.Insert([]byte(k), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trie.Lookup([]byte("abcdefa"))
	}
}

func BenchmarkTriePrefix(b *testing.B) {
	trie := New()
	keys := generateKeys(6, "")

	for i, k := range keys {
		trie.Insert([]byte(k), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trie.Prefix([]byte("abcd"))
	}
}

//
// Helpers and stuff.
//
func generateKeys(level int, prefix string) []string {
	keys := make([]string, 0)

	for i := 'a'; i < 'h'; i++ {
		key := prefix + string(i)
		keys = append(keys, key)

		if level > 0 {
			keys = append(keys, generateKeys(level-1, key)...)
		}
	}

	return keys
}

func setupTrie() Trie {
	trie := New()
	trie.Insert([]byte("20140901"), "20140901")
	trie.Insert([]byte("20140902"), "20140902")
	trie.Insert([]byte("20140903"), "20140903")
	trie.Insert([]byte("20140904"), "20140904")
	trie.Insert([]byte("20140905"), "20140905")
	trie.Insert([]byte("20140906"), "20140906")
	trie.Insert([]byte("20140907"), "20140907")
	trie.Insert([]byte("20140908"), "20140907")
	trie.Insert([]byte("20140909"), "20140907")
	trie.Insert([]byte("20140910"), "20140910")
	trie.Insert([]byte("20140911"), "20140911")
	return trie
}
