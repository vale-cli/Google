package main

import "os"

// cdf is registered with go-cmdtest so .ct cases can change directory into a
// fixture before invoking vale.
func cdf() int {
	os.Chdir(os.Args[1])
	return 0
}

func main() {}
