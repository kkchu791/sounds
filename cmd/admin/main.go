package main

import (
	"fmt"
	"os"

	"github.com/kkchu791/sounds/internal/app/admin"
)

func main() {
	if err := admin.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
