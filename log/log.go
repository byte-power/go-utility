package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	fieldLoggerKey = "logger"

	fieldTraceID       = "trace_id"
	fieldTraceDuration = "trace_dur"
)

type Logger struct {
	level   Level
	trace   Trace
	outputs []Output
}

type Output interface {
	Level() Level
	LogModuleAndPairs(l Level, subject string, pairs []LogPair)
}

func NewLogger(outputs ...Output) Logger {
	return Logger{outputs: outputs, level: LevelDebug}
}

func NewLoggerWithTrace(logger *Logger, traceTime time.Time, tracePairs ...LogPair) Logger {
	traceID := NewTraceID().String()
	trace := NewTrace(traceID, traceTime, tracePairs...)
	return logger.WithTrace(trace)
}

// DO NOT USE IN SERVER
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func MakeConsoleOutput(name string, fmt LocalFormat, level Level, stream ConsoleStream) Output {
	writer := newZapConsoleWriter(stream.stream())
	return newZapLogger(name, fmt, level, writer)
}

func MakeFileOutput(name string, fmt LocalFormat, level Level, location string, rotation FileRotation) Output {
	writer := newZapFileWriter(location, rotation)
	return newZapLogger(name, fmt, level, writer)
}

func WithContext(ctx context.Context) Logger {
	l, _ := ctx.Value(fieldLoggerKey).(Logger)
	return l
}

func (l Logger) WrapContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, fieldLoggerKey, l)
}

func (l Logger) WithTrace(trace Trace) Logger {
	l.trace = l.trace.merge(trace)
	return l
}

func (l Logger) WithTraceLogs(pairs ...LogPair) Logger {
	if len(pairs) > 0 {
		l.trace = l.trace.merge(Trace{pairs: pairs})
	}
	return l
}

// 得到有哪些output可输出指定level的log
func (l Logger) marchLevelOutputs(level Level) []Output {
	outputs := make([]Output, 0, len(l.outputs))
	for _, it := range l.outputs {
		if level >= it.Level() {
			outputs = append(outputs, it)
		}
	}
	return outputs
}

// 产生需要log的数据
// 如有起始时间，将包含起始时间至此时的时长
// 如已有pairs与参数pairs之间有重复的key，它们的值将会被合并
func (l Logger) producePairs(pairs []LogPair) []LogPair {
	if l.trace.isEmpty() {
		return pairs
	}
	dedup := make(map[string]any, l.trace.length()+len(pairs)+1)
	for _, v := range l.trace.pairs {
		dedup[v.key] = v.value
	}
	if !l.trace.startTime.IsZero() {
		dedup[fieldTraceDuration] = time.Since(l.trace.startTime).String()
	}
	for _, v := range pairs {
		exists := dedup[v.key]
		if exists != nil {
			dedup[v.key] = fmt.Sprintf("%v,%v", exists, v.value)
		} else {
			dedup[v.key] = v.value
		}
	}
	toLog := make([]LogPair, 0, len(dedup))
	for k, v := range dedup {
		toLog = append(toLog, LogPair{key: k, value: v})
	}
	return toLog
}

func (l Logger) logPairs(level Level, subject string, pairs []LogPair) {
	if !(level >= l.level) {
		return
	}
	toOutputs := l.marchLevelOutputs(level)
	if len(toOutputs) == 0 {
		return
	}
	toLog := l.producePairs(pairs)
	for _, it := range toOutputs {
		it.LogModuleAndPairs(level, subject, toLog)
	}
}

func (l Logger) Debug(subject string, pairs ...LogPair) {
	l.logPairs(LevelDebug, subject, pairs)
}

func (l Logger) Info(subject string, pairs ...LogPair) {
	l.logPairs(LevelInfo, subject, pairs)
}

func (l Logger) Warn(subject string, pairs ...LogPair) {
	l.logPairs(LevelWarn, subject, pairs)
}

func (l Logger) Error(subject string, pairs ...LogPair) {
	l.logPairs(LevelError, subject, pairs)
}

func (l Logger) CanOutput(level Level) bool {
	for _, it := range l.outputs {
		if level >= it.Level() {
			return true
		}
	}
	return false
}

func (l Logger) IsEmpty() bool { return len(l.outputs) == 0 }

type LogPair struct {
	key   string
	value interface{}
}

func (p LogPair) String() string {
	return fmt.Sprintf("%s: %v", p.key, p.value)
}

func Any(k string, v interface{}) LogPair {
	return LogPair{key: k, value: v}
}

func String(k, v string) LogPair {
	return LogPair{key: k, value: v}
}

func JsonString(k string, v any) LogPair {
	vBytes, _ := json.Marshal(v)
	return LogPair{key: k, value: string(vBytes)}
}

// func Number[T utility.Number](k string, v T) LogPair {
// 	return LogPair{key: k, value: v}
// }

func Error(err error) LogPair {
	return LogPair{key: "error", value: err}
}

func Stack(stack []byte) LogPair {
	return LogPair{key: "stack", value: string(stack)}
}

func ConvertStrMapToLogPairs(values map[string]interface{}) []LogPair {
	pairs := make([]LogPair, 0, len(values))
	for key, value := range values {
		pairs = append(pairs, LogPair{key: key, value: value})
	}
	return pairs
}
