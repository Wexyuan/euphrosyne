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
	// Level is the minimum log level (debug/info/warn/error, case-insensitive).
	Level string
	// Format is the log encoding format (console/json, case-insensitive).
	Format string
	// Output is the output stream (stdout/stderr, case-insensitive).
	Output string
}

// New creates the logger from the given options.
func New(opts Options) (*Logger, error) {
	level, err := opts.parseLogLevel()
	if err != nil {
		return nil, err
	}

	encoder, err := opts.buildEncoder()
	if err != nil {
		return nil, err
	}

	writer, err := opts.resolveOutput()
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writer, level)
	base := zap.New(core, zap.AddCaller())
	return &Logger{SugaredLogger: base.Sugar()}, nil
}

// parseLogLevel parses the configured level.
func (o Options) parseLogLevel() (zapcore.Level, error) {
	switch strings.ToLower(o.Level) {
	case "":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
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
	cfg := zap.NewDevelopmentEncoderConfig()
	switch strings.ToLower(o.Format) {
	case "":
		return zapcore.NewConsoleEncoder(cfg), nil
	case "console":
		return zapcore.NewConsoleEncoder(cfg), nil
	case "json":
		return zapcore.NewJSONEncoder(cfg), nil
	default:
		return nil, fmt.Errorf("[logger] invalid format %q: must be console or json", o.Format)
	}
}

// resolveOutput resolves the writer from the configured output.
func (o Options) resolveOutput() (zapcore.WriteSyncer, error) {
	switch strings.ToLower(o.Output) {
	case "":
		return os.Stdout, nil
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return nil, fmt.Errorf("[logger] invalid output %q: must be stdout or stderr", o.Output)
	}
}
