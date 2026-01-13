package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	colorGray   = "\x1b[90m"
	colorReset  = "\x1b[0m"
	jsonIndent  = "  "
	fieldModule = "module"
	fieldComp   = "component"
)

var (
	log        zerolog.Logger
	prettyMode bool
)

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

func Get() *zerolog.Logger       { return &log }
func Sub(module string) zerolog.Logger { return log.With().Str(fieldModule, module).Logger() }
func Writer() io.Writer          { return log }
func RawWriter() io.Writer       { return os.Stdout }

func PrettyJSON(v interface{}) string {
	if !prettyMode {
		b, _ := json.Marshal(v)
		return string(b)
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", jsonIndent)
	_ = enc.Encode(v)
	return colorGray + buf.String() + colorReset
}

func Component(name string) *zerolog.Event     { return log.Info().Str(fieldComp, name) }
func WarnComponent(name string) *zerolog.Event { return log.Warn().Str(fieldComp, name) }
func ErrorComponent(name string) *zerolog.Event { return log.Error().Str(fieldComp, name) }
func DebugComponent(name string) *zerolog.Event { return log.Debug().Str(fieldComp, name) }
func WithError(err error) *zerolog.Event       { return log.Error().Err(err) }

func Debug(msg string) { log.Debug().Msg(msg) }
func Info(msg string)  { log.Info().Msg(msg) }
func Warn(msg string)  { log.Warn().Msg(msg) }
func Error(msg string) { log.Error().Msg(msg) }

func Debugf(format string, v ...interface{}) { log.Debug().Msgf(format, v...) }
func Infof(format string, v ...interface{})  { log.Info().Msgf(format, v...) }
func Warnf(format string, v ...interface{})  { log.Warn().Msgf(format, v...) }
func Errorf(format string, v ...interface{}) { log.Error().Msgf(format, v...) }
