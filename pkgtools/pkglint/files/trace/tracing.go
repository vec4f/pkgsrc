// Package trace traces function calls and steps in-between.
package trace

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
)

// Tracer produces a hierarchical structure of log events.
type Tracer struct {
	Tracing bool
	Out     io.Writer
	depth   int
}

// Result marks an argument of Tracer.Call as a result of that function.
// It is only logged when exiting from the function but not when entering it.
type Result struct {
	pointer interface{}
}

func (t *Tracer) Stepf(format string, args ...interface{}) {
	if t.Tracing {
		msg := fmt.Sprintf(format, args...)
		_, _ = io.WriteString(t.Out, fmt.Sprintf("TRACE: %s  %s\n", t.traceIndent(), msg))
	}
}

func (t *Tracer) Step1(format string, arg0 string) {
	t.Stepf(format, arg0)
}

func (t *Tracer) Step2(format string, arg0, arg1 string) {
	t.Stepf(format, arg0, arg1)
}

func (t *Tracer) Call0() func() {
	return t.traceCall()
}

func (t *Tracer) Call1(arg1 string) func() {
	return t.traceCall(arg1)
}

func (t *Tracer) Call2(arg1, arg2 string) func() {
	return t.traceCall(arg1, arg2)
}

// Call records a function call in the tracing log, both when entering and
// when leaving the function.
//
// Usage:
//  if trace.Tracing {
//      defer trace.Call(arg1, arg2, trace.Result(result1), trace.Result(result2))()
// }
func (t *Tracer) Call(args ...interface{}) func() {
	return t.traceCall(args...)
}

// http://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
func isNil(a interface{}) bool {
	defer func() {
		recover()
	}()
	return a == nil || reflect.ValueOf(a).IsNil()
}

func argsStr(args []interface{}) string {
	rv := ""
	for _, arg := range args {
		if rv != "" {
			rv += ", "
		}
		if str, ok := arg.(fmt.Stringer); ok && !isNil(str) {
			rv += str.String()
		} else {
			rv += fmt.Sprintf("%#v", arg)
		}
	}
	return rv
}

func (t *Tracer) traceIndent() string {
	var indent strings.Builder
	for i := 0; i < t.depth; i++ {
		_, _ = fmt.Fprintf(&indent, "%d ", i+1)
	}
	return indent.String()
}

func (t *Tracer) traceCall(args ...interface{}) func() {
	if !t.Tracing {
		panic("Internal pkglint error: calls to trace.Call must only occur in tracing mode")
	}

	funcname := "?"
	if pc, _, _, ok := runtime.Caller(2); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcname = strings.TrimPrefix(fn.Name(), "netbsd.org/pkglint.")
		}
	}
	indent := t.traceIndent()
	_, _ = fmt.Fprintf(t.Out, "TRACE: %s+ %s(%s)\n", indent, funcname, argsStr(withoutResults(args)))
	t.depth++

	return func() {
		t.depth--
		_, _ = fmt.Fprintf(t.Out, "TRACE: %s- %s(%s)\n", indent, funcname, argsStr(withResults(args)))
	}
}

// Result marks an argument as a result and is only logged when the function returns.
func (t *Tracer) Result(rv interface{}) Result {
	if reflect.ValueOf(rv).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("Result must be called with a pointer to the result, not %#v.", rv))
	}
	return Result{rv}
}

func withoutResults(args []interface{}) []interface{} {
	for i, arg := range args {
		if _, ok := arg.(Result); ok {
			return args[0:i]
		}
	}
	return args
}

func withResults(args []interface{}) []interface{} {
	for i, arg := range args {
		if _, ok := arg.(Result); ok {
			var awr []interface{}
			awr = append(awr, args[0:i]...)
			awr = append(awr, "=>")
			for _, res := range args[i:] {
				pointer := reflect.ValueOf(res.(Result).pointer)
				actual := reflect.Indirect(pointer).Interface()
				awr = append(awr, actual)
			}
			return awr
		}
	}
	return args
}
