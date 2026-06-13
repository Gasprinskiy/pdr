package z_logger

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string, err error)
	Fatal(msg string, err error)
}
