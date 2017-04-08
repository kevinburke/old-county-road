# Go html boilerplate

This is a starter pack for doing web development with Go, with support for some
of the things you'll usually want to add to an HTML web server:

- Adding templates and rendering them
- Regex matching for routes
- Logging requests responses
- Serving static content
- Watching/restarting the server after changes to CSS/templates
- Loading configuration from a config file
- Flash success and error messages

Feel free to adapt as you see fit.

To get started, run `go get ./...` and then `make serve` to start a server on
port 7065.

Templates go in the "templates" folder; you can see how they're loaded by
examining the `init` function in main.go.

Static files go in the "static" folder. Run "make assets" to recompile them into
the binary.
