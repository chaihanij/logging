package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	globalLevel  log.Level = log.InfoLevel
	globalOut    io.Writer = os.Stderr
	globalFormat log.Formatter
	// logs keeps track of all created logs so that we can apply log globally
	// @LOCK should we get rid of this or keep it minimal?
	loggers []*Logger
	// hooks
	hookLogger *Logger
	// common or aggregate log
	commonLogger *Logger
	// cdr log
	cdrLogger *Logger
	// stats
	statsLogger *Logger
)

func init() {
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	// Set format
	// log.SetFormatter(&log.JSONFormatter{})
	// create common log
	l := log.New().WithField("module", "common")
	l.Logger.Level = log.InfoLevel
	l.Logger.Out = os.Stderr
	commonLogger = &Logger{"common", l}
	hookLogger = &Logger{"hook", l}
	// create crd log
	c := log.New().WithField("module", "cdr")
	c.Logger.Level = log.InfoLevel
	c.Logger.Out = os.Stderr
	cdrLogger = &Logger{"cdr", c}
	// create crd log
	s := log.New().WithField("module", "stats")
	s.Logger.Level = log.InfoLevel
	s.Logger.Out = os.Stderr
	statsLogger = &Logger{"stats", s}
	// commonLogger = log.New()
	// hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	// if err == nil {
	// 	commonLogger.Hooks.Add(hook)
	// } else {
	// 	fmt.Fprintf(os.Stdout, "fail to create log hook: %v", err)
	// }
}

// Log appends line, file and function context to the logger
func Log() *log.Entry {
	if pc, f, line, ok := runtime.Caller(1); ok {
		fnName := runtime.FuncForPC(pc).Name()
		file := strings.Split(f, "mobilebid")[1]
		caller := fmt.Sprintf("%s:%v %s", file, line, fnName)
		return log.WithField("caller", caller)
	}
	return &log.Entry{}
}
func CommonLogger() *Logger {
	return commonLogger
}
func SetCommonLoggerLevel(level log.Level) {
	commonLogger.SetLevel(level)
}
func SetCommonLoggerFormat(format log.Formatter) {
	globalFormat = format
	commonLogger.SetFormat(format)
}
func SetCommonLoggerOut(out io.Writer) {
	commonLogger.SetOut(out)
}
func CDRLogger() *Logger {
	return cdrLogger
}
func SetCDRLoggerLevel(level log.Level) {
	cdrLogger.SetLevel(level)
}
func SetCDRLoggerFormat(format log.Formatter) {
	globalFormat = format
	cdrLogger.SetFormat(format)
}
func SetCDRLoggerOut(out io.Writer) {
	cdrLogger.SetOut(out)
}
func GetCDRLoggerOut() io.Writer {
	return cdrLogger.Logger.Out
}
func StatsLogger() *Logger {
	return statsLogger
}
func SetStatsLoggerLevel(level log.Level) {
	statsLogger.SetLevel(level)
}
func SetStatsLoggerFormat(format log.Formatter) {
	globalFormat = format
	statsLogger.SetFormat(format)
}
func SetStatsLoggerOut(out io.Writer) {
	statsLogger.SetOut(out)
}
func SetLoggerLevel(level log.Level) {
	globalLevel = level
	for _, logger := range loggers {
		logger.SetLevel(globalLevel)
	}
}
func SetLoggerFormat(format log.Formatter) {
	globalFormat = format
	for _, logger := range loggers {
		logger.SetFormat(globalFormat)
	}
}
func SetLoggerOut(out io.Writer) {
	globalOut = out
	for _, logger := range loggers {
		logger.SetOut(globalOut)
	}
}
func AddHook(hook log.Hook) {
	// log.AddHook(hook)
	for _, logger := range loggers {
		logger.AddHook(hook)
	}
}

type Logger struct {
	name string
	*log.Entry
}

func NewLogger(name string, fields ...log.Fields) *Logger {
	l := log.New().WithField("module", name)
	for _, f := range fields {
		l = l.WithFields(f)
	}
	l.Logger.Level = globalLevel
	l.Logger.Out = globalOut
	if globalFormat != nil {
		// l.Logger.Formatter = &log.JSONFormatter{}
		l.Logger.Formatter = globalFormat
	}
	logger := &Logger{name, l}
	loggers = append(loggers, logger)
	return logger
}

// Log is here to make compatible with gokit log
func (l *Logger) Log(keyvals ...interface{}) error {
	l.Infoln(keyvals...)
	return nil
}

// SetLevel sets log leve
func (l *Logger) SetLevel(level log.Level) {
	l.Logger.Level = level
}
func (l *Logger) SetFormat(format log.Formatter) {
	l.Logger.Formatter = format
}
func (l *Logger) SetOut(out io.Writer) {
	l.Logger.Out = out
}
func (l *Logger) AddHook(hook log.Hook) {
	l.Logger.Hooks.Add(hook)
}
