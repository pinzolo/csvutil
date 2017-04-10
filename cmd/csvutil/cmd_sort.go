package main

var cmdSort = &Command{
	Run:       runSort,
	UsageLine: "sort ",
	Short:     "[Not implemented yet]",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdSort.Flag.BoolVar(&flagA, "a", false, "")
}

// runSort executes sort command and return exit code.
func runSort(args []string) int {

	return 0
}
