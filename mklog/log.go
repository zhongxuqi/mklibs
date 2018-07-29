package mklog

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	isatty "github.com/mattn/go-isatty"
	"github.com/zhongxuqi/mklibs/common"
)

type Level int

const (
	LevelDebug Level = 0
	LevelInfo  Level = 1
	LevelError Level = 2

	ContextLog = "mklog-instance"

	colorNone   = "\x1b[0m"
	colorRed    = "\x1b[1;31m"
	colorGreen  = "\x1b[1;32m"
	colorYellow = "\x1b[1;33m"
	colorBlue   = "\x1b[1;34m"
	colorPurple = "\x1b[1;35m"
)

var (
	LevelMap = map[Level]string{
		LevelDebug: "Debug",
		LevelInfo:  "Info",
		LevelError: "Error",
	}

	LevelColorMap = map[Level]string{
		LevelDebug: colorYellow,
		LevelInfo:  colorGreen,
		LevelError: colorRed,
	}
	NoColor = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

type Logger interface {
	SetOutput(writer io.Writer)
	SetLevel(lvl Level)
	Context() context.Context
	GetLogID() string
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

type logger struct {
	ctx    context.Context
	writer io.Writer // output io
	level  Level     // output level
	logID  string    // log id
}

func New() Logger {
	b := uuid.New()
	return &logger{
		ctx:   context.TODO(),
		level: LevelDebug,
		logID: base64.StdEncoding.EncodeToString(b[:]),
	}
}

func NewWithReq(req *http.Request) Logger {
	b := uuid.New()
	logID := req.Header.Get(common.HttpLogID)
	if logID == "" {
		logID = base64.URLEncoding.EncodeToString(b[:])
		req.Header.Set(common.HttpLogID, logID)
	}
	return &logger{
		ctx:   context.TODO(),
		level: LevelDebug,
		logID: logID,
	}
}

func NewWithContext(ctx context.Context) Logger {
	v := ctx.Value(ContextLog)
	if v != nil {
		if ml, ok := v.(Logger); ok {
			return ml
		}
	}
	b := uuid.New()
	return &logger{
		ctx:   context.TODO(),
		level: LevelDebug,
		logID: base64.StdEncoding.EncodeToString(b[:]),
	}
}

func (s *logger) getPrefix(level Level) string {
	_, file, line, _ := runtime.Caller(3)
	index := strings.LastIndex(file, "src")
	if index > 0 {
		index += 4
	}
	if NoColor {
		return fmt.Sprintf("%s[%s][%s]%s:%d:", time.Now().Format(time.RFC3339), LevelMap[level],
			s.logID, file[index:], line)
	}
	return fmt.Sprintf("%s%s%s[%s]%s[%s]%s%s:%d:%s", colorBlue, time.Now().Format(time.RFC3339), LevelColorMap[level], LevelMap[level],
		colorPurple, s.logID, colorBlue, file[index:], line, colorNone)
}

func (s *logger) SetOutput(writer io.Writer) {
	s.writer = writer
}

func (s *logger) SetLevel(lvl Level) {
	s.level = lvl
}

func (s *logger) Context() context.Context {
	return context.WithValue(s.ctx, ContextLog, s)
}

func (s *logger) GetLogID() string {
	return s.logID
}

func (s *logger) Debugf(format string, v ...interface{}) {
	s.writeLog(LevelDebug, format, v...)
}

func (s *logger) Infof(format string, v ...interface{}) {
	s.writeLog(LevelInfo, format, v...)
}

func (s *logger) Errorf(format string, v ...interface{}) {
	s.writeLog(LevelError, format, v...)
}

func (s *logger) writeLog(lvl Level, format string, v ...interface{}) {
	if lvl < s.level {
		return
	}
	out := s.writer
	if out == nil {
		out = os.Stdout
	}
	args := make([]interface{}, 0, len(v)+1)
	args = append(args, s.getPrefix(lvl))
	args = append(args, v...)
	fmt.Fprintf(out, "%s"+format+"\n", args...)
}
