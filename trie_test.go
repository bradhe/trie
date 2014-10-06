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
