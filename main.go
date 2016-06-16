package main

import (
	"os"

	"github.com/ks888/okdoc/commands"
)

func main() {
	if err := commands.OkdocCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
