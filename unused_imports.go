package main

import "unused"

// This declaration marks the import as used by referencing an
// item from the package.
var _ = unused.Item  // TODO: Delete before committing!

func main() {
	debugData := debug.Profile()
	_ = debugData // Used only during debugging.
	....
}
