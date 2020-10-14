package resolver

const (
	ErrorFirstLastMissing = "you must provide a `first` or `last` value to properly paginate the objects"
	ErrorFirstLastInvalid = "you must provide a valid `first` or `last` value in range 1-100"
	ErrorCursorInvalid    = "cursor value is invalid"
)
