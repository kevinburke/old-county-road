# go html templates

This is a starter pack for doing web development with Go; rendering templates,
serving static content, and watching/restarting after changes. Feel free to
adapt as you see fit.

To get started, run `go get ./...` and then `make serve`.

Templates go in the "templates" folder; you can see how they're loaded by
examining the `init` function in main.go.

Static files go in the "static" folder. Run "make assets" to recompile them into
the binary.
