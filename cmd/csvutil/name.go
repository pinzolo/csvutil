package main

var cmdName = &Command{
	Run:       runName,
	UsageLine: "name ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdName.Flag.BoolVar(&flagA, "a", false, "")
}

// runName executes name command and return exit code.
func runName(args []string) int {

	return 0
}
