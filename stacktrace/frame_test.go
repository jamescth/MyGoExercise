package stacktrace

import (
	"fmt"
	"testing"
)

func TestParsePackageName(t *testing.T) {
	var (
		name             = "stacktrace.captureFunc"
		expectedPackage  = "stacktrace"
		expectedFunction = "captureFunc"
	)

	pack, funcName := parseFunctionName(name)
	fmt.Println("TestParsePackageName", pack, funcName)

	if pack != expectedPackage {
		t.Fatal("expected package %q but received %q", expectedPackage, pack)
	}
	if funcName != expectedFunction {
		t.Fatal("expected function %q but received %q", expectedFunction, funcName)
	}
}
