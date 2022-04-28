package filter

// Signature for each type of filter
// takes in a line of input to determine if input is accepted
type FilterFunc func(line []byte) (accepted bool, err error)

func MakeFilters() []FilterFunc {
	// TODO: accept commandline args and make list of FilterFunc
	return []FilterFunc{Prefix([]byte("hello"))}
}
