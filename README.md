# lookup

Simple utility for looking up a value in a JSON document using a simple JSON path expression.

You can do things like `x.store.book[0].title`.

`x` is the root node.

This is a spartan implementation of JSON paths. Only `.` and `[` `]` are supported so far.
