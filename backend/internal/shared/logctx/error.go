package logctx

type errorWithLogContext struct {
	err    error
	fields Fields
}

func (e *errorWithLogContext) Error() string {
	return e.err.Error()
}

func (e *errorWithLogContext) Unwrap() error {
	return e.err
}
