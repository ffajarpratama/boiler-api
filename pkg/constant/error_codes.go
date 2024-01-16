package constant

// default error Code
const (
	DefaultUnhandledError = 1000 + iota
	DefaultNotFoundError
	DefaultBadRequestError
	DefaultUnauthorizedError
	DefaultDuplicateDataError
)

var ErrorMessageMap = map[int]string{
	DefaultUnhandledError: "Something went wrong with our side, please try again",
	DefaultNotFoundError:  "Data not found",
}
