package csvutil

var supportedSortDirections = []string{"asc", "desc"}

// SortOption is option holder for Sort.
type SortOption struct {
	// Source file does not have header line. (default false)
	NoHeader bool
	// Encoding of source file. (default utf8)
	Encoding string
	// Encoding for output.
	OutputEncoding string
	// Column symbol of target column
	Column string
	// Type is sort key's data type
	Type string
	// Sort direction
	Direction string
}
