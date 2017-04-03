package main

var cmdGenerate = &Command{
	Run:       runGenerate,
	UsageLine: "generate ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdGenerate.Flag.BoolVar(&flagA, "a", false, "")
}

// runGenerate executes generate command and return exit code.
func runGenerate(args []string) int {

	return 0
}
