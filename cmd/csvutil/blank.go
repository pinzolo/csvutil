package main

var cmdBlank = &Command{
	Run:       runBlank,
	UsageLine: "blank ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdBlank.Flag.BoolVar(&flagA, "a", false, "")
}

// runBlank executes blank command and return exit code.
func runBlank(args []string) int {

	return 0
}
