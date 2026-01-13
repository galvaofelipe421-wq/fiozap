package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// ANSI color codes
const (
	colorGray  = "\x1b[90m"
	colorReset = "\x1b[0m"
)

var (
	log        zerolog.Logger
	prettyMode bool
)

// Init initializes the global logger with the specified level and format
func Init(level string, pretty bool) {
	prettyMode = pretty

	var output io.Writer = os.Stdout
	if pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	log = zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Logger()
}

// Get returns the global logger instance
func Get() *zerolog.Logger {
	return &log
}

// Sub creates a sublogger with a module field
func Sub(module string) zerolog.Logger {
	return log.With().Str("module", module).Logger()
}

// Writer returns the logger as an io.Writer
func Writer() io.Writer {
	return log
}

// RawWriter returns os.Stdout for raw output (like QR codes)
func RawWriter() io.Writer {
	return os.Stdout
}

// PrettyJSON formats data as JSON with optional pretty printing and color
func PrettyJSON(v interface{}) string {
	if !prettyMode {
		b, _ := json.Marshal(v)
		return string(b)
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
	return colorGray + buf.String() + colorReset
}

// Component returns an info event with component field
func Component(name string) *zerolog.Event {
	return log.Info().Str("component", name)
}

// WarnComponent returns a warn event with component field
func WarnComponent(name string) *zerolog.Event {
	return log.Warn().Str("component", name)
}

// ErrorComponent returns an error event with component field
func ErrorComponent(name string) *zerolog.Event {
	return log.Error().Str("component", name)
}

// WithError returns an error event with the error attached
func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

// Legacy simple logging
func Info(msg string)  { log.Info().Msg(msg) }
func Warn(msg string)  { log.Warn().Msg(msg) }
func Error(msg string) { log.Error().Msg(msg) }

// Legacy formatted logging
func Infof(format string, v ...interface{})  { log.Info().Msgf(format, v...) }
func Warnf(format string, v ...interface{})  { log.Warn().Msgf(format, v...) }
func Errorf(format string, v ...interface{}) { log.Error().Msgf(format, v...) }
func Debugf(format string, v ...interface{}) { log.Debug().Msgf(format, v...) }
