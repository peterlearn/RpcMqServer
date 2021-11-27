package amqprpc

//Logger is our contract for the logger
type Logger interface {
	// Debug uses fmt.Sprint to construct and log a message.
	Debug(args ...interface{})

	// Info uses fmt.Sprint to construct and log a message.
	Info(args ...interface{})

	// Warn uses fmt.Sprint to construct and log a message.
	Warn(args ...interface{})

	// Error uses fmt.Sprint to construct and log a message.
	Error(args ...interface{})

	// Panic uses fmt.Sprint to construct and log a message, then panics.
	Panic(args ...interface{})

	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...interface{})

	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})

	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})

	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})

	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})

	// Panicf uses fmt.Sprintf to log a templated message, then panics.
	Panicf(template string, args ...interface{})

	// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
	Fatalf(template string, args ...interface{})
}

type LogWrapper struct {
	log Logger
}

func NewLogWrapper(log Logger) *LogWrapper {
	return &LogWrapper{log}
}

func (l *LogWrapper) Debug(args ...interface{}) {
	if l.log != nil {
		l.log.Debug(args...)
	}
}

func (l *LogWrapper) Info(args ...interface{}) {
	if l.log != nil {
		l.log.Info(args...)
	}
}

func (l *LogWrapper) Warn(args ...interface{}) {
	if l.log != nil {
		l.log.Warn(args...)
	}
}

func (l *LogWrapper) Error(args ...interface{}) {
	if l.log != nil {
		l.log.Error(args...)
	}
}

func (l *LogWrapper) Panic(args ...interface{}) {
	if l.log != nil {
		l.log.Panic(args...)
	}
}

func (l *LogWrapper) Fatal(args ...interface{}) {
	if l.log != nil {
		l.log.Fatal(args...)
	}
}

func (l *LogWrapper) Debugf(template string, args ...interface{}) {
	if l.log != nil {
		l.log.Debugf(template, args...)
	}
}

func (l *LogWrapper) Infof(template string, args ...interface{}) {
	if l.log != nil {
		l.log.Infof(template, args...)
	}
}

func (l *LogWrapper) Warnf(template string, args ...interface{}) {
	if l.log != nil {
		l.log.Warnf(template, args...)
	}
}

func (l *LogWrapper) Errorf(template string, args ...interface{}) {
	if l.log != nil {
		l.log.Errorf(template, args...)
	}
}

func (l *LogWrapper) Panicf(template string, args ...interface{}) {
	if l.log != nil {
		l.log.Panicf(template, args...)
	}
}

func (l *LogWrapper) Fatalf(template string, args ...interface{}) {
	if l.log != nil {
		l.log.Fatalf(template, args...)
	}
}
