package util

import (
	"flag"
	"fmt"
)

// ParseFlags ...
func ParseFlags() {
	var identity, filter string
	flag.StringVar(&identity, "i", "~/.ssh/id_rsa", "path to identity file")
	flag.StringVar(&filter, "f", "", "string to filter instances")
	flag.Parse()
	fmt.Println(identity)
}
