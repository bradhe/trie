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

func TestTrieRangeOutside(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Range([]byte("test3"), []byte("test4"))

	if len(vals) != 0 {
		t.Fatalf(`Expected length of val to be 0, got %d.`, len(vals))
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

func TestTrieRangeUnbalancedBeginning(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	// NOTE: This is invalid as endKey technically comes before startKey
	vals := trie.Range([]byte("test1"), []byte("test"))

	if len(vals) != 0 {
		t.Fatalf(`Expected length of val to be 0, got %d.`, len(vals))
	}
}

func TestTriePrefix(t *testing.T) {
	trie := New()
	trie.Insert([]byte("test1"), "Hello")
	trie.Insert([]byte("test2"), "World")

	vals := trie.Prefix([]byte("test"))

	if len(vals) != 2 {
		t.Fatalf(`Expected length of val to be 2, got %d.`, len(vals))
	}

	// We also want to test that it got the write keys back.
	if vals["test1"] != "Hello" {
		t.Logf("%v", vals)
		t.Fatalf(`Expected "test1" to be "Hello", got %s.`, vals["test1"])
	}

	if vals["test2"] != "World" {
		t.Fatalf(`Expected "test2" to be "World", got %s.`, vals["test2"])
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
