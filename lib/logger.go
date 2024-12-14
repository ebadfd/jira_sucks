package lib

import (
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	APP_ENV   = "ENVIRONMENT"
	LOG_LEVEL = "LOG_LEVEL"
)

var globalLog *Logger
var zapLogger *zap.Logger

// Logger structure
type Logger struct {
	*zap.SugaredLogger
}

// FxLogger logger for go-fx [subbed from main logger]
type FxLogger struct {
	*Logger
}

// GetLogger gets the global instance of the logger
func GetLogger() Logger {
	if globalLog != nil {
		return *globalLog
	}
	globalLog := newLogger()
	return *globalLog
}

func newLogger() *Logger {
	env := os.Getenv(APP_ENV)
	logLevel := os.Getenv(LOG_LEVEL)
	var config zap.Config
	var level zapcore.Level

	if env == "production" {
		config = zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.InfoLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(zapcore.DebugLevel)
	}

	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zap.PanicLevel
	}
	config.Level.SetLevel(level)
	zapLogger, _ = config.Build()
	globalLog := zapLogger.Sugar()

	return &Logger{
		SugaredLogger: globalLog,
	}

}

func newSugaredLogger(logger *zap.Logger) *Logger {
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// GetFxLogger gets logger for go-fx
func (l *Logger) GetFxLogger() fxevent.Logger {
	logger := zapLogger.WithOptions(
		zap.WithCaller(false),
	)
	return &FxLogger{Logger: newSugaredLogger(logger)}
}

func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Debug("OnStart hook executing: ",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Debug("OnStart hook failed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Logger.Debug("OnStart hook executed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Debug("OnStop hook executing: ",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Debug("OnStop hook failed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Logger.Debug("OnStop hook executed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		l.Logger.Debug("supplied: ", zap.String("type", e.TypeName), zap.Error(e.Err))
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Debug("provided: ", e.ConstructorName, " => ", rtype)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Debug("decorated: ",
				zap.String("decorator", e.DecoratorName),
				zap.String("type", rtype),
			)
		}
	case *fxevent.Invoking:
		l.Logger.Debug("invoking: ", e.FunctionName)
	case *fxevent.Started:
		if e.Err == nil {
			l.Logger.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err == nil {
			l.Logger.Debug("initialized: custom fxevent.Logger -> ", e.ConstructorName)
		}
	}
}

// Printf prints go-fx logs
func (l FxLogger) Printf(str string, args ...any) {
	if len(args) > 0 {
		l.Debugf(str, args)
	}
	l.Debug(str)
}
