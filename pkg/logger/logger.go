package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger provides leveled logging.
type Logger struct {
	*zap.SugaredLogger
}

// Options holds the logger settings.
type Options struct {
	Level  string // minimum log level (debug/info/warn/error)
	Format string // log encoding format (console/json)
	Output string // output stream (stdout/stderr)
}

func New(opts Options) (*Logger, error) {
	level, err := opts.parseLogLevel()
	if err != nil {
		return nil, err
	}

	encoder, err := opts.buildEncoder()
	if err != nil {
		return nil, err
	}

	writer, err := opts.resolveWriter()
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writer, level)
	base := zap.New(core, zap.AddCaller())
	return &Logger{SugaredLogger: base.Sugar()}, nil
}

// parseLogLevel parses the level from the configured value.
func (o Options) parseLogLevel() (zapcore.Level, error) {
	switch strings.ToLower(o.Level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "", "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return 0, fmt.Errorf("[logger] invalid level %q: must be debug, info, warn or error", o.Level)
	}
}

// buildEncoder builds the encoder from the configured format.
func (o Options) buildEncoder() (zapcore.Encoder, error) {
	switch strings.ToLower(o.Format) {
	case "", "console":
		return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), nil
	case "json":
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), nil
	default:
		return nil, fmt.Errorf("[logger] invalid format %q: must be console or json", o.Format)
	}
}

// resolveWriter resolves the writer from the configured output.
func (o Options) resolveWriter() (zapcore.WriteSyncer, error) {
	switch strings.ToLower(o.Output) {
	case "", "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return nil, fmt.Errorf("[logger] invalid output %q: must be stdout or stderr", o.Output)
	}
}
