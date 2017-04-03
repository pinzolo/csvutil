package main

var cmdAddress = &Command{
	Run:       runAddress,
	UsageLine: "address ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdAddress.Flag.BoolVar(&flagA, "a", false, "")
}

// runAddress executes address command and return exit code.
func runAddress(args []string) int {

	return 0
}
