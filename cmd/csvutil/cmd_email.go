package main

var cmdEmail = &Command{
	Run:       runEmail,
	UsageLine: "email ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdEmail.Flag.BoolVar(&flagA, "a", false, "")
}

// runEmail executes email command and return exit code.
func runEmail(args []string) int {

	return 0
}
