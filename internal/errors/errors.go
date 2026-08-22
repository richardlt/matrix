// Package errors adds stack traces to the standard error wrapping idiom.
//
// Errorf mirrors fmt.Errorf, so %w, Is, As and Unwrap behave exactly as they do
// in the standard library. On top of that the returned error carries the stack
// of the call site, which %+v prints.
//
// A stack is captured once, by the deepest call that creates one. Wrapping an
// error that already carries a stack keeps the original, so %+v shows where the
// failure started rather than where it was last annotated. Every wrap therefore
// takes a message: it is the message chain, not the stack, that says what the
// code was doing and to which file, layer or software it was doing it.
//
// Is, As, Unwrap and Join are re-exported so this package fully replaces the
// standard errors package and callers never need to import both.
package errors

import (
	stderrors "errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strconv"
)

const maxDepth = 32

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target.
func As(err error, target any) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling Unwrap on err.
func Unwrap(err error) error { return stderrors.Unwrap(err) }

// Join wraps stderrors.Join, returning an error that is any of the given errors.
func Join(errs ...error) error { return stderrors.Join(errs...) }

// New returns a plain error carrying no stack, for package-level sentinels that
// callers match with Is. Sentinels are built once at init, so a stack taken here
// would record the init goroutine and, being the first stack in the chain, would
// stop Errorf recording the site the failure actually happened at. Use Errorf for
// any error a function returns.
func New(text string) error {
	return stderrors.New(text)
}

// Errorf formats an error exactly like fmt.Errorf, %w included, and records the
// stack of the call site. Wrapping an error that already carries a stack keeps
// the original.
func Errorf(format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	if stackOf(err) != nil {
		return &withStack{err: err}
	}
	return &withStack{err: err, stack: callers()}
}

type withStack struct {
	err   error
	stack []uintptr
}

func (w *withStack) Error() string         { return w.err.Error() }
func (w *withStack) Unwrap() error         { return w.err }
func (w *withStack) StackTrace() []uintptr { return w.stack }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, Details(w))
			return
		}
		_, _ = io.WriteString(s, w.Error())
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = io.WriteString(s, strconv.Quote(w.Error()))
	}
}

// Details renders err's message followed by its stack, if any error in the
// chain carries one. It works whatever wraps the error, so a log site can call
// it without knowing how the error was built.
func Details(err error) string {
	if err == nil {
		return "<nil>"
	}
	out := err.Error()
	pcs := stackOf(err)
	if pcs == nil {
		return out
	}
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		out += "\n\t" + f.Function + "\n\t\t" + filepath.Base(f.File) + ":" + strconv.Itoa(f.Line)
		if !more {
			break
		}
	}
	return out
}

// stackOf returns the stack of the first error in the chain that carries one.
func stackOf(err error) []uintptr {
	for err != nil {
		if st, ok := err.(interface{ StackTrace() []uintptr }); ok {
			if pcs := st.StackTrace(); pcs != nil {
				return pcs
			}
		}
		err = stderrors.Unwrap(err)
	}
	return nil
}

func callers() []uintptr {
	pcs := make([]uintptr, maxDepth)
	// skip runtime.Callers, callers and Errorf itself
	n := runtime.Callers(3, pcs)
	return pcs[:n]
}
