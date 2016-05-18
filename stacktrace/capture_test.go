package stacktrace

import (
	"fmt"
	"testing"
)

func captureFunc() Stacktrace {
	return Capture(0)
}

func TestCaptureTestFunc(t *testing.T) {
	stack := captureFunc()

	if len(stack.Frames) == 0 {
		t.Fatal("expected stack frame to be returned")
	}

	// the first frame is the caller
	frame := stack.Frames[0]
	fmt.Println("TestCaptureTestFunc", frame)
	if expected := "captureFunc"; frame.Function != expected {
		t.Fatal("expected function %q but received %q", expected, frame.Function)
	}
	if expected := "stacktrace"; frame.Package != expected {
		t.Fatal("expected function %q but received %q", expected, frame.Package)
	}
	if expected := "capture_test.go"; frame.File != expected {
		t.Fatal("expected function %q but received %q", expected, frame.File)
	}
}
