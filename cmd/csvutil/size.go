package main

var cmdSize = &Command{
	Run:       runSize,
	UsageLine: "size ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdSize.Flag.BoolVar(&flagA, "a", false, "")
}

// runSize executes size command and return exit code.
func runSize(args []string) int {

	return 0
}
