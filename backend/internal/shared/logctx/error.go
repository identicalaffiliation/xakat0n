package logctx

type logContextError struct {
	err    error
	fields Fields
}

func (e *logContextError) Error() string {
	return e.err.Error()
}

func (e *logContextError) Unwrap() error {
	return e.err
}
