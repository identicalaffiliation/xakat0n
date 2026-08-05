package ports

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}
