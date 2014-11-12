package logex

import (
	"encoding/json"
	"fmt"
	goLog "log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var DebugLevel = 1

type Logger struct {
	depth int
	reqid string
}

func NewLogger(l int) Logger {
	return Logger{l, ""}
}

var goLogStd = goLog.New(os.Stderr, "", goLog.LstdFlags)
var std = Logger{1, ""}
var Std = Logger{1, ""}
var (
	Println     = std.Println
	Infof       = std.Infof
	Info        = std.Info
	Debug       = std.Debug
	Debugf      = std.Debugf
	Trace       = std.Trace
	Tracef      = std.Tracef
	TraceEnter  = std.TraceEnter
	TraceEnterf = std.TraceEnterf
	TraceExit   = std.TraceExit
	TraceExitf  = std.TraceExitf
	Error       = std.Error
	Errorf      = std.Errorf
	Warn        = std.Warn
	Warnf       = std.Warnf
	PrintStack  = std.PrintStack
	Stack       = std.Stack
	Panic       = std.Panic
	Panicf      = std.Panicf
	Fatal       = std.Fatal
	Fatalf      = std.Fatalf
	Struct      = std.Struct
	Pretty      = std.Pretty
	Todo        = std.Todo
)

var (
	INFO   = "[\x1b[1;33mINFO\x1b[0m] "
	ERROR  = "[\x1b[0;35mERROR\x1b[0m] "
	PANIC  = "[PANIC] "
	DEBUG  = "[DEBUG] "
	TRACE  = "[TRACE] "
	WARN   = "[WARN] "
	FATAL  = "[FATAL] "
	STRUCT = "[STRUCT] "
	PRETTY = "[PRETTY] "
	TODO   = "[" + color("35", "TODO") + "] "
)

func color(col, s string) string {
	return "\x1b[0;" + col + "m" + s + "\x1b[0m"
}

func init() {
	if os.Getenv("DEBUG") != "" {
		DebugLevel = 0
	}
}

func D(i int) Logger {
	return std.D(i - 1)
}

func (l Logger) D(i int) Logger {
	return Logger{l.depth + i, l.reqid}
}

// Pretty ----------------------------------------------------------------------

func (l Logger) Pretty(os ...interface{}) {
	content := ""
	colors := []string{"31", "32", "33", "35"}
	for i, o := range os {
		if ret, err := json.MarshalIndent(o, "", "\t"); err == nil {
			content += color(colors[i%len(colors)], string(ret)) + "\n"
		}
	}
	l.Output(2, PRETTY, content)
}

// Print -----------------------------------------------------------------------

func (l Logger) Print(o ...interface{}) {
	l.Output(2, "", sprint(o))
}
func (l Logger) Printf(layout string, o ...interface{}) {
	l.Output(2, "", sprintf(layout, o))
}
func (l Logger) Println(o ...interface{}) {
	l.Output(2, "", sprint(o))
}

// Info ------------------------------------------------------------------------

func (l Logger) Info(o ...interface{}) {
	l.Output(2, INFO, sprint(o))
}
func (l Logger) Infof(f string, o ...interface{}) {
	l.Output(2, INFO, sprintf(f, o))
}

// Debug -----------------------------------------------------------------------

func (l Logger) Debug(o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	l.Output(2, DEBUG, sprint(o))
}
func (l Logger) Debugf(f string, o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	l.Output(2, DEBUG, sprintf(f, o))
}

// Trace -----------------------------------------------------------------------

func (l Logger) Trace(o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	l.Output(2, TRACE, sprint(o))
}

func (l Logger) Tracef(f string, o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	l.Output(2, TRACE, sprintf(f, o))
}

func (l Logger) TraceEnter(o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return
	}
	func_ := runtime.FuncForPC(pc)
	s := fmt.Sprintf("enter %s", func_.Name())
	s1 := sprint(o)
	if len(s1) > 0 {
		s += ", " + s1
	}
	l.Output(2, TRACE, s)
}

func (l Logger) TraceEnterf(f string, o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return
	}
	func_ := runtime.FuncForPC(pc)
	s := fmt.Sprintf("enter %s", func_.Name())
	s1 := sprintf(f, o)
	if len(s1) > 0 {
		s += ", " + s1
	}
	l.Output(2, TRACE, s)
}

