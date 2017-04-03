package main

var cmdCount = &Command{
	Run:       runCount,
	UsageLine: "count ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdCount.Flag.BoolVar(&flagA, "a", false, "")
}

// runCount executes count command and return exit code.
func runCount(args []string) int {

	return 0
}
