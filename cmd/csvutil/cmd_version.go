package main

import (
	"fmt"
	"os"

	"github.com/pinzolo/csvutil"
)

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "バージョン表示",
	Long: `DESCRIPTION
        現在の csvutil のバージョンを表示します。
	`,
}

// runVersion executes Version command and return exit code.
func runVersion(args []string) int {
	fmt.Fprintf(os.Stdout, "csvutil version %s\n", csvutil.Version)
	return 0
}
