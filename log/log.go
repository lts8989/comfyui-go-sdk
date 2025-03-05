package log

import "fmt"

// Logger logger 对象
var Logger LoggerInterface

func init() {
	Logger = &DefaultLogger{}
}

// InitLogger 初始化logger对象
func InitLogger(ILogger LoggerInterface) {
	Logger = ILogger
}

// LoggerInterface 日志接口
type LoggerInterface interface {
	Debugf(format string, params ...interface{})

	Infof(format string, params ...interface{})

	Warnf(format string, params ...interface{})

	Errorf(format string, params ...interface{})

	Debug(v ...interface{})

	Info(v ...interface{})

	Warn(v ...interface{})

	Error(v ...interface{})
}

// Debugf debug 格式化
func Debugf(format string, params ...interface{}) {
	Logger.Debugf(format, params...)
}

// Infof 打印info
func Infof(format string, params ...interface{}) {
	Logger.Infof(format, params...)
}

// Warnf warn格式化
func Warnf(format string, params ...interface{}) {
	Logger.Warnf(format, params...)
}

// Errorf error格式化
func Errorf(format string, params ...interface{}) {
	Logger.Errorf(format, params...)
}

// Debug 打印debug
func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

// Info 打印Info
func Info(v ...interface{}) {
	Logger.Info(v...)
}

// Warn 打印Warn
func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

// Error 打印Error
func Error(v ...interface{}) {
	Logger.Error(v...)
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
}

// Debugf debug 格式化
func (d *DefaultLogger) Debugf(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params...)
}

// Infof 打印info
func (d *DefaultLogger) Infof(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params...)
}

// Warnf warn格式化
func (d *DefaultLogger) Warnf(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params...)
}

// Errorf error格式化
func (d *DefaultLogger) Errorf(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params...)
}

// Debug 打印debug
func (d *DefaultLogger) Debug(v ...interface{}) {
	fmt.Println(v)
}

// Info 打印Info
func (d *DefaultLogger) Info(v ...interface{}) {
	fmt.Println(v)
}

// Warn 打印Warn
func (d *DefaultLogger) Warn(v ...interface{}) {
	fmt.Println(v)
}

// Error 打印Error
func (d *DefaultLogger) Error(v ...interface{}) {
	fmt.Println(v)
}
