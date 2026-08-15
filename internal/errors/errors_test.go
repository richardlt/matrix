package errors

import (
	stderrors "errors"
	"fmt"
	"strings"
	"testing"
)

var errBase = stderrors.New("base")

var errSentinel = New("api is closed")

func TestNilHandling(t *testing.T) {
	if got := Details(nil); got != "<nil>" {
		t.Errorf("Details(nil) = %q", got)
	}
}

func TestSingleStackNotDuplicated(t *testing.T) {
	inner := Errorf("inner: %w", errBase)
	outer := Errorf("outer: %w", inner)
	deep := Errorf("deep: %w", outer)

	// exactly one frame set in the whole chain
	n := 0
	for e := error(deep); e != nil; e = stderrors.Unwrap(e) {
		if st, ok := e.(interface{ StackTrace() []uintptr }); ok && st.StackTrace() != nil {
			n++
		}
	}
	if n != 1 {
		t.Errorf("want 1 stack in the chain, got %d", n)
	}
	// and it points at the innermost call
	if !strings.Contains(Details(deep), "TestSingleStackNotDuplicated") {
		t.Error("stack should reach the test function")
	}
}

func TestStackComesFromOrigin(t *testing.T) {
	origin := func() error { return Errorf("origin: %w", errBase) }
	err := Errorf("wrapped: %w", origin())
	d := Details(err)
	// the closure frame must be present: the stack is the origin's, not the wrapper's
	if !strings.Contains(d, "func1") {
		t.Errorf("stack should be captured at the origin, got:\n%s", d)
	}
}

func TestStdlibInterop(t *testing.T) {
	err := Errorf("ctx: %w", errBase)
	if !stderrors.Is(err, errBase) {
		t.Error("errors.Is broken")
	}
	if stderrors.Unwrap(stderrors.Unwrap(err)) != errBase {
		t.Error("unwrap chain broken")
	}
	// still works when a plain fmt.Errorf sits on top
	if !stderrors.Is(fmt.Errorf("top: %w", err), errBase) {
		t.Error("errors.Is broken under fmt.Errorf")
	}
}

func TestVerbs(t *testing.T) {
	err := Errorf("ctx: %w", errBase)
	if got := fmt.Sprintf("%v", err); got != "ctx: base" {
		t.Errorf("%%v = %q", got)
	}
	if got := fmt.Sprintf("%s", err); got != "ctx: base" {
		t.Errorf("%%s = %q", got)
	}
	if got := fmt.Sprintf("%q", err); got != `"ctx: base"` {
		t.Errorf("%%q = %q", got)
	}
	if !strings.Contains(fmt.Sprintf("%+v", err), "TestVerbs") {
		t.Error("verbose format should include the stack")
	}
}

func TestErrorfWithoutWrapping(t *testing.T) {
	err := Errorf("plain %d", 42)
	if err.Error() != "plain 42" {
		t.Errorf("got %q", err.Error())
	}
	if !strings.Contains(Details(err), "TestErrorfWithoutWrapping") {
		t.Error("should still carry a stack")
	}
}

func TestReExports(t *testing.T) {
	target := stderrors.New("target")
	err := Errorf("ctx: %w", target)

	if !Is(err, target) {
		t.Error("Is re-export broken")
	}
	var se interface{ Error() string }
	if !As(err, &se) {
		t.Error("As re-export broken")
	}
	if Unwrap(Unwrap(err)) != target {
		t.Error("Unwrap re-export broken")
	}
	joined := Join(target, stderrors.New("other"))
	if !Is(joined, target) {
		t.Error("Join re-export broken")
	}
}

func TestNewCarriesNoStack(t *testing.T) {
	// New is for sentinels: a stack taken at init would be the init goroutine's.
	if stackOf(New("sentinel")) != nil {
		t.Error("New must not capture a stack")
	}
}

func TestSentinelDoesNotSuppressCallSiteStack(t *testing.T) {
	// Wrapping a sentinel must still record where the failure happened.
	wrapped := Errorf("sending request: %w", errSentinel)
	if !strings.Contains(Details(wrapped), "TestSentinelDoesNotSuppressCallSiteStack") {
		t.Errorf("stack should point at the failure site, got:\n%s", Details(wrapped))
	}
}