func (l Logger) TraceExit(o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return
	}
	func_ := runtime.FuncForPC(pc)
	s := fmt.Sprintf("exit %s", func_.Name())
	s1 := sprint(o)
	if len(s1) > 0 {
		s += ", " + s1
	}
	l.Output(2, TRACE, s)
}

func (l Logger) TraceExitf(f string, o ...interface{}) {
	if DebugLevel > 0 {
		return
	}
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return
	}
	func_ := runtime.FuncForPC(pc)
	s := fmt.Sprintf("exit %s", func_.Name())
	s1 := sprintf(f, o)
	if len(s1) > 0 {
		s += ", " + s1
	}
	l.Output(2, TRACE, s)
}

// Todo ------------------------------------------------------------------------

func (l Logger) Todo(o ...interface{}) {
	l.Output(2, TODO, sprint(o))
}

// Error -----------------------------------------------------------------------

func (l Logger) Error(o ...interface{}) {
	l.Output(2, ERROR, sprint(o))
}
func (l Logger) Errorf(f string, o ...interface{}) {
	l.Output(2, ERROR, sprintf(f, o))
}

// Warn ------------------------------------------------------------------------

func (l Logger) Warn(o ...interface{}) {
	l.Output(2, WARN, sprint(o))
}
func (l Logger) Warnf(f string, o ...interface{}) {
	l.Output(2, WARN, sprintf(f, o))
}

// Panic -----------------------------------------------------------------------

func (l Logger) Panic(o ...interface{}) {
	l.Output(2, PANIC, sprint(o))
	panic(o)
}
func (l Logger) Panicf(f string, o ...interface{}) {
	info := sprintf(f, o)
	l.Output(2, PANIC, info)
	panic(info)
}

// Fatal -----------------------------------------------------------------------

func (l Logger) Fatal(o ...interface{}) {
	l.Output(2, FATAL, sprint(o))
	os.Exit(1)
}
func (l Logger) Fatalf(f string, o ...interface{}) {
	l.Output(2, FATAL, sprintf(f, o))
	os.Exit(1)
}

// Struct ----------------------------------------------------------------------

func (l Logger) Struct(o ...interface{}) {
	items := make([]interface{}, 0, len(o)*2)
	for _, item := range o {
		items = append(items, item, item)
	}
	layout := strings.Repeat(", %T(%+v)", len(o))
	if len(layout) > 0 {
		layout = layout[2:]
	}
	l.Output(2, STRUCT, sprintf(layout, items))
}

// Stack -----------------------------------------------------------------------

func (l Logger) PrintStack() {
	Info(string(l.Stack()))
}

func (l Logger) Stack() []byte {
	a := make([]byte, 1024*1024)
	n := runtime.Stack(a, true)
	return a[:n]
}

func (l Logger) Output(calldepth int, level, msg string) error {
	calldepth += l.depth + 1
	return goLogStd.Output(calldepth, level+msg+" \x1b[0;37m"+l.makePrefix(calldepth)+"\x1b[0m")
}

func (l Logger) makePrefix(calldepth int) string {
	pc, f, line, _ := runtime.Caller(calldepth)
	name := runtime.FuncForPC(pc).Name()
	name = path.Base(name) // only use package.funcname
	f = path.Base(f)       // only use filename

	tags := make([]string, 0, 3)

	pos := name + ":" + f + ":" + strconv.Itoa(line)
	tags = append(tags, pos)
	if l.reqid != "" {
		tags = append(tags, l.reqid)
	}
	return "[" + strings.Join(tags, "][") + "]"
}

func sprint(o []interface{}) string {
	decodeTrackError(o)
	return joinInterface(o, " ")
}
func sprintf(f string, o []interface{}) string {
	decodeTrackError(o)
	return fmt.Sprintf(f, o...)
}

func decodeTrackError(o []interface{}) {
	for idx, obj := range o {
		if te, ok := obj.(*TrackError); ok {
			o[idx] = te.StackError()
		}
	}
}
