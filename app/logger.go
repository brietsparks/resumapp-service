package app

type Logger interface {
	Info(args ...interface{})
	Fatal(args ...interface{})
}
