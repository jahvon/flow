package io

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"

	"github.com/jahvon/flow/internal/io/styles"
)

type Format int

const (
	HumanReadable Format = iota
	Storage
)

type Logger struct {
	handler *log.Logger
	format  Format

	background  bool
	pendingRead bool
	data        string
	mu          *sync.RWMutex
}

func NewLogger(format Format, background bool) *Logger {
	logger := &Logger{
		format:     format,
		background: background,
		mu:         &sync.RWMutex{},
	}

	handler := log.NewWithOptions(logger, log.Options{ReportTimestamp: true, ReportCaller: false})
	switch format {
	case HumanReadable:
		handler.SetLevel(log.InfoLevel)
		applyHumanReadableFormat(handler)
	case Storage:
		handler.SetLevel(log.DebugLevel)
		applyStorageFormat(handler)
	default:
		panic("unknown format")
	}
	logger.handler = handler
	return logger
}

func applyHumanReadableFormat(handler *log.Logger) {
	handler.SetFormatter(log.TextFormatter)
	handler.SetTimeFormat(time.Kitchen)
	handler.SetStyles(styles.LoggerStyles())
	handler.SetColorProfile(termenv.ColorProfile())
}

func applyStorageFormat(handler *log.Logger) {
	handler.SetFormatter(log.JSONFormatter)
	handler.SetTimeFormat(time.RFC822)
	handler.SetStyles(log.DefaultStyles())
}

func (l *Logger) SetBackground(v bool) {
	l.background = v
}

func (l *Logger) ReadAllData() string {
	if !l.background {
		return ""
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.pendingRead = false
	return l.data
}

// PendingRead returns true if there is pending data to be read.
// This is useful for background loggers that need to be read
// in a loop before exiting.
func (l *Logger) PendingRead() bool {
	if !l.background {
		return false
	}
	return l.pendingRead
}

func (l *Logger) WriteStr(data string) {
	if data == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.background {
		_, err := os.Stdout.Write([]byte(data))
		if err != nil {
			panic(err)
		}
		return
	}
	if l.data == "" {
		l.data = data
		return
	}
	l.data += strings.Join([]string{l.data, data}, "\n")
	l.pendingRead = true
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.background {
		return os.Stdout.Write(p)
	}
	if l.data == "" {
		l.data = string(p)
		return len(p), nil
	}
	l.data += string(p)
	l.pendingRead = true
	return len(p), nil
}

func (l *Logger) SetLevel(level int) {
	switch level {
	case -3:
		l.handler.SetLevel(log.FatalLevel)
	case -2:
		l.handler.SetLevel(log.ErrorLevel)
	case -1:
		l.handler.SetLevel(log.WarnLevel)
	case 0:
		l.handler.SetLevel(log.InfoLevel)
	case 1:
		l.handler.SetLevel(log.DebugLevel)
	default:
		l.handler.SetLevel(log.InfoLevel)
	}
}

func (l *Logger) AsPlainText(exec func()) {
	l.handler.SetFormatter(log.TextFormatter)
	exec()
	l.handler.SetFormatter(log.LogfmtFormatter)
}

func (l *Logger) Infof(msg string, args ...any) {
	l.handler.Infof(msg, args...)
	l.pendingRead = true
}

func (l *Logger) Debugf(msg string, args ...any) {
	l.handler.Debugf(msg, args...)
	l.pendingRead = true
}

func (l *Logger) Error(err error, msg string) {
	if msg == "" {
		l.Errorf(err.Error())
		return
	}
	l.Errorx(err.Error(), "err", err)
	l.pendingRead = true
}

func (l *Logger) Errorf(msg string, args ...any) {
	l.handler.Errorf(msg, args...)
	l.pendingRead = true
}

func (l *Logger) Warnf(msg string, args ...any) {
	l.handler.Warnf(msg, args...)
	l.pendingRead = true
}

func (l *Logger) FatalErr(err error) {
	l.Fatalf(err.Error())
	l.pendingRead = true
}

func (l *Logger) Fatalf(msg string, args ...any) {
	if l.format == HumanReadable {
		l.handler.Fatalf(msg, args...)
		return
	}
	l.handler.Errorf(msg, args...)
	l.pendingRead = true
}

func (l *Logger) Infox(msg string, kv ...any) {
	l.handler.Info(msg, kv...)
	l.pendingRead = true
}

func (l *Logger) Debugx(msg string, kv ...any) {
	l.handler.Debug(msg, kv...)
	l.pendingRead = true
}

func (l *Logger) Errorx(msg string, kv ...any) {
	l.handler.Error(msg, kv...)
	l.pendingRead = true
}

func (l *Logger) Warnx(msg string, kv ...any) {
	l.handler.Warn(msg, kv...)
	l.pendingRead = true
}

func (l *Logger) Fatalx(msg string, kv ...any) {
	if l.format == HumanReadable {
		l.handler.Fatal(msg, kv...)
		return
	}
	l.handler.Error(msg, kv...)
	l.pendingRead = true
}

func (l *Logger) PlainTextInfo(msg string) {
	if l.format == HumanReadable {
		_, _ = fmt.Fprintln(l, styles.RenderInfo(msg))
		return
	}
	_, _ = fmt.Fprintln(l, msg)
}

func (l *Logger) PlainTextSuccess(msg string) {
	if l.format == HumanReadable {
		_, _ = fmt.Fprintln(l, styles.RenderSuccess(msg))
		return
	}
	_, _ = fmt.Fprintln(l, msg)
}
