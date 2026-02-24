package random

type ZeroLengthError struct {
}

func (e *ZeroLengthError) Error() string {
	return "length cannot be zero"
}
