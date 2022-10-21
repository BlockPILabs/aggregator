package log

import (
	"context"
	"fmt"
	"github.com/inconshreveable/log15"
	"sync"
)

type Logger interface {
	log15.Logger
	Trace(msg string, ctx ...interface{})
	NewLogger(ctx ...interface{}) Logger
	SetLevel(lvl log15.Lvl) Logger
	SetLevelString(lvlString string) Logger
	NewContextLogger(c context.Context, ctx ...interface{}) (context.Context, Logger)
}

type loggerCtx struct {
}

type logger struct {
	log15.Logger
}

func (l *logger) GetHandler() log15.Handler {
	return l.Logger.GetHandler()
}

func (l *logger) SetHandler(h log15.Handler) {
	l.Logger.SetHandler(h)
}
func (l *logger) NewLogger(ctx ...interface{}) Logger {
	return &logger{Logger: l.Logger.New(ctx...)}
}
func (l *logger) NewContextLogger(c context.Context, ctx ...interface{}) (context.Context, Logger) {
	nl := l.NewLogger(ctx...)
	return WithContext(c, nl), nl
}
func (l *logger) SetLevel(lvl log15.Lvl) Logger {
	l.SetHandler(log15.LvlFilterHandler(lvl, l.GetHandler()))
	return l
}
func (l *logger) SetLevelString(lvlString string) Logger {
	lvl, err := log15.LvlFromString(lvlString)
	if err != nil {
		return l
	}
	return l.SetLevel(lvl)
}
func (l *logger) New(ctx ...interface{}) log15.Logger {
	return l.NewLogger(ctx...)
}

func (l *logger) Trace(msg string, ctx ...interface{}) {
	l.Logger.Error(msg, ctx...)
}

var root *logger
var moduleLogs sync.Map

func init() {
	root = &logger{Logger: log15.Root()}
	root.SetLevelString("debug")
}

func Module(module string) Logger {
	if module == "root" {
		return Root()
	}
	logI, ok := moduleLogs.Load(module)
	if !ok {
		log := newModule(module)
		moduleLogs.Store(module, log)
		return log
	}
	log, ok := logI.(Logger)
	if !ok {
		log = newModule(module)
		moduleLogs.Store(module, log)
		return log
	}
	return log
}

func Module15(module string) log15.Logger {
	return log15.New("module", module)
}

func newModule(module string) Logger {
	log := Root().NewLogger("module", module)
	return log
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...).(*logger)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	root.Debug(msg, ctx...)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	root.Info(msg, ctx...)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.Warn(msg, ctx...)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.Error(msg, ctx...)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	root.Crit(msg, ctx...)
}

// Trace is a convenient alias for Root().Trace
func Trace(msg string, ctx ...interface{}) {
	root.Trace(msg, ctx...)
}

func StackError(msg string, ctx error) {
	root.Error(fmt.Sprintf("%s   => error : %+v", msg, ctx))
}
func SetHandler(h log15.Handler) {
	root.SetHandler(h)
}

func GetHandler() log15.Handler {
	return root.GetHandler()
}

func SetLevel(lvl log15.Lvl) Logger {
	return Root().SetLevel(lvl)
}
func SetLevelString(lvlString string) Logger {
	return Root().SetLevelString(lvlString)
}

// WithContext context log
func WithContext(ctx context.Context, l Logger) context.Context {
	if l == nil {
		l = Root()
	}
	k := loggerCtx{}
	ctx = context.WithValue(ctx, k, l)
	return ctx
}
