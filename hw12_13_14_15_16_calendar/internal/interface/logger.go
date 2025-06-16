package interfaces

//go:generate mockgen -source=logger.go -package=mocks -destination=../../mocks/mock_logger.go
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}
