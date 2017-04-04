package main

var cmdTel = &Command{
	Run:       runTel,
	UsageLine: "tel ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdTel.Flag.BoolVar(&flagA, "a", false, "")
}

// runTel executes tel command and return exit code.
func runTel(args []string) int {

	return 0
}
