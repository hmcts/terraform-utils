package main

import (
	"flag"
	"fmt"
	"os"
	"tf-utils/tfwhitelist"
)

var whitelist = flag.Bool("whitelist", false, "whitelist terraform modules + resources")

func main() {
	flag.Parse()

	var infraPath string
	var whitelistPaths []string

	if *whitelist {
		if flag.NArg() < 2 {
			_, _ = fmt.Fprintf(os.Stderr,
				"error: tf-utils --whitelist <terraform-infra-dir-path> <whitelist-file-path>...<whitelist-file-path>\n")
			os.Exit(1)
		} else {
			infraPath = flag.Arg(0)
			for i := 0; i < flag.NArg(); i++ {
				whitelistPaths = append(whitelistPaths, flag.Arg(i))
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Currently only whitelist operation is available\n")
		os.Exit(1)
	}

	err := tfwhitelist.LoadAndMatchAll(infraPath, whitelistPaths)
	if err != nil {
		os.Exit(1)
	}
}
