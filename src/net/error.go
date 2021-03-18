package net

// Error type holds HTTP status code
type Error interface {
	error

	Status() int
}

// StatusError ...
type StatusError struct {
	Err  error
	Code int
}

// Status returns HTTP status code
func (e StatusError) Status() int {
	return e.Code
}

// Error interface implementation
func (e StatusError) Error() string {
	return e.Err.Error()
}
