package zaplogger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

func capitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelCapitalString(l))
}

// CapitalString returns an all-caps ASCII representation of the log level.
func levelCapitalString(l zapcore.Level) string {
	// Printing levels in all-caps is common enough that we should export this
	// functionality.
	switch l {
	case zapcore.DebugLevel:
		return "[DEBUG]"
	case zapcore.InfoLevel:
		return "[INFO]"
	case zapcore.WarnLevel:
		return "[WARN]"
	case zapcore.ErrorLevel:
		return "[ERROR]"
	case zapcore.DPanicLevel:
		return "[DPANIC]"
	case zapcore.PanicLevel:
		return "[PANIC]"
	case zapcore.FatalLevel:
		return "[FATAL]"
	default:
		return fmt.Sprintf("[LEVEL(%d)]", l)
	}
}
