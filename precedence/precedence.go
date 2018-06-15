package precedence

const (
	Lowest = iota + 1
	Equal
	LteGte
	Add
	Multi
	Prefix
	Call
)
