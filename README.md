# lookup

This package provides a way to search and manipulate JSON files using simple JSON path expressions.

## Utilities

A utility named `loookup` is included, that can be used for finding a value in a JSON document when given a simple JSON path expression.

You can do things like `x.store.book[0].title`.

`x` is the root node.

This is a spartan implementation of JSON paths. Only `.` and `[` `]` are supported so far.
