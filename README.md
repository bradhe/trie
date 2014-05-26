# Trie

This is a quick and dirty implementation of a Trie. It supports a few
operations, namely `Prefix`, `Range`, and `Lookup`. It's best intended for when
you have a large list of strings that you need to efficiently search through.

## Usage

```go
package main

import (
  "github.com/bradhe/trie"
)

func main() {
  t := trie.New()
  t.Insert([]byte("test1"), "First Object")
  t.Insert([]byte("test2"), "Second Object")
  t.Insert([]byte("test3"), "Third Object")

  // Get a single object.
  _ = t.Lookup([]byte("test1"))

  // Get a range of objects.
  _ = t.Range([]byte("test2"), []byte("test3"))

  // Get all objects that have a prefix.
  _ = t.Prefix([]byte("test"))
}
```

## Suggested Improvements

A few ways that this implementation could be more efficient:

1. Sort the list of child nodes to make lookups slightly faster.
1. More benchmarks comparing it to other methodologies.
