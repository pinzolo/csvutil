package main

var cmdStruct = &Command{
	Run:       runStruct,
	UsageLine: "struct ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdStruct.Flag.BoolVar(&flagA, "a", false, "")
}

// runStruct executes struct command and return exit code.
func runStruct(args []string) int {

	return 0
}
